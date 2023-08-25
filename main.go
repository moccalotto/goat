package main

import (
	h "goat/glhelp"
	m "goat/motor" // gg = goat motor
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func main() {

	// start the mainthread system, allowing us to make calls on the main thread later
	h.StartMainThreadSystem(func() {

		// call actualMain() and ensure that it runs on the main thread
		// because it contains tons of OpenGL calls, etc.
		h.RunOnMain(actualMain)
	})
}

type PolyMap map[*h.Polygon]*h.Polygon

type Weapon struct {
	Rate     float32           // shots per second
	LastShot float32           // when was last shot fired
	NextShot float32           // when can you fire the next shot
	Fire     func() []m.Thing // Spawn all the necessary shots
}

const (
	protBaseAngle = 270 * h.Degrees
	maxBankAngle  = 5 * h.Degrees
	sceneW        = 30
	sceneH        = 20
	WinResFactor  = 100        // pixels per "square"
	maxY          = sceneH / 2 // (screen  height in px) / ( 100 * 2)
	maxX          = sceneW / 2 // (screen  width in px) / ( 100 * 2)
	shotDelay     = 0.1        // this should be dynamic, depending on what weapon the player has
	shotSpeed     = 15         // this should depend on the player's weapon
	maxProtShots  = 20         // this should depend on the player's wepaon
)

var (
	bgScrollSpeed float32 = 0.08
	maxProtSpeed  float32 = 10
	camera        *h.Camera
	window        *glfw.Window
	background    *m.SpriteThing
	spriteSheet   *m.SpriteRenderable

	protagonist      *m.SpriteThing
	protShotTemplate *m.SpriteThing

	protDown     bool
	protUp       bool
	protShoot    bool
	lastProtShot float32 = -1.0e10

	wireframe bool

	windowOptions *WindowOptions = &WindowOptions{
		Title:     "GOAT",
		Width:     sceneW * 100,
		Height:    sceneH * 100,
		Resizable: false,
	}
)

// This is the actual main function
// must run on mainthread
func actualMain() {

	Setup()

	h.EnableBlending()

	//
	// timing and bookkeeping variables

	for !window.ShouldClose() {

		h.ClearScreenF(0, 0, 0, 0)

		m.Machine.Tick()

		Update()
		Draw()

		window.SwapBuffers()
		h.AssertGLOK("EndOfDraw")
		glfw.PollEvents()

	}
}

// //////////////////////////////////////////////////////////////////
// UPDATE
// //////////////////////////////////////////////////////////////////
func Update() {

	// scroll background
	bgMove := bgScrollSpeed * m.Machine.Delta
	background.UniSubTexPos = background.UniSubTexPos.Add(mgl32.Vec4{bgMove, 0, bgMove, 0})

	neutralAngle := float32(protBaseAngle)
	minAngle := neutralAngle - maxBankAngle
	maxAngle := neutralAngle + maxBankAngle
	protagonist.Velocity.VelY = 0
	snapToHigh := neutralAngle + 1*h.Degrees
	snapToLow := neutralAngle - 1*h.Degrees
	bankRotSpeed := float32(maxBankAngle) * 3
	returnRotSpeed := bankRotSpeed * 6

	if protUp && !protDown {
		// move up
		protagonist.Velocity.VelY = maxProtSpeed
		protagonist.Velocity.VelR = bankRotSpeed
	} else if protDown && !protUp {
		// move down
		protagonist.Velocity.VelY = -maxProtSpeed
		protagonist.Velocity.VelR = -bankRotSpeed
	} else if protagonist.Trans.GetAll().R > snapToHigh {
		protagonist.Velocity.VelR = -returnRotSpeed
	} else if protagonist.Trans.GetAll().R < snapToLow {
		protagonist.Velocity.VelR = returnRotSpeed
	} else {
		protagonist.Velocity.VelR = 0
		protagonist.Trans.SetRotation(neutralAngle)
	}

	if protagonist.Trans.GetAll().R > maxAngle {
		protagonist.Trans.SetRotation(maxAngle)
	} else if protagonist.Trans.GetAll().R < minAngle {
		protagonist.Trans.SetRotation(minAngle)
	}

	// change rotation graduately to point in the direction of "travel"
	if protShoot {
		protagonistShoots()
	}

	m.Machine.UpdateThings()

	{
		// do a bunch of gofunc
		// Update all shots (in a goroutine)
		// Update all enemies (in a goroutine)
		// Update all enemy shots
	}

	// Cull dead shots
	// Cull dead enemies
	// Spawn new enemies (random die roll, chance decreases the more enemies on screen
	// Spawn obstacles (maybe)
	// allow enemies to shoot
	// calculate and handle all shot collisions
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
	//
	// Here you can draw unmanaged things
	//
	m.Machine.DrawThings()
}

// //////////////////////////////////////////////////////////////////
// SETUP
// //////////////////////////////////////////////////////////////////
func Setup() {

	m.Start()

	_, _window, err := initGlfw(windowOptions)
	h.GlPanicIfErrNotNil(err)
	window = _window

	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	InitializeSpritesheet()

	initializeCamera(windowOptions)
	initializeBackground()
	initializeProtagonist()
	initializeProtShotTemplate()
	initializeKeyboardHandler()
}

func InitializeSpritesheet() {
	verts, texCoords, indeces := h.PolygonCoords(4, 45*h.Degrees, math.Sqrt2, math.Sqrt2)

	m.Machine.LoadShader("main", "shaders/polygon.vert", "shaders/polygon.frag")
	m.Machine.LoadTextureAtlas("assets/Spritesheet/sheet.xml")
	spriteSheet = m.CreateSprite("main", "sheet.png", verts, texCoords, indeces)
	spriteSheet.Finalize()
	m.Machine.Renderables["main"] = spriteSheet
}

// //////////////////////////////////////////////////////////////////
// Keyboard handler
// //////////////////////////////////////////////////////////////////
func initializeKeyboardHandler() {
	glfw.GetCurrentContext().SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
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
			if keydown {
				wireframe = !wireframe
				h.Wireframe(wireframe)
			}
		}
	})
}

func protagonistShoots() {
	if m.Machine.Now < lastProtShot+shotDelay {
		return
	}

	// TODO: More logic about if prot can fire

	spawnProtagonistShot()
	lastProtShot = m.Machine.Now
}

func spawnProtagonistShot() *m.SpriteThing {
	p := protShotTemplate

	// must have a clone feature
	shot := &m.SpriteThing{
		Renderable:   p.Renderable,
		Camera:       p.Camera,
		Trans:        p.Trans,
		UniSubTexPos: p.UniSubTexPos,
		UniColor:     p.UniColor,
		UniColorMix:  p.UniColorMix,
	}

	loc := protagonist.Trans.GetAll()

	si, co := h.Sincos(loc.R - 270*h.Degrees)

	offsetX, offsetY := loc.ScaleX*co, loc.ScaleY*si
	shot.Trans.SetPos(loc.X+offsetX, loc.Y+offsetY)
	shot.Trans.SetRotation(loc.R)
	shot.Velocity.VelX = shotSpeed * co
	shot.Velocity.VelY = shotSpeed * si
	shot.Velocity.VelR = loc.R - protBaseAngle
	shot.Velocity.VelTR = loc.R - protBaseAngle

	m.Machine.Things = append(m.Machine.Things, shot)

	return shot
}

func initializeProtShotTemplate() {
	dims := m.Machine.SubTextures["sheet.png/laserBlue01.png"].GetDims()
	shotAspect := dims[0] / dims[1]
	shot := &m.SpriteThing{
		Renderable:   spriteSheet,
		Camera:       camera,
		UniSubTexPos: dims,
	}
	const shotSize = 0.2
	shot.Trans.SetScale(shotSize, shotSize*shotAspect)
	shot.UniColor = mgl32.Vec4{1, 1, 1, 0}
	shot.UniColorMix = 0.1
	shot.Velocity.VelX = 10
	m.Machine.Named["protLaser"] = shot
	protShotTemplate = shot
}

func initializeProtagonist() {

	spriteTemplate := &m.SpriteThing{
		Renderable: spriteSheet,
		Camera:     camera,
	}
	m.Machine.Named["spriteTemplate"] = spriteTemplate

	protagonist = spriteTemplate.Clone().(*m.SpriteThing)

	protagonist.UniSubTexPos = m.Machine.SubTextures["sheet.png/playerShip1_blue.png"].GetDims()
	protagonist.UniColorMix = 0.5
	protagonist.UniColor = mgl32.Vec4{1, 1, 1, 0.5}

	protagonist.Trans.SetScale(1, 1)
	protagonist.Trans.SetRotation(protBaseAngle)
	protagonist.Trans.SetPos(-maxX+2, 0.0)

	m.Machine.Named["protagonist"] = protagonist
	m.Machine.Things = append(m.Machine.Things, protagonist)
}

func initializeBackground() {
	verts, texCoords, indeces := h.PolygonCoords(4, 45*h.Degrees, math.Sqrt2/2, math.Sqrt2/2)
	m.Machine.LoadTexture("background", "assets/Backgrounds/purple.png")

	bgRenderable := m.CreateSprite("main", "background", verts, texCoords, indeces)
	bgRenderable.Texture.SetRepeatS()
	bgRenderable.Finalize()

	background = &m.SpriteThing{
		Renderable: bgRenderable,
		Camera:     camera,
	}
	background.Trans.SetPos(0, 0)
	background.Trans.SetScale(sceneW-2, sceneH-2)
	background.UniColor = mgl32.Vec4{1, 1, 1, 1}
	background.UniColorMix = 0.0
	background.UniSubTexPos = mgl32.Vec4{0, 0, 1, 1}

	m.Machine.Named["background"] = background
	m.Machine.Things = append(m.Machine.Things, background)
}

func initializeCamera(options *WindowOptions) {

	camera = h.CreateCamera()
	camera.SetFrameSize(sceneW, sceneH)
	camera.SetPosition(0, 0)

	m.Machine.Cameras["main"] = camera
}
