package main

import (
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func (dm *Drawing) ProcessEvents(maxEventsToProcess ...int) bool {
	// Only process a few events at a time, unless otherwise specified.
	max := 3
	if len(maxEventsToProcess) > 0 {
		max = maxEventsToProcess[0]
	}

	// Sanity check. Max should never be too large
	if max > 100 {
		log.Printf("The »max« arg is too large. It is %d, but it cannot be larger than %d\n", max, 100)
		max = 100
	}

	for i := 0; i < max; i++ {
		event := sdl.PollEvent()

		if event == nil {
			return true // no events om the event buffer, let's leave
		}

		switch t := event.(type) {
		case *sdl.KeyboardEvent:
			if t.State == sdl.PRESSED && t.Repeat == 0 {
				ke := CreateKeyEvent(t)
				dm.onKeydown(ke)
			}
			// if t.State == sdl.PRESSED && t.Repeat > 0 {
			// ke := CreateKeyEvent(t)
			// dm.onKeyRepeat(ke)
			// }
			if t.State == sdl.RELEASED {
				ke := CreateKeyEvent(t)
				dm.onKeyup(ke)
			}

			return true
		case *sdl.QuitEvent:
			log.Printf("Received Quit event: %+v", event)
			dm.ForceQuit(0)
			return false
		}

		if i == max-1 {
			log.Println("Event buffer full")
		}
	}

	return true
}

func (dm *Drawing) ForceQuit(code ...int) {
	sdl.Quit()

	exitCode := 0
	if len(code) > 0 {
		exitCode = code[0]
	}

	os.Exit(exitCode)
}

func (dm *Drawing) Quit() {

	quitEvent := sdl.QuitEvent{
		Type:      sdl.QUIT,
		Timestamp: uint32(sdl.GetTicks64()),
	}

	sdl.PushEvent(&quitEvent)
}

func (dm *Drawing) GetCanvasSize() (int32, int32) {
	w, h, _ := dm.renderer.GetOutputSize()

	return w, h
}

func (dm *Drawing) GetViewSize() (int32, int32) {
	vp := dm.renderer.GetViewport()

	return vp.W, vp.H
}

// Sleep until the SDL tickcounter == end
// every chunk_size_ms we wake up and process events.
func (dm *Drawing) SleepUntil(end uint64, chunk_size_ms uint64) {

	now := sdl.GetTicks64()
	delta := end - now

	if delta <= uint64(chunk_size_ms) {
		sdl.Delay(uint32(delta))
		return
	}

	dm.ProcessEvents()

	// we don't want to wake up for event processing more
	// than 50 times per sleep session
	safe_chunk_size := delta / 50

	// Chunk size must never decrease.
	// this is important! if you change it
	// recursion might break or explode
	if chunk_size_ms < safe_chunk_size {
		chunk_size_ms = safe_chunk_size
	}

	dm.SleepUntil(end, chunk_size_ms)
}

func (dm *Drawing) Sleep(ms uint64, relative ...bool) {
	if len(relative) > 0 && relative[0] == true {
		return
	}
	dm.SleepUntil(sdl.GetTicks64()+ms, 100)
}

func (dm *Drawing) Push() {
	dm_copy := *dm
	dm.stack = append(dm.stack, &dm_copy)
}

func (dm *Drawing) Pop() {
	// before := *dm
	stack := dm.stack
	x, _ := stack[len(stack)-1], stack[:len(stack)-1]
	*dm = *x

	// bookkeeping: settings have been changed and we need to notify SDL
	dm.renderer.SetScale(dm.scaleX, dm.scaleY)
}
