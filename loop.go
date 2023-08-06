package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func loop(drawing *Drawing, sys *SystemSettings, luaState *lua.LState) {
	drawParam := luar.New(luaState, drawing)
	sysParam := luar.New(luaState, sys)
	drawFunc := luaState.GetGlobal("Draw")
	renderer := drawing.renderer

	luaState.SetGlobal("Diller", luar.New(luaState, func() {
		print("diller")
	}))

	for handlingEvents() {
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		renderer.SetDrawColor(255, 0, 0, 255)

		drawing.updateLocalVariables()

		callLuaFunc(luaState, drawFunc, drawParam, sysParam)

		renderer.Present()
	}
}

func callLuaFunc(lua_state *lua.LState, function lua.LValue, args ...lua.LValue) {
	err := lua_state.CallByParam(lua.P{
		Fn:      function,
		NRet:    0,
		Protect: true,
		Handler: &lua.LFunction{},
	}, args...)

	if err != nil {
		panic(err)
	}
}

func handlingEvents() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			return false
		case *sdl.KeyboardEvent:
			fmt.Printf("%+v\n", t)
		}
	}

	return true
}
