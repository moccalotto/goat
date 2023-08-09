package main

import "github.com/veandco/go-sdl2/sdl"

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
	luaInvokeFunc("Draw()", dm.script, dm.drawFunc)

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
