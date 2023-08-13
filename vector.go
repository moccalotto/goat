package main

import (
	"fmt"
	"math"
)

type Vector struct {
	components []float64
}

// Create an empty vector the the given dimensionality
func CreateVectorD(dimensions int) Vector {
	return Vector{
		components: make([]float64, dimensions),
	}
}

// Create new vector - variadic
func CreateVectorV(components ...float64) Vector {
	return Vector{
		components: components,
	}
}

// Create new vector from a list of components
func CreateVectorL(components []float64) Vector {
	length := len(components)
	retVec := CreateVectorD(length)

	// copy components into the retVec vector
	count := copy(retVec.components, components)
	if count != length {
		panic(fmt.Errorf(
			"Something went wrong when copying the component list. Expected %d elements, but got %d",
			length,
			count))
	}

	return retVec
}

// Create a vector from polar coordinates
func CreatePolarVector(direction float64, magnitude ...float64) Vector {
	y, x := math.Sincos(direction)

	_scale := 1.0

	if len(magnitude) > 0 {
		_scale = magnitude[0]
	}

	return Vector{
		components: []float64{
			x * _scale,
			y * _scale,
		},
	}
}

// Get the x-component of a vector
func (vec Vector) X(update ...float64) float64 {
	const idx int = 0
	retval := vec.components[idx]

	if len(update) > 0 {
		vec.components[idx] = update[idx]
	}

	return retval
}

// Get the y-component of a vector
func (vec Vector) Y(update ...float64) float64 {
	const idx int = 1
	retval := vec.components[idx]

	if len(update) > 0 {
		vec.components[idx] = update[idx]
	}

	return retval
}

// Return a clone of this vector
func (vec Vector) Clone() Vector {
	return CreateVectorL(vec.components)
}

// Return a new vector equal to (vec - other)
func (vec Vector) Sub(other Vector) Vector {
	if len(other.components) != len(vec.components) {
		panic(
			fmt.Errorf(
				"cannot sub vectors of different dimensions [%d dims vs %d dims]",
				len(other.components),
				len(vec.components),
			),
		)
	}

	ret := vec.Clone()

	for i, v := range other.components {
		ret.components[i] -= v
	}

	return ret
}

// Return a new vector equal to (vec + other)
func (vec Vector) Add(other Vector, inline ...bool) Vector {
	if len(other.components) != len(vec.components) {

		panic(
			fmt.Errorf(
				"cannot add vectors of different dimensions [%d dims vs %d dims]",
				len(other.components),
				len(vec.components),
			),
		)
	}

	ret := vec.Clone()

	for i, v := range other.components {
		ret.components[i] += v
	}

	return ret
}

// Create a new version of vec, scaled by magnitude
func (vec Vector) Scale(magnitude float64) Vector {
	ret := CreateVectorD(len(vec.components))

	for i, v := range vec.components {
		ret.components[i] = v * magnitude
	}

	return ret
}

// Get the length of a vector
func (vec Vector) Len() float64 {
	sum := 0.0

	for _, v := range vec.components {
		sum += v * v
	}

	return math.Sqrt(sum)
}

// Get the angle of a vector (in radians)
func (vec Vector) Angle() float64 {
	x, y := vec.components[0], vec.components[1]

	return math.Atan(y / x)
}

// Overwrite the components of a vector
func (vec Vector) Set(vals ...float64) Vector {

	if len(vals) != len(vec.components) {
		panic("Cannot change a vector's dimensions!")
	}

	copy(vec.components, vals)

	return vec
}

// ....
func (vec Vector) String() string {

	return fmt.Sprintf("%+v", vec.components)
}
