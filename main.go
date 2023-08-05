package main

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	lua_state := lua.NewState()
	defer lua_state.Close()

	if err := lua_state.DoFile("script.lua"); err != nil {
		panic(err)
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	defer sdl.Quit()

	err := lua_state.CallByParam(lua.P{
		Fn:      lua_state.GetGlobal("Setup"),
		NRet:    0,
		Protect: true,
	})
	if nil != err {
		panic("Err")
	}

	window, err := sdl.CreateWindow(
		"test",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		800,
		600,
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE,
	)

	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(
		window,
		-1,
		sdl.RENDERER_ACCELERATED|
			sdl.RENDERER_PRESENTVSYNC,
	)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	drawing := CreateDrawing(renderer)

	loops(drawing, lua_state)

	os.Exit(0)
}
