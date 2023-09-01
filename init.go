package main

import (
	m "goat/motor"
	"goat/util"

	"github.com/go-gl/glfw/v3.3/glfw"
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
		0*util.Degrees,
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

	gCamera, _ = m.Machine.GetCamera(CAMERA_ID)
	gCamera.SetFrameSize(SCENE_W, SCENE_H)
	gCamera.SetXY(0, 0)
}

// ||=================================================================
// ||
// || Handle Keyboard events
// ||
// ||=================================================================
func initKeyboardHandler() {
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
