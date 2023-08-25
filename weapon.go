package main

import m "goat/motor"

type Weapon struct {
	Rate     float32          // shots per second
	LastShot float32          // when was last shot fired
	NextShot float32          // when can you fire the next shot
	Fire     func() []m.Thing // Spawn all the necessary shots
}
