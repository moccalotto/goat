package motor

import (
	"github.com/go-gl/mathgl/mgl32"
)

type SimpleBehavior struct {
	updateFunc func(*SpriteEnt)
}

func (SB *SimpleBehavior) Update(e *SpriteEnt) {
	SB.updateFunc(e)
}

func CreateSimpleBehavior(f func(e *SpriteEnt)) *SimpleBehavior {
	return &SimpleBehavior{
		updateFunc: f,
	}
}

type Behavior interface {
	Update(e *SpriteEnt)
}

type SpriteEnt struct {
	SpriteGl *SpriteGl
	Camera   *Camera
	Position

	UniSubTexPos mgl32.Vec4
	UniColor     mgl32.Vec4
	UniColorMix  float32
	Components   uint64
	Behavior     Behavior

	Deleted bool
}

func CreateBasicEnt(spriteGl *SpriteGl, camera *Camera) *SpriteEnt {
	return &SpriteEnt{
		SpriteGl: spriteGl,
		Camera:   camera,
		Behavior: nil,
	}
}

func (E *SpriteEnt) Draw() {
	if E.Deleted {
		return
	}
	camMatrix := E.Camera.GetMatrix()
	thingMatrix := E.GetMatrix()

	E.SpriteGl.UniSubTexPos = E.UniSubTexPos
	E.SpriteGl.UniColorMix = E.UniColorMix
	E.SpriteGl.UniColor = E.UniColor
	E.SpriteGl.Draw(camMatrix, thingMatrix)
}

func (E *SpriteEnt) Update() {
	if E.Deleted {
		return
	}
	if E.Behavior != nil {
		E.Behavior.Update(E)
	}
}

func (E *SpriteEnt) Clone() *SpriteEnt {
	return &SpriteEnt{
		SpriteGl:     E.SpriteGl,
		Camera:       E.Camera,
		Position:     E.Position,
		UniSubTexPos: E.UniSubTexPos,
		UniColor:     E.UniColor,
		UniColorMix:  E.UniColorMix,
		Behavior:     E.Behavior,
	}
}
