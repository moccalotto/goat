package main

import (
	m "goat/motor"
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Protagonist struct {
	sprite *m.SpriteRenderable
	camera *m.Camera

	Loc        m.Transformation // Location, Rotation, and Scale
	MinX, MaxX float32          // max and min values of X location
	MinY, MaxY float32          // max and min values of Y location
	Vel        m.PseudoForce    // Velocity vector, rotational velocity, and turning velocity (change in velocity direction)
	Acc        m.PseudoForce    // aceleration vector

	UniSubTexPos mgl32.Vec4
	UniColor     mgl32.Vec4
	UniColorMix  float32

	Weapon *Weapon

	LastShot       float64
	Up, Down, Fire bool
}

func CreateProtagonist(sprite *m.SpriteRenderable, cam *m.Camera) *Protagonist {

	prot := Protagonist{
		sprite:       sprite,
		camera:       cam,
		Vel:          m.PseudoForce{},
		Acc:          m.PseudoForce{},
		UniSubTexPos: mgl32.Vec4{0, 0, 1, 1},
		UniColor:     mgl32.Vec4{1, 1, 1, 1},
		UniColorMix:  0.5,
		Up:           false,
		Down:         false,
		Fire:         false,
	}

	prot.Vel.Max = 0.5
	prot.Acc.Max = 7

	return &prot
}

func (P *Protagonist) Draw() {
	P.sprite.UniColor = P.UniColor
	P.sprite.UniColorMix = P.UniColorMix
	P.sprite.UniSubTexPos = P.UniSubTexPos
	P.sprite.Draw(P.camera.GetMatrix(), P.Loc.GetMatrix())
}

func (P *Protagonist) Update() {
	ctx := glfw.GetCurrentContext()

	P.Up = ctx.GetKey(glfw.KeyUp) != glfw.Release
	P.Down = ctx.GetKey(glfw.KeyDown) != glfw.Release
	P.Fire = ctx.GetKey(glfw.KeySpace) != glfw.Release

	P.handleMovement()
}

func (P *Protagonist) handleMovement() {

	const standStillThreshold = 0.2

	goUp := (P.Up && !P.Down)
	goDown := (P.Down && !P.Up)
	hasVelocity := P.Vel.Vec.Len() > standStillThreshold

	switch true {
	case goUp:
		P.Acc.Vec[1] = P.Acc.Max
	case goDown:
		P.Acc.Vec[1] = -P.Acc.Max
	case hasVelocity:
		P.Acc.Vec = P.Vel.Vec.Normalize().Mul(-P.Acc.Max)
	default:
		P.Vel.Vec[1] = 0
		P.Acc.Vec[1] = 0

		// Do not apply forces that are near zero anywayy. Just return
		return
	}

	log.Printf("%+v\n", P.Loc.GetAll())

	P.applyForces()
}

func (P *Protagonist) applyForces() {
	dx, dy := P.Vel.XY()
	P.Vel.Apply(&P.Acc, m.Machine.Delta)
	P.Loc.RestrictedMove(dx, dy, P.MinX, P.MinY, P.MaxX, P.MaxY)
}

func (P *Protagonist) UseSubSprite(s string) {
	P.UniSubTexPos = m.Machine.SubTextures[s].GetDims()
}

func (P *Protagonist) Clone() m.Thing {
	return &Protagonist{
		sprite:       P.sprite,
		camera:       P.camera,
		Loc:          P.Loc,
		Vel:          P.Vel,
		UniSubTexPos: P.UniSubTexPos,
		UniColor:     P.UniColor,
		UniColorMix:  P.UniColorMix,
	}
}
