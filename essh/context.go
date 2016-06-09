package essh

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

type Context struct {
	Key           string
	DataDir       string
	LoadedModules map[string]*Module
	Type          int
}

// alias
type Registry Context

const (
	ContextTypeGlobal = 0
	ContextTypeLocal = 1
)

var CurrentContext *Context
var ContextMap map[string]*Context = map[string]*Context{}

func NewContext(dataDir string, contextType int) *Context {
	ctx := &Context{
		Key:           fmt.Sprintf("%x", sha256.Sum256([]byte(dataDir))),
		DataDir:       dataDir,
		LoadedModules: map[string]*Module{},
		Type:          contextType,
	}

	return ctx
}

func (ctx *Context) ModulesDir() string {
	return filepath.Join(ctx.DataDir, "modules")
}

func (ctx *Context) TmpDir() string {
	return filepath.Join(ctx.DataDir, "tmp")
}

func (ctx *Context) MkDirs() error {
	if _, err := os.Stat(ctx.ModulesDir()); os.IsNotExist(err) {
		err = os.MkdirAll(ctx.ModulesDir(), os.FileMode(0755))
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(ctx.TmpDir()); os.IsNotExist(err) {
		err = os.MkdirAll(ctx.TmpDir(), os.FileMode(0755))
		if err != nil {
			return err
		}
	}

	return nil
}

func (ctx *Context) TypeString() string {
	if ctx.Type == ContextTypeGlobal {
		return "global"
	} else if ctx.Type == ContextTypeLocal {
		return "local"
	}

	return "Unknown"
}
