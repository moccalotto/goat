package motor

import (
	"fmt"
	"goat/glhelp"

	"github.com/go-gl/mathgl/mgl32"
)

type UncachedTransformation struct {
	x,
	y,
	scaleX,
	scaleY,
	r float32 // rotation
	allowNegative bool
}

// Cache transformation
type Transformation struct {
	trans       UncachedTransformation
	cacheValid  bool
	matrixCache mgl32.Mat3
}

func (T *UncachedTransformation) GetMatrix() mgl32.Mat3 {

	if !T.allowNegative && T.scaleX <= 0 {
		glhelp.GlPanic(fmt.Errorf("scale must be > 0. But [%.2f, %.2f] given", T.scaleX, T.scaleY))
	}
	if !T.allowNegative && T.scaleY <= 0 {
		glhelp.GlPanic(fmt.Errorf("scale must be > 0. But [%.2f, %.2f] given", T.scaleX, T.scaleY))
	}
	translate := mgl32.Translate2D(T.x, T.y)
	rotate := mgl32.HomogRotate2D(T.r)
	scale := mgl32.Scale2D(T.scaleX, T.scaleY)

	return glhelp.MatMulX3(translate, rotate, scale)
}

func (CT *Transformation) InvalidateCache() {
	CT.cacheValid = false
}

func CreateCachedTransformation() Transformation {
	return Transformation{
		trans: UncachedTransformation{
			x: 1, y: 1, scaleX: 1, scaleY: 1, r: 0,
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

	CT.cacheValid = CT.cacheValid && (x == CT.trans.x) && (y == CT.trans.y)

	CT.trans.x = x
	CT.trans.y = y
}

func (CT *Transformation) Move(dist mgl32.Vec2) {
	CT.SetPos(
		CT.trans.x+dist[0],
		CT.trans.y+dist[1],
	)
}

func (CT *Transformation) SetRotation(r float32) {

	CT.cacheValid = CT.cacheValid && (r == CT.trans.r)

	CT.trans.r = r
}

func (CT *Transformation) Rotate(r float32) {
	CT.SetRotation(
		CT.trans.r + r,
	)
}

func (CT *Transformation) SetScale(sx, sy float32) {

	CT.cacheValid = CT.cacheValid && (sx == CT.trans.scaleX) && (sy == CT.trans.scaleY)

	CT.trans.scaleX = sx
	CT.trans.scaleY = sy
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
	x := CT.trans.x + dx
	y := CT.trans.y + dy

	x = mgl32.Clamp(x, minX, maxX)
	y = mgl32.Clamp(y, minY, maxY)

	CT.SetPos(x, y)
}
