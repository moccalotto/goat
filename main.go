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
	cfg := NewConfig()

	script := luaLoadScript(cfg)

	window, renderer := setupSDL(cfg)

	dm := CreateDrawing(renderer, script)

	dm.bgColor = cfg.Background

	dm.injectFunctions()

	for dm.ProcessEvents(100) {

		renderer.SetDrawColor(dm.bgColor.R, dm.bgColor.G, dm.bgColor.B, dm.bgColor.A)
		renderer.Clear()

		dm.draw()

		renderer.Present()
	}

	script.Close()
	renderer.Destroy()
	window.Destroy()
	sdl.Quit()

	return 0
}
