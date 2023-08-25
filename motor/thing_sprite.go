package motor

import (
	"goat/glhelp"

	"github.com/go-gl/mathgl/mgl32"
)

type SpriteThing struct {
	Renderable   *SpriteRenderable
	Camera       *Camera
	Trans        Transformation // create struct. Contains x,y, rot, scalex, scaley, and can calc and cache the translation matrix
	Velocity     PseudoForce
	UniSubTexPos mgl32.Vec4
	UniColor     mgl32.Vec4
	UniColorMix  float32
}

func (P *SpriteThing) Draw() {
	camMatrix := P.Camera.GetMatrix()
	thingMatrix := P.Trans.GetMatrix()

	P.Renderable.UniSubTexPos = P.UniSubTexPos
	P.Renderable.UniColorMix = P.UniColorMix
	P.Renderable.UniColor = P.UniColor
	P.Renderable.Draw(camMatrix, thingMatrix)
}

func (P *SpriteThing) Update() {
	if P.Velocity.TR*P.Velocity.TR > 0.00001 {
		si, co := glhelp.Sincos(P.Velocity.TR * Machine.Delta)

		x1, y1 := P.Velocity.Vec[0], P.Velocity.Vec[1]

		P.Velocity.Vec[0] = co*x1 - si*y1
		P.Velocity.Vec[1] = si*x1 + co*y1
	}

	P.Trans.Move(mgl32.Vec2{P.Velocity.Vec[0], P.Velocity.Vec[1]}.Mul(Machine.Delta))
	P.Trans.Rotate(P.Velocity.R * Machine.Delta)
}

func (P *SpriteThing) Clone() Thing {
	return &SpriteThing{
		Renderable: P.Renderable,
		Camera:     P.Camera,
		Trans:      P.Trans,
	}
}
