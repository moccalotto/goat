package main

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
)

type Config struct {
	Width      int32
	Height     int32
	Title      string
	CanResize  bool
	Scale      float32
	Background sdl.Color
}

func NewConfig() *Config {
	return &Config{
		Title:     "Goaty McWindow",
		Width:     800,
		Height:    600,
		CanResize: false,
		Scale:     1,
		Background: sdl.Color{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		},
	}
}

func (cfg *Config) configSanityCheck() error {
	if cfg.Width < 50 {
		return errors.New("width too low")
	}
	if cfg.Height < 50 {
		return errors.New("height too low")
	}
	if len(cfg.Title) > 100 {
		return errors.New("title longer than 100 characters")
	}
	if cfg.Scale < 1 || cfg.Scale > 50 {
		return errors.New("scale must be between 1 and 50")
	}

	return nil
}

func (cfg *Config) canResizeSdlFlag() uint32 {
	if !cfg.CanResize {
		return 0
	}
	return sdl.WINDOW_RESIZABLE
}
