package motor

import (
	"goat/glhelp"
)

type Force struct {
	Vec         glhelp.V2
	Max         float32 // Maximum length of vec
	Min         float32 // Minimum length of vec
	Rot         float32 // Rotation around own axis (radians per second)
	ApplyLimits bool
}

func (F *Force) AddTo(other Force, delta float32) Force {
	if F.ApplyLimits {
		F.Vec = F.Vec.ClampLen(F.Min, F.Max)
	}

	other.Vec = other.Vec.Add(F.Vec.Scale(delta))

	if other.ApplyLimits {
		other.Vec = other.Vec.ClampLen(other.Min, other.Max)
	}

	return other
}

func (F *Force) ApplyToPosition(p *Position, delta float32) {
	if F.ApplyLimits {
		F.Vec = F.Vec.ClampLen(F.Min, F.Max)
	}

	p.Move(F.Vec.Scale(delta))
	p.Rotate(F.Rot * delta)
}
