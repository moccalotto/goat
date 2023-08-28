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
	SQRT2_HALF       = math.Sqrt2 / 2
	SCENE_W          = 20
	SCENE_H          = SCENE_W * 9 / 16
	MARGIN           = 0.5
	MIN_X            = -SCENE_W / 2
	MAX_X            = SCENE_W / 2
	MIN_Y            = -SCENE_H / 2
	MAX_Y            = SCENE_H / 2
	HERO_BASE_SPEED  = SCENE_H * 3 / 4
	PX_FACTOR        = 150 // pixels per "square"
	CAMERA_ID        = "mainCamera"
	SHADER_ID        = "shaders/sprite"
	BACKGROUND_TEX   = "Backgrounds/purple.png"
	ATLAS_ID         = "Spritesheet/sheet.xml"
	HERO_TEX_ID      = "playerShip1_blue.png"
	LASER_BLUE_03_ID = "laserBlue03.png"
)

var (
	gBgScrollSpeed    float32 = 0.08
	gCamera           *m.Camera
	gWindow           *glfw.Window
	gBackgroundEntity *m.BasicEntity
	gGodSprite        *m.Sprite
	gHero             *m.BasicEntity
	gShots            []*m.BasicEntity = make([]*m.BasicEntity, 0, 200)
	gEnemies          []*m.BasicEntity = make([]*m.BasicEntity, 0, 200)

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

// /////////////////////////////////////////////////////////////////
//
// # UPDATE
//
// *
// *
// *
// *
// *
// *
// *
// *
// /////////////////////////////////////////////////////////////////
func Update() {

	gHero.Update()

    for i, shot := range gShots {
        if shot == nil {
            continue
        }
        shot.Update()


        x, y, _ := shot.Get()
        if x > MAX_X || x < MIN_X || y > MAX_Y || y < MIN_Y {
            gShots[i] = nil
        }
    }

    if m.Machine.TickCount % 60 == 0 {
        print(len(gShots))
        print("\n")
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

	///
	/// scroll background
	//////////////////////
	bgDist := gBgScrollSpeed * m.Machine.Delta
	gBackgroundEntity.UniSubTexPos = gBackgroundEntity.UniSubTexPos.Add(mgl32.Vec4{bgDist, 0, bgDist, 0})

	gBackgroundEntity.Draw()
	gHero.Draw()
    for _, shot := range gShots {
        if shot == nil {
            continue
        }

        shot.Draw()
    }
}

// //////////////////////////////////////////////////////////////////
// SETUP
// //////////////////////////////////////////////////////////////////
func Setup() {

	m.Start()

	h.EnableBlending()

	m.Machine.AssetPath = "assets"

	initializeCamera() // Must be called fairly early
	initializeKeyboardHandler()
	initializeGodSprite()
	initializeBackground()
	initializeHero()
}

// /
// /
// / BACKGROUND
// //////////////////////////////////////////////
func initializeBackground() {
	verts, texCoords, indeces := h.SquareCoords()

	bgSprite := m.CreateSprite(SHADER_ID, BACKGROUND_TEX, verts, texCoords, indeces)
	bgSprite.Texture.SetRepeatS()
	bgSprite.Finalize()

	gBackgroundEntity = &m.BasicEntity{
		Sprite:       bgSprite,
		Camera:       gCamera,
		UniColor:     mgl32.Vec4{},
		UniColorMix:  0.0,
		UniSubTexPos: mgl32.Vec4{0, 0, 1, 1},
	}

	gBackgroundEntity.SetScale(SCENE_W-MARGIN*2, SCENE_H-MARGIN*2)
}

// /
// /
// / CAMERA
// //////////////////////////////////////////////
func initializeCamera() {

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
	glfw.GetCurrentContext().SetKeyCallback(func(_ /* key */ *glfw.Window, key glfw.Key, _ /* scancode */ int, action glfw.Action, _ /* mods */ glfw.ModifierKey) {
		if action == glfw.Repeat {
			return
		}
		keydown := action == glfw.Press

		if keydown && key == glfw.KeyEscape {
			glfw.GetCurrentContext().SetShouldClose(true)
		}
	})
}
