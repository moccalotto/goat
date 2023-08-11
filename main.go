package main

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	exitCode := 0
	defer os.Exit(exitCode)

	// Tell SDL that we're currently in the main
	// thread, and that our main program code should
	// continue inside it.
	// This way, SDL can find its way back to the main thread
	// when it needs to.
	sdl.Main(func() {
		exitCode = run()
	})
}

// This is the actual main function
func run() int {
	script := luaLoadScript()

	// todo: switch to sfml https://www.sfml-dev.org/
	// 	     https://github.com/SFML/SFML
	window, renderer := setupSDL()

	dm := CreateDrawing(renderer, script)

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
		* Rectange (fill color, stroke thickness, stroke color, corner_radius, location, size_a, size_b, rotation)
		* Image
	*/

	dm.exportFunctions()

	dm.setup()

	for dm.ProcessEvents(100) {

		dm.draw()

	}

	script.Close()
	renderer.Destroy()
	window.Destroy()
	sdl.Quit()

	return 0
}
