package main

import (
	"time"

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
	W        int32
	H        int32
}

func CreateDrawing(renderer *sdl.Renderer) *Drawing {
	renderer.SetScale(10, 10)
	return &Drawing{
		renderer: renderer,
		W:        renderer.GetViewport().W,
		H:        renderer.GetViewport().H,
	}
}

func (dm *Drawing) updateLocalVariables() {
	dm.W = dm.renderer.GetViewport().W
	dm.H = dm.renderer.GetViewport().H
}

func (dm *Drawing) Scale(scale float32) {
	dm.renderer.SetScale(scale, scale)
}

/**
 * Draw a line
 *
 * TODO:
 * A line is a rect with a given color and size
 * It is rotated and translated to match the given
 * start and end points. Then it is copied into the
 * screen via the renderer
 */
func (dm *Drawing) Line(x1, y1, x2, y2 float32) {
	dm.renderer.DrawLineF(x1, y1, x2, y2)
}

// Get the viewport
func (dm *Drawing) Size() sdl.Rect {
	return dm.renderer.GetViewport()
}

func (dm *Drawing) Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
