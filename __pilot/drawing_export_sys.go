/*****************************************
 * Contains all exported
 * functions that are "system" related.
 *****************************************/

package pilot

import (
	h "goat/shed"
	"log"
	"os"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func (dm *Drawing) ProcessEvents(maxEventsToProcess ...int) bool {
	// Only process a few events at a time, unless otherwise specified.
	glfw.PollEvents()
	return !dm.window.ShouldClose()
}

func (dm *Drawing) ForceQuit(code ...int) {
	exitCode := 0
	if len(code) > 0 {
		exitCode = code[0]
	}
	h.RunOnMain(func() {
		os.Exit(exitCode)
	})
}

func (dm *Drawing) Quit() {
	// add quit event to glfw even queue if possible
}

// Sleep until the SDL tickcounter == end
// every chunk_size_ms we wake up and process events.
func (dm *Drawing) SleepUntil(end float64, chunkSize float64) {

	now := glfw.GetTime()

	delta := (end - now)

	if delta <= chunkSize {
		time.Sleep(time.Duration(delta * 1e9))
		return
	}

	dm.ProcessEvents()

	// we don't want to wake up for event processing more
	// than 50 times per sleep session
	safe_chunk_size := delta / 50

	// Chunk size must never decrease.
	// this is important! if you change it
	// recursion might break or explode
	if chunkSize < safe_chunk_size {
		chunkSize = safe_chunk_size
	}

	// a bit of tail recursion never hurt anyone.
	dm.SleepUntil(end, chunkSize)
}

func (dm *Drawing) Sleep(sec float64, relative ...bool) {

	startTime := dm.nowTime

	if len(relative) == 0 || !relative[0] {
		startTime = glfw.GetTime()
	}

	endTime := startTime + sec

	const defaultChunkSize = 0.010 // 10 msec

	dm.SleepUntil(endTime, defaultChunkSize)
}

// Push all settings onto a the stack.
func (dm *Drawing) Push() {
	dm_copy := *dm
	dm.stack = append(dm.stack, &dm_copy)
}

// Pop all settings from stack.
func (dm *Drawing) Pop() {
	// before := *dm
	stack := dm.stack
	x, _ := stack[len(stack)-1], stack[:len(stack)-1]
	*dm = *x
}

func (dm *Drawing) WinSize(x, y int, center ...bool) {
	if x > 0 && y > 0 {
		dm.window.SetSize(x, y)
	}

	if len(center) > 0 && center[0] {
		return // TODO something
	}
}

func (dm *Drawing) WinTitle(title string) {
	dm.window.SetTitle(title)
}

func (dm *Drawing) Pause() {
	log.Print("Pause()")
	dm.paused = true
}

func (dm *Drawing) Resume() {
	log.Print("Resume()")
	dm.paused = false
}

func (dm *Drawing) FrameRateCap(val float64) {
	dm.frameRateCap = val
}

func (dm *Drawing) Dump(x ...interface{}) {
	for i, v := range x {
		log.Printf("Dump%3d: %+v", i, v)
	}
}

func (dm *Drawing) HasKey(name string) bool {
	return false
}
