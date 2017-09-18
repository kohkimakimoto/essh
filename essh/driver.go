package essh

import (
	"bytes"
	"fmt"
	"github.com/yuin/gopher-lua"
	"runtime"
	"strings"
	"text/template"
)

type Driver struct {
	Name     string
	Props    map[string]interface{}
	Engine   func(*Driver) (string, error)
	Registry *Registry
	Group    *Group
	LValues  map[string]lua.LValue
	Parent   *Driver
	Child    *Driver
}

var Drivers map[string]*Driver

var DefaultDriver *Driver
var DefaultDriverName = "default"

func NewDriver() *Driver {
	return &Driver{
		Props:   map[string]interface{}{},
		LValues: map[string]lua.LValue{},
	}
}

func (driver *Driver) MapLValuesToLTable(tb *lua.LTable) {
	for key, value := range driver.LValues {
		tb.RawSetString(key, value)
	}
}

func (driver *Driver) GenerateRunnableContent(sshConfigPath string, task *Task, host *Host) (string, error) {
	for key, value := range driver.LValues {
		driver.Props[key] = toGoValue(value)
	}

	if driver.Engine == nil {
		return "", fmt.Errorf("invalid driver '%s'. The engine was not defined.", driver.Name)
	}

	templateText, err := driver.Engine(driver)
	if err != nil {
		return "", err
	}

	scripts := []map[string]string{}
	if task.File != "" {
		tContent, err := GetContentFromPath(task.File)
		if err != nil {
			return "", err
		}
		scripts = append(scripts, map[string]string{"code": string(tContent)})
	} else {
		scripts = task.Script
	}

	funcMap := template.FuncMap{
		"ShellEscape":  ShellEscape,
		"ToUpper":      strings.ToUpper,
		"ToLower":      strings.ToLower,
		"EnvKeyEscape": EnvKeyEscape,
		"Add": func(x, y int) int {
			return x + y
		},
	}

	dict := map[string]interface{}{
		"GOARCH":        runtime.GOARCH,
		"GOOS":          runtime.GOOS,
		"Debug":         debugFlag,
		"Driver":        driver,
		"Task":          task,
		"Host":          host,
		"Scripts":       scripts,
		"SSHConfigPath": sshConfigPath,
	}

	baseTempl, err := template.New("base").Funcs(funcMap).Parse(templateText)
	if err != nil {
		return "", err
	}

	tmpl, err := baseTempl.Parse(EnvironmentTemplate)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, dict)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

const EnvironmentTemplate = `{{define "environment" -}}
export ESSH_TASK_NAME={{.Task.Name | ShellEscape}}
export ESSH_SSH_CONFIG={{.SSHConfigPath}}
export ESSH_DEBUG="{{if .Debug}}1{{end}}"
{{range $key, $value := .Task.Props -}}
export ESSH_TASK_PROPS_{{$key | ToUpper | EnvKeyEscape}}={{$value | ShellEscape }}
{{end -}}
{{range $index, $value := .Task.Args -}}
export ESSH_TASK_ARGS_{{Add $index 1 }}={{$value | ShellEscape }}
{{end -}}
{{if .Host -}}
export ESSH_HOSTNAME={{.Host.Name | ShellEscape}}
export ESSH_HOST_HOSTNAME={{.Host.Name | ShellEscape}}
{{range $i, $kvpair := .Host.SortedSSHConfig -}}
{{range $key, $value := $kvpair -}}
export ESSH_HOST_SSH_{{$key | ToUpper}}={{$value | ShellEscape }}
{{end -}}
{{end -}}
{{range $key, $value := .Host.Props -}}
export ESSH_HOST_PROPS_{{$key | ToUpper | EnvKeyEscape}}={{$value | ShellEscape }}
{{end -}}
{{range $i, $value := .Host.Tags -}}
export ESSH_HOST_TAGS_{{$value | ToUpper | EnvKeyEscape}}=1
{{end -}}
{{end -}}
{{end}}
`

func removeDriverInGlobalSpace(driver *Driver) {
	d := Drivers[driver.Name]
	if d == driver {
		if d.Child != nil {
			newDriver := d.Child
			Drivers[newDriver.Name] = newDriver
			newDriver.Parent = nil
		} else {
			delete(Drivers, d.Name)
		}
	}
}

func esshDriver(L *lua.LState) int {
	first := L.CheckAny(1)
	if tb, ok := toLTable(first); ok {
		name := DefaultDriverName
		d := registerDriver(L, name)
		setupDriver(L, d, tb)
		L.Push(newLDriver(L, d))

		return 1
	}

	name := L.CheckString(1)
	if L.GetTop() == 1 {
		// object or DSL style
		d := registerDriver(L, name)
		L.Push(newLDriver(L, d))

		return 1
	} else if L.GetTop() == 2 {
		// function style
		tb := L.CheckTable(2)
		d := registerDriver(L, name)
		setupDriver(L, d, tb)
		L.Push(newLDriver(L, d))

		return 1
	}

	panic("driver requires 1 or 2 arguments")
}

func registerDriver(L *lua.LState, name string) *Driver {
	if debugFlag {
		fmt.Printf("[essh debug] register driver: %s\n", name)
	}

	d := NewDriver()
	d.Name = name
	d.Registry = CurrentRegistry

	if driver := Drivers[d.Name]; driver != nil {
		// detect same name driver
		d.Child = driver
		driver.Parent = d
	}

	if EvaluatingModule != nil {
		EvaluatingModule.Drivers = append(EvaluatingModule.Drivers, d)
	}

	Drivers[d.Name] = d

	return d
}

func setupDriver(L *lua.LState, driver *Driver, config *lua.LTable) {
	config.ForEach(func(k, v lua.LValue) {
		if kstr, ok := toString(k); ok {
			updateDriver(L, driver, kstr, v)
		}
	})
}

func updateDriver(L *lua.LState, driver *Driver, key string, value lua.LValue) {
	driver.LValues[key] = value

	switch key {
	case "engine":
		if engineFn, ok := value.(*lua.LFunction); ok {
			driver.Engine = func(driver *Driver) (string, error) {
				err := L.CallByParam(lua.P{
					Fn:      engineFn,
					NRet:    1,
					Protect: true,
				}, newLDriver(L, driver))
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
		} else if engineStr, ok := toString(value); ok {
			driver.Engine = func(driver *Driver) (string, error) {
				return engineStr, nil
			}
		} else {
			L.RaiseError("driver 'engine' have to be a function or string.")
		}
	}
}

const LDriverClass = "Driver*"

func registerDriverClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LDriverClass)
	mt.RawSetString("__call", L.NewFunction(driverCall))
	mt.RawSetString("__index", L.NewFunction(driverIndex))
	mt.RawSetString("__newindex", L.NewFunction(driverNewindex))
}

func newLDriver(L *lua.LState, driver *Driver) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = driver
	L.SetMetatable(ud, L.GetTypeMetatable(LDriverClass))
	return ud
}

func checkDriver(L *lua.LState) *Driver {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Driver); ok {
		return v
	}
	L.ArgError(1, "Driver object expected")
	return nil
}

func driverCall(L *lua.LState) int {
	driver := checkDriver(L)
	tb := L.CheckTable(2)

	setupDriver(L, driver, tb)

	L.Push(L.CheckUserData(1))
	return 1
}

func driverIndex(L *lua.LState) int {
	driver := checkDriver(L)
	index := L.CheckString(2)

	if index == "name" {
		L.Push(L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LString(driver.Name))
			return 1
		}))
		return 1
	}

	v, ok := driver.LValues[index]
	if v == nil || !ok {
		v = lua.LNil
	}

	L.Push(v)
	return 1
}

func driverNewindex(L *lua.LState) int {
	driver := checkDriver(L)
	index := L.CheckString(2)
	value := L.CheckAny(3)

	updateDriver(L, driver, index, value)

	return 0
}
