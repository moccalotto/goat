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
	scriptFile   string        // filename of the lua script
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

func CreateDrawing(renderer *sdl.Renderer, scriptFile string) *Drawing {

	script := lua.NewState()

	ticks := sdl.GetTicks64()

	dm := &Drawing{
		script:       script,
		scriptFile:   scriptFile,
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

	dm.loadLuaScript()

	dm.scaleX, dm.scaleY = renderer.GetScale()
	dm.keydownCallback = luaFuncOrNil(script.GetGlobal("Keydown"))
	dm.keyupCallback = luaFuncOrNil(script.GetGlobal("Keyup"))

	return dm
}

func (dm *Drawing) Destroy() {
	dm.script.Close()
}

// Call the CallSetupFunc() function before the first frame is framed
func (dm *Drawing) CallSetupFunc() {
	if dm.frameCount > 0 {
		panic("This should never happen")
	}

	setupFunc := dm.script.GetGlobal("Setup")
	luaInvokeFunc("Setup()", dm.script, setupFunc)
}

func (dm *Drawing) loadLuaScript() {

	dm.setupLuaFunctions()

	if err := dm.script.DoFile(dm.scriptFile); err != nil {
		panic(err)
	}
	dm.drawFunc = dm.script.GetGlobal("Draw")
}

// Do the CallDrawFunc phase of the game loop.
func (dm *Drawing) CallDrawFunc() {

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

	//
	//********************************************
	// Draw all "things"
	//********************************************
	if len(dm.stack) > 0 {
		panic("You must call Pop() as many times as you have called Push()")
	}

	//
	//********************************************
	// Start framerate limit logic
	//
	// Important: setting a low frameRateCap
	// results on slow message handling.
	// longer delays should be handled
	// with the Sleep() method.
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
func (dm *Drawing) assignVarToScript(name string, value interface{}) {
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
func (dm *Drawing) setupLuaFunctions() {

	fun := func(name string, value interface{}) {
		dm.script.SetGlobal(name, luar.New(dm.script, value))
	}

	fun("WinSize", dm.WinSize)
	fun("WinTitle", dm.WinTitle)
	fun("ProcessEvents", dm.ProcessEvents)
	fun("FrameRateCap", dm.FrameRateCap)

	fun("GridSize", dm.GridSize)
	fun("CanvasSize", dm.CanvasSize)
	fun("Sleep", dm.Sleep)
	fun("Dump", dm.Dump)
	fun("Log", log.Printf)
	fun("Quit", dm.Quit)
	fun("ForceQuit", dm.ForceQuit)
	fun("FrameCount", func() uint64 { return dm.frameCount })
	fun("Pause", dm.Pause)
	fun("Resume", dm.Resume)

	fun("Scale", dm.Scale)
	fun("Color", dm.Color)
	fun("Background", dm.Background)
	fun("Dot", dm.Dot)
	fun("Push", dm.Push)
	fun("Pop", dm.Pop)
	fun("Line", dm.Line)
	fun("Rectangle", dm.Rectangle)
	fun("Polygon", dm.Polygon)
	fun("PolarVector", CreateVectorPolar)
	fun("Vector", func(composants ...float64) Vector { return Vector{components: composants} })

	/*****************************************
	 * TEST FUNCTIONS
	 ****************************************/

	fun("HasKey", dm.HasKey)
}
