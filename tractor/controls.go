package tractor

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	Controls *ControlsType = &ControlsType{}
)

// =========================================================================
// ||
// || Primarily used to interface with glfw
// ||
// =========================================================================
type ControlsType struct {
	E *EngineType
}

func (C *ControlsType) lazyInit() {
	if C.E == nil {
		C.E = Engine
	}
}

// =========================================================================
// || Is the given key pressed?
// ||
// || If engine is not set, the main engine (defined in a global variable)
// || will be used.
// =========================================================================
func (C *ControlsType) KeyPressed(key KeyCode, engine ...*EngineType) bool {
	C.lazyInit()

	k := glfw.Key(key)
	e := *C.E
	if len(engine) > 0 {
		e = *engine[0]
	}

	return e.Window.GetKey(k) != Release
}

func (C *ControlsType) HandleKeys(kh KeyboardHandler) {
	C.lazyInit()

	//
	glfw.GetCurrentContext().SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

		kev := KeyEvent{
			Key:      KeyCode(key),
			ScanCode: scancode,
			Pressed:  action == glfw.Press,
			Repeated: action == glfw.Repeat,
			Released: action == glfw.Release,
			Ctrl:     mods&glfw.ModControl != 0,
			Shift:    mods&glfw.ModShift != 0,
			Alt:      mods&glfw.ModAlt != 0,
			Gui:      mods&glfw.ModSuper != 0,
			L_Ctrl:   key == glfw.KeyLeftControl,
			R_Ctrl:   key == glfw.KeyRightControl,
			L_Alt:    key == glfw.KeyLeftAlt,
			R_Alt:    key == glfw.KeyRightAlt,
			AltGr:    key == glfw.KeyRightAlt,
			L_Shift:  key == glfw.KeyLeftShift,
			R_Shift:  key == glfw.KeyRightShift,

			Up:    key == glfw.KeyDown,
			Down:  key == glfw.KeyDown,
			Left:  key == glfw.KeyLeft,
			Right: key == glfw.KeyRight,

			Escape: key == glfw.KeyEscape,
		}

		// ModCapsLock ModifierKey = C.GLFW_MOD_CAPS_LOCK
		// ModNumLock  ModifierKey = C.GLFW_MOD_NUM_LOCK

		kh(&kev)
	})
}
