package main

import (
	"log"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// My dreap is to do all of this:: https://github.com/fogleman/gg

type Drawing struct {
	renderer   *sdl.Renderer
	script     *lua.LState
	drawFunc   lua.LValue
	prevTicks  uint64
	nowTicks   uint64
	deltaTicks uint64
}

func CreateDrawing(renderer *sdl.Renderer, script *lua.LState) *Drawing {
	sdl.GetTicks64()
	renderer.SetScale(1, 1)

	ticks := sdl.GetTicks64()
	drawing := &Drawing{
		script:     script,
		renderer:   renderer,
		drawFunc:   script.GetGlobal("Draw"),
		nowTicks:   ticks,
		prevTicks:  ticks,
		deltaTicks: 0,
		// keypresses, etc
	}

	return drawing
}

func (dm *Drawing) draw() {
	dm.prevTicks = dm.nowTicks
	dm.nowTicks = sdl.GetTicks64()
	dm.deltaTicks = dm.nowTicks - dm.prevTicks

	// aim for around 60 fps
	if dm.deltaTicks < 16 {
		delay := uint32(16 - dm.deltaTicks)
		sdl.Delay(delay)
		log.Printf("Delay: %d\n", delay)
	}

	dm.injectVariables()

	invokeLuaFunc("Draw()", dm.script, dm.drawFunc)
}

func (dm *Drawing) injectVariables() {
	x, y, err := dm.renderer.GetOutputSize()
	dm.setGlobalScriptEntry("Width", x)
	dm.setGlobalScriptEntry("Height", y)
	dm.setGlobalScriptEntry("Ticks", dm.nowTicks)
	dm.setGlobalScriptEntry("PrevTicks", dm.prevTicks)
	dm.setGlobalScriptEntry("DeltaTicks", dm.deltaTicks)

	if err != nil {
		panic(err)
	}

	// keyboard events
}

func (dm *Drawing) injectFunctions() {
	dm.setGlobalScriptEntry("Line", dm.Line)
	dm.setGlobalScriptEntry("Scale", dm.Scale)
	dm.setGlobalScriptEntry("GetSize", dm.GetSize)
	dm.setGlobalScriptEntry("SleepMs", dm.SleepMs)
	dm.setGlobalScriptEntry("Log", log.Printf)
	dm.setGlobalScriptEntry("ForceQuit", dm.ForceQuit)
	dm.setGlobalScriptEntry("Quit", dm.ForceQuit)
	dm.setGlobalScriptEntry("Background", dm.Background)
	dm.setGlobalScriptEntry("Color", dm.Color)
}

func (dm *Drawing) setGlobalScriptEntry(name string, value interface{}) {
	dm.script.SetGlobal(
		name,
		luar.New(dm.script, value),
	)
}

func (dm *Drawing) ForceQuit(code int, why string) {
	log.Printf("ForceQuit() called from lua. Code %d, reason: %s\n", code, why)
	os.Exit(code)
}

func (dm *Drawing) Color(r, g, b, a uint8) {
	dm.renderer.SetDrawColor(r, g, b, a)
}
func (dm *Drawing) Background(c ...uint8) {
	_x, _y, _z, _a, _ := dm.renderer.GetDrawColor()

	switch len(c) {
	case 1:
		dm.renderer.SetDrawColor(c[0], c[0], c[0], 255)
	case 3:
		dm.renderer.SetDrawColor(c[0], c[1], c[2], 255)
	case 4:
		dm.renderer.SetDrawColor(c[0], c[1], c[2], c[3])
	default:
		panic("Background() takes 1, 3, or 4 arguments.")

	}
	dm.renderer.Clear()
	dm.renderer.SetDrawColor(_x, _y, _z, _a)
}

func (dm *Drawing) Quit() {
	sdl.PushEvent(
		&sdl.QuitEvent{
			Type:      sdl.QUIT,
			Timestamp: uint32(sdl.GetTicks64()),
		},
	)
}

func (dm *Drawing) Scale(scale float32) {
	dm.renderer.SetScale(scale, scale)
}

func (dm *Drawing) Line(x1, y1, x2, y2 float32) {
	dm.renderer.DrawLineF(x1, y1, x2, y2)
}

func (dm *Drawing) GetSize() sdl.Rect {
	return dm.renderer.GetViewport()
}

func (dm *Drawing) SleepMs(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
