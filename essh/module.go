package essh

import (
	"github.com/yuin/gopher-lua"
)

type Module struct {
	// Name is url that is used as go-getter src.
	// examples:
	//   github.com/aaa/bbb
	//   git::github.com/aaa/bbb.git
	Name string
	// LValues
	LValues map[string]lua.LValue
}

var RootModules []*Module = []*Module{}

func NewModule(name string) *Module {
	return &Module{
		Name:    name,
		LValues: map[string]lua.LValue{},
	}
}
