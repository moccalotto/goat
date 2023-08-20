package glhelp

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Quad struct {
	buffers      []uint32  // handle to vertex buffer
	vao          uint32    // handle to vertex array
	vertShaderId uint32    // handle to the compiled vertex shader
	fragShaderId uint32    // handle to the compiled fragment shader
	programId    uint32    // handle to the linked shader program
	Scale        float32   // How much to scale the model
	verts        []float32 // the model's vertices (3 floats per UNIQUE vertex)
	indeces      []uint32  // indeces into the verts array (one index per vertex)
	readyToDraw  bool
}

func CreateQuad() *Quad {

	return &Quad{
		buffers: make([]uint32, 2),

		verts: []float32{
			1, 1, 0, // top right
			-1, 1, 0, // bot right
			-1, -1, 0, // bot left
			1, -1, 0, // top left
		},
		indeces: []uint32{
			0, 2, 3,
			0, 1, 2,
		},
		readyToDraw: false,
	}
}

func (QD *Quad) Initialize() {

	var err error
	QD.vertShaderId, err = compileShader(gl.VERTEX_SHADER, "shaders/quad.vert")
	if err != nil {
		panic(err)
	}

	QD.fragShaderId,
		err = compileShader(gl.FRAGMENT_SHADER, "shaders/quad.frag")
	if err != nil {
		panic(err)
	}

	QD.programId = gl.CreateProgram()
	gl.AttachShader(QD.programId, QD.vertShaderId)
	gl.AttachShader(QD.programId, QD.fragShaderId)
	gl.LinkProgram(QD.programId)

	var success int32
	gl.GetProgramiv(QD.programId, gl.LINK_STATUS, &success)
	if success != gl.TRUE {
		logStr := GetProgramLog(QD.programId)
		panic(logStr)
	}
	gl.ValidateProgram(QD.programId)
	gl.GetProgramiv(QD.programId, gl.VALIDATE_STATUS, &success)
	if success != gl.TRUE {
		logStr := GetProgramLog(QD.programId)
		panic(logStr)
	}

	gl.DetachShader(QD.programId, QD.vertShaderId)
	gl.DetachShader(QD.programId, QD.fragShaderId)
	gl.DeleteShader(QD.vertShaderId)
	gl.DeleteShader(QD.fragShaderId)

	gl.UseProgram(QD.programId)

	AssertGLOK("CreateQuad - Shader Setup")

	// Create VAO
	gl.GenVertexArrays(1, &QD.vao)
	gl.BindVertexArray(QD.vao)

	// Prepare buffers
	gl.GenBuffers(2, &QD.buffers[0])

	// Create VBO for vertices
	gl.BindBuffer(gl.ARRAY_BUFFER, QD.buffers[0])
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(QD.verts), gl.Ptr(QD.verts), gl.STATIC_DRAW)

	// Buffer for elements
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, QD.buffers[1])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(QD.indeces), gl.Ptr(QD.indeces), gl.STATIC_DRAW)

	/* Specify that our coordinate data is going into attribute index 0, and contains two floats per vertex */
	vertexArrayLocation := uint32(0)
	gl.VertexAttribPointer(vertexArrayLocation, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(vertexArrayLocation)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	AssertGLOK("Quad : Initialize")
}

func (QD *Quad) Draw() {
	gl.UseProgram(QD.programId)
	gl.BindVertexArray(QD.vao)
	AssertGLOK("Quad: Draw")

	gl.DrawElements(gl.TRIANGLES, int32(len(QD.indeces)), gl.UNSIGNED_INT, nil)
	AssertGLOK("Quad: Draw")
}

func (QD *Quad) Destroy() {
	gl.DeleteBuffers(2, &QD.buffers[0])
	gl.DeleteProgram(QD.programId)
	AssertGLOK("Triangle.Destroy")
}
