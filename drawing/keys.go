package drawing

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Keyboard event type
// This object is sent to lua whenever
// a key is pressed, released, or repeated
// TODO: add locale aware characters
type KeyEvent struct {
	Code glfw.Key

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
