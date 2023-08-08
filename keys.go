package main

import "github.com/veandco/go-sdl2/sdl"

// Keyboard event type
// This object is sent to lua whenever
// a key is pressed, released, or repeated
// TODO: add locale aware characters
type KeyEvent struct {
	Code sdl.Keycode

	Str  string
	Name string

	Ctrl   bool
	L_Ctrl bool
	R_Ctrl bool

	Alt   bool
	L_Alt bool
	R_Alt bool
	AltGr bool // same as ralt

	Shift   bool
	R_Shift bool
	L_Shift bool

	Gui   bool
	L_Gui bool
	R_Gui bool

	Pressed  bool
	Released bool
	Repeated bool

	Up    bool
	Down  bool
	Left  bool
	Right bool

	Escape bool
}

// Create and initialize a KeyEvent Struct from an SDL keyboard event.
func CreateKeyEvent(k *sdl.KeyboardEvent) *KeyEvent {
	ke := KeyEvent{}
	sym := k.Keysym.Sym
	mod := k.Keysym.Mod

	ke.Name = sdl.GetScancodeName(k.Keysym.Scancode)
	ke.Str = ke.Name
	ke.Code = sym

	switch k.State {
	case sdl.RELEASED:
		ke.Released = true
	case sdl.PRESSED:
		ke.Pressed = true
	default:
		if k.Repeat > 0 {
			ke.Repeated = true
		} else {
			panic("Keys shoulde be pressed, released or repeated")
		}
	}

	if k.Keysym.Sym == sdl.K_ESCAPE {
		ke.Escape = true
	}

	switch sym {
	case sdl.K_DOWN:
		ke.Down = true
	case sdl.K_UP:
		ke.Up = true
	case sdl.K_LEFT:
		ke.Left = true
	case sdl.K_RIGHT:
		ke.Right = true
	}

	if mod&sdl.KMOD_LALT != 0 {
		ke.L_Alt = true
		ke.Alt = true
	}
	if mod&sdl.KMOD_RALT != 0 {
		ke.R_Alt = true
		ke.Alt = true
		ke.AltGr = true
	}

	if mod&sdl.KMOD_LCTRL != 0 {
		ke.L_Ctrl = true
		ke.Ctrl = true
	}
	if mod&sdl.KMOD_RCTRL != 0 {
		ke.R_Ctrl = true
		ke.Ctrl = true
	}

	if mod&sdl.KMOD_LSHIFT != 0 {
		ke.L_Shift = true
		ke.Shift = true
	}
	if mod&sdl.KMOD_RSHIFT != 0 {
		ke.R_Shift = true
		ke.Shift = true
	}

	if mod&sdl.KMOD_LGUI != 0 {
		ke.L_Gui = true
		ke.Gui = true
	}
	if mod&sdl.KMOD_RGUI != 0 {
		ke.R_Gui = true
		ke.Gui = true
	}

	return &ke
}
