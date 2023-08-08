package main

import (
	"log"

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

// Functions sare injected into the script only once.
func (dm *Drawing) injectFunctions() {
	dm.setGlobalScriptEntry("Line", dm.Line)
	dm.setGlobalScriptEntry("Scale", dm.Scale)
	dm.setGlobalScriptEntry("GetViewSize", dm.GetViewSize)
	dm.setGlobalScriptEntry("GetCanvasSize", dm.GetCanvasSize)
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
}
