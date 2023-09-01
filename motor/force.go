package motor

import (
	"goat/util"
)

type Force struct {
	Vec util.V2
	Rot float32 // Rotation around own axis (radians per second)
}
