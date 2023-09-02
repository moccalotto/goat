package tractor

import (
	"goat/shed"
)

// Draw a sprite on screen
type Sprite struct {
	Renderer *TexQuadRenderer
	Camera   *Camera
	Position

	UniSubTexPos shed.V4
	UniColor     shed.V4
	UniColorMix  float32

	Deleted bool
}

func CreateSpriteAdv(renderer *TexQuadRenderer, camera *Camera) *Sprite {
	return &Sprite{
		Renderer:     renderer,
		Camera:       camera,
		UniSubTexPos: shed.V4{C1: 0, C2: 0, C3: 1, C4: 1},
		UniColorMix:  0,
	}
}

func (E *Sprite) Draw() {
	if E.Deleted {
		return
	}
	camMatrix := E.Camera.GetMatrix()
	thingMatrix := E.GetMatrix()

	E.Renderer.UniSubTexPos = E.UniSubTexPos
	E.Renderer.UniColorMix = E.UniColorMix
	E.Renderer.UniColor = E.UniColor
	E.Renderer.Draw(camMatrix, thingMatrix)
}

func (E *Sprite) Update() {
	if E.Deleted {
		return
	}
}

func (E *Sprite) Clone() *Sprite {
	return &Sprite{
		Renderer:     E.Renderer,
		Camera:       E.Camera,
		Position:     E.Position,
		UniSubTexPos: E.UniSubTexPos,
		UniColor:     E.UniColor,
		UniColorMix:  E.UniColorMix,
	}
}
