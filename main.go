package main

import (
	m "goat/motor"
	h "goat/util"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func main() {

	// start the mainthread system, allowing us to make calls on the main thread later
	// h.StartMainThreadSystem(actualMain)
    actualMain()
}

const (
	SCENE_W = 3000
	SCENE_H = SCENE_W * 9 / 16
	MARGIN  = 0.5
	MIN_X   = -SCENE_W / 2
	MAX_X   = SCENE_W / 2
	MIN_Y   = -SCENE_H / 2
	MAX_Y   = SCENE_H / 2

	PX_FACTOR     = 1 // pixels per "square"
	CAMERA_ID     = "main"
	SPRITE_SHADER = "shaders/sprite"
	RECT_SHADER   = "shaders/rect"
	BG_TEX_FN     = "Backgrounds/purple.png"
	ATLAS_FN      = "Spritesheet/sheet.xml"
	TEST_TEX_FN   = "playerShip1_blue.png"
)

var (
	gBgScrollSpeed       float32 = 0.08
	gCamera              *m.Camera
	gWindow              *glfw.Window
	gBackgroundSprite    *m.Sprite
	gMainTexQuadRenderer *m.TexQuadRenderer
	gMainRectRenderer    *m.BasicRectRenderer
	gMainSprite          *m.Sprite
	gMainRect            *m.BasicRect
	gMainLine            *m.BasicLine

	gWindowOptions *WindowOptions = &WindowOptions{
		Title:     "GOAT",
		Width:     SCENE_W * PX_FACTOR,
		Height:    SCENE_H * PX_FACTOR,
		Resizable: false,
	}
)

// ||========================================================
// ||
// || ACTUAL MAIN FUNC
// ||
// ||========================================================
func actualMain() {

	// || MUST BE THE FIRST THING WE DO
	// ||
	// ||
	// || TODO: Move to Machine
	// ||
	// ||=======================================
	_, _window, err := initGlfw(gWindowOptions)
	h.GlPanicIfErrNotNil(err)
	gWindow = _window

	Setup()

	//
	// timing and bookkeeping variables

	for !gWindow.ShouldClose() {
		h.Clear()

		m.Machine.Tick()

		Update()
		Draw()

		gWindow.SwapBuffers()

		h.AssertGLOK("EndOfDraw")

		glfw.PollEvents()

	}
}

// ||========================================================
// ||
// || Update
// ||
// ||========================================================

func Update() {
	sin, cos := h.Sincos(m.Machine.Now)
	// Background
	// =================
	bgDist := gBgScrollSpeed * m.Machine.Delta
	gBackgroundSprite.UniSubTexPos = gBackgroundSprite.UniSubTexPos.Add(mgl32.Vec4{bgDist, 0, bgDist, 0})

	// Rotate main sprite
	gMainSprite.Rotate(m.Machine.Delta)

	// Rotate and color main rect
	gMainRect.Rotate(-m.Machine.Delta * 2)
	gMainRect.Color.Y = 0.5 + sin*0.5
	gMainRect.Color.Z = 0.5 + cos*0.5
}

// ||========================================================
// ||
// || DRAW
// ||
// ||========================================================
func Draw() {

	gBackgroundSprite.Draw()

	// Main Sprite
	// ================
	gMainSprite.Draw()

	// Main Rect
	// ================
	gMainRect.Draw()

	gMainLine = m.CreateBasicLine(
		-250, 0, // pt1
		250, 0, // pt2
		50, // thickness
		gCamera,
		gMainRectRenderer,
	)
	gMainLine.Draw()
}

// ||========================================================
// ||
// || SETUP
// ||
// ||========================================================
func Setup() {

	m.Start()

	h.EnableBlending()

	m.Machine.AssetPath = "assets"

	initCamera() // Must be called fairly early
	initKeyboardHandler()
	initBackground()
	initMainSprite()
	initBasicRect()
}
