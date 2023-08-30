package main

import (
	m "goat/motor"
	"math/rand"
	"regexp"
)

var (
	gEnemyPool  []*m.SpriteEnt = make([]*m.SpriteEnt, 20)
	gEnemyTypes []*EnemyType   = make([]*EnemyType, 0, 50)
	_enemySlots []*m.SpriteEnt = make([]*m.SpriteEnt, 20)
)

type EnemyType struct {
	name     string
	weapon   *Weapon
	template *m.SpriteEnt
}

type EnemyBehavior struct {
	typ        *EnemyType
	lastShotAt float64
}

func (EB *EnemyBehavior) Clone() *EnemyBehavior {
	return &EnemyBehavior{
		typ:        EB.typ,
		lastShotAt: EB.lastShotAt,
	}
}

func (EB *EnemyBehavior) Update(e *m.SpriteEnt) {

	if EB.lastShotAt > m.Machine.Now64 {
		return
	}

	if m.Machine.Now64 > EB.lastShotAt && rand.Float32() > m.Machine.Delta {
		EB.typ.weapon.Fire(e)
		EB.lastShotAt = m.Machine.Now64 + EB.typ.weapon.cooldown
	}
}

func (E *EnemyType) Spawn() bool {

	var ty float32
	var step float32 = 1.2
	tx := float32(MAX_X - 1)
	emptySpots := []float32{}
	for ty = MIN_Y + 1; ty < MAX_Y; ty += step {
		if isSpotTakenByEnemy(tx, ty) {
			continue
		}
		emptySpots = append(emptySpots, ty)

	}

	if len(emptySpots) == 0 {
		return false
	}

	ty = emptySpots[rand.Int()%len(emptySpots)]
	nmy := E.template.Clone()
	nmy.SetXY(MAX_X-1, ty)
	nmy.Behavior = nmy.
		Behavior.(*EnemyBehavior).
		Clone()
	addEnemyToPool(nmy)

	return true
}

func addEnemyToPool(nmy *m.SpriteEnt) {
	for i, _e := range gEnemyPool {
		if _e == nil || _e.Deleted {
			gEnemyPool[i] = nmy
			return
		}
	}

	gEnemyPool = append(gEnemyPool, nmy)
}

func isSpotTakenByEnemy(x, y float32) bool {

	for _, e := range gEnemyPool {
		if e == nil || e.Deleted {
			continue
		}

		ex, ey, _ := e.GetXYA()

		if ey == y && x == ex {
			return true
		}
	}
	return false
}

func initializeEnemies() {

	nameParser := regexp.MustCompile(`^enemy([a-zA-Z]+)(\d+)`)

	for _, st := range m.Machine.AtlasDescriptors[ATLAS_ID].SubTextures {

		isEnemy := nameParser.MatchString(st.Name)
		// This SubTexture is not a weapon, skip it.
		if !isEnemy {
			continue
		}

		template := m.CreateBasicEnt(gSpriteSheet, gCamera)
		template.UniSubTexPos = m.Machine.GetDimsForSubtexture(ATLAS_ID, st.Name)
		template.SetScale(1, 1)
		template.SetAngleOffset(90 * Degrees)
		template.SetAngle(180 * Degrees)
		template.SetXY(MAX_X-2, 0)
		// template.Behavior = m.CreateSimpleBehavior(func(e m.Ent) {
		// 	nmy := e.(*m.BasicEnt)
		// 	rnd := rand.Float32()
		// 	if rnd < m.Machine.Delta {
		// 		gCurWeapon.Fire(nmy)
		// 	}
		// })

		enemy := EnemyType{
			name:     st.Name,
			weapon:   gAllWeapons[rand.Int()%len(gAllWeapons)],
			template: template,
		}

		template.Behavior = &EnemyBehavior{
			lastShotAt: 0,
			typ:        &enemy,
		}

		gEnemyTypes = append(gEnemyTypes, &enemy)
	}
}
