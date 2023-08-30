package main

import (
	h "goat/glhelp"
	m "goat/motor" // gg = goat motor

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func initializeHero() {
	gHero = m.CreateBasicEnt(gSpriteSheet, gCamera)

	gHero.UniColor = mgl32.Vec4{1, 1, 1, 1}
	gHero.UniColorMix = 0.5
	gHero.UniSubTexPos = m.Machine.GetDimsForSubtexture(ATLAS_ID, HERO_TEX_ID)

	gHero.SetScale(1, 1)
	gHero.SetAngle(0 * h.Degrees)
	gHero.SetAngleOffset(-90 * Degrees)
	gHero.LimitAngle(-5*Degrees, 5*Degrees)
	gHero.LimitLocation(MIN_X+MARGIN, MIN_Y+MARGIN, MIN_X+2, MAX_Y-MARGIN)
	gHero.SetXY(MIN_X+2, 0)

	DownForce := m.Force{
		Vec: h.V2{Y: -HERO_BASE_SPEED},
		Rot: -0.3,
	}
	UpForce := m.Force{
		Vec: h.V2{Y: HERO_BASE_SPEED},
		Rot: 0.3,
	}

	gHero.Behavior = m.CreateSimpleBehavior(func(_ *m.SpriteEnt) {

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
			gHero.RotateTowards(0*h.Degrees, m.Machine.Delta*10)
			gHero.SnapAngleTo(0, 0.5*h.Degrees)
		}

		if m.KeyPressed(glfw.KeySpace) && m.Machine.Now64 > gCooldownUntil {
			gCurWeapon.Fire(gHero)
			gCooldownUntil = m.Machine.Now64 + gCurWeapon.cooldown
		}
	})
}
