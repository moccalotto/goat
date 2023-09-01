package motor

import (
	"fmt"
	u "goat/util"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// ||
// || Textured Quad Renderer
// ||
// || Render a quad with a texture
// || and map into subtextures.
// ||=============================
type TexQuadRenderer struct {
	Shader  *u.ShaderProgram
	Texture *u.TextureWrapper
	Atlas   *u.AtlasDescriptor // may be nil

	// Uniform variables to send to the shader
	UniColor     mgl32.Vec4
	UniSubTexPos mgl32.Vec4
	UniColorMix  float32

	// Buffer initialization stuff
	buffersReady bool
	vaoHandle    uint32
	bufferHandle uint32
}

// ||
// || Create a textured quad that contains a spritesheet / texture atlas
// ||
// || shaderFileBasename: filename shader files. It is assumed that all shaders have the same base fílename,
// || and that only the extension changes such that: frag shaders end in .frag, vert shaders in .vert, etc.
// ||
// || atlas: filename of the texture atlas (xml file)
// || subTexName the "filename" of the subtexture, as described in the texture atlas descriptor
// ||
// ||  verts the vertices of the quad
// ||  texCoords the overall texture coordinates of the entire sheet. should be [0, 0, 1, 1] because it will be modified by the coordinates of the subtextures
// ||  indexes: indeces used in the element array
// || ==================================================================
func CreateTexAtlasRenderer(shaderFileBasename, atlas, subTexName string) *TexQuadRenderer {
	shader, err := Machine.GetShader(shaderFileBasename)
	u.GlPanicIfErrNotNil(err)

	atlasDescriptor := Machine.LoadTextureAtlas(atlas)

	subTexInfo := atlasDescriptor.GetSubTexture(subTexName)

	w, h := atlasDescriptor.Texture.GetSize()

	s := TexQuadRenderer{
		Shader:       shader,
		Texture:      atlasDescriptor.Texture,
		Atlas:        atlasDescriptor,
		UniColor:     [4]float32{},
		UniSubTexPos: subTexInfo.GetDims(float32(w), float32(h)),
		UniColorMix:  0,
		buffersReady: false,
		vaoHandle:    0,
	}

	return &s
}

// ||
// || Create a Texture Quad for a non-atlassed texture
// ||
// ||
// || ===================================================
func CreateTexQuadRenderer(shaderAlias, textureAlias string) *TexQuadRenderer {

	tex, err := Machine.GetTexture(textureAlias)
	u.GlPanicIfErrNotNil(err)

	shader, err := Machine.GetShader(shaderAlias)
	u.GlPanicIfErrNotNil(err)

	T := TexQuadRenderer{
		Shader:       shader,
		Texture:      tex,
		UniColor:     mgl32.Vec4{1, 1, 1, 1},
		UniColorMix:  0.5, // mix tex and unicolor equally. good for debugging
		UniSubTexPos: mgl32.Vec4{1, 1, 1, 1},
		buffersReady: false,
	}

	return &T
}

func (R *TexQuadRenderer) Finalize() {
	R.Shader.Use()

	R.Texture.Finalize()

	R.Shader.SetUniformAttr("uniColorMix", R.UniColorMix)
	R.Shader.SetUniformAttr("uniColor", R.UniColor)
	R.Shader.SetUniformAttr("uniSubTexPos", R.UniSubTexPos)

	if R.buffersReady {
		return
	}

	// vertex and texture buffers are interleaved
	const vt_bytes_pr_stride = u.F32_SIZE * 5   // 4 bytes per float and 5 floats per stride/vert
	const vt_floats_total = 4 * 5               // 4 strides and 5 floats per stride
	const vt_len = u.F32_SIZE * vt_floats_total // total length (in bytes) of the VT buffer
	const Z, HI, LO = 1.0, 0.5, -0.5            // convenience
	vt_buffer := [vt_floats_total]float32{
		HI, HI, Z /* <== Vert || tex ==> */, 1, 1,
		LO, HI, Z /* <== Vert || tex ==> */, 0, 1,
		LO, LO, Z /* <== Vert || tex ==> */, 0, 0,
		HI, LO, Z /* <== Vert || tex ==> */, 1, 0,
	}
	vt_ptr := u.GlPtr32f(&vt_buffer[0])

	// Create buffer
	gl.GenBuffers(1, &R.bufferHandle)

	//
	// Vertex Array Object
	gl.GenVertexArrays(1, &R.vaoHandle)
	gl.BindVertexArray(R.vaoHandle)

	//
	// Vertex Buffer Object
	gl.BindBuffer(gl.ARRAY_BUFFER, R.bufferHandle)
	gl.BufferData(gl.ARRAY_BUFFER, vt_len, vt_ptr, gl.STATIC_DRAW)

	R.Shader.EnableVertexAttribArray("iVert")
	R.Shader.VertexAttribPointer("iVert", 3, gl.FLOAT, false, vt_bytes_pr_stride, 0)
	defer R.Shader.DisableVertexAttribArray("iVert")

	//
	// Texture Buffers
	R.Texture.Bind()
	R.Shader.EnableVertexAttribArray("iTexCoord")
	R.Shader.VertexAttribPointer("iTexCoord", 2, gl.FLOAT, false, vt_bytes_pr_stride, 3*u.F32_SIZE)
	defer R.Shader.DisableVertexAttribArray("iTexCoord")

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	R.Texture.Unbind()
}

func (R *TexQuadRenderer) MustUseSubTextureByName(name string) {
	if R.Atlas == nil {
		u.GlPanic(fmt.Errorf("this texture quad does not have subtextures"))
	}
}

func (R *TexQuadRenderer) Draw(camMatrix, objTranslationMatrix mgl32.Mat3) {

	R.Shader.Use()
	gl.BindVertexArray(R.vaoHandle)
	R.Texture.Bind()

	trMatrix := camMatrix.Mul3(objTranslationMatrix)

	u.GlPanicIfErrNotNil(R.Shader.SetUniformAttr("uniColor", R.UniColor))
	u.GlPanicIfErrNotNil(R.Shader.SetUniformAttr("uniColorMix", R.UniColorMix))
	u.GlPanicIfErrNotNil(R.Shader.SetUniformAttr("uniSubTexPos", R.UniSubTexPos))
	u.GlPanicIfErrNotNil(R.Shader.SetUniformAttr("uniTransformation", trMatrix))

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

	u.AssertGLOK("SpriteRenderable.Draw", R.Shader, 22)
}

func (R *TexQuadRenderer) Clone() *TexQuadRenderer {

	return &TexQuadRenderer{
		Shader:       R.Shader,
		Texture:      R.Texture,
		UniColor:     R.UniColor,
		UniSubTexPos: R.UniSubTexPos,
		UniColorMix:  R.UniColorMix,
		buffersReady: R.buffersReady,
		vaoHandle:    R.vaoHandle,
		bufferHandle: R.bufferHandle,
	}
}
