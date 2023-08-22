package glhelp

/**
This file contains some unorthodox error handling that
borders on exception handling.
But its fast and easy.
*/

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/go-gl/gl/v4.6-core/gl"
)

var (
	AlwaysPanic = true
	AlwaysLog   = true
	logger      = log.New(os.Stderr, "GLH", log.LstdFlags)
)

// Log and then panic
func GlPanic(err error) {
	logger.Writer().Write(debug.Stack())
	logger.Println(err)
	panic(err)
}

func GlPanicIfErrNotNil(err error) {
	if err != nil {
		GlPanic(err)
	}
}

// Panic if AlwaysPanic == true
// Always log if AlwaysLog == true
// if AlwaysLog and AlwaysPanic are false,
// this function does nothing
func GlProbablePanic(err error) error {

	if AlwaysPanic {
		GlPanic(err)
	}

	if AlwaysLog {
		logger.Println(err)
	}

	return err
}

// Checks if there are any opengl errors in the queue and panics if necessary
func AssertGLOK(values ...interface{}) error {
	errCode := gl.GetError()

	if errCode == gl.NO_ERROR {
		return nil
	}

	if len(values) == 0 {
		return GlProbablePanic(fmt.Errorf("openGL error. Code: %d", errCode))
	}

	return GlProbablePanic(fmt.Errorf("[%s] OpenGL Error: Code %d", values[0], errCode))
}

// Get hte program error log
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

// Get the shader error log
func GetShaderInfoLog(shaderId uint32) string {
	var log_length int32
	gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &log_length)

	if log_length == 0 {
		return ""
	}

	buffer := make([]byte, log_length+1)
	gl.GetShaderInfoLog(shaderId, int32(len(buffer)), nil, &buffer[0])

	return string(buffer)
}
