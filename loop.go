package main

import (
	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
)

func loops(window *sdl.Window, renderer *sdl.Renderer, lua_state *lua.LState) {
	drawFunc := lua_state.GetGlobal("Draw")
	updateFunc := lua_state.GetGlobal("Setup")
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}

		if err := lua_state.CallByParam(lua.P{
			Fn:      drawFunc,
			NRet:    0,
			Protect: true,
		}, lua.LString("Draw()")); err != nil {
			panic(err)
		}
		if err := lua_state.CallByParam(lua.P{
			Fn:      updateFunc,
			NRet:    0,
			Protect: true,
		}, lua.LString("Update()")); err != nil {
			panic(err)
		}

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.DrawLine(0, 0, renderer.GetViewport().W, renderer.GetViewport().H)
		renderer.Present()
	}
}
