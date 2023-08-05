package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func loops(drawing *Drawing, lua_state *lua.LState) {
	dm_param := luar.New(lua_state, drawing)
	drawFunc := lua_state.GetGlobal("Draw")
	updateFunc := lua_state.GetGlobal("Update")
	renderer := drawing.renderer

	for handlingEvents() {
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		renderer.SetDrawColor(255, 0, 0, 255)

		if err := lua_state.CallByParam(lua.P{
			Fn:      updateFunc,
			NRet:    0,
			Protect: true,
		}, lua.LString("Params for update")); err != nil {
			//			panic(err)
		}
		err := lua_state.CallByParam(lua.P{
			Fn:      drawFunc,
			NRet:    0,
			Protect: true,
		}, luar.New(lua_state, dm_param))

		if err != nil {
			fmt.Printf("Error: %v", err)
		}

		renderer.Present()
	}
}

func handlingEvents() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			return false
		}
	}
}
