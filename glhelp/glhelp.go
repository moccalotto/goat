package glhelp

/**
Wrapper for all calls to opengl.

- make importing opengl less stupid
- streamline calls so they return better errors
- wrap stuff in sane data structures


*/
import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.6-core/gl"
)

func Init() error {
	ret := gl.Init()

	return ret
}

func Str(s string) *uint8 {
	return gl.Str(s + "\x00")
}

func WindowViewport2X(x, y, w, h int32) error {

	var vp [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &vp[0])
	log.Printf("%+v\n", vp)

	gl.Viewport(
		0,
		0,
		w*2,
		h*2,
	)

	if errCode := gl.GetError(); errCode != 0 {
		return fmt.Errorf("cound not intialize viewport. Error code: %v", errCode)
	}

	return nil
}

func ClearF(r, g, b, a float32) {
	// gl.ClearColor(r, g, b, a)
	// ClearX()
}

func ClearI(r, g, b, a uint8) {
	ClearF(
		float32(r)/255.0,
		float32(g)/255.0,
		float32(b)/255.0,
		float32(a)/255.0,
	)
}

func ClearX() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

var (
	program_id uint32 = 0
)

func GetShaderLog(shaderId uint32) string {
	var log_length int32
	gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &log_length)

	if log_length == 0 {
		return ""
	}

	buffer := make([]byte, log_length+1)
	gl.GetShaderInfoLog(shaderId, int32(len(buffer)), nil, &buffer[0])

	return string(buffer)
}

func GetProgramLog(programId uint32) string {
	var log_length int32

	gl.GetProgramiv(programId, gl.INFO_LOG_LENGTH, &log_length)

	if log_length == 0 {
		return ""
	}

	buffer := make([]byte, log_length+1)
	gl.GetProgramInfoLog(programId, int32(len(buffer)), nil, &buffer[0])

	return string(buffer)
}

func CreateProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {

	// VERTEX SHADER
	vert_shader_id, err := compileShader(gl.VERTEX_SHADER, "vertex", vertexShaderSource)

	if err != nil {
		return 0, err
	}

	frag_shader_id, err := compileShader(gl.FRAGMENT_SHADER, "fragment", fragmentShaderSource)

	if err != nil {
		return 0, err
	}

	// CREATE PROGRAM
	program_id := gl.CreateProgram()
	gl.AttachShader(program_id, vert_shader_id)
	gl.AttachShader(program_id, frag_shader_id)

	gl.LinkProgram(program_id)

	var link_status int32
	gl.GetProgramiv(program_id, gl.LINK_STATUS, &link_status)

	if link_status != gl.TRUE {

		log := GetProgramLog(program_id)

		return 0, fmt.Errorf("could not link shaders. %s", log)
	}

	// SOME SAY I SHOULD DO THIS HERE????
	// Seems weird
	gl.DetachShader(program_id, vert_shader_id)
	gl.DetachShader(program_id, frag_shader_id)
	gl.DeleteShader(vert_shader_id)
	gl.DeleteShader(frag_shader_id)

	return program_id, nil
}
