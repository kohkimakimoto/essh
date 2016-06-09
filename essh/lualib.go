package essh

import (
	"fmt"
	"github.com/cjoudrey/gluahttp"
	"github.com/kohkimakimoto/essh/gluamapper"
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
	L.SetGlobal("driver", L.NewFunction(esshDriver))
	// L.SetGlobal("remote_task", L.NewFunction(esshRemoteTask))
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
		"driver":  esshDriver,
		"require": esshRequire,
		"debug":   esshDebug,
	})
}

func esshDebug(L *lua.LState) int {
	msg := L.CheckString(1)
	if debugFlag {
		fmt.Printf("[essh debug] %s\n", msg)
	}

	return 0
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

func esshDriver(L *lua.LState) int {
	first := L.CheckAny(1)
	if tb, ok := first.(*lua.LTable); ok {
		registerDriverByTable(L, tb)
		return 0
	}

	name := L.CheckString(1)

	// procedural style
	if L.GetTop() == 2 {
		tb := L.CheckTable(2)
		registerDriver(L, name, tb)

		return 0
	}

	// DSL style
	L.Push(L.NewFunction(func(L *lua.LState) int {
		tb := L.CheckTable(1)
		registerDriver(L, name, tb)

		return 0
	}))

	return 1
}

func registerDriverByTable(L *lua.LState, tb *lua.LTable) {
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

			registerDriver(L, string(name), config)
		})
	} else { // array
		for i := 1; i <= maxn; i++ {
			value := tb.RawGetInt(i)
			valueTb, ok := value.(*lua.LTable)
			if !ok {
				return
			}
			registerDriverByTable(L, valueTb)
		}
	}
}

func registerDriver(L *lua.LState, name string, config *lua.LTable) {
	driver := NewDriver()
	driver.Name = name
	driver.Context = CurrentContext

	if debugFlag {
		fmt.Printf("[essh debug] register driver: %s\n", name)
	}

	engine := config.RawGetString("engine")
	if engine == lua.LNil {
		L.RaiseError("driver 'engine' is must.")
	} else {
		if engineFn, ok := engine.(*lua.LFunction); ok {
			driver.Engine = func(driver *Driver) (string, error) {
				err := L.CallByParam(lua.P{
					Fn:      engineFn,
					NRet:    1,
					Protect: true,
				}, driver.Config)
				if err != nil {
					return "", err
				}

				ret := L.Get(-1) // returned value
				L.Pop(1)

				if retStr, ok := toString(ret); ok {
					return retStr, nil
				} else {
					return "", fmt.Errorf("driver engine has to return a string.")
				}
			}
		} else if engineStr, ok := toString(engine); ok {
			driver.Engine = func(driver *Driver) (string, error) {
				return engineStr, nil
			}
		} else {
			L.RaiseError("driver 'engine' have to be a function or string.")
		}
	}

	driver.Config = config

	mapper := gluamapper.NewMapper(gluamapper.Option{
		NameFunc: func(s string) string {
			return s
		},
	})
	mapper.Map(driver.Config, &driver.Props)

	Drivers[driver.Name] = driver
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
		Name:    name,
		Config:  newConfig,
		Props:   map[string]string{},
		Hooks:   map[string][]interface{}{},
		Tags:    []string{},
		Context: CurrentContext,
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
		err := registerHook(L, h, "before_connect", hookTb.RawGetString("before_connect"))
		if err != nil {
			panic(err)
		}

		err = registerHook(L, h, "after_connect", hookTb.RawGetString("after_connect"))
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

	extend := config.RawGetString("extend")
	if extendStr, ok := toString(extend); ok {
		h.Extend = extendStr
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
			hooks := host.Hooks[hookPoint]
			hooks = append(hooks, hookFn)
			host.Hooks[hookPoint] = hooks
		} else if hookString, ok := toString(hook); ok {
			hooks := host.Hooks[hookPoint]
			hooks = append(hooks, hookString)
			host.Hooks[hookPoint] = hooks
		} else if tb, ok := toLTable(hook); ok {
			maxn := tb.MaxN()
			if maxn == 0 { // table
				return fmt.Errorf("invalid hook type '%v'. hook must be string, function or table of array.", hook)
			}

			for i := 1; i <= maxn; i++ {
				if err := registerHook(L, host, hookPoint, tb.RawGetInt(i)); err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("invalid hook type '%v'. hook must be string, function or table of array.", hook)
		}
	}

	return nil
}

func registerTask(L *lua.LState, name string, config *lua.LTable) {
	task := NewTask()
	task.Name = name
	task.Context = CurrentContext

	description := config.RawGetString("description")
	if descStr, ok := toString(description); ok {
		task.Description = descStr
	}

	pty := config.RawGetString("pty")
	if ptyBool, ok := toBool(pty); ok {
		task.Pty = ptyBool
	}

	driver := config.RawGetString("driver")
	if driverStr, ok := toString(driver); ok {
		task.Driver = driverStr
	}

	parallel := config.RawGetString("parallel")
	if parallelBool, ok := toBool(parallel); ok {
		task.Parallel = parallelBool
	}

	privileged := config.RawGetString("privileged")
	if privilegedBool, ok := toBool(privileged); ok {
		task.Privileged = privilegedBool
	}

	disabled := config.RawGetString("disabled")
	if disabledBool, ok := toBool(disabled); ok {
		task.Disabled = disabledBool
	}

	hidden := config.RawGetString("hidden")
	if hiddenBool, ok := toBool(hidden); ok {
		task.Hidden = hiddenBool
	}

	lock := config.RawGetString("lock")
	if lockBool, ok := toBool(lock); ok {
		task.Lock = lockBool
	}

	script, err := toScript(L, config.RawGetString("script"))
	if err != nil {
		L.RaiseError("%v", err)
	}
	task.Script = script

	file := config.RawGetString("file")
	if fileStr, ok := toString(file); ok {
		task.File = fileStr
	}

	if task.File != "" && len(task.Script) > 0 {
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

	// configure is deprecated. this is a old version of override_config
	configure := config.RawGetString("configure")
	if configure != lua.LNil {
		if configureFn, ok := configure.(*lua.LFunction); ok {
			task.Configure = func() error {
				err := L.CallByParam(lua.P{
					Fn:      configureFn,
					NRet:    0,
					Protect: false,
				})
				if err != nil {
					return err
				}

				return nil
			}
		} else {
			L.RaiseError("configure have to be a function.")
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
					Protect: false,
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
			L.RaiseError("prepare have to be a function.")
		}
	}

	Tasks = append(Tasks, task)
}

func toScript(L *lua.LState, value lua.LValue) ([]map[string]string, error) {
	ret := []map[string]string{}

	if tb, ok := toLTable(value); ok {
		maxn := tb.MaxN()
		if maxn == 0 { // table
			if tb.RawGetString("code") == lua.LNil {
				return nil, fmt.Errorf("if a 'script' entry is table, it has to have 'code' property.")
			}

			m := map[string]string{}
			tb.ForEach(func(k, v lua.LValue) {
				vs, ok := toString(v)
				if !ok {
					vb, ok := toBool(v)
					if !ok {
						panic("if a 'script' entry is table, it's value has to be string or bool.")
					}
					if vb {
						vs = "true"
					} else {
						vs = "false"
					}
				}
				ks, ok := toString(k)
				if !ok {
					panic("if a 'script' entry is table, it's property has to be string.")
				}
				m[ks] = vs
			})

			ret = append(ret, m)
		} else { // array
			for i := 1; i <= maxn; i++ {
				value := tb.RawGetInt(i)
				if fn, ok := toLFunction(value); ok {
					err := L.CallByParam(lua.P{
						Fn:      fn,
						NRet:    1,
						Protect: false,
					})
					if err != nil {
						panic(err)
					}
					funcRet := L.Get(-1)
					L.Pop(1)

					script, err := toScript(L, funcRet)
					if err != nil {
						return nil, err
					}
					ret = append(ret, script...)
				} else {
					script, err := toScript(L, value)
					if err != nil {
						return nil, err
					}
					ret = append(ret, script...)
				}
			}
		}
		return ret, nil
	} else if str, ok := toString(value); ok {
		return []map[string]string{
			map[string]string{"code": str},
		}, nil
	}

	return nil, fmt.Errorf("'scrpt' got a invalid value.")
}

func esshRequire(L *lua.LState) int {
	name := L.CheckString(1)

	module := CurrentContext.LoadedModules[name]
	if module == nil {
		module = NewModule(name)

		update := updateFlag
		if CurrentContext.Type == ContextTypeGlobal && noGlobalFlag {
			update = false
		}

		err := module.Load(update)
		if err != nil {
			L.RaiseError("%v", err)
		}

		indexFile := module.IndexFile()
		if _, err := os.Stat(indexFile); err != nil {
			L.RaiseError("invalid module: %v", err)
		}
		if err := L.DoFile(indexFile); err != nil {
			L.RaiseError("%v", err)
		}

		// get a module return value
		ret := L.Get(-1)
		module.Value = ret

		// register loaded module.
		CurrentContext.LoadedModules[name] = module

		return 1
	}

	L.Push(module.Value)
	return 1
}

// This code inspired by https://github.com/yuin/gluamapper/blob/master/gluamapper.go
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
