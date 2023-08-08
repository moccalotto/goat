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

	// dm.beforeFirstDraw() -- a function that is called just before the loop starts,
	// but after setup, so that the renderer and alle the settings are ready to rock.

	for dm.ProcessEvents(100) {

		renderer.SetDrawColor(dm.bgColor.R, dm.bgColor.G, dm.bgColor.B, dm.bgColor.A)
		renderer.Clear()

		dm.draw()

		if dm.autorender == AUTORENDER_ALWAYS {
			renderer.Present()
		}
		if dm.autorender == AUTORENDER_SKIP_ONCE {
			dm.autorender = AUTORENDER_ALWAYS
		}

	}

	script.Close()
	renderer.Destroy()
	window.Destroy()
	sdl.Quit()

	return 0
}
