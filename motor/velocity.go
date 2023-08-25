package motor

import (
	"github.com/go-gl/mathgl/mgl32"
)

type PseudoForce struct {
	Vec mgl32.Vec2
	Max float32
	R   float32 // rotation around own axis (radians per second)
	TR  float32 // how fast the object's trajecory is changing (radians per second)
}

func (PF *PseudoForce) Apply(other *PseudoForce, delta float32) {
	tmp := PF.Vec.Add(other.Vec.Mul(delta))

	if tmp.Len() > PF.Max {
		tmp = tmp.Normalize().Mul(PF.Max)
	}

	PF.Vec = tmp
}

func (PF *PseudoForce) XY() (float32, float32) {
	return PF.Vec[0], PF.Vec[1]
}
