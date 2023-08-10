/*****************************************
* Contains all exported functions that
* are related to graphics and rendering.
******************************************/

package main

import "github.com/veandco/go-sdl2/sdl"

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
func (dm *Drawing) Background(c ...uint8) {
	switch len(c) {
	case 1:
		dm.Background(c[0], c[0], c[0], 255)
	case 3:
		dm.Background(c[0], c[1], c[2], 255)
	case 4:
		dm.bgColor = sdl.Color{R: c[0], G: c[1], B: c[2], A: c[3]}
	default:
		panic("Background() takes 1, 3, or 4 arguments.")
	}
}

func (dm *Drawing) Scale(scale float32) {
	dm.scaleX = scale
	dm.scaleY = scale
	dm.renderer.SetScale(scale, scale)
}

func (dm *Drawing) Line(x1, y1, x2, y2 float32) {
	dm.applySettingsToRenderer()
	dm.renderer.DrawLineF(x1, y1, x2, y2)
}

func (dm *Drawing) Dot(x, y float32) {
	dm.applySettingsToRenderer()
	dm.renderer.DrawPointF(x, y)
}

func (dm *Drawing) Rectangle(x1, y1, x2, y2 float32) {
	dm.applySettingsToRenderer()

	vertices := []sdl.Vertex{
		{sdl.FPoint{x1, y1}, dm.fgColor, sdl.FPoint{0, 0}},
		{sdl.FPoint{x2, y1}, dm.fgColor, sdl.FPoint{0, 0}},
		{sdl.FPoint{x2, y2}, dm.fgColor, sdl.FPoint{0, 0}},

		{sdl.FPoint{x2, y2}, dm.fgColor, sdl.FPoint{0, 0}},
		{sdl.FPoint{x1, y2}, dm.fgColor, sdl.FPoint{0, 0}},
		{sdl.FPoint{x1, y1}, dm.fgColor, sdl.FPoint{0, 0}},
	}
	dm.renderer.RenderGeometry(nil, vertices, nil)

	// dm.renderer.DrawRectF(&sdl.FRect{
	// 	X: x1,
	// 	Y: y1,
	// 	H: y2 - y1,
	// 	W: x2 - x1,
	// })
}
