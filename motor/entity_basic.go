package motor

import (
	"log"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type BasicEntity struct {
	Renderable   *Sprite
	Camera       *Camera
	Position     // create struct. Contains x,y, rot, scalex, scaley, and can calc and cache the translation matrix
	UniSubTexPos mgl32.Vec4
	UniColor     mgl32.Vec4
	UniColorMix  float32
	forces       map[string]Force

	// Deleted bool
}

func (E *BasicEntity) Draw() {
	camMatrix := E.Camera.GetMatrix()
	thingMatrix := E.GetMatrix()

	log.Printf("%+v\n", thingMatrix)
	time.Sleep(100 * time.Millisecond)

	E.Renderable.UniSubTexPos = E.UniSubTexPos
	E.Renderable.UniColorMix = E.UniColorMix
	E.Renderable.UniColor = E.UniColor
	E.Renderable.Draw(camMatrix, thingMatrix)
}

func (E *BasicEntity) Update() {
	for _, force := range E.forces {
		E.ApplyForce(force, Machine.Delta)
	}
}

func (E *BasicEntity) Clone() Entity {
	return &BasicEntity{
		Renderable:   E.Renderable,
		Camera:       E.Camera,
		Position:     E.Position,
		UniSubTexPos: E.UniSubTexPos,
		UniColor:     E.UniColor,
		UniColorMix:  E.UniColorMix,
	}
}
