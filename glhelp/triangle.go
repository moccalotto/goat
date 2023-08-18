package glhelp

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Triangle struct {
	vbo          uint32 // handle to vertex buffer
	vao          uint32 // handle to vertex array
	vertShaderId uint32 // handle to the compiled vertex shader
	fragShaderId uint32 // handle to the compiled fragment shader
	programId    uint32 // handle to the linked shader program
	Scale        float32
	verts        []float32
	readyToDraw  bool // have all the OpenGL calls been made so we're ready to call Draw()
}

func CreateTriangle() *Triangle {

	return &Triangle{
		vao:   0,
		vbo:   0,
		Scale: 1.0,
		verts: []float32{
			0, 1, 0, // top center
			-1, -1, 0, // low left
			1, -1, 0, // low right
		},
		readyToDraw: false,
	}
}

func (TR *Triangle) Initialize() {
	var err error
	TR.vertShaderId, err = compileShader(gl.VERTEX_SHADER, "shaders/triangle.vert")
	if err != nil {
		panic(err)
	}

	TR.fragShaderId, err = compileShader(gl.FRAGMENT_SHADER, "shaders/triangle.frag")
	if err != nil {
		panic(err)
	}

	TR.programId = gl.CreateProgram()
	gl.AttachShader(TR.programId, TR.vertShaderId)
	gl.AttachShader(TR.programId, TR.fragShaderId)
	gl.LinkProgram(TR.programId)

	var success int32
	gl.GetProgramiv(TR.programId, gl.LINK_STATUS, &success)
	if success != gl.TRUE {
		logStr := GetProgramLog(TR.programId)
		panic(logStr)
	}
	gl.ValidateProgram(TR.programId)
	gl.GetProgramiv(TR.programId, gl.VALIDATE_STATUS, &success)
	if success != gl.TRUE {
		logStr := GetProgramLog(TR.programId)
		panic(logStr)
	}

	gl.DetachShader(TR.programId, TR.vertShaderId)
	gl.DetachShader(TR.programId, TR.fragShaderId)
	gl.DeleteShader(TR.vertShaderId)
	gl.DeleteShader(TR.fragShaderId)

	AssertGLOK("CreateTriangle - Shader Setup")

	// Create VAO
	gl.GenVertexArrays(1, &TR.vao)
	gl.BindVertexArray(TR.vao)

	// Create VBO for vertices
	gl.GenBuffers(1, &TR.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, TR.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(TR.verts), gl.Ptr(TR.verts), gl.STATIC_DRAW)

	/* Specify that our coordinate data is going into attribute index 0, and contains two floats per vertex */
	vertexArrayLocation := uint32(0)
	gl.VertexAttribPointer(vertexArrayLocation, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(vertexArrayLocation)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	AssertGLOK("CreateTriangle - Buffer Setup")

	TR.readyToDraw = true
}

func (TR *Triangle) Draw() {
	if !TR.readyToDraw {
		panic("You forgot to call Initialize()")
	}

	gl.UseProgram(TR.programId)
	gl.BindVertexArray(TR.vao)
	gl.Uniform1f(0, TR.Scale)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(TR.verts)))
	// pass
}

func (TR *Triangle) Destroy() {
	gl.DeleteBuffers(1, &TR.vao)
	gl.DeleteBuffers(1, &TR.vbo)
	gl.DeleteProgram(TR.programId)
	AssertGLOK("Triangle.Destroy")
}
