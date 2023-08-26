package motor

import (
	"goat/glhelp"

	"github.com/go-gl/mathgl/mgl32"
)

func Rotate(v mgl32.Vec2, r float32) mgl32.Vec2 {
	x, y := v[0], v[1]
	sin, cos := glhelp.Sincos(r)
	return mgl32.Vec2{
		cos*x - sin*y,
		sin*x + cos*y,
	}
}

type Velocity struct {
	v        mgl32.Vec2
	Max      float32 // Maximum velocity (world units per second)
	Rotation float32 // Rotation around own axis (radians per second)
}

func (PF *Velocity) Apply(other *Velocity, delta float32) {
}
