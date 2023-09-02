/*****************************************
* Contains all exported functions that
* are related to graphics and rendering.
******************************************/

package pilot

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

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

func (dm *Drawing) Polygon(centerX, centerY, radius, angle float64, edgeVertexCount int) []mgl32.Vec3 {

	if edgeVertexCount < 3 {
		panic("verticeCount must be at least 3")
	}

	const tau = 2 * math.Pi
	radiansPerSlice := tau / float64(edgeVertexCount)
	triangleCount := edgeVertexCount - 2

	tmp := make([]mgl32.Vec3, edgeVertexCount)
	for i := 0; i < edgeVertexCount; i++ {
		currentAngle := radiansPerSlice * float64(i)
		_sin, _cos := math.Sincos(currentAngle)
		tmp[i] = mgl32.Vec3{
			/* x: */ float32(_cos * currentAngle),
			/* y: */ float32(_sin * currentAngle),
			/* z: */ 0.0,
		}
	}

	result := make([]mgl32.Vec3, 0)
	for i := 0; i < triangleCount; i++ {
		result = append(result,
			tmp[i],
			tmp[i+1],
			tmp[i+2],
		)
	}

	return result
}
