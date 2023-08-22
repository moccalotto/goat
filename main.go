package main

import (
	"goat/glhelp"
	"math/rand"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	// start the mainthread system, allowing us to make calls on the main thread later
	StartMainThreadSystem(func() {

		// call actualMain() and ensure that it runs on the main thread
		// because it contains tons of OpenGL calls, etc.
		RunOnMain(actualMain)
	})
}

const (
	protBaseAngle = 270 * glhelp.Degrees
	maxBankAngle  = 10 * glhelp.Degrees
	maxY          = 10  // (screen  height in px) / ( 100 * 2)
	maxX          = 14  // (screen  width in px) / ( 100 * 2)
	shotDelay     = 0.1 // this should be dynamic, depending on what weapon the player has
	shotSpeed     = 15  // this should depend on the player's weapon
	maxProtShots  = 20  // this should depend on the player's wepaon
)

type PolyMap map[*glhelp.Polygon]*glhelp.Polygon

var (
	camera     *glhelp.Camera
	window     *glfw.Window
	background *glhelp.Polygon

	enemySpawnRate    float32 = 3
	alive             bool    = true
	enemies           PolyMap = make(PolyMap)
	enemyShots        PolyMap = make(PolyMap)
	enemyTemplate     *glhelp.Polygon
	enemyShotTemplate *glhelp.Polygon

	protagonist      *glhelp.Polygon
	protShots        PolyMap = make(PolyMap)
	protShotTemplate *glhelp.Polygon
	now              float32 = 0
	prev             float32 = 0
	delta            float32 = 0
	tickCount        uint64  = 0

	protSpeed    float32 = 10
	protDown     bool
	protUp       bool
	protShoot    bool
	lastProtShot float32 = -1.0e10

	wireframe bool

	windowOptions *WindowOptions = &WindowOptions{
		Title:     "GOAT",
		Width:     3000,
		Height:    2000,
		Resizable: false,
	}
)

// This is the actual main function
// must run on mainthread
func actualMain() {

	Setup()

	glhelp.EnableBlending()

	//
	// timing and bookkeeping variables

	for !window.ShouldClose() {
		tickCount = tickCount + 1
		prev = now
		now = float32(glfw.GetTime())
		delta = now - prev

		glhelp.ClearScreenF(0, 0, 0, 0)

		if alive {
			Update()
		}
		Draw()

		window.SwapBuffers()
		glhelp.AssertGLOK("EndOfDraw")
		glfw.PollEvents()

	}
}

// //////////////////////////////////////////////////////////////////
// UPDATE
// //////////////////////////////////////////////////////////////////
func Update() {
	angle := float32(0)
	distance := float32(0)
	_, y := protagonist.GetPos()

	if protUp && !protDown {
		angle = maxBankAngle
		distance = mgl32.Clamp(
			protSpeed*delta,
			0,
			maxY-y-1,
		)
	}
	if protDown && !protUp {
		angle = -maxBankAngle
		distance = -mgl32.Clamp(
			protSpeed*delta,
			0,
			maxY+y-1,
		)
	}

	if distance != 0 {
		protagonist.Move(0, distance)
	}

	// change rotation graduately to point in the direction of "travel"

	rot := protagonist.GetRotation()
	newAngle := glhelp.Lerp(rot, protBaseAngle+angle, delta*5)
	angleDiff := mgl32.Abs(rot - newAngle)

	if angleDiff > 0.05*glhelp.Degrees {
		// set the angle normally if we still have some way to go to reach destination angle
		protagonist.SetRotation(newAngle)
	} else if angleDiff > 0.01 {
		protagonist.SetRotation(protBaseAngle + angle)
	}
	if protShoot {
		protagonistShoots()
	}

	{ // do a bunch of gofunc
		donsky := make(chan bool)

		// Update all shots (in a goroutine)
		go func() {
			for _, shot := range protShots {
				angle := shot.GetRotation()
				si, co := glhelp.Sincos(angle)

				// the equation is a bit wonky because the model is rotated
				shot.Move(-si*shotSpeed*delta, co*shotSpeed*delta)
			}

			donsky <- true
		}()

		// Update all enemies (in a goroutine)
		go func() {
			for _, enemy := range enemies {
				mix := enemy.GetColorMix()
				if mix > 0.1 {
					newMix := mgl32.Clamp(
						mix-delta*2, // go from fully invisible to full visible in .5 second
						0,
						1,
					)
					enemy.SetColorMix(newMix)
				}
			}
			donsky <- true
		}()

		// Update all enemy shots
		go func() {
			for _, eshot := range enemyShots {
				angle := eshot.GetRotation()
				si, co := glhelp.Sincos(angle)

				eshot.Move(-si*shotSpeed*delta, co*shotSpeed*delta)
			}
			donsky <- true
		}()

		// Wait for the shot- and enemy routines to be done
		<-donsky
		<-donsky
		<-donsky
	}

	// Roll a die and try to spawn a new enemy
	if rand.Float32() < delta*enemySpawnRate {
		nmy := trySpawnEnemy()
		if nmy != nil {
			enemies[nmy] = nmy
		}
	}

	for _, enemy := range enemies {
		if rand.Float32() < delta {
			enemyShoots(enemy)
		}
	}

	for _, shot := range protShots {
		if shouldCull(shot) {
			delete(protShots, shot)
			continue
		}

		for _, enemy := range enemies {
			if collides(shot, enemy) {
				delete(protShots, shot)
				delete(enemies, enemy)
				continue
			}
		}
	}

	for _, eShot := range enemyShots {
		if collides(eShot, protagonist) {
			background.SetColor(mgl32.Vec4{1, 0, 0, 1})
			background.SetColorMix(1)
			alive = false
		}
	}

}

// //////////////////////////////////////////////////////////////////
// DRAW
// Order of operations:
// * Background
// * Protagonist
// * Enemies
// * Enemy shots
// * Protagonist shots
// * Asteroids and environment
// * Explosions and effects
// //////////////////////////////////////////////////////////////////
func Draw() {
	camMatrix := camera.GetTransformationMatrix() // can omit. in a shootemup the camera never moves

	background.SetWireframe(wireframe)
	background.Draw(camMatrix)

	protagonist.SetWireframe(wireframe)
	protagonist.Draw(camMatrix)

	for _, enemy := range enemies {
		enemy.SetWireframe(wireframe)
		enemy.Draw(camMatrix)
	}

	for _, eshot := range enemyShots {
		eshot.SetWireframe(wireframe)
		eshot.Draw(camMatrix)
	}

	for _, shot := range protShots {
		shot.SetWireframe(wireframe)
		shot.Draw(camMatrix)
	}
}

// //////////////////////////////////////////////////////////////////
// SETUP
// //////////////////////////////////////////////////////////////////
func Setup() {
	_, _window, err := initGlfw(windowOptions)
	window = _window
	glhelp.GlPanicIfErrNotNil(err)

	glfw.GetCurrentContext().SetKeyCallback(KeyHandler)
	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	initializeCamera(windowOptions)
	initializeBackground()
	initializeProtagonist()
	initializeProtShotTemplate()
	initializeEnemyTemplate()
	initializeEnemyShotTemplate()
}

// //////////////////////////////////////////////////////////////////
// Keyboard handler
// //////////////////////////////////////////////////////////////////
func KeyHandler(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Repeat {
		return
	}

	keydown := action == glfw.Press

	switch key {
	case glfw.KeyEscape:
		glfw.GetCurrentContext().SetShouldClose(true)
	case glfw.KeyUp:
		protUp = keydown
	case glfw.KeyDown:
		protDown = keydown
	case glfw.KeySpace:
		protShoot = keydown
	case glfw.KeyF1:
		if !keydown {
			break
		}
		wireframe = !wireframe
		glhelp.Wireframe(wireframe)

	case glfw.KeyF2:
		if !keydown {
			break
		}
		for _, e := range enemies {
			enemyShoots(e)
		}
	case glfw.KeyEnter:
		if !keydown {
			break
		}
		tmp := trySpawnEnemy()
		if tmp != nil {
			enemies[tmp] = tmp
		}
	}

}

func enemyShoots(enemy *glhelp.Polygon) {
	shot := enemyShotTemplate.Copy()
	x, y := enemy.GetPos()

	offsetX, offsetY := glhelp.Sincos(enemy.GetRotation())

	shot.SetPosition(x+offsetX*2, y+offsetY*2) // todo - move away from the enemy, using its angle

	enemyShots[shot] = shot
}

func protagonistShoots() {
	if now < lastProtShot+shotDelay {
		return
	}

	if len(protShots) >= maxProtShots {
		return
	}

	lastProtShot = now

	shot := spawnProtagonistShot()

	protShots[shot] = shot
}

func spawnProtagonistShot() *glhelp.Polygon {
	protX, protY := protagonist.GetPos()

	shot := protShotTemplate.Copy()

	shot.SetRotation(protagonist.GetRotation())
	shot.SetPosition(protX+1, protY)

	return shot
}

func trySpawnEnemy() *glhelp.Polygon {

	enemy := enemyTemplate.Copy()

	enemy.SetPosition(
		maxX-2,
		-maxY+rand.Float32()*maxY*2,
	)

	for _, e := range enemies {
		if collides(enemy, e, 1.5) {
			return nil
		}
	}

	return enemy
}

func shouldCull(p *glhelp.Polygon) bool {
	x, y := p.GetPos()

	return y > maxY ||
		y < -maxY ||
		x > maxX ||
		x < -maxX
}

func collides(a, b *glhelp.Polygon, scaleFactor ...float32) bool {
	sa := a.GetScaleV()
	sb := b.GetScaleV()
	scale := float32(1)
	if len(scaleFactor) > 0 {
		scale = scaleFactor[0]
	}

	maxDist := glhelp.Max(sa.Len(), sb.Len())

	length := a.GetPosV().Sub(b.GetPosV()).Len()

	return length < maxDist*scale
}

func initializeEnemyTemplate() {
	enemyTemplate = glhelp.CreateSprite("assets/PNG/Enemies/enemyRed4.png")
	enemyTemplate.Initialize()

	enemyTemplate.SetScale(2, 2)
	enemyTemplate.SetRotation(270 * glhelp.Degrees)
	enemyTemplate.SetColorMix(1)
	enemyTemplate.SetColor(mgl32.Vec4{1, 1, 1, 0})
}
func initializeEnemyShotTemplate() {
	enemyShotTemplate = glhelp.CreateSprite("assets/PNG/Lasers/laserRed16.png")
	enemyShotTemplate.Initialize()

	enemyShotTemplate.SetScale(0.4, 1.4)
	enemyShotTemplate.SetRotation(90 * glhelp.Degrees)
	enemyShotTemplate.SetColorMix(0.1)
	enemyShotTemplate.SetColor(mgl32.Vec4{1, 1, 1, 0})
}

func initializeProtShotTemplate() {
	protShotTemplate = glhelp.CreateSprite("assets/PNG/Lasers/laserBlue01.png")
	protShotTemplate.Initialize()

	protShotTemplate.SetScale(0.4, 1.4)
	protShotTemplate.SetRotation(-90 * glhelp.Degrees)
	protShotTemplate.SetColorMix(0.1)
	protShotTemplate.SetColor(mgl32.Vec4{1, 1, 1, 0})
}

func initializeProtagonist() {
	protagonist = glhelp.CreateSprite("assets/PNG/playerShip1_green.png")
	protagonist.Initialize()

	protagonist.SetScale(2, 2)
	protagonist.SetPosition(-12, 0)
	protagonist.SetRotation(protBaseAngle)

	protagonist.SetColorMix(0.1)
	protagonist.SetColor(mgl32.Vec4{1, 1, 1, 0})
}

func initializeBackground() {
	background = glhelp.CreateSprite("assets/Backgrounds/purple.png")
	background.Initialize()

	w, h := camera.GetFrameSize()
	background.SetScale(w-1, h-1)
	background.SetColorMix(0.0)
	background.SetColor(mgl32.Vec4{0, 0, 0, 0})
}

func initializeCamera(options *WindowOptions) {

	camera = glhelp.CreateCamera()
	camera.SetFrameSize(float32(options.Width)/100, float32(options.Height)/100)
	camera.SetPosition(0, 0)
}
