/// lua utils

package main

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
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
func luaInvokeFunc(context string, script *lua.LState, funcCandidate lua.LValue, args ...lua.LValue) error {

	fn := luaFuncOrNil(funcCandidate)
	if nil == fn {
		return fmt.Errorf("%s : fn was not a function", context)
	}

	err := script.CallByParam(lua.P{
		Fn:      fn,
		NRet:    0,
		Protect: false,
		Handler: &lua.LFunction{},
	}, args...)

	if err != nil {
		panic(fmt.Sprintf("Lua error (%s) - %s", context, err))
	}

	return nil
}
