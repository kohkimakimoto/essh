package zssh

import (
	"github.com/yuin/gopher-lua"
	"text/template"
	"bytes"
	"fmt"
	"sort"
)

type Host struct {
	Name string
	Config *lua.LTable
	Hooks map[string]func() error
	Description string
	Hidden bool
}

func (h *Host) Values() []map[string]interface{} {

	values := []map[string]interface{}{}

	var names []string

	h.Config.ForEach(func(k lua.LValue, v lua.LValue) {
		if keystr, ok := toString(k); ok {
			names = append(names, keystr)
		}

	})

	sort.Strings(names)

	for _, name := range names {
		lvalue := h.Config.RawGetString(name)
		value := map[string]interface{}{name: toGoValue(lvalue)}
		values = append(values, value)
	}

	return values
}

var Hosts []*Host = []*Host{}

func GetHost(hostname string) *Host {
	for _, host := range Hosts {
		if host.Name == hostname {
			return host
		}
	}

	return nil
}

var hostsTemplate = `# Generated from '{{.ConfigFile}}' by using https://github.com/kohkimakimoto/zssh
# Don't edit this file manually.
{{range $i, $host := .Hosts}}
Host {{$host.Name}}{{range $ii, $value := $host.Values}}{{range $k, $v := $value}}
    {{$k}} {{$v}}{{end}}{{end}}
{{end}}
`

func GenHostsConfig() ([]byte, error) {
	tmpl, err := template.New("T").Parse(hostsTemplate)
	if err != nil {
		return nil, err
	}

	input := map[string]interface{}{"Hosts": Hosts, "ConfigFile": ConfigFile}
	var b bytes.Buffer
	if err := tmpl.Execute(&b, input); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// This code refers to https://github.com/yuin/gluamapper/blob/master/gluamapper.go
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
			ret := make(map[interface{}]interface{})
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
