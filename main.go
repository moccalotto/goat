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

	thing := glhelp.CreatePolygon(4, []float32{}, "textures/cat.png")
	thing.Initialize()

	zero := float32(0)
	time := zero
	frameCount := zero
	pulse := zero + 0.01

	camera := glhelp.CreateCamera()
	camera.SetFrameSize(10, 10)
	camera.SetPosition(0, 0)
	camera.SetRotation(45 * glhelp.Degrees)
	camera.Zoom(0.1)

	thing.SetPosition(0, 0)

	thing.SetScale(2, 2)

	thing2 := thing.Copy()
	thing2.SetPosition(0, -3)
	thing2.SetScale(3, 1)

	for !window.ShouldClose() {
		frameCount++

		glhelp.ClearScreenF(0.1, 0.1, 0.1, 1.0)

		time += pulse
		camera.Zoom(1 + pulse*3)

		if time > 0.9 {
			pulse = -pulse
		}
		if time < 0 {
			pulse = -pulse
		}

		thing.SetRotation(time)

		thing.Draw(camera.GetTransformationMatrix())
		thing2.Draw(camera.GetTransformationMatrix())

		glfw.PollEvents()

		glhelp.AssertGLOK("EndOfDraw")
		window.SwapBuffers()
	}
}
