package main

import (
	"fmt"
	"math"
)

type VectorType struct {
	elements []float64
}

func PrepareVector(count int) VectorType {
	return VectorType{
		elements: make([]float64, count),
	}
}
func CreateVector(elements ...float64) VectorType {
	return VectorType{
		elements: elements,
	}
}

func (vec VectorType) Clone() VectorType {
	ret := PrepareVector(len(vec.elements))
	copy(ret.elements, vec.elements)

	return ret
}

func (vec VectorType) Partner() VectorType {
	return PrepareVector(len(vec.elements))
}

func (vec VectorType) Invert() VectorType {
	ret := vec.Partner()

	for i, v := range vec.elements {
		ret.elements[i] = -v
	}

	return ret
}

func (vec VectorType) Sub(other VectorType) VectorType {
	if len(other.elements) != len(vec.elements) {
		panic(
			fmt.Errorf(
				"cannot sub vectors of different dimensions [%d dims vs %d dims]",
				len(other.elements),
				len(vec.elements),
			),
		)
	}

	ret := vec.Clone()

	for i, v := range other.elements {
		ret.elements[i] -= v
	}

	return ret
}

func (vec VectorType) Add(other VectorType) VectorType {
	if len(other.elements) != len(vec.elements) {

		panic(
			fmt.Errorf(
				"cannot add vectors of different dimensions [%d dims vs %d dims]",
				len(other.elements),
				len(vec.elements),
			),
		)
	}

	ret := vec.Clone()
	for i, v := range other.elements {
		ret.elements[i] += v
	}

	return ret
}

func (vec VectorType) Scale(magnitude float64) VectorType {
	ret := vec.Partner()

	for i, _ := range vec.elements {
		ret.elements[i] *= magnitude
	}

	return ret
}

func (vec VectorType) Len() float64 {
	sum := 0.0

	for _, v := range vec.elements {
		sum += v * v
	}

	return math.Sqrt(sum)
}

func (vec VectorType) X() float64 {
	return vec.elements[0]
}

func (vec VectorType) Y() float64 {
	return vec.elements[1]
}
func (vec VectorType) Z() float64 {
	return vec.elements[2]
}
func (vec VectorType) W() float64 {
	return vec.elements[3]
}

func (vec VectorType) Angle() float64 {
	x, y := vec.elements[0], vec.elements[1]

	return math.Atan(y / x)
}

func PolarVector(direction float64, scale ...float64) VectorType {
	x, y := math.Sincos(direction)

	_scale := 1.0

	if len(scale) > 0 {
		_scale = scale[0]
	}

	return VectorType{
		elements: []float64{
			x * _scale,
			y * _scale,
		},
	}
}

func (vec VectorType) Set(vals ...float64) VectorType {

	if len(vals) != len(vec.elements) {
		panic("Cannot change a vector's dimensions!")
	}

	copy(vec.elements, vals)

	return vec
}

func (vec VectorType) ToString() string {

	// return fmt.Sprintf("%+v", vec.elements)

	ret := fmt.Sprintf("[%.2f", vec.elements[0])

	for i := 1; i < len(vec.elements); i++ {
		ret = fmt.Sprintf("%s, %.2f", ret, vec.elements[i])
	}

	return ret + "]"
}
