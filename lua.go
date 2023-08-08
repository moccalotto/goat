/// lua utils

package main

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func luaFuncOrNil(fn lua.LValue) *lua.LFunction {
	if nil == fn {
		return nil
	}
	if lua.LNil == fn {
		return nil
	}
	if fn.Type() != lua.LTFunction {
		return nil
	}

	return fn.(*lua.LFunction)
}

// Call a function inside the lua script
func luaInvokeFunc(context string, script *lua.LState, fn lua.LValue, args ...lua.LValue) error {

	_func := luaFuncOrNil(fn)
	if nil == _func {
		return fmt.Errorf("%s : fn was not a function", context)
	}

	err := script.CallByParam(lua.P{
		Fn:      _func,
		NRet:    0,
		Protect: false,
		Handler: &lua.LFunction{},
	}, args...)

	if err != nil {
		panic(fmt.Sprintf("Lua error (%s) - %s", context, err))
	}

	return nil
}

func luaLoadScript(cfg *Config) *lua.LState {
	script := lua.NewState()

	if err := script.DoFile("script.lua"); err != nil {
		panic(err)
	}

	// Call the setup function.
	// In that function we set parameters needed to set the initial size and
	// title of the window, etc.
	setupFunc := script.GetGlobal("Setup")

	script.SetGlobal("cfg", luar.New(script, cfg))

	if setupFunc != lua.LNil {
		luaInvokeFunc("Setup()", script, setupFunc, luar.New(script, cfg))
	}

	if err := cfg.configSanityCheck(); err != nil {
		panic(err)
		// log.Fatalf("Invalid configuration: %v\n", err, err)
	}

	return script
}
