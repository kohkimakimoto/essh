package zssh

import (
	"github.com/yuin/gopher-lua"
)

func LoadFunctions(L *lua.LState) {
	L.SetGlobal("Host", L.NewFunction(coreHost))
}

func coreHost(L *lua.LState) int {
	name := L.CheckString(1)

	if L.GetTop() == 2 {
		tb := L.CheckTable(2)
		registerHost(L, name, tb)

		return 0
	}

	L.Push(L.NewFunction(func(L *lua.LState) int {
		tb := L.CheckTable(1)
		registerHost(L, name, tb)

		return 0
	}))

	return 1
}

func registerHost(L *lua.LState, name string, config *lua.LTable) {
	newConfig := L.NewTable()
	config.ForEach(func(k lua.LValue, v lua.LValue) {
		if k.String() != "hooks" && k.String() != "description" && k.String() != "hidden" {
			newConfig.RawSet(k, v)
		}
	})

	h := &Host{
		Name: name,
		Config: newConfig,
		Hooks: map[string]func() error{},
	}

	hooks := config.RawGetString("hooks")
	if hookTb, ok := toLTable(hooks); ok {
		hookBefore := hookTb.RawGetString("before")
		if hookBeforeFn, ok := toLFunction(hookBefore); ok {
			h.Hooks["before"] = func() error {
				err := L.CallByParam(lua.P{
					Fn:      hookBeforeFn,
					NRet:    0,
					Protect: true,
				})
				return err
			}
		}

		hookAfter := hookTb.RawGetString("after")
		if hookAfterFn, ok := toLFunction(hookAfter); ok {
			h.Hooks["after"] = func() error {
				err := L.CallByParam(lua.P{
					Fn:      hookAfterFn,
					NRet:    0,
					Protect: true,
				})
				return err
			}
		}
	}

	description := config.RawGetString("description")
	if descStr, ok := toString(description); ok {
		h.Description = descStr
	}

	hidden := config.RawGetString("hidden")
	if hiddenBool, ok := toBool(hidden); ok {
		h.Hidden = hiddenBool
	}

	Hosts = append(Hosts, h)
}

