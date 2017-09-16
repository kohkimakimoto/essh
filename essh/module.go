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

func (m *Module) IndexFile() string {
	return path.Join(m.Dir(), "index.lua")
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
	indexFile := m.IndexFile()
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
