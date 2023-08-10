package main

import (
	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
)

// My dream is to do all of this:: https://github.com/fogleman/gg

type Drawing struct {
	renderer     *sdl.Renderer // The sdl renderer we're using. It's most likely a hardware renderer.
	script       *lua.LState   // The lua script we're running
	drawFunc     lua.LValue    // The lua function to call to draw
	prevTicks    uint64        // ticks when previous draw started
	nowTicks     uint64        // ticks when current draw started
	deltaTicks   uint64        // number of ticks the last render cycle took
	fgColor      sdl.Color     // Foreground color. The color of the "ink" if you will
	bgColor      sdl.Color     // Background color. The color of the "paper"
	scaleX       float32       // How big (in pixels) are the virtual pixels in the X direction.
	scaleY       float32       // How big (in pixels) are the virtual pixels in the Y direction.
	frameRateCap float32       // Maximum allowed number of updates per second. - During this delay no events are processed.
	stack        []*Drawing    // The stack that allows us to store and recall colors, scales, and other such settings.
	frameCount   uint64        // The number of calls to draw(). Starts at 1

	keydownCallback *lua.LFunction
	keyupCallback   *lua.LFunction
}

func CreateDrawing(renderer *sdl.Renderer, script *lua.LState) *Drawing {

	ticks := sdl.GetTicks64()

	dm := &Drawing{
		script:       script,
		renderer:     renderer,
		drawFunc:     script.GetGlobal("Draw"),
		nowTicks:     ticks,
		prevTicks:    ticks,
		deltaTicks:   0,
		fgColor:      sdl.Color{R: 0, G: 0, B: 0, A: 255},
		bgColor:      sdl.Color{R: 255, G: 255, B: 255, A: 255},
		scaleX:       1,
		scaleY:       1,
		frameRateCap: 0,
		stack:        make([]*Drawing, 0),
	}

	dm.scaleX, dm.scaleY = renderer.GetScale()
	dm.keydownCallback = luaFuncOrNil(script.GetGlobal("Keydown"))
	dm.keyupCallback = luaFuncOrNil(script.GetGlobal("Keyup"))

	return dm
}

// Call the setup() function before the first frame is framed
func (dm *Drawing) setup() {
	if dm.frameCount > 0 {
		return
	}

	setupFunc := dm.script.GetGlobal("Setup")
	luaInvokeFunc("Setup()", dm.script, setupFunc)
}
