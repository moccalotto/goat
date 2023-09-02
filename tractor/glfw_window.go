package tractor

import (
	"goat/shed"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type WindowOptions struct {
	Width     int
	Height    int
	Title     string
	Resizable bool
}

func glfwBool(b bool) int {
	if b {
		return glfw.True
	}

	return glfw.False
}

func glfwCreateWin(O *WindowOptions) (freeFunc func(), window *glfw.Window, err error) {

	err = glfw.Init()
	shed.GlPanicIfErrNotNil(err)

	glfw.WindowHint(glfw.Visible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfwBool(O.Resizable))
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err = glfw.CreateWindow(O.Width, O.Height, O.Title, nil, nil) // last two args are about multiple windows
	shed.GlPanicIfErrNotNil(err)

	window.MakeContextCurrent()

	err = gl.Init()
	shed.GlPanicIfErrNotNil(err)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	shed.GlLog("OpenGL version: %s\n", version)

	freeFunc = func() {
		if window != nil {
			window.Destroy()
		}
		glfw.Terminate()
	}

	return freeFunc, window, nil
}
