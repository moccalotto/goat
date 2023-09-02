package tractor

import (
	"fmt"
	u "goat/shed"

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
	UniColor     u.V4
	UniSubTexPos u.V4
	UniColorMix  float32

	// Buffer initialization stuff
	buffersReady bool
	vaoHandle    uint32
	bufferHandle uint32

	finalized bool
}

// || ========================================================================================================================================================================
// ||
// || SPRITE WITH ATLAS
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
// || ========================================================================================================================================================================
func CreateTexAtlasRenderer(shaderFileBasename, atlas, subTexName string) *TexQuadRenderer {
	shader, err := Engine.GetShader(shaderFileBasename)
	u.GlPanicIfErrNotNil(err)

	atlasDescriptor := Engine.LoadTextureAtlas(atlas)

	subTexInfo := atlasDescriptor.GetSubTexture(subTexName)

	w, h := atlasDescriptor.Texture.GetSize()

	s := TexQuadRenderer{
		Shader:       shader,
		Texture:      atlasDescriptor.Texture,
		Atlas:        atlasDescriptor,
		UniColor:     u.V4{},
		UniSubTexPos: subTexInfo.GetDims(float32(w), float32(h)),
		UniColorMix:  0,
		buffersReady: false,
		vaoHandle:    0,
	}

	return &s
}

// || ===================================================
// ||
// || Create a Texture Quad for a non-atlassed texture
// ||
// || ===================================================
func CreateTexQuadRenderer(shaderAlias, textureAlias string) *TexQuadRenderer {

	tex, err := Engine.GetTexture(textureAlias)
	u.GlPanicIfErrNotNil(err)

	shader, err := Engine.GetShader(shaderAlias)
	u.GlPanicIfErrNotNil(err)

	T := TexQuadRenderer{
		Shader:       shader,
		Texture:      tex,
		UniColor:     u.OPAQ_WHITE(),
		UniColorMix:  0.5, // mix tex and unicolor equally. good for debugging
		UniSubTexPos: u.V4{C1: 0, C2: 0, C3: 1, C4: 1},
		buffersReady: false,
	}

	return &T
}

func (R *TexQuadRenderer) Finalize() {

	if R.finalized {
		return
	}

	R.Shader.Use()

	R.Texture.Finalize()

	if R.buffersReady {
		return
	}

	// Create interleaved array of verts and
	// texture coordinates.
	const vt_bytes_pr_stride = u.F32_SIZE * 5   // 4 bytes per float and 5 floats per stride/vert
	const vt_floats_total = 4 * 5               // 4 strides and 5 floats per stride
	const vt_len = u.F32_SIZE * vt_floats_total // total length (in bytes) of the VT buffer
	const Z, HI, LO = 1.0, 0.5, -0.5            // convenience
	vt_buffer := [vt_floats_total]float32{
		HI, HI, Z /* <== Vert | Tex ==> */, 1, 1,
		LO, HI, Z /* <== Vert | Tex ==> */, 0, 1,
		LO, LO, Z /* <== Vert | Tex ==> */, 0, 0,
		HI, LO, Z /* <== Vert | Tex ==> */, 1, 0,
	}
	vt_ptr := u.GlPtr32f(&vt_buffer[0])

	// Create buffer
	gl.GenBuffers(1, &R.bufferHandle)

	//
	// Vertex Array Object
	gl.GenVertexArrays(1, &R.vaoHandle)
	gl.BindVertexArray(R.vaoHandle)
	defer gl.BindVertexArray(0)

	//
	// Vertex Buffer Object
	gl.BindBuffer(gl.ARRAY_BUFFER, R.bufferHandle)
	defer gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BufferData(gl.ARRAY_BUFFER, vt_len, vt_ptr, gl.STATIC_DRAW)

	R.Shader.EnableVertexAttribArray("iVert")
	R.Shader.VertexAttribPointer("iVert", 3, gl.FLOAT, false, vt_bytes_pr_stride, 0)
	defer R.Shader.DisableVertexAttribArray("iVert")

	//
	// Texture Buffers
	R.Texture.Bind()
	defer R.Texture.Unbind()
	R.Shader.EnableVertexAttribArray("iTexCoord")
	R.Shader.VertexAttribPointer("iTexCoord", 2, gl.FLOAT, false, vt_bytes_pr_stride, 3*u.F32_SIZE)
	defer R.Shader.DisableVertexAttribArray("iTexCoord")

	R.finalized = true
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
