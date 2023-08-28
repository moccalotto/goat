package glhelp

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// /
// /
// /
// / TYPES
// /////////////////////////////////7
type V4 struct{ X, Y, Z, W float32 }
type V3 struct{ X, Y, Z float32 }
type V2 struct{ X, Y float32 }
type Color V4

// /
// /
// /
// / V2 Implementation
// //////////////////////////////////////////
func Vec2FromMgl(mvec mgl32.Vec2) (vec V2) {
	vec.X = mvec[0]
	vec.Y = mvec[1]
	return
}

func Vec2(x, y float32) V2 {
	return V2{x, y}
}

func PolarV2(radians, length float32) V2 {

	sin, cos := Sincos(radians)

	return V2{
		X: cos * length,
		Y: sin * length,
	}
}

func (vec V2) ToComplex() complex64 {
	return complex(vec.X, vec.Y)
}

// Return the polar coordinates of the vector
func (vec V2) ToPolar() (angle float32, length float32) {
	angle = vec.Angle()
	length = vec.Len()

	return // Using named output args improves intellisense
}

func (vec V2) Len() float32 {
	return Hypot(vec.X, vec.Y)
}

func (vec V2) Angle() float32 {
	return float32(math.Atan2(
		float64(vec.Y), float64(vec.X),
	))
}

func (vec V2) Add(other V2) V2 {
	vec.X += other.X
	vec.Y += other.Y

	return vec
}
func (vec V2) Sub(other V2) V2 {
	vec.X -= other.X
	vec.Y -= other.Y

	return vec
}

func (vec V2) Mul(other V2) V2 {
	vec.X *= other.X
	vec.Y *= other.Y

	return vec
}

func (vec V2) Div(other V2) V2 {
	vec.X /= other.X
	vec.Y /= other.Y

	return vec
}

func (vec V2) Rotate(radians float32) V2 {
	sin, cos := Sincos(radians)

	x, y := vec.X, vec.Y

	return V2{
		X: cos*x - sin*y,
		Y: sin*x + cos*y,
	}
}

// Multiply the vector's length by a given factor
func (vec V2) Scale(factor float32) V2 {

	vec.X *= factor
	vec.Y *= factor

	return vec
}

// Turn the vector to point in a different direction
func (vec V2) SetAngle(radians float32) V2 {
	return PolarV2(radians, vec.Len())
}

func (vec V2) SetLength(newLen float32) V2 {
	curLen := vec.Len()
	factor := newLen / curLen

	vec.X *= factor
	vec.Y *= factor

	return vec
}

// Set the vector's length to 1
func (vec V2) Normalize() V2 {
	curLen := vec.Len()

	vec.X /= curLen
	vec.Y /= curLen

	return vec
}

// Clamp the vector's length within two values
func (vec V2) ClampLen(min, max float32) V2 {

	if max < min {
		max, min = min, max
	}

	l := vec.Len()

	if l > max {
		return vec.SetLength(max)
	}
	if l < min {
		return vec.SetLength(min)
	}
	return vec
}

func (vec V2) ClampAngle(min, max float32) V2 {
	if max < min {
		max, min = min, max
	}

	angle, length := vec.ToPolar()

	if angle > max {
		return PolarV2(max, length)
	}

	if angle < min {
		return PolarV2(min, length)
	}

	return vec
}

func (vec V2) ToVec3(z float32) V3 {
	return V3{vec.X, vec.Y, z}
}

func (vec V2) ToVec4(z, w float32) V4 {
	return V4{vec.X, vec.Y, z, w}
}

// Turn into array, useful for sending data to shaders
func (vec V2) ToArray() [2]float32 {
	return [2]float32{vec.X, vec.Y}
}

// /
// /
// /
// / V3 Implementation
// //////////////////////////////////////////