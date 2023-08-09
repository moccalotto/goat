package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	luar "layeh.com/gopher-luar"
)

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

func (dm *Drawing) SetWinSize(x, y int32, center ...bool) {
	w, err := dm.renderer.GetWindow()

	if err != nil {
		panic(err)
	}

	if x > 0 && y > 0 {
		w.SetSize(x, y)
	}

	if len(center) > 0 && center[0] {
		w.SetPosition(
			sdl.WINDOWPOS_CENTERED,
			sdl.WINDOWPOS_CENTERED,
		)
	}
}

func (dm *Drawing) SetWinTitle(title string) {
	w, err := dm.renderer.GetWindow()

	if err != nil {
		panic(err)
	}

	w.SetTitle(title)
}

// Functions sare injected into the script only once.
func (dm *Drawing) injectFunctions() {

	dm.setGlobalScriptEntry("SetWinSize", dm.SetWinSize)
	dm.setGlobalScriptEntry("SetWinTitle", dm.SetWinTitle)

	dm.setGlobalScriptEntry("Line", dm.Line)
	dm.setGlobalScriptEntry("Scale", dm.Scale)
	dm.setGlobalScriptEntry("GridSize", dm.GridSize)
	dm.setGlobalScriptEntry("CanvasSize", dm.CanvasSize)
	dm.setGlobalScriptEntry("Sleep", dm.Sleep)
	dm.setGlobalScriptEntry("Log", log.Printf)
	dm.setGlobalScriptEntry("Quit", dm.Quit)
	dm.setGlobalScriptEntry("ForceQuit", dm.ForceQuit)
	dm.setGlobalScriptEntry("Background", dm.Background)
	dm.setGlobalScriptEntry("Color", dm.Color)
	dm.setGlobalScriptEntry("Dot", dm.Dot)
	dm.setGlobalScriptEntry("Push", dm.Push)
	dm.setGlobalScriptEntry("Pop", dm.Pop)
	dm.setGlobalScriptEntry("ProcessEvents", dm.ProcessEvents)
	dm.setGlobalScriptEntry("Autorender", dm.Autorender)
	dm.setGlobalScriptEntry("Rectangle", dm.Rectangle)
	dm.setGlobalScriptEntry("Counter", func() uint64 { return dm.frameCount })

}
