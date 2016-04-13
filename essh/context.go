package essh

import (
	"path/filepath"
)

type Context struct {
	DataDir       string
	LoadedModules map[string]*Module
}

func (ctx *Context) ModulesDir() string {
	return filepath.Join(ctx.DataDir, "modules")
}

func (ctx *Context) LockDir() string {
	return filepath.Join(ctx.DataDir, "lock")
}
