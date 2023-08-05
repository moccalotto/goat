package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

// My dreap is to do all of this:: https://github.com/fogleman/gg

type Drawing struct {
	// width  int
	// height int
	// color  color.Color
	// fillPattern   Pattern
	// strokePattern Pattern
	// strokePath    raster.Path
	// fillPath      raster.Path
	// start         Point
	// current       Point
	// hasCurrent    bool
	// dashes        []float64
	// dashOffset    float64
	// lineWidth     float64
	// lineCap       LineCap
	// lineJoin      LineJoin
	// fillRule      FillRule
	// fontFace      font.Face
	// fontHeight    float64
	// stack       []*DrawingMachine
	// transformer *mat32.Mat3
	renderer *sdl.Renderer
}

func CreateDrawing(renderer *sdl.Renderer) *Drawing {
	return &Drawing{
		renderer: renderer,
	}
}

// Draw a line
func (dm *Drawing) Line(x1, y1, x2, y2 int32) {
	dm.renderer.DrawLine(x1, y1, x2, y2)
}

// Get viewport width
func (dm *Drawing) W() int32 {
	return dm.renderer.GetViewport().W
}

// Get viewport height
func (dm *Drawing) H() int32 {
	return dm.renderer.GetViewport().H
}
