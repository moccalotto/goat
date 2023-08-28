package main

import (
	// h "goat/glhelp"
	"goat/glhelp"
	m "goat/motor"
)

var (
	gShot *m.BasicEntity
)

// Initialize sprites needed for the hero's projectiles
func initializeHeroShot() {

	sprite := gGodSprite.Clone()

	gShot = m.CreateBasicEntity(sprite, gCamera)
	gShot.UniSubTexPos = m.Machine.GetDimsForSubtexture(ATLAS_ID, LASER_BLUE_03_ID)

	w, h := gShot.UniSubTexPos[2]-gShot.UniSubTexPos[0], gShot.UniSubTexPos[3]-gShot.UniSubTexPos[1]
	aspect := h / w
	gShot.SetScale(0.3, 0.3*aspect)
	gShot.SetRotation(glhelp.Tau * 3 / 4)
	gShot.SetPos(0, 0)
}

func heroShoots() {
	if gShot == nil {
		initializeHeroShot()
	}

	shot := gShot.Clone().(*m.BasicEntity)

	x, y, angle := gHero.Get()

	shot.SetPos(x, y)
	shot.SetRotation(angle)

	force := m.Force{Vec: glhelp.V2{X: 10}.Rotate(angle - 270*glhelp.Degrees)}

	shot.Behavior = func(e m.Entity) {
		s := e.(*m.BasicEntity)
		s.ApplyForce(force, m.Machine.Delta)
	}

	addShot(shot)
}

func addShot(shot *m.BasicEntity) {

	// is there an empty spot in the shots array ?
	for i, s := range gShots {
		if s == nil {
			// found empty slot. Use it and leave.
			gShots[i] = shot
			return
		}
	}

	// no empty slots found, add new element to end of slice
	gShots = append(gShots, shot)
}
