package tractor

import (
	u "goat/shed"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// ||=============================
// ||
// || Basic (Filled) Rect Renderer
// ||
// || Render a quad with a texture
// || and map into subtextures.
// ||=============================
type BasicRectRenderer struct {
	Shader *u.ShaderProgram

	// Uniform variables to send to the shader
	UniColor u.V4

	// Buffer initialization stuff
	buffersReady bool
	vaoHandle    uint32
	bufferHandle uint32 // we only have the vertex buffer.
}

func CreateBasicRectRenderer(shaderFileBaseName string) *BasicRectRenderer {

	shader, err := Engine.GetShader(shaderFileBaseName)
	u.GlPanicIfErrNotNil(err)

	return &BasicRectRenderer{
		Shader:   shader,
		UniColor: u.OPAQ_WHITE(),
	}
}

func (R *BasicRectRenderer) Finalize() {
	R.Shader.Use()

	R.Shader.SetUniformAttr("uniColor", R.UniColor)

	if R.buffersReady {
		return
	}

	const vt_floats_total = 4 * 5               // 4 strides and 5 floats per stride
	const vt_len = u.F32_SIZE * vt_floats_total // total length (in bytes) of the VT buffer
	const Z, HI, LO = 1.0, 0.5, -0.5            // convenience
	vt_buffer := [vt_floats_total]float32{
		HI, HI, Z,
		LO, HI, Z,
		LO, LO, Z,
		HI, LO, Z,
	}
	vt_ptr := u.GlPtr32f(&vt_buffer[0])

	// Create buffers
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
	R.Shader.VertexAttribPointer("iVert", 3, gl.FLOAT, false, 0, 0)
	defer R.Shader.DisableVertexAttribArray("iVert")

	// Cleanup
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	R.Shader.DisableVertexAttribArray("iVert")
}

func (R *BasicRectRenderer) Draw(camMatrix, objTranslationMatrix mgl32.Mat3, color u.V4) {

	R.Shader.Use()
	gl.BindVertexArray(R.vaoHandle)

	trMatrix := camMatrix.Mul3(objTranslationMatrix)

	u.GlPanicIfErrNotNil(R.Shader.SetUniformAttr("uniColor", color))
	u.GlPanicIfErrNotNil(R.Shader.SetUniformAttr("uniTransformation", trMatrix))

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

	u.AssertGLOK("BasicRectRenderer.Draw", R.Shader, 22)
}

func (R *BasicRectRenderer) Clone() *BasicRectRenderer {

	return &BasicRectRenderer{
		Shader:       R.Shader,
		UniColor:     R.UniColor,
		buffersReady: R.buffersReady,
		vaoHandle:    R.vaoHandle,
		bufferHandle: R.bufferHandle,
	}
}
