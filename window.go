package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type WindowOptions struct {
	Width     int
	Height    int
	Title     string
	Resizable bool
}

func (O *WindowOptions) ResizableInt() int {
	if O.Resizable {
		return glfw.True
	}

	return glfw.False
}

func initGlfw(options *WindowOptions) (func(), *glfw.Window, error) {

	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("glfw.Init failed: %v", err))
	}

	glfw.WindowHint(glfw.Visible, glfw.True)
	glfw.WindowHint(glfw.Resizable, options.ResizableInt())
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(options.Width, options.Height, options.Title, nil, nil) // utilize last argument to support multiple windows.
	if err != nil {
		panic(fmt.Errorf("CreateWindow failed: %v", err))
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	freeFuc := func() {
		window.Destroy()
		glfw.Terminate()
	}

	return freeFuc, window, nil
}
