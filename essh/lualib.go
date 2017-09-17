package essh

import (
	"fmt"
	"github.com/cjoudrey/gluahttp"
	"github.com/kohkimakimoto/gluaenv"
	"github.com/kohkimakimoto/gluafs"
	"github.com/kohkimakimoto/gluaquestion"
	"github.com/kohkimakimoto/gluatemplate"
	"github.com/kohkimakimoto/gluayaml"
	"github.com/otm/gluash"
	"github.com/yuin/gluare"
	"github.com/yuin/gopher-lua"
	gluajson "layeh.com/gopher-json"
	"net/http"
	"os"
	"path/filepath"
)

func InitLuaState(L *lua.LState) {
	// custom type.
	registerHostClass(L)
	registerTaskClass(L)
	registerDriverClass(L)
	registerHostQueryClass(L)
	registerRegistryClass(L)
	registerGroupClass(L)
	registerModuleClass(L)

	// global functions
	L.SetGlobal("host", L.NewFunction(esshHost))
	L.SetGlobal("task", L.NewFunction(esshTask))
	L.SetGlobal("driver", L.NewFunction(esshDriver))
	L.SetGlobal("group", L.NewFunction(esshGroup))
	L.SetGlobal("module", L.NewFunction(esshModule))

	// deprecated. for BC
	L.SetGlobal("import", L.NewFunction(esshImport))

	// modules
	L.PreloadModule("json", gluajson.Loader)
	L.PreloadModule("fs", gluafs.Loader)
	L.PreloadModule("yaml", gluayaml.Loader)
	L.PreloadModule("template", gluatemplate.Loader)
	L.PreloadModule("question", gluaquestion.Loader)
	L.PreloadModule("env", gluaenv.Loader)
	L.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	L.PreloadModule("re", gluare.Loader)
	L.PreloadModule("sh", gluash.Loader)

	// global variables
	lessh := L.NewTable()
	L.SetGlobal("essh", lessh)
	lessh.RawSetString("ssh_config", lua.LNil)
	lessh.RawSetString("version", lua.LString(Version))
	lessh.RawSetString("module", lua.LNil)

	L.SetFuncs(lessh, map[string]lua.LGFunction{
		// aliases global function.
		"host":   esshHost,
		"task":   esshTask,
		"driver": esshDriver,
		"group":  esshGroup,
		"module": esshModule,

		// utility functions
		"debug":            esshDebug,
		"select_hosts":     esshSelectHosts,
		"current_registry": esshCurrentRegistry,
	})
}

func esshDebug(L *lua.LState) int {
	msg := L.CheckString(1)
	if debugFlag {
		fmt.Printf("[essh debug] %s\n", msg)
	}

	return 0
}

func esshImport(L *lua.LState) int {
	name := L.CheckString(1)
	lessh, ok := toLTable(L.GetGlobal("essh"))
	if !ok {
		L.RaiseError("'essh' global variable is broken")
	}
	mod := lessh.RawGetString("package")
	if mod != lua.LNil {
		L.RaiseError("'essh.pkg' is existed. does not support nested pkg importing.")
	}

	pkg := CurrentRegistry.LoadedPackages[name]
	if pkg == nil {
		pkg = NewPackage(name)

		update := updateFlag
		if CurrentRegistry.Type == RegistryTypeGlobal && !withGlobalFlag {
			update = false
		}

		err := pkg.Load(update)
		if err != nil {
			L.RaiseError("%v", err)
		}

		indexFile := pkg.IndexFile()
		if _, err := os.Stat(indexFile); err != nil {
			L.RaiseError("invalid pkg: %v", err)
		}

		// init pkg variable
		modulevar := L.NewTable()
		modulevar.RawSetString("path", lua.LString(filepath.Dir(indexFile)))
		modulevar.RawSetString("import_path", lua.LString(name))
		lessh.RawSetString("package", modulevar)

		if err := L.DoFile(indexFile); err != nil {
			panic(err)
		}
		// remove pkg variable
		lessh.RawSetString("package", lua.LNil)

		// get a pkg return value
		ret := L.Get(-1)
		pkg.Value = ret

		// register loaded pkg.
		CurrentRegistry.LoadedPackages[name] = pkg

		return 1
	}

	L.Push(pkg.Value)
	return 1
}

func esshCurrentRegistry(L *lua.LState) int {
	L.Push(newLRegistry(L, CurrentRegistry))
	return 1
}

// This code inspired by https://github.com/yuin/gluamapper/blob/master/gluamapper.go
func toGoValue(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[string]interface{})
			v.ForEach(func(key, value lua.LValue) {
				keystr := fmt.Sprint(toGoValue(key))
				ret[keystr] = toGoValue(value)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, toGoValue(v.RawGetInt(i)))
			}
			return ret
		}
	default:
		return v
	}
}

func toBool(v lua.LValue) (bool, bool) {
	if lv, ok := v.(lua.LBool); ok {
		return bool(lv), true
	} else {
		return false, false
	}
}

func toString(v lua.LValue) (string, bool) {
	if lv, ok := v.(lua.LString); ok {
		return string(lv), true
	} else {
		return "", false
	}
}

func toMap(v lua.LValue) (map[string]interface{}, bool) {
	if lv, ok := toGoValue(v).(map[string]interface{}); ok {
		return lv, true
	} else {
		return nil, false
	}
}

func toSlice(v lua.LValue) ([]interface{}, bool) {
	gov := toGoValue(v)
	if lv, ok := gov.([]interface{}); ok {
		return lv, true
	} else if lv, ok := gov.(map[string]interface{}); ok {
		if len(lv) == 0 {
			return []interface{}{}, true
		}
		return nil, false
	} else {
		return nil, false
	}
}

func toLFunction(v lua.LValue) (*lua.LFunction, bool) {
	if lv, ok := v.(*lua.LFunction); ok {
		return lv, true
	} else {
		return nil, false
	}
}

func toLTable(v lua.LValue) (*lua.LTable, bool) {
	if lv, ok := v.(*lua.LTable); ok {
		return lv, true
	} else {
		return nil, false
	}
}

func toLUserData(v lua.LValue) (*lua.LUserData, bool) {
	if lv, ok := v.(*lua.LUserData); ok {
		return lv, true
	} else {
		return nil, false
	}
}

func toFloat64(v lua.LValue) (float64, bool) {
	if lv, ok := v.(lua.LNumber); ok {
		return float64(lv), true
	} else {
		return 0, false
	}
}
