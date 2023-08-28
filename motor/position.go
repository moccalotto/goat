package motor

import (
	"fmt"
	"goat/glhelp"

	"github.com/go-gl/mathgl/mgl32"
)

type UncachedPosition struct {
	// Location
	X             float32 // position on x-axis
	Y             float32 // position on y-axis
	MinX          float32 // min x value, only used if ClampLocation is true
	MaxX          float32 // max y-value, only used if ClampLocation is true
	MinY          float32 // same as minX
	MaxY          float32 // same as MaxX
	LimitLocation bool    // should coordinates, scale, and rotation be clamped to MinX, MaxX, etc.

	// Scale
	ScaleX        float32 // size in x-direction of the object as it was created, not as it is facing now
	ScaleY        float32 //
	MinScaleX     float32 // Minimum X-scale
	MaxScaleX     float32 // Maximum X-scale
	MinScaleY     float32 //
	MaxScaleY     float32 //
	LimitScale    bool    // Should we enforce the scale limits. Clamping will be performed when GetMatrix is called
	AllowNegScale bool    // Allow negative scales (which would flip rendered models) This is separate from ClampScale to improve default behavior when creating a Transformation where all values are zero/false

	// Rotation
	Rot           float32 // rotation
	MaxRot        float32 // Max angle float32 relative to X-axis
	MinRot        float32 // Min angle, relative to X-axis
	LimitRotation bool    // Should we calmp Rot to its min and max values
}

// Cache transformation
type Position struct {
	trans       UncachedPosition
	cacheValid  bool
	matrixCache mgl32.Mat3
}

func (P *UncachedPosition) GetMatrix() mgl32.Mat3 {

	if (!P.AllowNegScale) && (P.ScaleX <= 0 || P.ScaleY <= 0) {
		glhelp.GlPanic(fmt.Errorf("scale must be > 0. But [%.2f, %.2f] given", P.ScaleX, P.ScaleY))
	}

	translate := mgl32.Translate2D(P.X, P.Y)
	rotate := mgl32.HomogRotate2D(P.Rot)
	scale := mgl32.Scale2D(P.ScaleX, P.ScaleY)

	return translate.Mul3(rotate).Mul3(scale)
}

func (P *UncachedPosition) normalizeRotation() {
	P.Rot = glhelp.NormalizeAngle(P.Rot)
}

func (P *Position) GetNormalizedRotation() float32 {
	return glhelp.NormalizeAngle(P.trans.Rot)
}

// clampToLimits position and rotation into their restricted values
func (P *Position) clampToLimits() {

	t := P.trans
	if t.LimitLocation && (t.X > t.MaxX || t.X < t.MinX || t.Y > t.MaxY || t.Y < t.MinY) {
		t.X = mgl32.Clamp(t.X, t.MinX, t.MaxX)
		t.Y = mgl32.Clamp(t.Y, t.MinY, t.MaxY)
		P.cacheValid = false
	}

	if t.LimitScale && (t.ScaleX < t.MinScaleX || t.ScaleX > t.MaxScaleX || t.ScaleY < t.MinScaleY || t.ScaleY < t.MaxScaleY) {
		t.ScaleX = mgl32.Clamp(t.ScaleX, t.MinScaleX, t.MaxScaleX)
		t.ScaleY = mgl32.Clamp(t.ScaleY, t.MinScaleY, t.MaxScaleY)
		P.cacheValid = false
	}

	if t.LimitRotation && (t.Rot < t.MinRot || t.Rot > t.MaxRot) {
		t.normalizeRotation()
		t.Rot = mgl32.Clamp(t.Rot, t.MinRot, t.MaxRot)
		P.cacheValid = false
	}
}

func (P *Position) ApplyForce(f Force, delta float32) {
	P.cacheValid = false

	P.SetPos(
		P.trans.X+f.Vec.X*delta,
		P.trans.Y+f.Vec.Y*delta,
	)

	P.SetRotation(
		P.trans.Rot + f.Rot*delta,
	)
	P.trans.Rot += f.Rot * delta
}

func (P *Position) InvalidateCache() {
	P.cacheValid = false
}

func CreateCachedTransformation() Position {
	return Position{
		trans: UncachedPosition{
			X: 1, Y: 1, ScaleX: 1, ScaleY: 1, Rot: 0,
		},
		cacheValid:  true,
		matrixCache: mgl32.Ident3(),
	}
}

func (P *Position) SetTransformation(t UncachedPosition) {

	P.cacheValid = false
	P.trans = t
}

func (P *Position) SetPos(x, y float32) {

	P.cacheValid = P.cacheValid && (x == P.trans.X) && (y == P.trans.Y)

	if P.trans.LimitLocation {
		x = mgl32.Clamp(x, P.trans.MinX, P.trans.MaxX)
		y = mgl32.Clamp(y, P.trans.MinY, P.trans.MaxY)
	}

	P.trans.X = x
	P.trans.Y = y
}

func (P *Position) ApplyLimits(location, rotation, scale bool) {
	// location
	P.trans.LimitLocation = location

	// rotation
	P.trans.LimitRotation = rotation

	// scale
	P.trans.LimitScale = scale
}

func (P *Position) Move(dist glhelp.V2) {
	P.SetPos(
		P.trans.X+dist.X,
		P.trans.Y+dist.Y,
	)
}

func (P *Position) SetRotation(r float32) {

	r = glhelp.NormalizeAngle(r)

	if P.trans.LimitRotation {
		r = mgl32.Clamp(r, P.trans.MinRot, P.trans.MaxRot)
	}

	P.cacheValid = P.cacheValid && (r == P.trans.Rot)

	P.trans.Rot = r
}

func (P *Position) Rotate(r float32) {
	P.SetRotation(P.trans.Rot + r)
}

// combination of lerp and rotate
func (P *Position) RotateTowards(target, amount float32) {

	target = glhelp.NormalizeAngle(target)

	P.SetRotation(
		glhelp.Lerp(P.trans.Rot, target, amount),
	)
}

func (P *Position) Get() (x, y, rot float32) {
	x = P.trans.X
	y = P.trans.Y
	rot = P.trans.Rot

	return
}

// if P.trans.Rot is within [distance] radians of target, then set rotation = target
func (P *Position) SnapRotationTo(target, distance float32) {

	target = glhelp.NormalizeAngle(target)

	d := mgl32.Abs(target - P.trans.Rot)

	if d < mgl32.Epsilon {
		return // distance to target is functionally zero. We're done.
	}

	if d > distance {
		return // too far away from target. Can't snap
	}

	P.SetRotation(target)
}

func (P *Position) SetScale(sx, sy float32) {

	P.cacheValid = P.cacheValid && (sx == P.trans.ScaleX) && (sy == P.trans.ScaleY)

	P.trans.ScaleX = sx
	P.trans.ScaleY = sy
}

func (P *Position) LimitScale(minX, minY, maxX, maxY float32) {
	P.trans.MinScaleX = minX
	P.trans.MinScaleY = minY
	P.trans.MaxScaleX = maxX
	P.trans.MaxScaleY = maxY
	P.trans.LimitScale = true
}

func (P *Position) LimitLocation(minX, minY, maxX, maxY float32) {
	P.trans.MinX = minX
	P.trans.MinY = minY
	P.trans.MaxX = maxX
	P.trans.MaxY = maxY
	P.trans.LimitLocation = true
}

func (P *Position) LimitRotation(min, max float32) {
	P.trans.MinRot = min
	P.trans.MaxRot = max
	P.trans.LimitRotation = true
}

func (P *Position) GetMatrix() mgl32.Mat3 {
	P.clampToLimits()
	if !P.cacheValid {
		P.matrixCache = P.trans.GetMatrix()
	}

	return P.matrixCache
}
