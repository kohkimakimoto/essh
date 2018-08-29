package essh

import (
	"fmt"
	"github.com/hashicorp/go-getter"
	"github.com/kohkimakimoto/essh/support/color"
	"github.com/yuin/gopher-lua"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Module struct {
	// Name is url that is used as go-getter src.
	// examples:
	//   github.com/aaa/bbb
	//   git::github.com/aaa/bbb.git
	Name string
	// LValues
	LValues map[string]lua.LValue
	// L
	L *lua.LState
	// Evaluated
	Evaluated bool

	Hosts   []*Host
	Tasks   []*Task
	Drivers []*Driver
	Modules []*Module
	Parant  *Module
}

var Modules []*Module = []*Module{}

var EvaluatingModule *Module

var UpdatedModules map[string]*Module = map[string]*Module{}

func NewModule(L *lua.LState, name string) *Module {
	return &Module{
		Name:      name,
		LValues:   map[string]lua.LValue{},
		L:         L,
		Evaluated: false,
		Hosts:     []*Host{},
		Tasks:     []*Task{},
		Drivers:   []*Driver{},
		Modules:   []*Module{},
	}
}

func (m *Module) AllHosts() []*Host {
	all := m.Hosts

	for _, nm := range m.Modules {
		all = append(all, nm.AllHosts()...)
	}

	return all
}

func (m *Module) AllTasks() []*Task {
	all := m.Tasks

	for _, nm := range m.Modules {
		all = append(all, nm.AllTasks()...)
	}

	return all
}

func (m *Module) AllDrivers() []*Driver {
	all := m.Drivers

	for _, nm := range m.Modules {
		all = append(all, nm.AllDrivers()...)
	}

	return all
}

func (m *Module) MapLValuesToLTable(tb *lua.LTable) {
	for key, value := range m.LValues {
		tb.RawSetString(key, value)
	}
}

func (m *Module) Load(update bool) error {
	// If you usually use git with essh, you can set variable "GIT_SSH=essh".
	// But this setting cause a error in a module loading.
	// When we load a module, essh can git protocol, but essh hasn't generated ssh_config used by git command.
	gitssh := os.Getenv("GIT_SSH")
	if filepath.Base(gitssh) == "essh" {
		os.Setenv("GIT_SSH", "ssh")
		defer func() {
			os.Setenv("GIT_SSH", gitssh)
		}()
	}

	src := m.Src()
	dst := m.Dir()

	if UpdatedModules[m.Name] != nil {
		// already updated
		update = false
	}

	if !update {
		if _, err := os.Stat(dst); err == nil {
			// If the directory already exists, then we're done since
			// we're not updating.
			return nil
		} else if !os.IsNotExist(err) {
			// If the error we got wasn't a file-not-exist error, then
			// something went wrong and we should report it.
			return fmt.Errorf("Error reading directory: %s", err)
		}
	}

	if debugFlag {
		fmt.Printf("[essh debug] module src '%s'\n", src)
	}

	if update {
		if _, err := os.Stat(dst); err == nil {
			fmt.Fprintf(os.Stdout, "Updating module: '%s' (into %s)\n", color.FgYB(m.Name), color.FgBold(CurrentRegistry.DataDir))
		} else {
			fmt.Fprintf(os.Stdout, "Installing module: '%s' (into %s)\n", color.FgYB(m.Name), color.FgBold(CurrentRegistry.DataDir))
		}

		UpdatedModules[m.Name] = m
	} else {
		fmt.Fprintf(os.Stdout, "Installing module: '%s' (into %s)\n", color.FgYB(m.Name), color.FgBold(CurrentRegistry.DataDir))
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	client := &getter.Client{
		Src:  src,
		Dst:  dst,
		Pwd:  pwd,
		Mode: getter.ClientModeDir,
	}
	if err := client.Get(); err != nil {
		return err
	}

	return nil
}

func (m *Module) Src() string {
	src := m.Name

	return src
}

func (m *Module) IndexFile() (string, error) {
	idx := path.Join(m.Dir(), "esshmodule.lua")
	if _, err := os.Stat(idx); os.IsNotExist(err) {
		idx = path.Join(m.Dir(), "index.lua")
		if _, err := os.Stat(idx); os.IsNotExist(err) {
			return "", fmt.Errorf("not found esshmodule.lua or index.lua: %v", err)
		}
	}

	return idx, nil
}

func (m *Module) Dir() string {
	return path.Join(CurrentRegistry.ModulesDir(), m.Key())
}

func (m *Module) Key() string {
	return strings.Replace(strings.Replace(m.Name, "/", "-", -1), ":", "-", -1)
}

func (m *Module) Evaluate() error {
	if m.Evaluated {
		return nil
	}

	if debugFlag {
		fmt.Printf("[essh debug] module evaluating '%s'\n", m.Name)
	}

	EvaluatingModule = m

	err := m.Load(updateFlag)
	if err != nil {
		return err
	}

	L := m.L
	indexFile, err := m.IndexFile()
	if err != nil {
		return err
	}

	lessh, ok := toLTable(L.GetGlobal("essh"))
	if !ok {
		return fmt.Errorf("'essh' global variable is broken")
	}

	// configure module variable
	modulevar := L.NewTable()
	modulevar.RawSetString("path", lua.LString(filepath.Dir(indexFile)))
	params := L.NewTable()
	for kstr, v := range m.LValues {
		params.RawSetString(kstr, v)
	}
	modulevar.RawSetString("params", params)

	// deprecated. for BC
	modulevar.RawSetString("var", params)

	lessh.RawSetString("module", modulevar)

	if err := L.DoFile(indexFile); err != nil {
		return err
	}

	// remove pkg variable
	lessh.RawSetString("module", lua.LNil)
	m.Evaluated = true
	EvaluatingModule = nil

	if len(m.Modules) > 0 {
		// has nested modules
		for _, nm := range m.Modules {
			if err := nm.Evaluate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func esshModule(L *lua.LState) int {
	value := L.CheckAny(1)
	if tb, ok := toLTable(value); ok {
		modulesTb := L.NewTable()
		tb.ForEach(func(k, v lua.LValue) {
			name, ok := toString(k)
			if !ok {
				panic(fmt.Sprintf("expected string of module's name but got '%v'\n", k))
			}

			config, ok := toLTable(v)
			if !ok {
				panic(fmt.Sprintf("expected table of module's config but got '%v'\n", v))
			}

			m := registerModule(L, name)
			setupModule(L, m, config)
			modulesTb.RawSetString(name, newLModule(L, m))
		})

		L.Push(modulesTb)
		return 1
	} else if name, ok := toString(value); ok {
		if L.GetTop() == 1 {
			// object or DSL style
			m := registerModule(L, name)
			L.Push(newLModule(L, m))

			return 1
		} else if L.GetTop() == 2 {
			// function style
			tb := L.CheckTable(2)
			m := registerModule(L, name)
			setupModule(L, m, tb)
			L.Push(newLModule(L, m))

			return 1
		} else {
			panic("module requires 1 or 2 arguments")
		}
	} else {
		panic(fmt.Sprintf("expected table or string but got '%v'\n", value))
	}
	return 0
}

func registerModule(L *lua.LState, name string) *Module {
	m := NewModule(L, name)

	Modules = append(Modules, m)

	if EvaluatingModule != nil {
		m.Parant = EvaluatingModule
		EvaluatingModule.Modules = append(EvaluatingModule.Modules, m)
	}

	return m
}

func setupModule(L *lua.LState, h *Module, config *lua.LTable) {
	config.ForEach(func(k, v lua.LValue) {
		if kstr, ok := toString(k); ok {
			updateModule(L, h, kstr, v)
		}
	})
}

func updateModule(L *lua.LState, h *Module, key string, value lua.LValue) {
	h.LValues[key] = value
}

const LModuleClass = "Module*"

func registerModuleClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LModuleClass)
	mt.RawSetString("__call", L.NewFunction(moduleCall))
	mt.RawSetString("__index", L.NewFunction(moduleIndex))
	mt.RawSetString("__newindex", L.NewFunction(moduleNewindex))
}

func newLModule(L *lua.LState, module *Module) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = module
	L.SetMetatable(ud, L.GetTypeMetatable(LModuleClass))
	return ud
}

func checkModule(L *lua.LState) *Module {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Module); ok {
		return v
	}
	L.ArgError(1, "Module object expected")
	return nil
}

func moduleCall(L *lua.LState) int {
	module := checkModule(L)
	tb := L.CheckTable(2)

	setupModule(L, module, tb)

	L.Push(L.CheckUserData(1))
	return 1
}

func moduleIndex(L *lua.LState) int {
	module := checkModule(L)
	index := L.CheckString(2)

	v, ok := module.LValues[index]
	if v == nil || !ok {
		v = lua.LNil
	}

	L.Push(v)
	return 1
}

func moduleNewindex(L *lua.LState) int {
	module := checkModule(L)
	index := L.CheckString(2)
	value := L.CheckAny(3)

	updateModule(L, module, index, value)

	return 0
}
