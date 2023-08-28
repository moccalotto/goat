package motor

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Behavior func(e Entity)

type BasicEntity struct {
	Sprite       *Sprite
	Camera       *Camera
	Position     // create struct. Contains x,y, rot, scalex, scaley, and can calc and cache the translation matrix
	UniSubTexPos mgl32.Vec4
	UniColor     mgl32.Vec4
	UniColorMix  float32
	Components   uint64
	Behavior Behavior

	// Deleted bool
}

func CreateBasicEntity(renderable *Sprite, camera *Camera) *BasicEntity {
	return &BasicEntity{
		Sprite:    renderable,
		Camera:    camera,
        Behavior: nil,
	}
}

func (E *BasicEntity) Draw() {
	camMatrix := E.Camera.GetMatrix()
	thingMatrix := E.GetMatrix()

	E.Sprite.UniSubTexPos = E.UniSubTexPos
	E.Sprite.UniColorMix = E.UniColorMix
	E.Sprite.UniColor = E.UniColor
	E.Sprite.Draw(camMatrix, thingMatrix)
}

func (E *BasicEntity) Update() {
    if E.Behavior != nil {
        E.Behavior(E)
    }
}

func (E *BasicEntity) Clone() Entity {
	return &BasicEntity{
		Sprite:       E.Sprite,
		Camera:       E.Camera,
		Position:     E.Position,
		UniSubTexPos: E.UniSubTexPos,
		UniColor:     E.UniColor,
		UniColorMix:  E.UniColorMix,
	}
}
