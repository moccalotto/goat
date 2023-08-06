package main

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func main() {
	exitCode := 0

	sdl.Main(func() {
		exitCode = 0
		run()
	})

	os.Exit(exitCode)
}

func run() {
	sdl.Main(func() {})
	script := lua.NewState()
	defer script.Close()

	if err := script.DoFile("script.lua"); err != nil {
		panic(err)
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	defer sdl.Quit()

	sys := &SystemSettings{}

	// Call the setup function.
	// In that function we set parameters needed to set the initial size and
	// title of the window, etc.
	setupFunc := script.GetGlobal("Setup")
	if setupFunc.Type() != lua.LTNil {
		callLuaFunc(script, setupFunc, luar.New(script, sys))
	}

	window, err := sdl.CreateWindow(
		"test",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		800,
		600,
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE,
	)

	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(
		window,
		-1,
		sdl.RENDERER_ACCELERATED|
			sdl.RENDERER_PRESENTVSYNC,
	)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	drawing := CreateDrawing(renderer)

	loop(drawing, sys, script)

	os.Exit(0)
}
