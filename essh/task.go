package essh

import (
	"fmt"
	"github.com/yuin/gopher-lua"
)

type Task struct {
	Name        string
	Description string
	Props       map[string]string
	Prepare     func() error
	Driver      string
	Pty         bool
	Script      []map[string]string
	File        string
	Backend     string
	Targets     []string
	Filters     []string
	Parallel    bool
	Privileged  bool
	User        string
	// deprecated? use only hidden?
	Disabled  bool
	Hidden    bool
	Prefix    string
	UsePrefix bool
	Registry  *Registry
	Group     *Group
	Args      []string
	LValues   map[string]lua.LValue
	Parent    *Task
	Child     *Task
}

var Tasks map[string]*Task

var DefaultTaskName = "default"

var (
	DefaultPrefixLocal  = `[local:{{.Host.Name}}]{{HostnameAlignString " "}}`
	DefaultPrefixRemote = `[remote:{{.Host.Name}}]{{HostnameAlignString " "}}`
)

const (
	TASK_BACKEND_LOCAL  = "local"
	TASK_BACKEND_REMOTE = "remote"
)

func NewTask() *Task {
	return &Task{
		Targets: []string{},
		Filters: []string{},
		Backend: TASK_BACKEND_LOCAL,
		Script:  []map[string]string{},
		Args:    []string{},
		LValues: map[string]lua.LValue{},
	}
}

func (t *Task) MapLValuesToLTable(tb *lua.LTable) {
	for key, value := range t.LValues {
		tb.RawSetString(key, value)
	}
}

func (t *Task) PublicName() string {
	return t.Name
}

func (t *Task) IsRemoteTask() bool {
	if t.Backend == TASK_BACKEND_REMOTE {
		return true
	} else {
		return false
	}
}

func (t *Task) TargetsSlice() []string {
	if len(t.Targets) >= 1 {
		return t.Targets
	}

	return []string{}
}

func (t *Task) FiltersSlice() []string {
	if len(t.Filters) >= 1 {
		return t.Filters
	}

	return []string{}
}

func (t *Task) DescriptionOrDefault() string {
	if t.Description == "" {
		return t.Name + " task"
	}

	return t.Description
}

func esshTask(L *lua.LState) int {
	first := L.CheckAny(1)
	if tb, ok := toLTable(first); ok {
		name := DefaultTaskName
		j := registerTask(L, name)
		setupTask(L, j, tb)
		L.Push(newLTask(L, j))

		return 1
	}

	name := L.CheckString(1)
	if L.GetTop() == 1 {
		// object or DSL style
		t := registerTask(L, name)
		L.Push(newLTask(L, t))

		return 1
	} else if L.GetTop() == 2 {
		// function style
		tb := L.CheckTable(2)
		t := registerTask(L, name)
		setupTask(L, t, tb)
		L.Push(newLTask(L, t))

		return 1
	}

	panic("task requires 1 or 2 arguments")
}

func registerTask(L *lua.LState, name string) *Task {
	if debugFlag {
		fmt.Printf("[essh debug] register task: %s\n", name)
	}

	t := NewTask()
	t.Name = name
	t.Registry = CurrentRegistry

	if task := Tasks[t.Name]; task != nil {
		// detect same name task
		t.Child = task
		task.Parent = t
	}

	Tasks[t.Name] = t

	return t
}

func setupTask(L *lua.LState, t *Task, config *lua.LTable) {
	config.ForEach(func(k, v lua.LValue) {
		if kstr, ok := toString(k); ok {
			updateTask(L, t, kstr, v)
		}
	})
}

const LTaskClass = "Task*"

func registerTaskClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LTaskClass)
	mt.RawSetString("__call", L.NewFunction(taskCall))
	mt.RawSetString("__index", L.NewFunction(taskIndex))
	mt.RawSetString("__newindex", L.NewFunction(taskNewindex))
}

func newLTask(L *lua.LState, task *Task) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = task
	L.SetMetatable(ud, L.GetTypeMetatable(LTaskClass))
	return ud
}

func checkTask(L *lua.LState) *Task {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Task); ok {
		return v
	}
	L.ArgError(1, "Task object expected")
	return nil
}

func taskCall(L *lua.LState) int {
	task := checkTask(L)

	arg := L.CheckAny(2)
	if tb, ok := toLTable(arg); ok {
		setupTask(L, task, tb)
		L.Push(L.CheckUserData(1))
		return 1
	}

	if lv, ok := arg.(lua.LString); ok {
		updateTask(L, task, "script", lv)
		L.Push(L.CheckUserData(1))
		return 1
	}

	L.ArgError(2, "Table or string expected")
	return 0
}

func taskIndex(L *lua.LState) int {
	task := checkTask(L)
	index := L.CheckString(2)

	if index == "name" {
		L.Push(L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LString(task.Name))
			return 1
		}))
		return 1
	}

	v, ok := task.LValues[index]
	if v == nil || !ok {
		v = lua.LNil
	}

	L.Push(v)
	return 1
}

func taskNewindex(L *lua.LState) int {
	task := checkTask(L)
	index := L.CheckString(2)
	value := L.CheckAny(3)

	updateTask(L, task, index, value)

	return 0
}

func updateTask(L *lua.LState, task *Task, key string, value lua.LValue) {
	task.LValues[key] = value

	switch key {
	case "backend":
		if backendStr, ok := toString(value); ok {
			task.Backend = backendStr
			if backendStr != TASK_BACKEND_LOCAL && backendStr != TASK_BACKEND_REMOTE {
				L.RaiseError("backend must be '%s' or '%s'.", TASK_BACKEND_LOCAL, TASK_BACKEND_REMOTE)
			}
		}
	case "targets":
		if targetsStr, ok := toString(value); ok {
			task.Targets = []string{targetsStr}
		} else if targetsSlice, ok := toSlice(value); ok {
			task.Targets = []string{}

			for _, target := range targetsSlice {
				if targetStr, ok := target.(string); ok {
					task.Targets = append(task.Targets, targetStr)
				}
			}
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	case "filters":
		if filtersStr, ok := toString(value); ok {
			task.Filters = []string{filtersStr}
		} else if filtersSlice, ok := toSlice(value); ok {
			task.Filters = []string{}

			for _, filter := range filtersSlice {
				if filterStr, ok := filter.(string); ok {
					task.Filters = append(task.Filters, filterStr)
				}
			}
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	case "description":
		if descStr, ok := toString(value); ok {
			task.Description = descStr
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	case "pty":
		if ptyBool, ok := toBool(value); ok {
			task.Pty = ptyBool
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	case "driver":
		if driverStr, ok := toString(value); ok {
			task.Driver = driverStr
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	case "parallel":
		if parallelBool, ok := toBool(value); ok {
			task.Parallel = parallelBool
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	case "user":
		if userStr, ok := toString(value); ok {
			task.User = userStr
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	case "privileged":
		if privilegedBool, ok := toBool(value); ok {
			task.Privileged = privilegedBool
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	case "disabled":
		if disabledBool, ok := toBool(value); ok {
			task.Disabled = disabledBool
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	case "hidden":
		if hiddenBool, ok := toBool(value); ok {
			task.Hidden = hiddenBool
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	case "script":
		script, err := toScript(L, value)
		if err != nil {
			L.RaiseError("%v", err)
		}
		task.Script = script

		if task.File != "" && len(task.Script) > 0 {
			L.RaiseError("invalid task definition: can't use 'script_file' and 'script' at the same time.")
		}
	case "script_file":
		if fileStr, ok := toString(value); ok {
			task.File = fileStr
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}

		if task.File != "" && len(task.Script) > 0 {
			L.RaiseError("invalid task definition: can't use 'script_file' and 'script' at the same time.")
		}
	case "prefix":
		if prefixBool, ok := toBool(value); ok {
			task.UsePrefix = prefixBool
		} else if prefixStr, ok := toString(value); ok {
			task.UsePrefix = true
			task.Prefix = prefixStr
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	case "prepare":
		if prepareFn, ok := value.(*lua.LFunction); ok {
			task.Prepare = func() error {
				err := L.CallByParam(lua.P{
					Fn:      prepareFn,
					NRet:    1,
					Protect: false,
				}, newLTask(L, task))
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
	case "props":
		if propsTb, ok := toLTable(value); ok {
			// initialize
			task.Props = map[string]string{}

			propsTb.ForEach(func(propsKey lua.LValue, propsValue lua.LValue) {
				propsKeyStr, ok := toString(propsKey)
				if !ok {
					L.RaiseError("props table's key must be a string: %v", propsKey)
				}
				propsValueStr, ok := toString(propsValue)
				if !ok {
					L.RaiseError("props table's value must be a string: %v", propsValue)
				}

				task.Props[propsKeyStr] = propsValueStr
			})
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	case "args":
		if argsSlice, ok := toSlice(value); ok {
			task.Args = []string{}

			for _, arg := range argsSlice {
				if argStr, ok := arg.(string); ok {
					task.Args = append(task.Args, argStr)
				}
			}
		} else {
			panic("invalid value of a task's field '" + key + "'.")
		}
	default:
		panic("unsupported task's field '" + key + "'.")
	}
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

	return nil, fmt.Errorf("'script' got a invalid value.")
}
