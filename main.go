package main

import (
	"os"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	RunInMainthread(func() {
		run()
	})
}

// This is the actual main function
func run() int {

	options := &WindowOptions{
		Title:     "GOAT",
		Width:     800,
		Height:    600,
		Resizable: false,
	}

	freeGl, window, err := createGlWindow(options) // could switch to sfml, but cant get it to compile.
	if err != nil {
		panic(err)
	}
	defer freeGl()

	dm := CreateDrawing(window, "script.lua")
	defer dm.Destroy()

	dm.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action != glfw.Press {
			return
		}

		print(glfw.GetKeyName(key, scancode))

		if key == glfw.KeyEscape {
			os.Exit(0)
		}
	})

	/* 	NEXT STEPS
	ROTATION MATRIX
	does what it says on the tin

	ENTITIES
	an entity is a simple shape (initially)
	It has a location, one or two colors, and a few other parameters.
	It can be ephemeral - so it is removed right after renedering - or it can persist
	You can chose to execute its draw instructions right away, or when the Draw() is at an end
		* Dot (Size, Location, Color)
		* Line (Stroke Width, Stroke Color, Location, Length, Rotation)
		* Circle (fill color, stroke thickness, stroke color, location, size)
		* Ellipse (fill color, stroke thickness, stroke color location, size_a, size_b, rotation)
		* Rectangle (fill color, stroke thickness, stroke color, corner_radius, location, size_a, size_b, rotation)
		* Image
	*/

	dm.CallSetupFunc()

	for dm.ProcessEvents(100) {

		dm.CallDrawFunc()

	}

	return 0
}
