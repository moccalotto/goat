package main

import (
	"log"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type SystemSettings struct {
}

func (sys *SystemSettings) Sleep(ms uint32) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (sys *SystemSettings) Log(args ...any) {
	log.Printf("%+v", args)
}

func (sys *SystemSettings) ForceQuit(code int, why string) {
	log.Printf("ForceQuit() called from lua. Code %d, reason: %s\n", code, why)
	os.Exit(code)
}

func (sys *SystemSettings) Quit() {
	sdl.PushEvent(
		&sdl.QuitEvent{
			Type:      sdl.QUIT,
			Timestamp: uint32(sdl.GetTicks64()),
		},
	)
}
