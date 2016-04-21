package essh

import (
	"path/filepath"
	"os"
)

type Context struct {
	DataDir       string
	LoadedModules map[string]*Module
	Type int
}

const (
	ContextTypeUserData = 0
	ContextTypeWorkingData = 1

)

func (ctx *Context) ModulesDir() string {
	return filepath.Join(ctx.DataDir, "modules")
}

func (ctx *Context) LockDir() string {
	return filepath.Join(ctx.DataDir, "lock")
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

	if _, err := os.Stat(ctx.LockDir()); os.IsNotExist(err) {
		err = os.MkdirAll(ctx.LockDir(), os.FileMode(0755))
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
	if ctx.Type == ContextTypeUserData {
		return "UserData"
	} else if ctx.Type == ContextTypeWorkingData {
		return "WorkingData"
	}

	return "Unknown"
}