package motor

import (
	h "goat/glhelp"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Renderable interface {
	// GetShader() *h.ShaderProgram
	// GetVerts() []float32
	// GetTexCoords() []float32
	// GetUniVariables() map[string]any

	Prepare()
}

type SpriteRenderable struct {
	Shader  *h.ShaderProgram
	Texture *h.Texture

	// Geometry
	Verts     []float32
	TexCoords []float32
	Indeces   []uint32

	// Uniform variables to send to the shader
	UniColor     mgl32.Vec4
	UniSubTexPos mgl32.Vec4
	UniColorMix  float32

	// Override colors in wireframe mode
	Wireframe      bool
	WireframeColor mgl32.Vec4

	// Buffer initialization stuff
	buffersReady  bool
	vaoHandle     uint32
	bufferHandles []uint32
}

func CreateSprite(shaderAlias, textureAlias string, verts, texCoords []float32, indeces []uint32) *SpriteRenderable {

	sprite := SpriteRenderable{
		Shader:         Machine.Shaders[shaderAlias],
		Texture:        Machine.Textures[textureAlias],
		Verts:          verts,
		TexCoords:      texCoords,
		Indeces:        indeces,
		UniColor:       mgl32.Vec4{1, 1, 1, 1},
		UniColorMix:    0.5, // mix tex and unicolor equally. good for debugging
		UniSubTexPos:   mgl32.Vec4{1, 1, 1, 1},
		Wireframe:      false,
		WireframeColor: mgl32.Vec4{1, 1, 1, 1},
		buffersReady:   false,
	}

	return &sprite
}

func (R *SpriteRenderable) Finalize() {
	R.Shader.Use()

	R.Texture.Finalize()

	R.Shader.SetUniformAttr("uniColorMix", R.UniColorMix)
	R.Shader.SetUniformAttr("uniColor", R.UniColor)
	R.Shader.SetUniformAttr("uniSubTexPos", R.UniSubTexPos)

	if R.buffersReady {
		return
	}

	// Create out buffers
	// VBO: 0	vertex
	// EBO: 1	element
	// TBO: 2	texture
	R.bufferHandles = make([]uint32, 3)
	gl.GenBuffers(int32(len(R.bufferHandles)), &R.bufferHandles[0])

	//
	// Vertex Array Object
	gl.GenVertexArrays(1, &R.vaoHandle)
	gl.BindVertexArray(R.vaoHandle)

	//
	// Vertex Buffer Object
	gl.BindBuffer(gl.ARRAY_BUFFER, R.bufferHandles[0])
	gl.BufferData(gl.ARRAY_BUFFER, h.Float32Size*len(R.Verts), gl.Ptr(R.Verts), gl.STATIC_DRAW)

	R.Shader.EnableVertexAttribArray("iVert")
	R.Shader.VertexAttribPointer("iVert", 3, gl.FLOAT, false, 0, nil)
	defer R.Shader.DisableVertexAttribArray("iVert")

	//
	// Element/Index Buffer Object
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, R.bufferHandles[1])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, h.Int32Size*len(R.Indeces), gl.Ptr(R.Indeces), gl.STATIC_DRAW)

	//
	// Texture Buffers
	R.Texture.Bind()
	if R.Shader.HasAttrib("iTexCoord") {
		h.AssertGLOK("102")
		gl.BindBuffer(gl.ARRAY_BUFFER, R.bufferHandles[2])
		h.AssertGLOK("104")
		gl.BufferData(gl.ARRAY_BUFFER, h.Float32Size*len(R.TexCoords), gl.Ptr(R.TexCoords), gl.STATIC_DRAW)
		h.AssertGLOK("106")
		R.Shader.EnableVertexAttribArray("iTexCoord")
		h.AssertGLOK("108")
		R.Shader.VertexAttribPointer("iTexCoord", 2, gl.FLOAT, false, 0, nil)
		defer R.Shader.DisableVertexAttribArray("iTexCoord")
	}

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	R.Texture.Unbind()
}

func (R *SpriteRenderable) Draw(camMatrix, objTranslationMatrix mgl32.Mat3) {

	R.Shader.Use()
	gl.BindVertexArray(R.vaoHandle)
	R.Texture.Bind()

	trMatrix := h.MatMulMany(
		camMatrix,
		objTranslationMatrix,
	)
	if R.Wireframe {
		h.GlPanicIfErrNotNil(R.Shader.SetUniformAttr("uniColor", R.WireframeColor))
		h.GlPanicIfErrNotNil(R.Shader.SetUniformAttr("uniColorMix", float32(1)))
	} else {
		h.GlPanicIfErrNotNil(R.Shader.SetUniformAttr("uniColor", R.UniColor))
		h.GlPanicIfErrNotNil(R.Shader.SetUniformAttr("uniColorMix", R.UniColorMix))
	}

	R.Shader.SetUniformAttr("uniSubTexPos", R.UniSubTexPos)
	h.GlPanicIfErrNotNil(R.Shader.SetUniformAttr("uTransformation", trMatrix))

	gl.DrawElements(gl.TRIANGLES, int32(len(R.Indeces)), gl.UNSIGNED_INT, nil)

	h.AssertGLOK("SpriteRenderable.Draw", R.Shader, 22)
}
