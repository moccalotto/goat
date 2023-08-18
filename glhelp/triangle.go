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

	verts []float32
}

var (
	blam int = 0
)

func CreateTriangle() *Triangle {

	blam += 1
	TR := Triangle{
		vao: 0,
		vbo: 0,
		verts: []float32{
			0, 1, 0, // top center
			-1, -1, 0, // low left
			1, -1, 0, // low right
		},
	}

	if blam%2 == 0 {
		for i, v := range TR.verts {
			TR.verts[i] = v * -1
		}
	}

	var err error
	TR.vertShaderId, err = compileShader(gl.VERTEX_SHADER, "triangle.vert")
	if err != nil {
		panic(err)
	}

	TR.fragShaderId, err = compileShader(gl.FRAGMENT_SHADER, "triangle.frag")
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

	gl.UseProgram(TR.programId)

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

	return &TR
}

func (TR *Triangle) Draw() {
	gl.UseProgram(TR.programId)

	gl.BindVertexArray(TR.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(TR.verts)))
	// pass
}
