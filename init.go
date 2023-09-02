package main

import (
	"goat/shed"
	m "goat/tractor"
)

// ||=================================================================
// ||
// || BASIC RECT
// ||
// ||=================================================================
func initBasicRect() {
	gMainRectRenderer = m.CreateBasicRectRenderer(RECT_SHADER)
	gMainRectRenderer.Finalize()
	gMainRect = m.CreateBasicRect(
		0, 0, // x, y
		500, 500, // w, h
		0*shed.Degrees,
		gCamera,
		gMainRectRenderer,
	)
}

// ||=================================================================
// ||
// || SPRITE
// ||
// ||=================================================================
func initMainSprite() {

	gMainTexQuadRenderer = m.CreateTexAtlasRenderer(SPRITE_SHADER, ATLAS_FN, TEST_TEX_FN)
	gMainTexQuadRenderer.Finalize()
	gMainSprite = m.CreateSpriteAdv(gMainTexQuadRenderer, gCamera)

	gMainSprite.SetXY(MIN_X, MIN_Y)
	gMainSprite.SetScale(4, 4)
}

// ||=================================================================
// ||
// || BACKGROUND
// ||
// ||=================================================================
func initBackground() {
	bgQuad := m.CreateTexQuadRenderer(SPRITE_SHADER, BG_TEX_FN)
	bgQuad.Texture.SetRepeatS()
	bgQuad.Finalize()

	gBackgroundSprite = m.CreateSpriteAdv(bgQuad, gCamera)

	gBackgroundSprite.SetScale(SCENE_W-MARGIN*2, SCENE_H-MARGIN*2)
}

// ||=================================================================
// ||
// || CAMERA
// ||
// ||=================================================================
func initCamera() {

	gCamera, _ = m.Engine.GetCamera(CAMERA_ID)
	gCamera.SetFrameSize(SCENE_W, SCENE_H)
	gCamera.SetXY(0, 0)
}

// ||=================================================================
// ||
// || Handle Keyboard events
// ||
// ||=================================================================
func initKeyboardHandler() {

	m.Controls.HandleKeys(func(kev *m.KeyEvent) {
		if kev.Escape {
			m.Engine.GracefulShutdown()
		}
	})

}
