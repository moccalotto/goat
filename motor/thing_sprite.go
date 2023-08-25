package motor

import (
	"goat/glhelp"

	"github.com/go-gl/mathgl/mgl32"
)

type SpriteThing struct {
	Renderable   *SpriteRenderable
	Camera       *glhelp.Camera
	Trans        CachedTransformation // create struct. Contains x,y, rot, scalex, scaley, and can calc and cache the translation matrix
	Velocity     Velocity
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
	if P.Velocity.VelTR*P.Velocity.VelTR > 0.00001 {
		si, co := glhelp.Sincos(P.Velocity.VelTR * Machine.Delta)

		x1, y1 := P.Velocity.VelX, P.Velocity.VelY

		P.Velocity.VelX = co*x1 - si*y1
		P.Velocity.VelY = si*x1 + co*y1
	}

	P.Trans.Move(mgl32.Vec2{P.Velocity.VelX, P.Velocity.VelY}.Mul(Machine.Delta))
	P.Trans.Rotate(P.Velocity.VelR * Machine.Delta)
}

func (P *SpriteThing) Clone() Thing {
	return &SpriteThing{
		Renderable: P.Renderable,
		Camera:     P.Camera,
		Trans:      P.Trans,
	}
}
