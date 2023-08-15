package main

import (
	"fmt"
	"math"
)

type Vector struct {
	x float64
	y float64
}

func CreateVector(x, y float64) Vector {
	return Vector{x: x, y: y}
}

// Create a vector from polar coordinates
func CreatePolarVector(radians float64, magnitude ...float64) Vector {
	y, x := math.Sincos(radians)

	var length float64
	if len(magnitude) > 0 {
		length = magnitude[0]
	} else {
		length = 1
	}

	return Vector{
		x: x * length,
		y: y * length,
	}
}

// Get the x-component of a vector
func (vec Vector) X() float64 {
	return vec.x
}

// Get the y-component of a vector
func (vec Vector) Y() float64 {
	return vec.y // return the original value of x
}

// Return a clone of this vector
func (vec Vector) Clone() Vector {
	return CreateVector(vec.x, vec.y)
}

// Return a new vector equal to (vec - other)
func (vec Vector) Sub(other Vector) Vector {
	return CreateVector(
		vec.x-other.x,
		vec.y-other.y,
	)
}

// Return a new vector equal to (vec + other)
func (vec Vector) Add(other Vector, inline ...bool) Vector {
	return CreateVector(
		vec.x+other.x,
		vec.y+other.y,
	)
}

// Return a new vector equal to (vec · magnitude)
func (vec Vector) Scale(magnitude float64) Vector {
	return CreateVector(
		vec.x*magnitude,
		vec.y*magnitude,
	)
}

// Get the length of a vector
func (vec Vector) Len() float64 {
	return math.Sqrt(vec.x*vec.x + vec.y*vec.y)
}

// Get the angle of a vector (in radians)
func (vec Vector) Angle() float64 {
	return math.Atan(vec.y / vec.x)
}

func (vec Vector) Rotate(radians float64) Vector {
	_sin, _cos := math.Sincos(radians)
	return CreateVector(
		vec.x*_cos-vec.y*_sin,
		vec.x*_sin+vec.y*_cos,
	)
}

// ToString
func (vec Vector) String() string {

	return fmt.Sprintf("[%.2f, %.2f]", vec.x, vec.y)
}
