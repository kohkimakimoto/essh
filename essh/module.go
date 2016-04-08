package essh

import (
	"fmt"
	"github.com/hashicorp/go-getter"
	"github.com/kohkimakimoto/essh/color"
	"os"
	"path"
	"strings"
)

type Module struct {
	// Name is url that is used as go-getter src.
	// examples:
	//   github.com/aaa/bbb
	//   git::github.com/aaa/bbb.git
	Name string
}

func NewModule(name string) *Module {
	return &Module{
		Name: name,
	}
}

func (m *Module) GetModule(update bool) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	source, err := getter.Detect(m.Name, wd, getter.Detectors)
	if err != nil {
		return err
	}

	dir := m.Dir()

	if !update {
		if _, err := os.Stat(dir); err == nil {
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
		fmt.Printf("[essh debug] module src '%s'", source)
	}

	fmt.Fprintf(color.StdoutWriter, "Getting module: '%s'\n", color.FgYB(m.Name))

	err = getter.Get(dir, source)
	if err != nil {
		os.RemoveAll(dir)
		return err
	}

	return nil
}

func (m *Module) Dir() string {
	return path.Join(ModulesDir, m.Key())
}

func (m *Module) Key() string {
	return strings.Replace(strings.Replace(m.Name, "/", "-", -1), ":", "-", -1)
}
