package essh

import (
	"crypto/sha256"
	"fmt"
	"github.com/yuin/gopher-lua"
)

type Registry struct {
	Key            string
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
		Type:           registryType,
	}

	return reg
}

//func (reg *Registry) PackagesDir() string {
//	return filepath.Join(reg.DataDir, "packages")
//}
//
//func (reg *Registry) ModulesDir() string {
//	return filepath.Join(reg.DataDir, "modules")
//}
//
//func (reg *Registry) LibDir() string {
//	return filepath.Join(reg.DataDir, "lib")
//}
//
//func (ctx *Registry) CacheDir() string {
//	return filepath.Join(ctx.DataDir, "cache")
//}
//
//func (reg *Registry) MkDirs() error {
//	if _, err := os.Stat(reg.PackagesDir()); os.IsNotExist(err) {
//		err = os.MkdirAll(reg.PackagesDir(), os.FileMode(0755))
//		if err != nil {
//			return err
//		}
//	}
//
//	if _, err := os.Stat(reg.ModulesDir()); os.IsNotExist(err) {
//		err = os.MkdirAll(reg.ModulesDir(), os.FileMode(0755))
//		if err != nil {
//			return err
//		}
//	}
//
//	if _, err := os.Stat(reg.LibDir()); os.IsNotExist(err) {
//		err = os.MkdirAll(reg.LibDir(), os.FileMode(0755))
//		if err != nil {
//			return err
//		}
//	}
//
//	if _, err := os.Stat(reg.CacheDir()); os.IsNotExist(err) {
//		err = os.MkdirAll(reg.CacheDir(), os.FileMode(0755))
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

func (reg *Registry) TypeString() string {
	if reg.Type == RegistryTypeGlobal {
		return "global"
	} else if reg.Type == RegistryTypeLocal {
		return "local"
	}

	panic("unknown context")
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
		"type":        registryType,
	}))
}

func registryType(L *lua.LState) int {
	reg := checkRegistry(L)
	L.Push(lua.LString(reg.TypeString()))
	return 1
}
