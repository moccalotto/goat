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

func (dm *Drawing) Color(c ...uint8) (uint8, uint8, uint8, uint8) {
	switch len(c) {
	case 0:
		return dm.fgColor.R, dm.fgColor.G, dm.fgColor.B, dm.fgColor.A
	case 1:
		return dm.Color(c[0], c[0], c[0], 255)
	case 3:
		return dm.Color(c[0], c[1], c[2], 255)
	case 4:
		dm.fgColor = Color{R: c[0], G: c[1], B: c[2], A: c[3]}
		return dm.Color()
	default:
		panic("Color() takes 1, 3, or 4 arguments.")
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
		dm.bgColor = Color{R: c[0], G: c[1], B: c[2], A: c[3]}
		return dm.Background()
	default:
		panic("Background() takes 1, 3, or 4 arguments.")
	}
}

func (dm *Drawing) Scale(scale float32) {
}

func (dm *Drawing) Line(x1, y1, x2, y2 float64) {
}

func (dm *Drawing) Dot(x, y float64) {
}

func (dm *Drawing) Rectangle(x1, y1, x2, y2 float64) {
}

func (dm *Drawing) Polygon(centerX, centerY, radius, angle float64, vertices int) {
	/*
		https://github.com/go-gl/mathgl/blob/e426c0894fa41bc41ac04704eef3ac4011b19ecf/mgl32/shapes.go#L34
	*/
}
