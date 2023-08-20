package glhelp

import (
	"errors"
	"fmt"
	"log"
	"math"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Represents a shader program
type ShaderProgram struct {
	uniforms     map[string]int32
	attribs      map[string]uint32
	vertShaderId uint32
	fragShaderId uint32
	programId    uint32
}

func CreateShaderProgramFromFiles(vertPath, fragPath string) *ShaderProgram {

	S := ShaderProgram{
		uniforms:     make(map[string]int32),
		attribs:      make(map[string]uint32),
		vertShaderId: 0,
		fragShaderId: 0,
		programId:    0,
	}
	var err error

	if S.vertShaderId, err = compileShader(gl.VERTEX_SHADER, vertPath); err != nil {
		GlPanic(err)
	}

	if S.fragShaderId, err = compileShader(gl.FRAGMENT_SHADER, fragPath); err != nil {
		GlPanic(err)
	}

	S.programId = gl.CreateProgram()
	gl.AttachShader(S.programId, S.vertShaderId)
	gl.AttachShader(S.programId, S.fragShaderId)
	gl.LinkProgram(S.programId)

	if err := S.getLinkError(); err != nil {
		GlPanic(fmt.Errorf("could not link shaders. %v", err))
	}

	gl.DetachShader(S.programId, S.vertShaderId)
	gl.DetachShader(S.programId, S.fragShaderId)
	gl.DeleteShader(S.vertShaderId)
	gl.DeleteShader(S.fragShaderId)

	AssertGLOK("CreateShaderFromFile")
	return &S
}

func (S *ShaderProgram) getAttribLocation(name string) (uint32, error) {
	loc, found := S.attribs[name]

	if found {
		return loc, nil
	}

	attr := gl.GetAttribLocation(S.programId, GlStr(name))

	if attr < 0 {
		err := fmt.Errorf("could not get location of attribute '%s'", name)
		return math.MaxUint32, GlProbablePanic(err)
	}

	S.attribs[name] = uint32(attr)

	return uint32(attr), nil
}

func (S *ShaderProgram) getUniformLocation(name string) (int32, error) {
	loc, found := S.uniforms[name]

	if found {
		return loc, nil
	}

	if loc == -1 {
		return -1, errors.New("not found")
	}

	attr := gl.GetUniformLocation(S.programId, GlStr(name))

	if attr < 0 {
		err := fmt.Errorf("could not get location of uniform '%s'", name)
		S.uniforms[name] = -1
		return math.MaxInt32, err
	}

	S.uniforms[name] = attr

	return attr, nil
}

func (s *ShaderProgram) SetUniformAttr(name string, value interface{}) error {
	loc, err := s.getUniformLocation(name)

	if err != nil {
		return err
	}

	switch typ := value.(type) {
	case int32:
		value := int32(value.(int32))
		gl.Uniform1iv(loc, 1, &value)
	case float32:
		value := value.(float32)
		gl.Uniform1fv(loc, 1, &value)
	case mgl32.Vec2:
		value := value.(mgl32.Vec2)
		gl.Uniform2fv(loc, 1, &value[0])
	case mgl32.Vec3:
		value := value.(mgl32.Vec3)
		gl.Uniform3fv(loc, 1, &value[0])
	case mgl32.Vec4:
		value := value.(mgl32.Vec4)
		gl.Uniform4fv(loc, 1, &value[0])
	case mgl32.Mat3:
		value := value.(mgl32.Mat3)
		gl.UniformMatrix3fv(loc, 1, false, &value[0])
	case *Texture:
		value := value.(*Texture).GetTextureUnit()
		gl.Uniform1iv(loc, 1, &value)
	default:
		return fmt.Errorf("unsupported data type: %v", typ)
	}
	return nil
}

func (S *ShaderProgram) DisableVertexAttribArray(name string) {

	location, err := S.getAttribLocation(name)

	if err != nil {
		GlPanic(fmt.Errorf("DisableVertexAttribArray: %s", err))
	}

	gl.EnableVertexAttribArray(location)

	AssertGLOK("DisableVertexAttribArray")
}

func (S *ShaderProgram) EnableVertexAttribArray(name string) {

	location, err := S.getAttribLocation(name)

	if err != nil {
		GlPanic(fmt.Errorf("EnableVertexAttribArray: %s", err))
	}

	gl.EnableVertexAttribArray(location)

	AssertGLOK("EnableVertexAttribArray", name)
}

func (S *ShaderProgram) VertexAttribPointer(name string, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
	pos, err := S.getAttribLocation(name)
	if err != nil {
		GlPanic(fmt.Errorf("VertexAttribPointer: %v", err))
	}

	gl.VertexAttribPointer(
		pos,        // Must match layout in shader
		size,       // size
		xtype,      // type of data in vector components (i think)
		normalized, // data is not normalized
		stride,     // stride (there are zero bytes in between the vertices in the array)
		pointer,
	)

	AssertGLOK("VertexAttribPointer", name)
}

func (S *ShaderProgram) Use() {
	gl.UseProgram(S.programId)
	AssertGLOK("Shader.Use")
}

func (S *ShaderProgram) Destroy() {
	gl.DeleteProgram(S.programId)
	AssertGLOK("Shader.Destroy")
}

func (S *ShaderProgram) getLinkError() error {

	var link_status int32

	gl.GetProgramiv(S.programId, gl.LINK_STATUS, &link_status)

	if link_status != gl.TRUE {
		logStr := GetProgramLog(S.programId)
		return GlProbablePanic(fmt.Errorf("linker Error: %v", logStr))
	}

	return nil
}

func compileShader(shaderType uint32, filePath string) (shader_id uint32, e error) {

	source, err := ReadFile(filePath)
	if err != nil {
		return 0, err
	}

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
		logStr := GetShaderInfoLog(shader_id)

		if logStr == "" {
			return 0, fmt.Errorf("cannot compile shader '%s'", filePath)
		}

		log.Println(logStr)
		return 0, fmt.Errorf("cannot compile shader '%s': %s", filePath, logStr)
	}
	return shader_id, nil
}
