package main

import (
	"goat/glhelp"
	"log"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// My dream is to do all of this:: https://github.com/fogleman/gg
// also https://github.com/faiface/pixel/

type Color struct{ R, G, B, A uint8 }

type Entity interface {
	Draw(dm *Drawing)
	GetId() uint64
	SetId(id uint64)
}

type Drawing struct {
	window       *glfw.Window      // The window we're drawing on
	script       *lua.LState       // The lua script we're running
	scriptFile   string            // filename of the lua script
	drawFunc     lua.LValue        // The lua function to call to draw
	prevTime     float64           // time when previous draw phase started (in seconds since program start)
	nowTime      float64           // time when current draw phase started (in seconds since program start)
	deltaTime    float64           // number of seconds last draw phase took
	fgColor      Color             // Foreground color. The color of the "ink" if you will
	bgColor      Color             // Background color. The color of the "paper"
	scaleX       float32           // How big (in pixels) are the virtual pixels in the X direction.
	scaleY       float32           // How big (in pixels) are the virtual pixels in the Y direction.
	frameRateCap float64           // Maximum allowed number of updates per second. - During this delay no events are processed.
	stack        []*Drawing        // The stack that allows us to store and recall colors, scales, and other such settings.
	frameCount   uint64            // The number of calls to draw(). Starts at 1
	paused       bool              // don't draw until Unpause() is called
	entities     map[uint64]Entity // list of entities to draw
	nextEntityId uint64            // total number of entities added
	createdAt    uint64

	keydownCallback *lua.LFunction
	keyupCallback   *lua.LFunction
}

func CreateDrawing(window *glfw.Window, scriptFile string) *Drawing {

	script := lua.NewState()

	glfw.SetTime(0)
	dm := &Drawing{
		window:       window,
		script:       script,
		scriptFile:   scriptFile,
		drawFunc:     script.GetGlobal("Draw"),
		nowTime:      0.0,
		prevTime:     0.0,
		deltaTime:    0.0,
		fgColor:      Color{R: 0, G: 0, B: 0, A: 255},
		bgColor:      Color{R: 255, G: 255, B: 255, A: 255},
		scaleX:       1,
		scaleY:       1,
		frameRateCap: -1.0,
		stack:        make([]*Drawing, 0),
		entities:     make(map[uint64]Entity),
		createdAt:    uint64(time.Now().UnixMilli()),
	}

	dm.loadLuaScript()

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
	/// RENDERER dm.renderer.Present() // TODO if autopresent

	// The number of times draw() has been called so far.
	dm.frameCount++

	// Conventional tick counters
	dm.prevTime = dm.nowTime                // Time when the previous update began
	dm.nowTime = glfw.GetTime()             // number of seconds since program started
	dm.deltaTime = dm.nowTime - dm.prevTime // number of seconds since last update

	glhelp.ClearI(dm.bgColor.R, dm.bgColor.G, dm.bgColor.B, dm.bgColor.A)

	glhelp.Triangle()

	//
	//********************************************
	// Call the Draw() function
	//********************************************
	if !dm.paused {
		luaInvokeFunc("Draw()", dm.script, dm.drawFunc)
	}

	for _, ent := range dm.entities {
		ent.Draw(dm)
	}

	dm.window.SwapBuffers()

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

	// TODO: implement delay for framerate cap
	now := glfw.GetTime()
	elapsed := now - dm.nowTime
	secPerFrame := 1.0 / dm.frameRateCap

	if elapsed > secPerFrame {
		return
	}

	// Sleep() ignores negative values, so no need to check for them here
	timeToSleep := secPerFrame - elapsed
	time.Sleep(time.Duration(timeToSleep) * time.Second)

	glfw.SwapInterval(1) // VSync:   1 = ON, 0 = OFF - negative values are allowed on certain GPUs
	dm.window.SwapBuffers()
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
	fun("Delta", func() float64 {
		return float64(dm.deltaTime) / 1000
	})

	fun("GridSize", "TODO")
	fun("CanvasSize", "TODO")
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
	fun("PolarVector", CreatePolarVector)
	fun("Vector", CreateVector)

	/*****************************************
	 * TEST FUNCTIONS
	 ****************************************/

	fun("HasKey", dm.HasKey)

	fun("ELine", func(x1, y1, x2, y2 float64) *ELine {
		el := CreateELine(x1, y1, x2, y2)
		return dm.AddEntity(el).(*ELine)
	})

	fun("ELineP", func(pivX, pivY, length, radians float64) *ELine {
		el := CreateELineX(pivX, pivY, length, radians)
		return dm.AddEntity(el).(*ELine)
	})

}
