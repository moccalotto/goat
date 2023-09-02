package main

import (
	"goat/shed"
	"goat/tractor"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {

	// start the mainthread system, allowing us to make calls on the main thread later
	// h.StartMainThreadSystem(actualMain)
	actualMain()
}

const (
	SCENE_W = 3000
	SCENE_H = SCENE_W * 9 / 16
	MARGIN  = 100
	MIN_X   = -SCENE_W / 2
	MAX_X   = SCENE_W / 2
	MIN_Y   = -SCENE_H / 2
	MAX_Y   = SCENE_H / 2

	PX_FACTOR       = 1 // pixels per "square"
	CAMERA_ID       = "main"
	SPRITE_SHADER   = "shaders/sprite"
	RECT_SHADER     = "shaders/rect"
	BG_TEX_FN       = "Backgrounds/purple.png"
	ATLAS_FN        = "Spritesheet/sheet.xml"
	TEST_TEX_FN     = "playerShip1_blue.png"
	BG_SCROLL_SPEED = 0.08
)

var (
	gCamera              *tractor.Camera
	gWindow              *glfw.Window
	gBackgroundSprite    *tractor.Sprite
	gMainTexQuadRenderer *tractor.TexQuadRenderer
	gMainRectRenderer    *tractor.BasicRectRenderer
	gMainSprite          *tractor.Sprite
	gMainRect            *tractor.BasicRect
	gMainLine            *tractor.BasicLine
)

// ||========================================================
// ||
// || ACTUAL MAIN FUNC
// ||
// ||========================================================
func actualMain() {

	tractor.StartMain(&tractor.WindowOptions{
		Title:     "GOAT",
		Width:     SCENE_W * PX_FACTOR,
		Height:    SCENE_H * PX_FACTOR,
		Resizable: false,
	})

	Setup()

	//
	// timing and bookkeeping variables

	tractor.Engine.Loop(func() {
		Update()
		Draw()
	})

	// Free/dispose all allocated resources
	tractor.Engine.Dispose()
}

// ||========================================================
// ||
// || Update
// ||
// ||========================================================

func Update() {

	// Background
	// =================
	bgDist := BG_SCROLL_SPEED * tractor.Engine.Delta
	gBackgroundSprite.UniSubTexPos = gBackgroundSprite.UniSubTexPos.Plus(shed.Vec4(bgDist, 0, bgDist, 0))

	// SPRITE
	// =====================================
	gMainSprite.Rotate(tractor.Engine.Delta)

	// RECTANGLE
	// =====================================
	sin0, _ := shed.Sincos(tractor.Engine.Now)
	sin120, _ := shed.Sincos(tractor.Engine.Now + shed.Tau/3)
	sin240, _ := shed.Sincos(tractor.Engine.Now + shed.Tau*2/3)
	//
	gMainRect.Rotate(-tractor.Engine.Delta * 2)
	gMainRect.Color.C1 = 0.5 + sin0*0.5
	gMainRect.Color.C2 = 0.5 + sin120*0.5
	gMainRect.Color.C3 = 0.5 + sin240*0.5
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

	p1 := shed.Vec2(-250, -250) // lower left
	p2 := shed.Vec2(250, -250)  // lower right
	gMainLine.SetColor(shed.RGBA(1, 1, 1, 1))
	gMainLine.SetPoints(p1, p2)
	gMainLine.Draw()

	p1 = p2
	p2 = shed.Vec2(250, 250) // upper right
	gMainLine.SetColor(shed.RGBA(1, 0, 1, 1))
	gMainLine.SetPoints(p1, p2)
	gMainLine.Draw()

	p1 = p2
	p2 = shed.Vec2(-250, 250) // upper left
	gMainLine.SetColor(shed.RGBA(1, 1, 0, 1))
	gMainLine.SetPoints(p1, p2)
	gMainLine.Draw()

	p1 = p2
	p2 = shed.Vec2(-250, -250) // lower left
	gMainLine.SetColor(shed.RGBA(0, 1, 1, 1))
	gMainLine.SetPoints(p1, p2)
	gMainLine.Draw()
}

// ||========================================================
// ||
// || SETUP
// ||
// ||========================================================
func Setup() {

	shed.EnableBlending()

	tractor.Engine.AssetPath = "assets"

	initCamera() // Must be called fairly early
	initKeyboardHandler()
	initBackground()
	initMainSprite()
	initBasicRect()
	gMainLine = tractor.CreateBasicLine(0, 0, 0, 0, 50, gCamera, gMainRectRenderer)

}
