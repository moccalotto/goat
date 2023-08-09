package main

import (
	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
)

// My dream is to do all of this:: https://github.com/fogleman/gg

type Drawing struct {
	renderer     *sdl.Renderer // The sdl renderer we're using. It's most likely a hardware renderer.
	script       *lua.LState   // The lua script we're running
	drawFunc     lua.LValue    // The lua function to call to draw
	prevTicks    uint64        // ticks when previous draw started
	nowTicks     uint64        // ticks when current draw started
	deltaTicks   uint64        // number of ticks the last render cycle took
	fgColor      sdl.Color     // Foreground color. The color of the "ink" if you will
	bgColor      sdl.Color     // Background color. The color of the "paper"
	scaleX       float32       // How big (in pixels) are the virtual pixels in the X direction.
	scaleY       float32       // How big (in pixels) are the virtual pixels in the Y direction.
	frameRateCap float32       // Maximum allowed number of updates per second. - During this delay no events are processed.
	stack        []*Drawing    // The stack that allows us to store and recall colors, scales, and other such settings.
	frameCount   uint64        // The number of calls to draw(). Starts at 1

	keydownCallback *lua.LFunction
	keyupCallback   *lua.LFunction
}

func CreateDrawing(renderer *sdl.Renderer, script *lua.LState) *Drawing {

	ticks := sdl.GetTicks64()

	dm := &Drawing{
		script:       script,
		renderer:     renderer,
		drawFunc:     script.GetGlobal("Draw"),
		nowTicks:     ticks,
		prevTicks:    ticks,
		deltaTicks:   0,
		fgColor:      sdl.Color{R: 0, G: 0, B: 0, A: 255},
		bgColor:      sdl.Color{R: 255, G: 255, B: 255, A: 255},
		scaleX:       1,
		scaleY:       1,
		frameRateCap: 0,
		stack:        make([]*Drawing, 0),
	}

	dm.scaleX, dm.scaleY = renderer.GetScale()
	dm.keydownCallback = luaFuncOrNil(script.GetGlobal("Keydown"))
	dm.keyupCallback = luaFuncOrNil(script.GetGlobal("Keyup"))

	return dm
}

// Call the setup() function before the first frame is framed
func (dm *Drawing) setup() {
	if dm.frameCount > 0 {
		return
	}

	setupFunc := dm.script.GetGlobal("Setup")
	luaInvokeFunc("Setup()", dm.script, setupFunc)
}

func (dm *Drawing) Color(c ...uint8) (uint8, uint8, uint8, uint8) {
	switch len(c) {
	case 0:
		return dm.fgColor.R, dm.fgColor.G, dm.fgColor.B, dm.fgColor.A
	case 1:
		return dm.Color(c[0], c[0], c[0], 255)
	case 3:
		return dm.Color(c[0], c[1], c[2], 255)
	case 4:
		dm.fgColor = sdl.Color{R: c[0], G: c[1], B: c[2], A: c[3]}
		return dm.Color()
	default:
		panic("Background() takes 1, 3, or 4 arguments.")
	}
}
func (dm *Drawing) Background(c ...uint8) {
	switch len(c) {
	case 1:
		dm.Background(c[0], c[0], c[0], 255)
	case 3:
		dm.Background(c[0], c[1], c[2], 255)
	case 4:
		dm.bgColor = sdl.Color{R: c[0], G: c[1], B: c[2], A: c[3]}
	default:
		panic("Background() takes 1, 3, or 4 arguments.")
	}
}

func (dm *Drawing) Scale(scale float32) {
	dm.scaleX = scale
	dm.scaleY = scale
	dm.renderer.SetScale(scale, scale)
}

func (dm *Drawing) Line(x1, y1, x2, y2 float32) {
	dm.applySettingsToRenderer()
	dm.renderer.DrawLineF(x1, y1, x2, y2)
}

func (dm *Drawing) Dot(x, y float32) {
	dm.applySettingsToRenderer()
	dm.renderer.DrawPointF(x, y)
}

func (dm *Drawing) Rectangle(x1, y1, x2, y2 float32) {
	dm.applySettingsToRenderer()

	vertices := []sdl.Vertex{
		{sdl.FPoint{x1, y1}, dm.fgColor, sdl.FPoint{0, 0}},
		{sdl.FPoint{x2, y1}, dm.fgColor, sdl.FPoint{0, 0}},
		{sdl.FPoint{x2, y2}, dm.fgColor, sdl.FPoint{0, 0}},

		{sdl.FPoint{x2, y2}, dm.fgColor, sdl.FPoint{0, 0}},
		{sdl.FPoint{x1, y2}, dm.fgColor, sdl.FPoint{0, 0}},
		{sdl.FPoint{x1, y1}, dm.fgColor, sdl.FPoint{0, 0}},
	}
	dm.renderer.RenderGeometry(nil, vertices, nil)

	// dm.renderer.DrawRectF(&sdl.FRect{
	// 	X: x1,
	// 	Y: y1,
	// 	H: y2 - y1,
	// 	W: x2 - x1,
	// })
}
