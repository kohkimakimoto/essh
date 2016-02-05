package zssh

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"unicode"
)

var (
	lzssh *lua.LTable
)

func InitLuaState(L *lua.LState) {
	L.SetGlobal("Host", L.NewFunction(coreHost))

	lzssh = L.NewTable()
	L.SetGlobal("zssh", lzssh)

	lzssh.RawSetString("ssh_config", lua.LNil)
}

func sshConfig(L *lua.LState) int {

	L.Push(lua.LString("aaaaaaa"))

	return 0
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
		Tags:   []string{},
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
		tagsTb.ForEach(func(_ lua.LValue, v lua.LValue) {
			if vs, ok := toString(v); ok {
				h.Tags = append(h.Tags, vs)
				//if hosts, ok := Tags[vs]; ok {
				//	hosts = append(hosts, h)
				//} else {
				//	hosts = []*Host{h}
				//	Tags[vs] = hosts
				//}
			} else {
				L.RaiseError("unsupported format of tags.")
			}
		})
	}

	Hosts = append(Hosts, h)
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
