package essh

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

type Registry struct {
	Key            string
	DataDir        string
	LoadedPackages map[string]*Package
	Type           int
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
		Key:            fmt.Sprintf("%x", sha256.Sum256([]byte(dataDir))),
		DataDir:        dataDir,
		LoadedPackages: map[string]*Package{},
		Type:           registryType,
	}

	return reg
}

func (reg *Registry) PackagesDir() string {
	return filepath.Join(reg.DataDir, "packages")
}

func (ctx *Registry) CacheDir() string {
	return filepath.Join(ctx.DataDir, "cache")
}

func (reg *Registry) MkDirs() error {
	if _, err := os.Stat(reg.PackagesDir()); os.IsNotExist(err) {
		err = os.MkdirAll(reg.PackagesDir(), os.FileMode(0755))
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(reg.CacheDir()); os.IsNotExist(err) {
		err = os.MkdirAll(reg.CacheDir(), os.FileMode(0755))
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
