package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
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
	paused       bool          // don't draw until Unpause() is called

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
		frameRateCap: -1.0,
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

// Do the draw phase of the game loop.
func (dm *Drawing) draw() {
	//
	//********************************************
	// Present last loop's renderings
	//********************************************
	dm.renderer.Present() // TODO if autopresent

	// The number of times draw() has been called so far.
	dm.frameCount++

	// The high resolution time at the beginning of the cycle.
	// It has sub-ms resolution and is usd to limit FPS
	startTicksHq := sdl.GetPerformanceCounter()

	// TODO: if dm.autoclear
	dm.renderer.SetDrawColor(dm.bgColor.R, dm.bgColor.G, dm.bgColor.B, dm.bgColor.A)
	dm.renderer.Clear()

	// Conventional tick counters
	dm.prevTicks = dm.nowTicks                 // Time when the previous update began
	dm.nowTicks = sdl.GetTicks64()             // Time the when the current update begins
	dm.deltaTicks = dm.nowTicks - dm.prevTicks // number of ms since last update

	//
	//********************************************
	// Call the Draw() function
	//********************************************
	if !dm.paused {
		luaInvokeFunc("Draw()", dm.script, dm.drawFunc)
	}

	if len(dm.stack) > 0 {
		panic("You must call Pop() as many times as you have called Push()")
	}

	//
	//********************************************
	// Start framerate limit logic
	//********************************************
	if dm.frameRateCap <= 0 {
		return
	}

	endTicksHq := sdl.GetPerformanceCounter() // The high resolution time of the end of the cycle.
	elapsedMs := float32(endTicksHq-startTicksHq) / float32(sdl.GetPerformanceFrequency()*1000)
	wantedFrameTimeMs := 1000.0 / dm.frameRateCap

	// use a raw SDL delay - no pausing to check for keyboard events.
	// so dm.frameRateCap should not be too low.
	sdl.Delay(uint32(wantedFrameTimeMs - elapsedMs))
}

// Before we draw something, ensure that the settings are applyed to the
// renderer. This is necessary for us to push/pop graphics settings.
func (dm *Drawing) applySettingsToRenderer() {
	dm.renderer.SetScale(dm.scaleX, dm.scaleY)
	dm.renderer.SetDrawColor(dm.fgColor.R, dm.fgColor.G, dm.fgColor.R, dm.fgColor.A)
}

// Helper function to inject a variable or function into the script
func (dm *Drawing) setGlobalScriptEntry(name string, value interface{}) {
	dm.script.SetGlobal(
		name,
		luar.New(dm.script, value),
	)
}

// triggered whenever our game loop receives a keydown event
func (dm *Drawing) onKeydown(ke *KeyEvent) {
	luaInvokeFunc("Keydown", dm.script, dm.keydownCallback, luar.New(dm.script, ke))
}

// triggered whenever our game loop receives a keyup event
func (dm *Drawing) onKeyup(ke *KeyEvent) {
	luaInvokeFunc("Keyup", dm.script, dm.keyupCallback, luar.New(dm.script, ke))
}

// Functions sare injected into the script only once.
func (dm *Drawing) exportFunctions() {

	dm.setGlobalScriptEntry("WinSize", dm.WinSize)
	dm.setGlobalScriptEntry("WinTitle", dm.WinTitle)
	dm.setGlobalScriptEntry("ProcessEvents", dm.ProcessEvents)
	dm.setGlobalScriptEntry("FrameRateCap", dm.FrameRateCap)

	dm.setGlobalScriptEntry("GridSize", dm.GridSize)
	dm.setGlobalScriptEntry("CanvasSize", dm.CanvasSize)
	dm.setGlobalScriptEntry("Sleep", dm.Sleep)
	dm.setGlobalScriptEntry("Log", log.Printf)
	dm.setGlobalScriptEntry("Quit", dm.Quit)
	dm.setGlobalScriptEntry("ForceQuit", dm.ForceQuit)
	dm.setGlobalScriptEntry("FrameCount", func() uint64 { return dm.frameCount })
	dm.setGlobalScriptEntry("Pause", dm.Pause)
	dm.setGlobalScriptEntry("Resume", dm.Resume)

	dm.setGlobalScriptEntry("Scale", dm.Scale)
	dm.setGlobalScriptEntry("Color", dm.Color)
	dm.setGlobalScriptEntry("Background", dm.Background)
	dm.setGlobalScriptEntry("Dot", dm.Dot)
	dm.setGlobalScriptEntry("Push", dm.Push)
	dm.setGlobalScriptEntry("Pop", dm.Pop)
	dm.setGlobalScriptEntry("Line", dm.Line)
	dm.setGlobalScriptEntry("Rectangle", dm.Rectangle)
	dm.setGlobalScriptEntry("Polygon", dm.Polygon)

	/*****************************************
	 * TEST FUNCTIONS
	 ****************************************/

	dm.setGlobalScriptEntry("PolarVector", PolarVector)
	dm.setGlobalScriptEntry("Vector", func(composants ...float64) VectorType {
		return VectorType{
			elements: composants,
		}
	})

}
