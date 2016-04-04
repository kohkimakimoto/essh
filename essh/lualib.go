package essh

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"unicode"
)

var (
	lessh *lua.LTable
)

func InitLuaState(L *lua.LState) {
	L.SetGlobal("Host", L.NewFunction(coreHost))
	L.SetGlobal("Task", L.NewFunction(coreTask))

	lessh = L.NewTable()
	L.SetGlobal("essh", lessh)

	lessh.RawSetString("ssh_config", lua.LNil)
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

func coreTask(L *lua.LState) int {
	name := L.CheckString(1)

	if L.GetTop() == 2 {
		tb := L.CheckTable(2)
		registerTask(L, name, tb)

		return 0
	}

	L.Push(L.NewFunction(func(L *lua.LState) int {
		tb := L.CheckTable(1)
		registerTask(L, name, tb)

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
		Hooks:  map[string]interface{}{},
		Tags:   []string{},
	}

	hooks := config.RawGetString("hooks")
	if hookTb, ok := toLTable(hooks); ok {
		// before depricated. use before_connect
		err := registerHook(L, h, "before", hookTb.RawGetString("before"))
		if err != nil {
			panic(err)
		}
		err = registerHook(L, h, "before_connect", hookTb.RawGetString("before_connect"))
		if err != nil {
			panic(err)
		}

		err = registerRemoteHook(L, h, "after_connect", hookTb.RawGetString("after_connect"))
		if err != nil {
			panic(err)
		}

		// after depricated. use after_disconnect
		err = registerHook(L, h, "after", hookTb.RawGetString("after"))
		if err != nil {
			panic(err)
		}

		err = registerHook(L, h, "after_disconnect", hookTb.RawGetString("after_disconnect"))
		if err != nil {
			panic(err)
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
			} else {
				L.RaiseError("unsupported format of tags.")
			}
		})
	}

	Hosts = append(Hosts, h)
}

func registerHook(L *lua.LState, host *Host, hookPoint string, hook lua.LValue) error {
	if hook != lua.LNil {
		if hookFn, ok := toLFunction(hook); ok {
			host.Hooks[hookPoint] = func() error {
				err := L.CallByParam(lua.P{
					Fn:      hookFn,
					NRet:    0,
					Protect: true,
				})
				return err
			}
		} else if hookString, ok := toString(hook); ok {
			host.Hooks[hookPoint] = hookString
		} else {
			return fmt.Errorf("invalid hook type %v", hook)
		}
	}
	return nil
}

func registerRemoteHook(L *lua.LState, host *Host, hookPoint string, hook lua.LValue) error {
	if hook != lua.LNil {
		if hookString, ok := toString(hook); ok {
			host.Hooks[hookPoint] = hookString
		} else {
			return fmt.Errorf("invalid hook type %v", hook)
		}
	}

	return nil
}

func registerTask(L *lua.LState, name string, config *lua.LTable) {
	task := &Task{
		Name:   name,
	}

	description := config.RawGetString("description")
	if descStr, ok := toString(description); ok {
		task.Description = descStr
	}

	Tasks = append(Tasks, task)
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
