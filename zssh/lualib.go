package zssh

import (
	"errors"
	"fmt"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
	"unicode"
)

func LoadFunctions(L *lua.LState) {
	L.SetGlobal("Host", L.NewFunction(coreHost))
	L.SetGlobal("Macro", L.NewFunction(coreMacro))
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

func coreMacro(L *lua.LState) int {
	name := L.CheckString(1)

	if L.GetTop() == 2 {
		tb := L.CheckTable(2)
		registerMacro(L, name, tb)

		return 0
	}

	L.Push(L.NewFunction(func(L *lua.LState) int {
		tb := L.CheckTable(1)
		registerMacro(L, name, tb)

		return 0
	}))

	return 1
}

func registerHost(L *lua.LState, name string, config *lua.LTable) {
	newConfig := L.NewTable()
	config.ForEach(func(k lua.LValue, v lua.LValue) {
		var firstChar rune
		for _, c := range k.String() {
			firstChar = c
			break
		}

		if unicode.IsUpper(firstChar) {
			newConfig.RawSet(k, v)
		}
	})

	h := &Host{
		Name:   name,
		Config: newConfig,
		Hooks:  map[string]func() error{},
		Tags:   map[string][]string{},
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

	tags := config.RawGetString("tags")
	if tagsTb, ok := tags.(*lua.LTable); ok {
		tagsTb.ForEach(func(k lua.LValue, v lua.LValue) {
			if ks, ok := toString(k); ok {
				if vs, ok := toString(v); ok {
					h.Tags[ks] = []string{vs}
				} else if vs, ok := toLTable(v); ok {
					var values = []string{}
					vs.ForEach(func(_ lua.LValue, vv lua.LValue) {
						if vvs, ok := toString(vv); ok {
							values = append(values, vvs)
						} else {
							L.RaiseError("unsupported format of tags.")
						}
					})
					h.Tags[ks] = values
				} else {
					L.RaiseError("unsupported format of tags.")
				}
			}
		})
	}

	Hosts = append(Hosts, h)
}

func registerMacro(L *lua.LState, name string, config *lua.LTable) {
	m := &Macro{
		OnTags:    map[string][]string{},
		OnServers: []string{},
	}

	if err := gluamapper.Map(config, m); err != nil {
		L.RaiseError("got a error when it is parsing the macro: %s", err)
	}

	m.Name = name

	on := config.RawGetString("on")
	if ontb, ok := on.(*lua.LTable); ok {
		m.RunLocally = false
		ontb.ForEach(func(k lua.LValue, v lua.LValue) {
			if ks, ok := toString(k); ok {
				// tags
				if vs, ok := toString(v); ok {
					m.OnTags[ks] = []string{vs}
				} else if vs, ok := toLTable(v); ok {
					var values = []string{}
					vs.ForEach(func(_ lua.LValue, vv lua.LValue) {
						if vvs, ok := toString(vv); ok {
							values = append(values, vvs)
						} else {
							L.RaiseError("unsupported format of tags.")
						}
					})
					m.OnTags[ks] = values
				} else {
					L.RaiseError("unsupported format of tags.")
				}
			} else {
				// servers
				if vs, ok := toString(v); ok {
					m.OnServers = append(m.OnServers, vs)
				}
			}
		})
	} else {
		m.RunLocally = true
	}

	confirm := config.RawGetString("confirm")
	switch confirmConverted := confirm.(type) {
	case lua.LBool:
		m.Confirm = bool(confirmConverted)
	case lua.LString:
		m.Confirm = true
		m.ConfirmText = string(confirmConverted)
	}

	command := config.RawGetString("command")
	if commandFn, ok := command.(*lua.LFunction); ok {
		m.CommandFunc = func(host *Host) (string, error) {
			lhost := L.NewUserData()
			lhost.Value = host
			L.SetMetatable(lhost, L.GetTypeMetatable(LHostClass))

			err := L.CallByParam(lua.P{
				Fn:      commandFn,
				NRet:    1,
				Protect: true,
			}, lhost)
			if err != nil {
				return "", err
			}

			ret := L.Get(-1) // returned value
			L.Pop(1)

			if ret == lua.LNil {
				return "", nil
			} else if retStr, ok := ret.(lua.LString); ok {
				return string(retStr), nil
			} else {
				return "", errors.New("return value must be string")
			}
		}
	} else if commandStr, ok := command.(lua.LString); ok {
		m.Command = string(commandStr)
	}

	Macros = append(Macros, m)
}

// This code refers to https://github.com/yuin/gluamapper/blob/master/gluamapper.go
func toGoValue(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[interface{}]interface{})
			v.ForEach(func(key, value lua.LValue) {
				keystr := fmt.Sprint(toGoValue(key))
				ret[keystr] = toGoValue(value)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, toGoValue(v.RawGetInt(i)))
			}
			return ret
		}
	default:
		return v
	}
}

func toBool(v lua.LValue) (bool, bool) {
	if lv, ok := v.(lua.LBool); ok {
		return bool(lv), true
	} else {
		return false, false
	}
}

func toString(v lua.LValue) (string, bool) {
	if lv, ok := v.(lua.LString); ok {
		return string(lv), true
	} else {
		return "", false
	}
}

func toLFunction(v lua.LValue) (*lua.LFunction, bool) {
	if lv, ok := v.(*lua.LFunction); ok {
		return lv, true
	} else {
		return nil, false
	}
}

func toLTable(v lua.LValue) (*lua.LTable, bool) {
	if lv, ok := v.(*lua.LTable); ok {
		return lv, true
	} else {
		return nil, false
	}
}
