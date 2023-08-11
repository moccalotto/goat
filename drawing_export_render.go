/*****************************************
* Contains all exported functions that
* are related to graphics and rendering.
*
*
* TODO: All calls to SDL must be queued,
* delayed and/or batched such that no SDL
* calls are made while the lua script is
* running - maybe except SDL_Delay et al.
*
* I should consider using entities for
* each object i draw on screen
*
******************************************/

package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

func sdlPoint64(x, y float64) sdl.FPoint {
	return sdl.FPoint{
		X: float32(x),
		Y: float32(y),
	}
}

func sdlVert64(x, y float64, color sdl.Color) sdl.Vertex {
	return sdl.Vertex{
		Position: sdl.FPoint{
			X: float32(x),
			Y: float32(y),
		},
		Color:    color,
		TexCoord: sdlPoint64(0, 0),
	}
}

func (dm *Drawing) Color(c ...uint8) (uint8, uint8, uint8, uint8) {
	switch len(c) {
	case 0:
		return dm.fgColor.R, dm.fgColor.G, dm.fgColor.B, dm.fgColor.A
	case 1:
		return dm.Color(c[0], c[0], c[0], 255)
	case 3:
		return dm.Color(c[0], c[1], c[2], 255)
	case 4:
		dm.fgColor = sdl.Color{R: c[0], G: c[1], B: c[2], A: c[3]}
		return dm.Color()
	default:
		panic("Background() takes 1, 3, or 4 arguments.")
	}
}
func (dm *Drawing) Background(c ...uint8) (uint8, uint8, uint8, uint8) {
	switch len(c) {
	case 0:
		return dm.bgColor.R, dm.bgColor.G, dm.bgColor.B, dm.bgColor.A
	case 1:
		return dm.Background(c[0], c[0], c[0], 255)
	case 3:
		return dm.Background(c[0], c[1], c[2], 255)
	case 4:
		dm.bgColor = sdl.Color{R: c[0], G: c[1], B: c[2], A: c[3]}
		return dm.Background()
	default:
		panic("Background() takes 1, 3, or 4 arguments.")
	}
}

func (dm *Drawing) Scale(scale float32) {
	dm.scaleX = scale
	dm.scaleY = scale
	dm.renderer.SetScale(scale, scale)
}

func (dm *Drawing) Line(x1, y1, x2, y2 float64) {
	dm.applySettingsToRenderer()
	dm.renderer.DrawLineF(
		float32(x1),
		float32(y1),
		float32(x2),
		float32(y2),
	)
}

func (dm *Drawing) Dot(x, y float64) {
	dm.applySettingsToRenderer()
	dm.renderer.DrawPointF(float32(x), float32(y))
}

func (dm *Drawing) Rectangle(x1, y1, x2, y2 float64) {
	dm.applySettingsToRenderer()

	vertices := []sdl.Vertex{
		sdlVert64(x1, y1, dm.fgColor),
		sdlVert64(x2, y1, dm.fgColor),
		sdlVert64(x2, y2, dm.fgColor),
		sdlVert64(x2, y2, dm.fgColor),
		sdlVert64(x1, y2, dm.fgColor),
		sdlVert64(x1, y1, dm.fgColor),
	}
	dm.renderer.RenderGeometry(nil, vertices, nil)
}

func (dm *Drawing) Polygon(centerX, centerY, radius, angle float64, vertices int) {

	totalAngle := math.Pi * 2
	trangleCount := vertices - 2
	radiansPerSlice := totalAngle / float64(vertices)

	edgeVertices := []sdl.Vertex{}

	for i := 0; i < vertices; i++ {

		_angle := angle + float64(i)*radiansPerSlice

		_x, _y := math.Sincos(_angle)
		x := _x*radius + centerX
		y := _y*radius + centerY

		vertex := sdlVert64(x, y, dm.fgColor)

		edgeVertices = append(edgeVertices, vertex)
	}

	resultVerts := []sdl.Vertex{}

	for i := 0; i < trangleCount; i++ {
		resultVerts = append(resultVerts,
			edgeVertices[0],
			edgeVertices[i+1],
			edgeVertices[i+2],
		)
	}

	dm.renderer.RenderGeometry(nil, resultVerts, nil)
}
