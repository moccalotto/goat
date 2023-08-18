package main

import (
	"goat/glhelp"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	// start the mainthread system, allowing us to make calls on the main thread later
	StartMainThreadSystem(func() {

		// call actualMain() and ensure that it runs on the main thread
		// because it contains tons of OpenGL calls, etc.
		RunOnMain(actualMain)
	})
}

// This is the actual main function
// must run on mainthread
func actualMain() {

	options := &WindowOptions{
		Title:     "GOAT",
		Width:     1500,
		Height:    1500,
		Resizable: false,
	}

	_, window, err := initGlfw(options) // could switch to sfml, but cant get it to compile.
	if err != nil {
		panic(err)
	}

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action != glfw.Press {
			return
		}

		if key == glfw.KeyEscape {
			glfw.GetCurrentContext().SetShouldClose(true)
		}
	})

	i := 0

	triangle := glhelp.CreateTriangle()
	triangle2 := glhelp.CreateTriangle()

	for !window.ShouldClose() {
		glhelp.ClearF(0.1, 0.1, 0.1, 1.0)

		triangle.Draw()
		triangle2.Draw()

		glfw.PollEvents()

		glhelp.AssertGLOK("EndOfDraw", i)
		window.SwapBuffers()

		i++
	}
	// dm := CreateDrawing(window, "script.lua")
	// defer dm.Destroy()

	// 	print(glfw.GetKeyName(key, scancode))

	// 	if key == glfw.KeyEscape {
	// 		os.Exit(0)
	// 	}
	// })
	// dm.CallSetupFunc()

	// for dm.ProcessEvents(100) {

	// 	dm.CallDrawFunc()
	// }
}
