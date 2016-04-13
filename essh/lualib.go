package essh

import (
	"fmt"
	"github.com/cjoudrey/gluahttp"
	"github.com/kohkimakimoto/gluafs"
	"github.com/kohkimakimoto/gluajson"
	"github.com/kohkimakimoto/gluaquestion"
	"github.com/kohkimakimoto/gluatemplate"
	"github.com/kohkimakimoto/gluayaml"
	"github.com/yuin/gopher-lua"
	"net/http"
	"os"
	"unicode"
)

func InitLuaState(L *lua.LState) {
	// custom type.
	// registerContextClass(L)
	registerTaskContextClass(L)

	// global functions
	L.SetGlobal("host", L.NewFunction(esshHost))
	L.SetGlobal("task", L.NewFunction(esshTask))
	// backend compatibility
	L.SetGlobal("Host", L.NewFunction(esshHost))
	L.SetGlobal("Task", L.NewFunction(esshTask))

	// modules
	L.PreloadModule("essh.json", gluajson.Loader)
	L.PreloadModule("essh.fs", gluafs.Loader)
	L.PreloadModule("essh.yaml", gluayaml.Loader)
	L.PreloadModule("essh.template", gluatemplate.Loader)
	L.PreloadModule("essh.question", gluaquestion.Loader)
	L.PreloadModule("essh.http", gluahttp.NewHttpModule(&http.Client{}).Loader)

	// global variables
	lessh := L.NewTable()
	L.SetGlobal("essh", lessh)
	lessh.RawSetString("ssh_config", lua.LNil)

	L.SetFuncs(lessh, map[string]lua.LGFunction{
		"host":    esshHost,
		"task":    esshTask,
		"require": esshRequire,
//		"reset":   esshReset,
	})
}

func esshHost(L *lua.LState) int {
	first := L.CheckAny(1)
	if tb, ok := first.(*lua.LTable); ok {
		registerHostByTable(L, tb)
		return 0
	}

	name := L.CheckString(1)

	// procedural style
	if L.GetTop() == 2 {
		tb := L.CheckTable(2)
		registerHost(L, name, tb)

		return 0
	}

	// DSL style
	L.Push(L.NewFunction(func(L *lua.LState) int {
		tb := L.CheckTable(1)
		registerHost(L, name, tb)

		return 0
	}))

	return 1
}

func registerHostByTable(L *lua.LState, tb *lua.LTable) {
	maxn := tb.MaxN()
	if maxn == 0 { // table
		tb.ForEach(func(key, value lua.LValue) {
			config, ok := value.(*lua.LTable)
			if !ok {
				return
			}
			name, ok := key.(lua.LString)
			if !ok {
				return
			}

			registerHost(L, string(name), config)
		})
	} else { // array
		for i := 1; i <= maxn; i++ {
			value := tb.RawGetInt(i)
			valueTb, ok := value.(*lua.LTable)
			if !ok {
				return
			}
			registerHostByTable(L, valueTb)
		}
	}
}

func esshTask(L *lua.LState) int {
	first := L.CheckAny(1)
	if tb, ok := first.(*lua.LTable); ok {
		registerTaskByTable(L, tb)
		return 0
	}

	name := L.CheckString(1)

	// procedural style
	if L.GetTop() == 2 {
		tb := L.CheckTable(2)
		registerTask(L, name, tb)

		return 0
	}

	// DSL style
	L.Push(L.NewFunction(func(L *lua.LState) int {
		tb := L.CheckTable(1)
		registerTask(L, name, tb)

		return 0
	}))

	return 1
}

func registerTaskByTable(L *lua.LState, tb *lua.LTable) {
	maxn := tb.MaxN()
	if maxn == 0 { // table
		tb.ForEach(func(key, value lua.LValue) {
			config, ok := value.(*lua.LTable)
			if !ok {
				return
			}
			name, ok := key.(lua.LString)
			if !ok {
				return
			}

			registerTask(L, string(name), config)
		})
	} else { // array
		for i := 1; i <= maxn; i++ {
			value := tb.RawGetInt(i)
			valueTb, ok := value.(*lua.LTable)
			if !ok {
				return
			}
			registerTaskByTable(L, valueTb)
		}
	}
}


func esshReset(L *lua.LState) int {
	if debugFlag {
		fmt.Println("[essh debug] reset tasks, hosts and modules.")
	}
	Tasks = []*Task{}
	Hosts = []*Host{}
	CurrentContext.LoadedModules = map[string]*Module{}

	return 0
}

func registerHost(L *lua.LState, name string, config *lua.LTable) {
	if debugFlag {
		fmt.Printf("[essh debug] register host: %s\n", name)
	}

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
		Props: map[string]string{},
		Hooks:  map[string]interface{}{},
		Tags:   []string{},
	}

	props := config.RawGetString("props")
	if propsTb, ok := toLTable(props); ok {
		propsTb.ForEach(func(propsKey lua.LValue, propsValue lua.LValue) {
			propsKeyStr, ok := toString(propsKey)
			if !ok {
				L.RaiseError("props table's key must be a string: %v", propsKey)
			}
			propsValueStr, ok := toString(propsValue)
			if !ok {
				L.RaiseError("props table's value must be a string: %v", propsValue)
			}

			h.Props[propsKeyStr] = propsValueStr
		})
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
	task := NewTask()
	task.Name = name

	description := config.RawGetString("description")
	if descStr, ok := toString(description); ok {
		task.Description = descStr
	}

	pty := config.RawGetString("pty")
	if ptyBool, ok := toBool(pty); ok {
		task.Pty = ptyBool
	}

	parallel := config.RawGetString("parallel")
	if parallelBool, ok := toBool(parallel); ok {
		task.Parallel = parallelBool
	}

	privileged := config.RawGetString("privileged")
	if privilegedBool, ok := toBool(privileged); ok {
		task.Privileged = privilegedBool
	}

	script := config.RawGetString("script")
	if scriptStr, ok := toString(script); ok {
		task.Script = scriptStr
	}

	file := config.RawGetString("file")
	if fileStr, ok := toString(file); ok {
		task.File = fileStr
	}

	if task.File != "" && task.Script != "" {
		L.RaiseError("invalid task definition: can't use 'file' and 'script' at the same time.")
	}

	on := config.RawGetString("on")
	if onStr, ok := toString(on); ok {
		task.On = []string{onStr}
	} else if onSlice, ok := toSlice(on); ok {
		for _, target := range onSlice {
			if targetStr, ok := target.(string); ok {
				task.On = append(task.On, targetStr)
			}
		}
	}

	foreach := config.RawGetString("foreach")
	if foreachStr, ok := toString(foreach); ok {
		task.Foreach = []string{foreachStr}
	} else if foreachSlice, ok := toSlice(foreach); ok {
		for _, target := range foreachSlice {
			if targetStr, ok := target.(string); ok {
				task.Foreach = append(task.Foreach, targetStr)
			}
		}
	}

	if len(task.Foreach) >= 1 && len(task.On) >= 1 {
		L.RaiseError("invalid task definition: can't use 'foreach' and 'on' at the same time.")
	}

	prefix := config.RawGetString("prefix")
	if prefixBool, ok := toBool(prefix); ok {
		if prefixBool {
			if task.IsRemoteTask() {
				task.Prefix = DefaultPrefixRemote
			} else {
				task.Prefix = DefaultPrefixLocal
			}
		}
	} else if prefixStr, ok := toString(prefix); ok {
		task.Prefix = prefixStr
	}

	configure := config.RawGetString("configure")
	if configure != lua.LNil {
		if configureFn, ok := configure.(*lua.LFunction); ok {
			task.Configure = func() error {
				err := L.CallByParam(lua.P{
					Fn:      configureFn,
					NRet:    0,
					Protect: true,
				})
				if err != nil {
					return err
				}

				return nil
			}
		} else {
			L.RaiseError("configure have to be function.")
		}
	}

	prepare := config.RawGetString("prepare")
	if prepare != lua.LNil {
		if prepareFn, ok := prepare.(*lua.LFunction); ok {
			task.Prepare = func(ctx *TaskContext) error {
				lctx := newLTaskContext(L, ctx)
				err := L.CallByParam(lua.P{
					Fn:      prepareFn,
					NRet:    1,
					Protect: true,
				}, lctx)
				if err != nil {
					return err
				}

				ret := L.Get(-1) // returned value
				L.Pop(1)

				if ret == lua.LNil {
					return nil
				} else if retB, ok := ret.(lua.LBool); ok {
					if retB {
						return nil
					} else {
						return fmt.Errorf("returned false from the prepare function.")
					}
				}

				return nil
			}
		} else {
			L.RaiseError("prepare have to be function.")
		}
	}

	Tasks = append(Tasks, task)
}

func esshRequire(L *lua.LState) int {
	name := L.CheckString(1)

	module := CurrentContext.LoadedModules[name]
	if module == nil {
		module = NewModule(name)
		err := module.Load(updateFlag)
		if err != nil {
			L.RaiseError("%v", err)
		}

		indexFile := module.IndexFile()
		if _, err := os.Stat(indexFile); err != nil {
			L.RaiseError("Could not load essh module: %v", err)
		}
		if err := L.DoFile(indexFile); err != nil {
			L.RaiseError("%v", err)
		}

		// get a module return value
		ret := L.Get(-1)
		module.Value = ret

		// register loaded module.
		CurrentContext.LoadedModules[name] = module
	}

	L.Push(module.Value)

	return 1
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
			ret := make(map[string]interface{})
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

func toMap(v lua.LValue) (map[string]interface{}, bool) {
	if lv, ok := toGoValue(v).(map[string]interface{}); ok {
		return lv, true
	} else {
		return nil, false
	}
}

func toSlice(v lua.LValue) ([]interface{}, bool) {
	if lv, ok := toGoValue(v).([]interface{}); ok {
		return lv, true
	} else {
		return nil, false
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

const LTaskContextClass = "TaskContext*"

func newLTaskContext(L *lua.LState, ctx *TaskContext) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = ctx
	L.SetMetatable(ud, L.GetTypeMetatable(LTaskContextClass))
	return ud
}

func registerTaskContextClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LTaskContextClass)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), taskContextMethods))
}

var taskContextMethods = map[string]lua.LGFunction{
	"payload": taskContextPayload,
}

func taskContextPayload(L *lua.LState) int {
	ctx := checkTaskContext(L)
	if L.GetTop() == 2 {
		ctx.Payload = L.CheckString(2)
		return 0
	}
	L.Push(lua.LString(ctx.Payload))
	return 1
}

func checkTaskContext(L *lua.LState) *TaskContext {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*TaskContext); ok {
		return v
	}
	L.ArgError(1, "TaskContext expected")
	return nil
}

//const LContextClass = "Context*"
//
//func newLContext(L *lua.LState, ctx *Context) *lua.LUserData {
//	ud := L.NewUserData()
//	ud.Value = ctx
//	L.SetMetatable(ud, L.GetTypeMetatable(LContextClass))
//	return ud
//}
//
//func registerContextClass(L *lua.LState) {
//	mt := L.NewTypeMetatable(LContextClass)
//	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), contextMethods))
//}
//
//var contextMethods = map[string]lua.LGFunction{
//	"datadir": contextDatadir,
//}
//
//func contextDatadir(L *lua.LState) int {
//
//	return 1
//}
//
//func checkContext(L *lua.LState) *TaskContext {
//	ud := L.CheckUserData(1)
//	if v, ok := ud.Value.(*Context); ok {
//		return v
//	}
//	L.ArgError(1, "Context expected")
//	return nil
//}
