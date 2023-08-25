package motor

import (
	"goat/glhelp"

	"github.com/go-gl/mathgl/mgl32"
)

type Transformation struct {
	X,
	Y,
	ScaleX,
	ScaleY,
	R float32 // rotation
}

type CachedTransformation struct {
	transformation Transformation
	cacheValid     bool
	matrixCache    mgl32.Mat3
}

func (T *Transformation) GetMatrix() mgl32.Mat3 {

	translate := mgl32.Translate2D(T.X, T.Y)
	rotate := mgl32.HomogRotate2D(T.R)
	scale := mgl32.Scale2D(T.ScaleX, T.ScaleY)

	return glhelp.MatMulMany(translate, rotate, scale)
}

func (CT *CachedTransformation) InvalidateCache() {
	CT.cacheValid = false
}

func CreateCachedTransformation() CachedTransformation {
	return CachedTransformation{
		transformation: Transformation{
			X: 1, Y: 1, ScaleX: 1, ScaleY: 1, R: 0,
		},
		cacheValid:  true,
		matrixCache: mgl32.Ident3(),
	}
}

func (CT *CachedTransformation) SetTransformation(t Transformation) {

	CT.cacheValid = false
	CT.transformation = t
}

func (CT *CachedTransformation) SetPos(x, y float32) {

	CT.cacheValid = CT.cacheValid && (x == CT.transformation.X) && (y == CT.transformation.Y)

	CT.transformation.X = x
	CT.transformation.Y = y
}

func (CT *CachedTransformation) Move(dist mgl32.Vec2) {
	CT.SetPos(
		CT.transformation.X+dist[0],
		CT.transformation.Y+dist[1],
	)
}

func (CT *CachedTransformation) SetRotation(r float32) {

	CT.cacheValid = CT.cacheValid && (r == CT.transformation.R)

	CT.transformation.R = r
}

func (CT *CachedTransformation) Rotate(r float32) {
	CT.SetRotation(
		CT.transformation.R + r,
	)
}

func (CT *CachedTransformation) SetScale(sx, sy float32) {

	CT.cacheValid = CT.cacheValid && (sx == CT.transformation.ScaleX) && (sy == CT.transformation.ScaleY)

	CT.transformation.ScaleX = sx
	CT.transformation.ScaleY = sy
}

func (CT *CachedTransformation) GetMatrix() mgl32.Mat3 {
	if !CT.cacheValid {
		CT.matrixCache = CT.transformation.GetMatrix()
	}

	return CT.matrixCache

}

func (CT *CachedTransformation) GetAll() Transformation {
	return CT.transformation
}
