package motor

import (
	"goat/glhelp"
)

type Force struct {
	Vec glhelp.V2
	Rot float32 // Rotation around own axis (radians per second)
}
