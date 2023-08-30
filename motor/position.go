package motor

import (
	"fmt"
	h "goat/glhelp"

	"github.com/go-gl/mathgl/mgl32"
)

// Cache transformation
type Position struct {
	cacheValid  bool
	matrixCache mgl32.Mat3
	// Location
	x             float32 // position on x-axis
	y             float32 // position on y-axis
	minX          float32 // min x value, only used if ClampLocation is true
	maxX          float32 // max y-value, only used if ClampLocation is true
	minY          float32 // same as minX
	maxY          float32 // same as MaxX
	limitLocation bool    // should coordinates, scale, and rotation be clamped to MinX, MaxX, etc.

	// Scale
	scaleX        float32 // size in x-direction of the object as it was created, not as it is facing now
	scaleY        float32 //
	minScaleX     float32 // Minimum X-scale
	maxScaleX     float32 // Maximum X-scale
	minScaleY     float32 //
	maxScaleY     float32 //
	limitScale    bool    // Should we enforce the scale limits. Clamping will be performed when GetMatrix is called
	allowNegScale bool    // Allow negative scales (which would flip rendered models) This is separate from ClampScale to improve default behavior when creating a Transformation where all values are zero/false

	// Rotation
	angleOffset float32 // This offsete is added when rottation matrix is calculated. Useful for compensating for rotated textures.
	angle       float32 // rotation
	minAngle    float32 // Min angle, relative to X-axis
	maxAngle    float32 // Max angle float32 relative to X-axis
	limitAngle  bool    // Should we calmp Rot to its min and max values
}

func (P *Position) createMatrix() mgl32.Mat3 {

	if (!P.allowNegScale) && (P.scaleX <= 0 || P.scaleY <= 0) {
		h.GlPanic(fmt.Errorf("scale must be > 0. But [%.2f, %.2f] given", P.scaleX, P.scaleY))
	}

	translate := mgl32.Translate2D(P.x, P.y)
	rotate := mgl32.HomogRotate2D(P.angle + P.angleOffset)
	scale := mgl32.Scale2D(P.scaleX, P.scaleY)

	return translate.Mul3(rotate).Mul3(scale)
}

// Make sure that X,Y, Scale, and Angle are all within the limits we (might) have set
// Thins function only invalidates cache if necessary
func (P *Position) clampToLimits() {

	if P.limitLocation && (P.x > P.maxX || P.x < P.minX || P.y > P.maxY || P.y < P.minY) {
		P.x = mgl32.Clamp(P.x, P.minX, P.maxX)
		P.y = mgl32.Clamp(P.y, P.minY, P.maxY)
		P.cacheValid = false
	}

	if P.limitScale && (P.scaleX < P.minScaleX || P.scaleX > P.maxScaleX || P.scaleY < P.minScaleY || P.scaleY < P.maxScaleY) {
		P.scaleX = mgl32.Clamp(P.scaleX, P.minScaleX, P.maxScaleX)
		P.scaleY = mgl32.Clamp(P.scaleY, P.minScaleY, P.maxScaleY)
		P.cacheValid = false
	}

	if P.limitAngle && (P.angle < P.minAngle || P.angle > P.maxAngle) {
		P.angle = mgl32.Clamp(P.angle, P.minAngle, P.maxAngle)
		P.cacheValid = false
	}
}

func (P *Position) ApplyForce(f Force, delta float32) {
	P.cacheValid = false

	P.SetXY(
		P.x+f.Vec.X*delta,
		P.y+f.Vec.Y*delta,
	)

	P.Rotate(f.Rot * delta)
	P.angle += f.Rot * delta
}

func CreatePosition() Position {
	return Position{
		x:           1,
		y:           1,
		scaleX:      1,
		scaleY:      1,
		angle:       0,
		cacheValid:  true,
		matrixCache: mgl32.Ident3(),
	}
}

func (P *Position) SetXY(x, y float32) {

	P.cacheValid = P.cacheValid && (x == P.x) && (y == P.y)

	if P.limitLocation {
		x = mgl32.Clamp(x, P.minX, P.maxX)
		y = mgl32.Clamp(y, P.minY, P.maxY)
	}

	P.x = x
	P.y = y
}

func (P *Position) SetLimits(location, rotation, scale bool) {
	// location
	P.limitLocation = location

	// rotation
	P.limitAngle = rotation

	// scale
	P.limitScale = scale
}

func (P *Position) Move(dist h.V2) {
	P.SetXY(
		P.x+dist.X,
		P.y+dist.Y,
	)
}

// This offset is added to angle for the purposes of calculating the rotation matrix.
// This is convenient if want to think of zero degrees as forward, even though it might be a different direction
// on the screen. It is also useful for rotating textures.
func (P *Position) SetAngleOffset(r float32) {
	P.cacheValid = P.cacheValid && (r == P.angleOffset)

	P.angleOffset = r
}

func (P *Position) SetAngle(r float32) {

	if P.limitAngle {
		r = mgl32.Clamp(r, P.minAngle, P.maxAngle)
	}

	P.cacheValid = P.cacheValid && (r == P.angle)

	P.angle = r
}

func (P *Position) Rotate(r float32) {
	P.SetAngle(P.angle + r)
}

// combination of lerp and rotate
func (P *Position) RotateTowards(target, amount float32) {

	P.SetAngle(
		h.Lerp(P.angle, target, amount),
	)
}

func (P *Position) GetXYA() (x, y, angle float32) {
	x = P.x
	y = P.y
	angle = P.angle

	return
}

// if P.Rot is within [distance] radians of target, then set rotation = target
func (P *Position) SnapAngleTo(target, distance float32) {

	d := mgl32.Abs(target - P.angle)

	if d < mgl32.Epsilon {
		return // distance to target is functionally zero. We're done.
	}

	if d > distance {
		return // too far away from target. Can't snap
	}

	P.SetAngle(target)
}

func (P *Position) SetScale(sx, sy float32) {

	P.cacheValid = P.cacheValid && (sx == P.scaleX) && (sy == P.scaleY)

	P.scaleX = sx
	P.scaleY = sy
}

func (P *Position) LimitScale(minX, minY, maxX, maxY float32) {
	P.minScaleX = minX
	P.minScaleY = minY
	P.maxScaleX = maxX
	P.maxScaleY = maxY
	P.limitScale = true
}

func (P *Position) LimitLocation(minX, minY, maxX, maxY float32) {
	P.minX = minX
	P.minY = minY
	P.maxX = maxX
	P.maxY = maxY
	P.limitLocation = true
}

func (P *Position) LimitAngle(min, max float32) {
	P.minAngle = min
	P.maxAngle = max
	P.limitAngle = true
}

func (P *Position) GetMatrix() mgl32.Mat3 {
	P.clampToLimits()
	if !P.cacheValid {
		P.matrixCache = P.createMatrix()
		P.cacheValid = true
	}

	return P.matrixCache
}
