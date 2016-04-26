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

const (
	ContextTypeUserData    = 0
	ContextTypeWorkingData = 1
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

func (ctx *Context) LockDir() string {
	return filepath.Join(ctx.DataDir, "lock")
}

func (ctx *Context) TmpDir() string {
	return filepath.Join(ctx.DataDir, "tmp")
}

func (ctx *Context) VarDir() string {
	return filepath.Join(ctx.DataDir, "var")
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

	if _, err := os.Stat(ctx.VarDir()); os.IsNotExist(err) {
		err = os.MkdirAll(ctx.VarDir(), os.FileMode(0755))
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
