package shed

import (
	"math"
)

// =======================================
// ||
// || Vector types
// ||
// =======================================
type V4 struct{ C1, C2, C3, C4 float32 }
type V3 struct{ X, Y, Z float32 }
type V2 struct{ X, Y float32 }

// ||=======================================
// ||
// ||V2: 2D Vector
// ||
// ||=======================================

// Constructor
func Vec2(x, y float32) V2 {
	return V2{x, y}
}

// Construct V2 from polar coords
func PolarV2(radians, length float32) V2 {

	sin, cos := Sincos(radians)

	return V2{
		X: cos * length,
		Y: sin * length,
	}
}

// Useful for converting into an OpenGL vector, which is just an array
func (vec V2) ToArray() [2]float32 {
	return [2]float32{vec.X, vec.Y}
}

// Return the polar coordinates of the vector
func (vec V2) ToPolar() (angle float32, length float32) {
	angle = vec.Angle()
	length = vec.Len()

	return // Using named output args improves intellisense
}

// Length / magnitude of the vector
func (vec V2) Len() float32 {
	return Hypot(vec.X, vec.Y)
}

// Angle of the vector, relative to the x-axis
func (vec V2) Angle() float32 {
	return float32(math.Atan2(
		float64(vec.Y), float64(vec.X),
	))
}

func (vec V2) Plus(other V2) V2 {
	vec.X += other.X
	vec.Y += other.Y

	return vec
}
func (vec V2) Minus(other V2) V2 {
	vec.X -= other.X
	vec.Y -= other.Y

	return vec
}

func (vec V2) MultipliedBy(other V2) V2 {
	vec.X *= other.X
	vec.Y *= other.Y

	return vec
}

func (vec V2) DividedBy(other V2) V2 {
	vec.X /= other.X
	vec.Y /= other.Y

	return vec
}

func (vec V2) Rotated(radians float32) V2 {
	sin, cos := Sincos(radians)

	x, y := vec.X, vec.Y

	return V2{
		X: cos*x - sin*y,
		Y: sin*x + cos*y,
	}
}

// Multiply the vector's length by a given factor
func (vec V2) Scaled(factor float32) V2 {

	vec.X *= factor
	vec.Y *= factor

	return vec
}

// Turn the vector to point in a different direction
func (vec V2) WithAnle(radians float32) V2 {
	return PolarV2(radians, vec.Len())
}

// Return new vector with same angle but different length
func (vec V2) WithLength(newLen float32) V2 {
	curLen := vec.Len()
	factor := newLen / curLen

	vec.X *= factor
	vec.Y *= factor

	return vec
}

// Set the vector's length to 1
func (vec V2) Normalized() V2 {
	curLen := vec.Len()

	vec.X /= curLen
	vec.Y /= curLen

	return vec
}

// Clamp the vector's length within two values
func (vec V2) ClampedLength(min, max float32) V2 {

	if max < min {
		max, min = min, max
	}

	l := vec.Len()

	if l > max {
		return vec.WithLength(max)
	}
	if l < min {
		return vec.WithLength(min)
	}
	return vec
}

// Return a version of vec where angle is clamped between min and max
func (vec V2) ClampedAngle(min, max float32) V2 {
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

// Treat vec1 and vec2 as points.
// Return a "point" right between vec1 and vec2
func (vec1 V2) Between(vec2 V2) V2 {
	return vec2.Minus(vec1).Scaled(0.5).Plus(vec1)
}

// ||=======================================
// ||
// ||V3: 3D Vector
// ||
// ||=======================================

// Useful for converting into an OpenGL vector, which is just an array
func (vec V3) ToArray() [3]float32 {
	return [3]float32{vec.X, vec.Y, vec.Z}
}

func VecXYZ(x, y, z float32) V3 {
	return V3{x, y, z}
}

// ||=======================================
// ||
// ||V4: 4D Vector
// ||
// ||=======================================

// Useful for converting into an OpenGL vector, which is just an array
func (vec V4) ToArray() [4]float32 {
	return [4]float32{vec.C1, vec.C2, vec.C3, vec.C4}
}

func (vec V4) Plus(other V4) V4 {
	return V4{
		vec.C1 + other.C1,
		vec.C2 + other.C2,
		vec.C3 + other.C3,
		vec.C4 + other.C4,
	}
}

func Vec4(c1, c2, c3, c4 float32) V4 {
	return V4{c1, c2, c3, c4}
}

func RGBA(r, g, b, a float32) V4 {
	return V4{r, g, b, a}
}

func OPAQ_WHITE() V4 {
	return V4{1, 1, 1, 1}
}
