package motor

import (
	"github.com/go-gl/mathgl/mgl32"
)

type BasicEntity struct {
	Renderable     *Sprite
	Camera         *Camera
	Transformation // create struct. Contains x,y, rot, scalex, scaley, and can calc and cache the translation matrix
	UniSubTexPos   mgl32.Vec4
	UniColor       mgl32.Vec4
	UniColorMix    float32
}

func (P *BasicEntity) Draw() {
	camMatrix := P.Camera.GetMatrix()
	thingMatrix := P.GetMatrix()

	P.Renderable.UniSubTexPos = P.UniSubTexPos
	P.Renderable.UniColorMix = P.UniColorMix
	P.Renderable.UniColor = P.UniColor
	P.Renderable.Draw(camMatrix, thingMatrix)
}

func (P *BasicEntity) Update() {
}

func (P *BasicEntity) Clone() Entity {
	return &BasicEntity{
		Renderable:     P.Renderable,
		Camera:         P.Camera,
		Transformation: P.Transformation,
		UniSubTexPos:   P.UniSubTexPos,
		UniColor:       P.UniColor,
		UniColorMix:    P.UniColorMix,
	}
}
