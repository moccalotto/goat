package shed

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Represents a shader program
type ShaderProgram struct {
	uniforms     map[string]int32
	attribs      map[string]int32
	vertShaderId uint32
	fragShaderId uint32
	programId    uint32
}

func CreateShaderProgramFromFiles(vertPath, fragPath string) *ShaderProgram {

	S := ShaderProgram{
		uniforms:     make(map[string]int32),
		attribs:      make(map[string]int32),
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

	if found && loc >= 0 {
		return uint32(loc), nil
	} else if found && loc < 0 {
		return 0, fmt.Errorf("could not get location of attribute '%s'", name)
	}

	loc = gl.GetAttribLocation(S.programId, GlStr(name))

	if loc < 0 {
		S.attribs[name] = -1
		return 0, fmt.Errorf("could not get location of attribute '%s'", name)
	}

	S.attribs[name] = loc

	return uint32(loc), nil
}

func (S ShaderProgram) HasAttrib(name string) bool {
	_, err := S.getAttribLocation(name)

	return err == nil
}

func (S *ShaderProgram) getUniformLocation(name string) (int32, error) {
	loc, found := S.uniforms[name]

	if found && loc >= 0 {
		return loc, nil
	} else if found && loc < 0 {
		return -1, fmt.Errorf("could not get location of uniform '%s'", name)
	}

	loc = gl.GetUniformLocation(S.programId, GlStr(name))

	if loc < 0 {
		S.uniforms[name] = -1
		return -1, fmt.Errorf("could not get location of uniform '%s'", name)
	}

	S.uniforms[name] = loc

	return loc, nil
}

func (s *ShaderProgram) SetUniformAttr(name string, value interface{}) error {
	loc, err := s.getUniformLocation(name)

	if err != nil {
		return err
	}

	switch typ := value.(type) {

	//Basic Types
	case bool:
		tmp := int32(0)
		if bool(value.(bool)) {
			tmp = 1
		}
		gl.Uniform1iv(loc, 1, &tmp)
	case int32:
		value := int32(value.(int32))
		gl.Uniform1iv(loc, 1, &value)
	case float32:
		value := value.(float32)
		gl.Uniform1fv(loc, 1, &value)

		// GOAT Types
	case V2:
		arr := value.(V2).ToArray()
		gl.Uniform2fv(loc, 1, &arr[0])
	case V3:
		arr := value.(V3).ToArray()
		gl.Uniform3fv(loc, 1, &arr[0])
	case V4:
		arr := value.(V4).ToArray()
		gl.Uniform4fv(loc, 1, &arr[0])

	case *TextureWrapper: // GOAT Texture Type
		value := int32(value.(*TextureWrapper).GetTextureUnit())
		gl.Uniform1iv(loc, 1, &value)

		// MGL Types
		//
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

	default:
		return GlProbablePanic(fmt.Errorf("unsupported data type: %v", typ))
	}
	return nil
}

func (S *ShaderProgram) DisableVertexAttribArray(name string) {

	location, err := S.getAttribLocation(name)

	if err != nil {
		GlPanic(fmt.Errorf("DisableVertexAttribArray: %s", err))
	}

	gl.DisableVertexAttribArray(location)

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

func (S *ShaderProgram) VertexAttribPointer(name string, size int32, xtype uint32, normalized bool, stride int32, pointer uintptr) {
	pos, err := S.getAttribLocation(name)
	if err != nil {
		GlPanic(fmt.Errorf("VertexAttribPointer: %v", err))
	}

	gl.VertexAttribPointerWithOffset(
		pos,        // Must match layout in shader
		size,       // number of components per vertex (1-4)
		xtype,      // type of data in vector components (i think)
		normalized, // false: data is not normalized
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

func compileShader(shaderType uint32, filePath string) (shader_id uint32, err error) {

	source, err := ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	if (shaderType != gl.VERTEX_SHADER) && (shaderType != gl.FRAGMENT_SHADER) {
		return 0, errors.New("invalid shader_type argument. Must be GL_FRAGMENT_SHADER or GL_VERTEX_SHADER")
	}

	shader_id = 0

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
