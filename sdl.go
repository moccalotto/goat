package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

func setupSDL() (*sdl.Window, *sdl.Renderer) {
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		"Goat Window",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		500,
		500,
		sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL,
	)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(
		window,
		-1, // TODO: select the opengl driver.
		sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC,
	)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	if err != nil {
		log.Fatalf("Could not create SDL renderer: %+v", err)
	}
	return window, renderer
}
