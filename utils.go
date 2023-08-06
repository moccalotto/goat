package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

type Config struct {
	Width     int32
	Height    int32
	Title     string
	CanResize bool
}

func (cfg *Config) CanResizeFlag() uint32 {
	if !cfg.CanResize {
		return 0
	}
	return sdl.WINDOW_RESIZABLE
}

// Call a function inside the lua script
func invokeLuaFunc(context string, lua_state *lua.LState, function lua.LValue, args ...lua.LValue) {

	err := lua_state.CallByParam(lua.P{
		Fn:      function,
		NRet:    0,
		Protect: true,
		Handler: &lua.LFunction{},
	}, args...)

	if err != nil {
		log.Fatalf("Error in %s: %+v", context, err)
	}
}

func setupLua(cfg *Config) *lua.LState {
	script := lua.NewState()

	if err := script.DoFile("script.lua"); err != nil {
		panic(err)
	}

	// Call the setup function.
	// In that function we set parameters needed to set the initial size and
	// title of the window, etc.
	setupFunc := script.GetGlobal("Setup")

	script.SetGlobal("Width", luar.New(script, func(val int32) {
		cfg.Width = val
	}))
	script.SetGlobal("Height", luar.New(script, func(val int32) {
		cfg.Height = val
	}))
	script.SetGlobal("Title", luar.New(script, func(val string) {
		cfg.Title = val
	}))
	script.SetGlobal("CanResize", luar.New(script, func(val ...bool) {
		if len(val) > 0 {
			cfg.CanResize = val[0]
			return
		}

		cfg.CanResize = true
	}))

	if setupFunc.Type() != lua.LTNil {
		invokeLuaFunc("Setup()", script, setupFunc)
	}

	return script
}

func setupSDL(cfg *Config) (*sdl.Window, *sdl.Renderer) {
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		cfg.Title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		cfg.Width,
		cfg.Height,
		sdl.WINDOW_SHOWN|cfg.CanResizeFlag(),
	)

	if err != nil {
		panic(err)
	}
	renderer, err := sdl.CreateRenderer(
		window,
		-1,
		sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC,
	)
	if err != nil {
		log.Fatalf("Could not create SDL renderer: %+v", err)
	}
	return window, renderer
}
