package motor

import (
	"goat/glhelp"

	"github.com/go-gl/mathgl/mgl32"
)

type UncachedTransformation struct {
	X,
	Y,
	ScaleX,
	ScaleY,
	R float32 // rotation
}

// Cache transformation
type Transformation struct {
	trans       UncachedTransformation
	cacheValid  bool
	matrixCache mgl32.Mat3
}

func (T *UncachedTransformation) GetMatrix() mgl32.Mat3 {

	translate := mgl32.Translate2D(T.X, T.Y)
	rotate := mgl32.HomogRotate2D(T.R)
	scale := mgl32.Scale2D(T.ScaleX, T.ScaleY)

	return glhelp.MatMulX3(translate, rotate, scale)
}

func (CT *Transformation) InvalidateCache() {
	CT.cacheValid = false
}

func CreateCachedTransformation() Transformation {
	return Transformation{
		trans: UncachedTransformation{
			X: 1, Y: 1, ScaleX: 1, ScaleY: 1, R: 0,
		},
		cacheValid:  true,
		matrixCache: mgl32.Ident3(),
	}
}

func (CT *Transformation) SetTransformation(t UncachedTransformation) {

	CT.cacheValid = false
	CT.trans = t
}

func (CT *Transformation) SetPos(x, y float32) {

	CT.cacheValid = CT.cacheValid && (x == CT.trans.X) && (y == CT.trans.Y)

	CT.trans.X = x
	CT.trans.Y = y
}

func (CT *Transformation) Move(dist mgl32.Vec2) {
	CT.SetPos(
		CT.trans.X+dist[0],
		CT.trans.Y+dist[1],
	)
}

func (CT *Transformation) SetRotation(r float32) {

	CT.cacheValid = CT.cacheValid && (r == CT.trans.R)

	CT.trans.R = r
}

func (CT *Transformation) Rotate(r float32) {
	CT.SetRotation(
		CT.trans.R + r,
	)
}

func (CT *Transformation) SetScale(sx, sy float32) {

	CT.cacheValid = CT.cacheValid && (sx == CT.trans.ScaleX) && (sy == CT.trans.ScaleY)

	CT.trans.ScaleX = sx
	CT.trans.ScaleY = sy
}

func (CT *Transformation) GetMatrix() mgl32.Mat3 {
	if !CT.cacheValid {
		CT.matrixCache = CT.trans.GetMatrix()
	}

	return CT.matrixCache

}

func (CT *Transformation) GetAll() UncachedTransformation {
	return CT.trans
}

func (CT *Transformation) RestrictedMove(dx, dy, minX, minY, maxX, maxY float32) {
	x := CT.trans.X + dx
	y := CT.trans.Y + dy

	x = mgl32.Clamp(x, minX, maxX)
	y = mgl32.Clamp(y, minY, maxY)

	CT.SetPos(x, y)
}
