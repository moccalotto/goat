package motor

type Velocity struct {
	VelX  float32 // Speed in X direction
	VelY  float32 // Speed in Y direction
	VelR  float32 // rotation around own axis (radians per second)
	VelTR float32 // how fast the object's trajecory is changing (radians per second)
}
