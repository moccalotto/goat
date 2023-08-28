package main

import (
	h "goat/glhelp"
	m "goat/motor" // gg = goat motor

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func initializeHero() {
	sprite := gGodSprite.Clone()

	gHero = m.CreateBasicEntity(sprite, gCamera)

	gHero.UniColor = mgl32.Vec4{1, 1, 1, 1}
	gHero.UniColorMix = 0.5
	gHero.UniSubTexPos = m.Machine.GetDimsForSubtexture(ATLAS_ID, HERO_TEX_ID)

	gHero.SetScale(1, 1)
	gHero.SetRotation(-90 * h.Degrees)
	gHero.LimitRotation(265*h.Degrees, 275*h.Degrees)
	gHero.LimitLocation(MIN_X+MARGIN, MIN_Y+MARGIN, MIN_X+2, MAX_Y-MARGIN)
	gHero.SetPos(MIN_X+2, 0)

	DownForce := m.Force{
		Vec: h.V2{Y: -HERO_BASE_SPEED},
		Rot: -0.3,
	}
	UpForce := m.Force{
		Vec: h.V2{Y: HERO_BASE_SPEED},
		Rot: 0.3,
	}

	gHero.Behavior = func(_ m.Entity) {

		revert := true
		if m.KeyPressed(glfw.KeyDown) {
			gHero.ApplyForce(DownForce, m.Machine.Delta)
			revert = false
		}

		if m.KeyPressed(glfw.KeyUp) {
			gHero.ApplyForce(UpForce, m.Machine.Delta)
			revert = false
		}

		if revert {
			gHero.RotateTowards(270*h.Degrees, m.Machine.Delta*10)
			gHero.SnapRotationTo(270, 0.5*h.Degrees)
		}

		if m.KeyPressed(glfw.KeySpace) {
			heroShoots()
		}
	}
}
