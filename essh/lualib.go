package essh

import (
	"fmt"
	"github.com/cjoudrey/gluahttp"
	"github.com/kohkimakimoto/essh/support/gluamapper"
	"github.com/kohkimakimoto/gluaenv"
	"github.com/kohkimakimoto/gluafs"
	"github.com/kohkimakimoto/gluaquestion"
	"github.com/kohkimakimoto/gluatemplate"
	"github.com/kohkimakimoto/gluayaml"
	gluajson "github.com/layeh/gopher-json"
	"github.com/otm/gluash"
	"github.com/yuin/gluare"
	"github.com/yuin/gopher-lua"
	"net/http"
	"os"
	"unicode"
)

func InitLuaState(L *lua.LState) {
	// custom type.
	registerTaskContextClass(L)
	registerHostClass(L)
	registerRegistryClass(L)

	// global functions
	L.SetGlobal("host", L.NewFunction(esshHost))
	L.SetGlobal("private_host", L.NewFunction(esshPrivateHost))
	L.SetGlobal("task", L.NewFunction(esshTask))
	L.SetGlobal("driver", L.NewFunction(esshDriver))

	// modules
	L.PreloadModule("json", gluajson.Loader)
	L.PreloadModule("fs", gluafs.Loader)
	L.PreloadModule("yaml", gluayaml.Loader)
	L.PreloadModule("template", gluatemplate.Loader)
	L.PreloadModule("question", gluaquestion.Loader)
	L.PreloadModule("env", gluaenv.Loader)
	L.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	L.PreloadModule("re", gluare.Loader)
	L.PreloadModule("sh", gluash.Loader)
	// for BC
	L.PreloadModule("glua.json", gluajson.Loader)
	L.PreloadModule("glua.fs", gluafs.Loader)
	L.PreloadModule("glua.yaml", gluayaml.Loader)
	L.PreloadModule("glua.template", gluatemplate.Loader)
	L.PreloadModule("glua.question", gluaquestion.Loader)
	L.PreloadModule("glua.env", gluaenv.Loader)
	L.PreloadModule("glua.http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	L.PreloadModule("glua.re", gluare.Loader)
	L.PreloadModule("glua.sh", gluash.Loader)

	// global variables
	lessh := L.NewTable()
	L.SetGlobal("essh", lessh)
	lessh.RawSetString("ssh_config", lua.LNil)

	L.SetFuncs(lessh, map[string]lua.LGFunction{
		// aliases global function.
		"host":         esshHost,
		"private_host": esshPrivateHost,
		"task":   esshTask,
		"driver": esshDriver,
		// utilities.
		"require":  esshRequire,
		"debug":    esshDebug,
		"gethosts": esshGethosts,
		"hosts":    esshGethosts,
		"registry": esshRegistry,
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
	name := L.CheckString(1)
	if L.GetTop() == 1 {
		// object or DSL style
		h := registerHost(L, name)
		L.Push(newLHost(L, h))

		return 1
	} else if L.GetTop() == 2 {
		// function style
		tb := L.CheckTable(2)
		h := registerHost(L, name)
		setupHost(L, h, tb)
		L.Push(newLHost(L, h))

		return 1
	}

	panic("host requires 1 or 2 arguments")
}

func esshPrivateHost(L *lua.LState) int {

	name := L.CheckString(1)
	if L.GetTop() == 1 {
		// object or DSL style
		h := registerHost(L, name)
		updateHost(L, h, "private", lua.LBool(true))
		L.Push(newLHost(L, h))

		return 1
	} else if L.GetTop() == 2 {
		// function style
		tb := L.CheckTable(2)
		h := registerHost(L, name)
		updateHost(L, h, "private", lua.LBool(true))
		setupHost(L, h, tb)
		L.Push(newLHost(L, h))

		return 1
	}

	panic("private_host requires 1 or 2 arguments")

	return 0
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

func registerHost(L *lua.LState, name string) *Host {
	if debugFlag {
		fmt.Printf("[essh debug] register host: %s\n", name)
	}

	h := NewHost()
	h.Name = name
	h.Registry = CurrentRegistry

	CurrentRegistry.Hosts[h.Name] = h

	return h
}

func setupHost(L *lua.LState, h *Host, config *lua.LTable) {
	config.ForEach(func(k, v lua.LValue) {
		if kstr, ok := toString(k); ok {
			updateHost(L, h, kstr, v)
		}
	})
}

func updateHost(L *lua.LState, h *Host, key string, value lua.LValue) {
	h.LValues[key] = value

	var firstChar rune
	for _, c := range key {
		firstChar = c
		break
	}

	if unicode.IsUpper(firstChar) {
		if valuestr, ok := toString(value); ok {
			h.SSHConfig[key] = valuestr
			return
		}

		panic("SSH property must be string")
	}

	switch key {
	case "props":
		if propsTb, ok := toLTable(value); ok {
			// initialize
			h.Props = map[string]string{}

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

	case "hooks":
		if hookTb, ok := toLTable(value); ok {
			// initialize
			h.Hooks = map[string][]interface{}{}

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
	case "description":
		if descStr, ok := toString(value); ok {
			h.Description = descStr
		}
	case "hidden":
		if hiddenBool, ok := toBool(value); ok {
			h.Hidden = hiddenBool
		}
	case "private":
		if privateBool, ok := toBool(value); ok {
			h.Private = privateBool
		}
	case "tags":
		if tagsTb, ok := toLTable(value); ok {
			// initialize
			h.Tags = []string{}

			tagsTb.ForEach(func(_ lua.LValue, v lua.LValue) {
				if vs, ok := toString(v); ok {
					h.Tags = append(h.Tags, vs)
				} else {
					L.RaiseError("unsupported format of tags.")
				}
			})
		}
	default:
		panic("unsupported host's field '" + key + "'.")

	}
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

func registerTask(L *lua.LState, name string, config *lua.LTable) *Task {
	task := NewTask()
	task.Name = name
	task.Registry = CurrentRegistry

	backend := config.RawGetString("backend")
	if backendStr, ok := toString(backend); ok {
		task.Backend = backendStr
		if backendStr != TASK_BACKEND_LOCAL && backendStr != TASK_BACKEND_REMOTE {
			L.RaiseError("backend must be '%s' or '%s'.", TASK_BACKEND_LOCAL, TASK_BACKEND_REMOTE)
		}
	}

	targets := config.RawGetString("targets")
	if targetsStr, ok := toString(targets); ok {
		task.Targets = []string{targetsStr}
	} else if targetsSlice, ok := toSlice(targets); ok {
		for _, target := range targetsSlice {
			if targetStr, ok := target.(string); ok {
				task.Targets = append(task.Targets, targetStr)
			}
		}
	}

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

	Tasks[task.Name] = task

	return task
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

	module := CurrentRegistry.LoadedModules[name]
	if module == nil {
		module = NewModule(name)

		update := updateFlag
		if CurrentRegistry.Type == RegistryTypeGlobal && !withGlobalFlag {
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
		CurrentRegistry.LoadedModules[name] = module

		return 1
	}

	L.Push(module.Value)
	return 1
}

func esshGethosts(L *lua.LState) int {
	lhosts := L.NewTable()

	for _, host := range SortedHosts() {
		lhost := newLHost(L, host)
		lhosts.Append(lhost)
	}

	L.Push(lhosts)
	return 1
}

func esshRegistry(L *lua.LState) int {
	L.Push(newLRegistry(L, CurrentRegistry))
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
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"payload": taskContextPayload,
	}))
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
	L.ArgError(1, "TaskContext object expected")
	return nil
}

const LHostClass = "Host*"

func newLHost(L *lua.LState, host *Host) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = host
	L.SetMetatable(ud, L.GetTypeMetatable(LHostClass))
	return ud
}

func checkHost(L *lua.LState) *Host {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Host); ok {
		return v
	}
	L.ArgError(1, "Host object expected")
	return nil
}

func registerHostClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LHostClass)
	mt.RawSetString("__call", L.NewFunction(hostCall))
	mt.RawSetString("__index", L.NewFunction(hostIndex))
	mt.RawSetString("__newindex", L.NewFunction(hostNewindex))
}

func hostCall(L *lua.LState) int {
	host := checkHost(L)
	tb := L.CheckTable(2)

	setupHost(L, host, tb)

	return 0
}

func hostIndex(L *lua.LState) int {
	host := checkHost(L)
	index := L.CheckString(2)

	v, ok := host.LValues[index]
	if v == nil || !ok {
		v = lua.LNil
	}

	L.Push(v)
	return 1
}

func hostNewindex(L *lua.LState) int {
	host := checkHost(L)
	index := L.CheckString(2)
	value := L.CheckAny(3)

	updateHost(L, host, index, value)

	return 0
}

const LRegistryClass = "Registry*"

func newLRegistry(L *lua.LState, ctx *Registry) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = ctx
	L.SetMetatable(ud, L.GetTypeMetatable(LRegistryClass))
	return ud
}

func checkRegistry(L *lua.LState) *Registry {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Registry); ok {
		return v
	}
	L.ArgError(1, "Registry object expected")
	return nil
}

func registerRegistryClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LRegistryClass)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"data_dir":    registryDataDir,
		"tmp_dir":     registryTmpDir,
		"modules_dir": registryModulesDir,
		"type":        registryType,
	}))
}

func registryDataDir(L *lua.LState) int {
	reg := checkRegistry(L)
	L.Push(lua.LString(reg.DataDir))
	return 1
}

func registryTmpDir(L *lua.LState) int {
	reg := checkRegistry(L)
	L.Push(lua.LString(reg.TmpDir()))
	return 1
}

func registryModulesDir(L *lua.LState) int {
	reg := checkRegistry(L)
	L.Push(lua.LString(reg.ModulesDir()))
	return 1
}

func registryType(L *lua.LState) int {
	reg := checkRegistry(L)
	L.Push(lua.LString(reg.TypeString()))
	return 1
}
