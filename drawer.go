package main

import "github.com/veandco/go-sdl2/sdl"

// My dreap is to do all of this:: https://github.com/fogleman/gg

type DrawingMachine struct {
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

func (dm DrawingMachine) Line(x1, y1, x2, y2 float32) {
	// TODO: transform points

	// TODO: consider rounding before converting to int.
	dm.renderer.DrawLine(
		int32(x1),
		int32(y1),
		int32(x2),
		int32(y2),
	)
}
