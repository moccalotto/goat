package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func loop(dm *Drawing) {
	renderer := dm.renderer
	dm.injectFunctions()

	renderer.SetDrawColor(255, 0, 0, 255)
	for handlingEvents() {

		dm.draw()

		renderer.Present()
	}
}

func handlingEvents() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			return false
		case *sdl.KeyboardEvent:
			keyCode := t.Keysym.Sym
			keys := ""

			// Modifier keys
			switch t.Keysym.Mod {
			case sdl.KMOD_LALT:
				keys += "Left Alt"
			case sdl.KMOD_LCTRL:
				keys += "Left Control"
			case sdl.KMOD_LSHIFT:
				keys += "Left Shift"
			case sdl.KMOD_LGUI:
				keys += "Left Meta or Windows key"
			case sdl.KMOD_RALT:
				keys += "Right Alt"
			case sdl.KMOD_RCTRL:
				keys += "Right Control"
			case sdl.KMOD_RSHIFT:
				keys += "Right Shift"
			case sdl.KMOD_RGUI:
				keys += "Right Meta or Windows key"
			case sdl.KMOD_NUM:
				keys += "Num Lock"
			case sdl.KMOD_CAPS:
				keys += "Caps Lock"
			case sdl.KMOD_MODE:
				keys += "AltGr Key"
			}

			if keyCode < 10000 {
				if keys != "" {
					keys += " + "
				}

				// If the key is held down, this will fire
				if t.Repeat > 0 {
					keys += string(keyCode) + " repeating"
				} else {
					if t.State == sdl.RELEASED {
						keys += string(keyCode) + " released"
					} else if t.State == sdl.PRESSED {
						keys += string(keyCode) + " pressed"
					}
				}

			}

			if keys != "" {
				fmt.Println(keys)
			}
		}
	}

	return true
}
