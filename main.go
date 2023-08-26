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

const (
	SCENE_W        = 20
	SCENE_H        = SCENE_W * 9 / 16
	PX_FACTOR      = 150 // pixels per "square"
	CAMERA_ID      = "mainCamera"
	SHADER_ID      = "shaders/sprite"
	BACKGROUND_TEX = "Backgrounds/purple.png"
	ATLAS_ID       = "Spritesheet/sheet.xml"
	HERO_TEX_ID    = "playerShip1_blue.png"
)

var (
	gScrollSpeed      float32 = 0.08
	gCamera           *m.Camera
	gWindow           *glfw.Window
	gBackgroundEntity *m.BasicEntity
	gGodSprite        *m.Sprite
	gHero             *m.BasicEntity

	gWindowOptions *WindowOptions = &WindowOptions{
		Title:     "GOAT",
		Width:     SCENE_W * PX_FACTOR,
		Height:    SCENE_H * PX_FACTOR,
		Resizable: false,
	}
)

// This is the actual main function
// must run on mainthread
func actualMain() {

	///
	/// MUST BE THE FIRST THING WE DO
	////////////////////////////////////////////
	_, _window, err := initGlfw(gWindowOptions)
	h.GlPanicIfErrNotNil(err)
	gWindow = _window

	Setup()

	//
	// timing and bookkeeping variables

	for !gWindow.ShouldClose() {

		h.ClearScreenF(0, 0, 0, 0)

		m.Machine.Tick()

		Update()
		Draw()

		gWindow.SwapBuffers()
		h.AssertGLOK("EndOfDraw")
		glfw.PollEvents()

	}
}

// //////////////////////////////////////////////////////////////////
// UPDATE
// //////////////////////////////////////////////////////////////////
func Update() {
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

	///
	/// scroll background
	//////////////////////
	bgMove := gScrollSpeed * m.Machine.Delta
	gBackgroundEntity.UniSubTexPos = gBackgroundEntity.UniSubTexPos.Add(mgl32.Vec4{bgMove, 0, bgMove, 0})

	gBackgroundEntity.Draw()
	gHero.Draw()
}

// //////////////////////////////////////////////////////////////////
// SETUP
// //////////////////////////////////////////////////////////////////
func Setup() {

	h.EnableBlending()

	m.Start()

	m.Machine.AssetPath = "assets"

	// initCamera must do that pretty early
	initializeCamera(gWindowOptions)

	initializeKeyboardHandler()
	initializeGodSprite()
	initializeBackground()
	initializeHero()
}

func initializeHero() {
	sprite := gGodSprite.Clone()

	gHero = &m.BasicEntity{
		Renderable:   sprite,
		Camera:       gCamera,
		UniColor:     mgl32.Vec4{1, 1, 1, 1},
		UniColorMix:  0.5,
		UniSubTexPos: m.Machine.GetDimsForSubtexture(ATLAS_ID, HERO_TEX_ID),
	}

	gHero.SetScale(1, 1)
	gHero.Rotate(-90 * h.Degrees)
}

// /
// /
// / BACKGROUND
// //////////////////////////////////////////////
func initializeBackground() {
	// verts, texCoords, indeces := h.SquareCoords()
	scale := math.Sqrt2 / 2
	verts, texCoords, indeces := h.PolygonCoords(4, 45*h.Degrees, scale, scale)

	bgSprite := m.CreateSprite(SHADER_ID, BACKGROUND_TEX, verts, texCoords, indeces)
	bgSprite.Texture.SetRepeatS()
	bgSprite.Finalize()

	gBackgroundEntity = &m.BasicEntity{
		Renderable:   bgSprite,
		Camera:       gCamera,
		UniColor:     mgl32.Vec4{},
		UniColorMix:  0.0,
		UniSubTexPos: mgl32.Vec4{0, 0, 1, 1},
	}
	gBackgroundEntity.SetScale(SCENE_W-2, SCENE_H-2)
}

// /
// /
// / CAMERA
// //////////////////////////////////////////////
func initializeCamera(options *WindowOptions) {

	gCamera = m.Machine.GetCamera(CAMERA_ID)
	gCamera.SetFrameSize(SCENE_W, SCENE_H)
	gCamera.SetPosition(0, 0)

	m.Machine.Cameras[CAMERA_ID] = gCamera
}

// /
// /
// / The main sprite that all other sprites are copied from
// //////////////////////////////////////////////
func initializeGodSprite() {
	verts, texCoords, indeces := h.SquareCoords()
	gGodSprite = m.CreateSpriteFromAtlas(SHADER_ID, ATLAS_ID, HERO_TEX_ID, verts, texCoords, indeces)
	gGodSprite.Finalize()
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

		if keydown && key == glfw.KeyEscape {
			glfw.GetCurrentContext().SetShouldClose(true)
		}
	})
}
