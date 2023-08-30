package main

import (
	"goat/glhelp"
	m "goat/motor"
	"regexp"
)

const (
	DEFAULT_COOLDOWN = 0.1
)

var (
	gAllWeapons    []*Weapon = []*Weapon{}
	gCurWeapon     *Weapon   = nil
	gProjScale     float32   = 0.5
	gProjSpeed     float32   = 10
	gCooldownUntil float64
	gProjPool      []*m.SpriteEnt = make([]*m.SpriteEnt, 100)
)

// ||
// ||
// ||
// ||Weapon contains the info necessary to spawn projectiles.
// ||It knows how fast their going, in which direction, and
// ||how they _behave_ while underway.
// ||=============================================================
type Weapon struct {
	name      string
	damage    float32
	speed     float32
	cooldown  float64
	projTempl *m.SpriteEnt
}

func (W *Weapon) Fire(shooter *m.SpriteEnt) *m.SpriteEnt {

	projectile := W.projTempl.Clone()

	entX, entY, entAngle := shooter.GetXYA()
	offsetV := glhelp.PolarV2(entAngle, 1)

	projectile.SetXY(entX+offsetV.X, entY+offsetV.Y)
	projectile.SetAngle(entAngle)

	projForce := m.Force{Vec: glhelp.V2{X: W.speed}.Rotate(entAngle)}

	// ||
	// ||
	// ||TODO:
	// || * Projectiles that rebound when they hit the edges and bounce around and make havoc
	// || * multi projectiles
	// || * curving projectiles
	// ||============================================================
	projectile.Behavior = m.CreateSimpleBehavior(func(e *m.SpriteEnt) {
		x, y, _ := e.GetXYA()

		if x > MAX_X || x < MIN_X || y > MAX_Y || y < MIN_Y {
			e.Deleted = true
			return
		}

		e.ApplyForce(projForce, m.Machine.Delta)

	})

	addProjectileToPool(projectile)

	return projectile
}
func addProjectileToPool(proj *m.SpriteEnt) {
	for i, _e := range gProjPool {
		if _e == nil || _e.Deleted {
			gProjPool[i] = proj
			return
		}
	}

	gProjPool = append(gProjPool, proj)
}

/*
TODO: Manually create each weapon to give it unique characteristics
*/
func initializeWeapons() {

	nameParser := regexp.MustCompile(`^laser[a-zA-Z]+\d+`)

	for _, st := range m.Machine.AtlasDescriptors[ATLAS_ID].SubTextures {

		matches := nameParser.MatchString(st.Name)
		if !matches {
			continue
		}

		aspect := m.Machine.GetAspectRatioForSubTexture(ATLAS_ID, st.Name)

		template := m.CreateBasicEnt(gSpriteSheet, gCamera)
		template.UniSubTexPos = m.Machine.GetDimsForSubtexture(ATLAS_ID, st.Name)
		template.SetScale(gProjScale*aspect, gProjScale)
		template.SetAngleOffset(-90 * Degrees)

		gAllWeapons = append(gAllWeapons, &Weapon{
			name:      st.Name,
			speed:     gProjSpeed,
			cooldown:  DEFAULT_COOLDOWN,
			projTempl: template,
		})
	}

	gCurWeapon = gAllWeapons[0]
}
