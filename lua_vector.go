package main

import (
	"fmt"
	"math"
)

/**
 * A 2d vector type for use inside lua scripts.
 * It uses float64, so we shouldnt mix it with the 32-bit vectors
 * we use for openGL
 */
type LuaVector struct {
	x float64
	y float64
}

func CreateLuaVector(x, y float64) LuaVector {
	return LuaVector{x: x, y: y}
}

// Create a vector from polar coordinates
func CreatePolarLuaVector(radians float64, magnitude ...float64) LuaVector {

	sin, cos := math.Sincos(radians)

	var length float64 = 1
	if len(magnitude) > 0 {
		length = magnitude[0]
	}

	return LuaVector{
		x: cos * length,
		y: sin * length,
	}
}

// Get the x-component of a vector
func (vec LuaVector) X() float64 {
	return vec.x
}

// Get the y-component of a vector
func (vec LuaVector) Y() float64 {
	return vec.y // return the original value of x
}

// Return a clone of this vector
func (vec LuaVector) Clone() LuaVector {
	return CreateLuaVector(vec.x, vec.y)
}

// Return a new vector equal to (vec - other)
func (vec LuaVector) Sub(other LuaVector) LuaVector {
	return CreateLuaVector(
		vec.x-other.x,
		vec.y-other.y,
	)
}

// Return a new vector equal to (vec + other)
func (vec LuaVector) Add(other LuaVector, inline ...bool) LuaVector {
	return CreateLuaVector(
		vec.x+other.x,
		vec.y+other.y,
	)
}

// Return a new vector equal to (vec · magnitude)
func (vec LuaVector) Scale(magnitude float64) LuaVector {
	return CreateLuaVector(
		vec.x*magnitude,
		vec.y*magnitude,
	)
}

// Get the length of a vector
func (vec LuaVector) Len() float64 {
	return math.Sqrt(vec.x*vec.x + vec.y*vec.y)
}

// Get the angle of a vector (in radians)
func (vec LuaVector) Angle() float64 {
	return math.Atan(vec.y / vec.x)
}

func (vec LuaVector) Rotate(radians float64) LuaVector {
	_sin, _cos := math.Sincos(radians)
	return CreateLuaVector(
		vec.x*_cos-vec.y*_sin,
		vec.x*_sin+vec.y*_cos,
	)
}

// ToString
func (vec LuaVector) String() string {

	return fmt.Sprintf("[%.2f, %.2f]", vec.x, vec.y)
}
