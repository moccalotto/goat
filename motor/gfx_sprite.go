package motor

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Draw a sprite on screen
type Sprite struct {
	Renderer *TexQuadRenderer
	Camera   *Camera
	Position

	UniSubTexPos mgl32.Vec4
	UniColor     mgl32.Vec4
	UniColorMix  float32
	Behavior     Behavior

	Deleted bool
}

func CreateSpriteAdv(renderer *TexQuadRenderer, camera *Camera) *Sprite {
	return &Sprite{
		Renderer:     renderer,
		Camera:       camera,
		Behavior:     nil,
		UniSubTexPos: mgl32.Vec4{0, 0, 1, 1},
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
	if E.Behavior != nil {
		E.Behavior.Update(E)
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
		Behavior:     E.Behavior,
	}
}
