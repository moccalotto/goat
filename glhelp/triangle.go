package glhelp

import (
	"math"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

var (
	_mainShader *Shader
	angle       float32 = 0
)

func Triangle() {
	if _mainShader == nil {

		var err error
		_mainShader, err = CreateProgramFromFiles(
			"shad_vert.vert",
			"shad_frag.frag",
		)

		if err != nil {
			panic(err)
		}

		_mainShader.Use()
		print("shaders compiled and loaded\n")
	}

	S := _mainShader

	angle = angle + 0.05
	if angle > 2*math.Pi {
		angle -= 2 * math.Pi
	}
	S.Uniform2fv("u_scale", []float32{0.6, 0.3})
	S.Uniform1f("u_rotAngle", angle)

	verts := []float32{
		-1.0, -1.0, 0.0,
		1.0, 1.0, 0.0,
		1.0, -1.0, 0.0,
	}

	var vert_buffer uint32
	gl.GenBuffers(1, &vert_buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vert_buffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, unsafe.Pointer(&verts[0]), gl.STATIC_DRAW)

	S.EnableVertexAttribArray("in_vert")

	S.VertexAttribPointer(
		"in_vert", // name of attribute
		3,         // size (vec3)
		gl.FLOAT,  // type of data in vector components (i think)
		false,     // data is not normalized
		0,         // stride (there are zero bytes in between the vertices in the array)
		unsafe.Pointer(nil),
	)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.DisableVertexAttribArray(0)

	if errCode := gl.GetError(); errCode != 0 {
		// panic(fmt.Sprintf("OpenGL error %d", errCode))
	}
}
