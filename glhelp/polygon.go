package glhelp

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Create a 1x1 square that is axis-aligned and centered on the origin
func PolygonCoords(sides int, angleOffset, width, height float64) (verts []float32, texCoords []float32, indeces []uint32) {
	anglePerSide := Tau / float64(sides)

	triangleCount := sides - 2

	maxX, maxY := float32(0), float32(0)

	for i := 0; i < sides; i++ {
		angle := anglePerSide*float64(i) + angleOffset

		y, x := math.Sincos(angle)
		fx, fy := float32(x*width), float32(y*height)

		maxY = Max(maxY, mgl32.Abs(fy))
		maxX = Max(maxX, mgl32.Abs(fy))

		texCoords = append(texCoords, fx, fy)

		verts = append(verts, fx, fy, 1.0)
	}

	spanX, spanY := 2*maxX, 2*maxY

	for i := 0; i < sides*2; i += 2 {
		/*
			If fx is positive then tx must be ≥ 0.5
			if fx == maxX then tx must be 1.0
			if fx == -maxX then tx must be 0.0

			in other words

			if fx + maxX == 1.0 * spanX then tx must be 1.0
			if fx + maxX == 0.5 * spanX then tx must be 0.5
			if fx + maxX == 0.0 * spanX then tx must be 0

			ipso facto:

			tx = (fx + maxX) / spanX
		*/
		fx, fy := texCoords[i], texCoords[i+1]

		tx := (fx + maxX) / spanX
		ty := (fy + maxY) / spanY

		texCoords[i] = tx
		texCoords[i+1] = 1 - ty
	}

	for i := 0; i < triangleCount; i++ {
		indeces = append(indeces,
			0,
			uint32(i+1),
			uint32(i+2),
		)
	}
	return
}
