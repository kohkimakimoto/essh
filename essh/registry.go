package essh

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

type Registry struct {
	Key           string
	DataDir       string
	LoadedModules map[string]*Module
	Hosts         map[string]*Host
	Tasks         map[string]*Task
	Drivers       map[string]*Driver
	Type          int
}

const (
	RegistryTypeGlobal = 0
	RegistryTypeLocal  = 1
)

var CurrentRegistry *Registry
var GlobalRegistry *Registry
var LocalRegistry *Registry

func NewRegistry(dataDir string, registryType int) *Registry {
	reg := &Registry{
		Key:           fmt.Sprintf("%x", sha256.Sum256([]byte(dataDir))),
		DataDir:       dataDir,
		LoadedModules: map[string]*Module{},
		Hosts:         map[string]*Host{},
		Tasks:         map[string]*Task{},
		Drivers:       map[string]*Driver{
			DefaultDriverName: DefaultDriver,
		},
		Type:          registryType,
	}

	return reg
}

func (reg *Registry) ModulesDir() string {
	return filepath.Join(reg.DataDir, "modules")
}

func (ctx *Registry) TmpDir() string {
	return filepath.Join(ctx.DataDir, "tmp")
}

func (reg *Registry) MkDirs() error {
	if _, err := os.Stat(reg.ModulesDir()); os.IsNotExist(err) {
		err = os.MkdirAll(reg.ModulesDir(), os.FileMode(0755))
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(reg.TmpDir()); os.IsNotExist(err) {
		err = os.MkdirAll(reg.TmpDir(), os.FileMode(0755))
		if err != nil {
			return err
		}
	}

	return nil
}

func (reg *Registry) TypeString() string {
	if reg.Type == RegistryTypeGlobal {
		return "global"
	} else if reg.Type == RegistryTypeLocal {
		return "local"
	}

	panic("unknown context")
}
