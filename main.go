package main

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	exitCode := 0

	// Tell SDL that we're currently in the main
	// thread, and that our main program code should
	// continue inside it.
	// This way, SDL can find its way back to the main thread
	// when it needs to.
	sdl.Main(func() {
		exitCode = run()
	})

	os.Exit(exitCode)
}

// This is the actual main function
func run() int {
	cfg := &Config{
		Title:     "Goaty McWindow",
		Width:     800,
		Height:    600,
		CanResize: false,
	}

	script := setupLua(cfg)

	window, renderer := setupSDL(cfg)

	drawing := CreateDrawing(renderer, script)

	loop(drawing)

	script.Close()
	renderer.Destroy()
	window.Destroy()
	sdl.Quit()

	return 0
}
