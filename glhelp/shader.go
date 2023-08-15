package glhelp

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Represents a shader program
type Shader struct {
	uniforms     map[string]int32
	attribs      map[string]uint32
	vertShaderId uint32
	fragShaderId uint32
	programId    uint32
	logger       *log.Logger
	panicLevel   int
}

func _readShaderFileOrPanic(filename string) string {

	bytes, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func (S *Shader) error(severity int, err error) error {
	S.logger.Print(err)

	if severity >= S.panicLevel {
		panic(err)
	}

	return err
}

func CreateProgramFromFiles(vertPath, fragPath string) (*Shader, error) {

	S := Shader{
		uniforms:     make(map[string]int32),
		attribs:      make(map[string]uint32),
		vertShaderId: 0,
		fragShaderId: 0,
		programId:    0,
		logger:       log.New(os.Stderr, "shader.go :> ", 0),
		panicLevel:   0, // errors with a "severity" level higher than panicLevel will cause panic()
	}
	var err error

	if S.vertShaderId, err = compileShader(gl.VERTEX_SHADER, vertPath, _readShaderFileOrPanic(vertPath)); err != nil {
		S.error(1, err)
		return nil, err
	}

	if S.fragShaderId, err = compileShader(gl.FRAGMENT_SHADER, fragPath, _readShaderFileOrPanic(fragPath)); err != nil {
		S.error(1, err)
		return nil, err
	}

	S.programId = gl.CreateProgram()
	gl.AttachShader(S.programId, S.vertShaderId)
	gl.AttachShader(S.programId, S.fragShaderId)
	gl.LinkProgram(S.programId)

	if err := S.getLinkError(); err != nil {
		S.CleanupShaders()
		log.Println("shader linking problem", err)
		return nil, fmt.Errorf("could not link shaders. %v", err)
	}

	S.CleanupShaders()

	return &S, nil
}

func (S *Shader) getAttribLocation(name string) (uint32, error) {
	loc, found := S.attribs[name]

	if found {
		return loc, nil
	}

	attr := gl.GetAttribLocation(S.programId, Str(name))

	if attr < 0 {
		err := fmt.Errorf("could not get location of attribute '%s'", name)
		return math.MaxUint32, S.error(1, err)
	}

	S.attribs[name] = uint32(attr)

	return uint32(attr), nil
}

func (S *Shader) getUniformLocation(name string) (int32, error) {
	loc, found := S.uniforms[name]

	if found {
		return loc, nil
	}

	attr := gl.GetUniformLocation(S.programId, Str(name))

	if attr < 0 {
		err := fmt.Errorf("could not get location of uniform '%s'", name)
		return math.MaxInt32, S.error(1, err)
	}

	S.uniforms[name] = attr

	return attr, nil
}

func (S *Shader) EnableVertexAttribArray(name string) error {

	location, err := S.getAttribLocation(name)

	if err != nil {
		return S.error(
			2,
			fmt.Errorf("EnableVertexAttribArray: %s", err),
		)
	}

	gl.EnableVertexAttribArray(location)

	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		return S.error(
			2,
			fmt.Errorf("EnableVertexAttribArray: %v", errCode),
		)
	}
	return nil
}
func (S *Shader) VertexAttribPointer(name string, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) error {
	pos, err := S.getAttribLocation(name)
	if err != nil {
		return S.error(
			2,
			fmt.Errorf("VertexAttribPointer: %v", err),
		)
	}

	gl.VertexAttribPointer(
		pos,        // Must match layout in shader
		size,       // size (vec3)
		xtype,      // type of data in vector components (i think)
		normalized, // data is not normalized
		stride,     // stride (there are zero bytes in between the vertices in the array)
		pointer,
	)

	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		return S.error(
			2,
			fmt.Errorf("VertexAttribPointer: GL Error Code: %v", errCode),
		)
	}

	return nil
}

func (S *Shader) Uniform1f(name string, value float32) error {
	location, err := S.getUniformLocation(name)

	if err != nil {
		return err
	}

	gl.Uniform1f(
		location,
		value,
	)

	if gl.GetError() != gl.NO_ERROR {
		panic("Could not Uniform1f")
	}

	return nil
}

func (S *Shader) Uniform2fv(name string, values []float32) error {
	location, err := S.getUniformLocation(name)

	if err != nil {
		return err
	}

	gl.Uniform2fv(
		location,
		int32(len(values)/2),
		&values[0],
	)

	if gl.GetError() != gl.NO_ERROR {
		panic("Could not Uniform2fv")
	}

	return nil
}

func (S *Shader) Use() {
	gl.UseProgram(S.programId)
}

func (S *Shader) Destroy() {
	S.CleanupShaders()
	gl.DeleteProgram(S.programId)
}

func (S *Shader) CleanupShaders() {
	gl.DetachShader(S.programId, S.vertShaderId)
	gl.DetachShader(S.programId, S.fragShaderId)
	gl.DeleteShader(S.vertShaderId)
	gl.DeleteShader(S.fragShaderId)
}

func (S *Shader) getLinkError() error {

	var link_status int32

	gl.GetProgramiv(S.programId, gl.LINK_STATUS, &link_status)

	if link_status != gl.TRUE {

		logStr := GetProgramLog(S.programId)

		return errors.New(logStr)
	}

	return nil
}

func compileShader(shaderType uint32, filenameHint, source string) (shader_id uint32, e error) {

	if (shaderType != gl.VERTEX_SHADER) && (shaderType != gl.FRAGMENT_SHADER) {
		return 0, errors.New("invalid shader_type argument. Must be GL_FRAGMENT_SHADER or GL_VERTEX_SHADER")
	}

	shader_id = gl.CreateShader(shaderType)
	source_bytes, free := gl.Strs(source)
	defer free()
	length := int32(len(source)) // len returns number of bytes in string, not number of chars.
	gl.ShaderSource(shader_id, 1, source_bytes, &length)
	gl.CompileShader(shader_id)

	var success int32
	gl.GetShaderiv(shader_id, gl.COMPILE_STATUS, &success)

	if success != gl.TRUE {
		logStr := GetShaderLog(shader_id)

		if logStr == "" {
			return 0, fmt.Errorf("cannot compile shader '%s'", filenameHint)
		}

		log.Println(logStr)
		return 0, fmt.Errorf("cannot compile shader '%s': %s", filenameHint, logStr)
	}
	return shader_id, nil
}
