package essh

import (
	"fmt"
	"github.com/hashicorp/go-getter"
	"github.com/kohkimakimoto/essh/color"
	"github.com/yuin/gopher-lua"
	"os"
	"path"
	"strings"
)

var LoadedModules = map[string]*Module{}

type Module struct {
	// Name is url that is used as go-getter src.
	// examples:
	//   github.com/aaa/bbb
	//   git::github.com/aaa/bbb.git
	Name string
	// Value is a lua value that is returned when a module's 'index.lua' file is evaluated.
	Value lua.LValue
}

func NewModule(name string) *Module {
	return &Module{
		Name: name,
	}
}

func (m *Module) Load(update bool) error {
	src := m.Name
	dst := m.Dir()

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

	fmt.Fprintf(color.StdoutWriter, "Getting module: '%s'\n", color.FgYB(m.Name))

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

func (m *Module) IndexFile() string {
	return path.Join(m.Dir(), "index.lua")
}

func (m *Module) Dir() string {
	return path.Join(ModulesDir(), m.Key())
}

func (m *Module) Key() string {
	return strings.Replace(strings.Replace(m.Name, "/", "-", -1), ":", "-", -1)
}
