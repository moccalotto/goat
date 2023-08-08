package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

func setupSDL(cfg *Config) (*sdl.Window, *sdl.Renderer) {
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		cfg.Title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		cfg.Width,
		cfg.Height,
		sdl.WINDOW_SHOWN|cfg.canResizeSdlFlag(),
	)

	if err != nil {
		panic(err)
	}
	renderer, err := sdl.CreateRenderer(
		window,
		-1,
		sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC,
	)
	if err != nil {
		log.Fatalf("Could not create SDL renderer: %+v", err)
	}

	renderer.SetScale(cfg.Scale, cfg.Scale)

	return window, renderer
}
