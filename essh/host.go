package essh

import (
	"bytes"
	"fmt"
	"github.com/yuin/gopher-lua"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

type Host struct {
	Name                 string
	Description          string
	Props                map[string]string
	HooksBeforeConnect   []interface{}
	HooksAfterConnect    []interface{}
	HooksAfterDisconnect []interface{}
	Hidden               bool
	Tags                 []string
	SSHConfig            map[string]string
	Registry             *Registry
	Group                *Group
	LValues              map[string]lua.LValue
	// If you define same name hosts in multi time, stores it in layered structure that uses Parent and Child.
	Parent *Host
	Child  *Host
}

var Hosts map[string]*Host

func NewHost() *Host {
	return &Host{
		Props:                map[string]string{},
		HooksBeforeConnect:   []interface{}{},
		HooksAfterConnect:    []interface{}{},
		HooksAfterDisconnect: []interface{}{},
		Tags:                 []string{},
		SSHConfig:            map[string]string{},
		LValues:              map[string]lua.LValue{},
	}
}

func (h *Host) MapLValuesToLTable(tb *lua.LTable) {
	for key, value := range h.LValues {
		tb.RawSetString(key, value)
	}
}

func (h *Host) SortedSSHConfig() []map[string]string {
	values := []map[string]string{}

	var names []string

	for name, _ := range h.SSHConfig {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		v := h.SSHConfig[name]
		value := map[string]string{name: v}
		values = append(values, value)
	}

	return values
}

func (h *Host) DescriptionOrDefault() string {
	if h.Description == "" {
		return h.Name + " host"
	}

	return h.Description
}

var hostsTemplate = `{{range $i, $host := .Hosts -}}
Host {{$host.Name}}{{range $ii, $param := $host.SortedSSHConfig}}{{range $k, $v := $param}}
    {{$k}} {{$v}}{{end}}{{end}}

{{end -}}`

func GenHostsConfig(enabledHosts []*Host) ([]byte, error) {
	tmpl, err := template.New("T").Parse(hostsTemplate)
	if err != nil {
		return nil, err
	}

	input := map[string]interface{}{"Hosts": enabledHosts}
	var b bytes.Buffer
	if err := tmpl.Execute(&b, input); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func GetTags(hosts map[string]*Host) []string {
	tagsMap := map[string]string{}
	tags := []string{}

	for _, host := range hosts {
		for _, t := range host.Tags {
			if _, exists := tagsMap[t]; !exists {
				tagsMap[t] = t
				tags = append(tags, t)
			}
		}
	}

	sort.Strings(tags)

	return tags
}

func HostnameAlignString(host *Host, hosts []*Host) func(string) string {
	var maxlen int
	for _, h := range hosts {
		size := len(h.Name)
		if maxlen < size {
			maxlen = size
		}
	}

	var namelen = len(host.Name)
	return func(s string) string {
		diff := maxlen - namelen
		return strings.Repeat(s, 1+diff)
	}
}

func esshHost(L *lua.LState) int {
	value := L.CheckAny(1)
	if tb, ok := toLTable(value); ok {

		hostsTb := L.NewTable()
		tb.ForEach(func(k, v lua.LValue) {
			name, ok := toString(k)
			if !ok {
				panic(fmt.Sprintf("expected string of host's name but got '%v'\n", k))
			}

			config, ok := toLTable(v)
			if !ok {
				panic(fmt.Sprintf("expected table of host's config but got '%v'\n", v))
			}

			h := registerHost(L, name)
			setupHost(L, h, config)
			hostsTb.RawSetString(name, newLHost(L, h))
		})

		L.Push(hostsTb)
		return 1
	} else if name, ok := toString(value); ok {
		if L.GetTop() == 1 {
			// object or DSL style
			h := registerHost(L, name)
			L.Push(newLHost(L, h))

			return 1
		} else if L.GetTop() == 2 {
			// function style
			tb := L.CheckTable(2)
			h := registerHost(L, name)
			setupHost(L, h, tb)
			L.Push(newLHost(L, h))

			return 1
		} else {
			panic("host requires 1 or 2 arguments")
		}
	} else {
		panic(fmt.Sprintf("expected table or string but got '%v'\n", value))
	}
}

func registerHost(L *lua.LState, name string) *Host {
	if debugFlag {
		fmt.Printf("[essh debug] register host: %s\n", name)
	}

	h := NewHost()
	h.Name = name
	h.Registry = CurrentRegistry

	if host := Hosts[h.Name]; host != nil {
		// detect same name host
		h.Child = host
		host.Parent = h
	}

	Hosts[h.Name] = h

	return h
}

func setupHost(L *lua.LState, h *Host, config *lua.LTable) {
	config.ForEach(func(k, v lua.LValue) {
		if kstr, ok := toString(k); ok {
			updateHost(L, h, kstr, v)
		}
	})
}

func updateHost(L *lua.LState, h *Host, key string, value lua.LValue) {
	h.LValues[key] = value

	var firstChar rune
	for _, c := range key {
		firstChar = c
		break
	}

	if unicode.IsUpper(firstChar) {
		if valuestr, ok := toString(value); ok {
			h.SSHConfig[key] = valuestr
			return
		}

		panic("SSH property must be string")
	}

	switch key {
	case "props":
		if propsTb, ok := toLTable(value); ok {
			// initialize
			h.Props = map[string]string{}

			propsTb.ForEach(func(propsKey lua.LValue, propsValue lua.LValue) {
				propsKeyStr, ok := toString(propsKey)
				if !ok {
					L.RaiseError("props table's key must be a string: %v", propsKey)
				}
				propsValueStr, ok := toString(propsValue)
				if !ok {
					L.RaiseError("props table's value must be a string: %v", propsValue)
				}

				h.Props[propsKeyStr] = propsValueStr
			})
		} else {
			panic("invalid value of a host's field '" + key + "'.")
		}
	case "hooks_before_connect":
		if tb, ok := toLTable(value); ok {
			maxn := tb.MaxN()
			hooks := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				hooks = append(hooks, toGoValue(tb.RawGetInt(i)))
			}

			h.HooksBeforeConnect = hooks
		} else {
			panic("invalid value of a host's field '" + key + "'.")
		}
	case "hooks_after_connect":
		if tb, ok := toLTable(value); ok {
			maxn := tb.MaxN()
			hooks := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				hooks = append(hooks, toGoValue(tb.RawGetInt(i)))
			}

			h.HooksAfterConnect = hooks
		} else {
			panic("invalid value of a host's field '" + key + "'.")
		}
	case "hooks_after_disconnect":
		if tb, ok := toLTable(value); ok {
			maxn := tb.MaxN()
			hooks := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				hooks = append(hooks, toGoValue(tb.RawGetInt(i)))
			}

			h.HooksAfterDisconnect = hooks
		} else {
			panic("invalid value of a host's field '" + key + "'.")
		}
	case "description":
		if descStr, ok := toString(value); ok {
			h.Description = descStr
		} else {
			panic("invalid value of a host's field '" + key + "'.")
		}

	case "hidden":
		if hiddenBool, ok := toBool(value); ok {
			h.Hidden = hiddenBool
		} else {
			panic("invalid value of a host's field '" + key + "'.")
		}

	case "tags":
		if tagsTb, ok := toLTable(value); ok {
			// initialize
			h.Tags = []string{}

			tagsTb.ForEach(func(_ lua.LValue, v lua.LValue) {
				if vs, ok := toString(v); ok {
					h.Tags = append(h.Tags, vs)
				} else {
					L.RaiseError("unsupported format of tags.")
				}
			})
		} else {
			panic("invalid value of a host's field '" + key + "'.")
		}

	default:
		panic("unsupported host's field '" + key + "'.")

	}
}

const LHostClass = "Host*"

func registerHostClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LHostClass)
	mt.RawSetString("__call", L.NewFunction(hostCall))
	mt.RawSetString("__index", L.NewFunction(hostIndex))
	mt.RawSetString("__newindex", L.NewFunction(hostNewindex))
}

func newLHost(L *lua.LState, host *Host) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = host
	L.SetMetatable(ud, L.GetTypeMetatable(LHostClass))
	return ud
}

func checkHost(L *lua.LState) *Host {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Host); ok {
		return v
	}
	L.ArgError(1, "Host object expected")
	return nil
}

func hostCall(L *lua.LState) int {
	host := checkHost(L)
	tb := L.CheckTable(2)

	setupHost(L, host, tb)

	L.Push(L.CheckUserData(1))
	return 1
}

func hostIndex(L *lua.LState) int {
	host := checkHost(L)
	index := L.CheckString(2)

	if index == "name" {
		L.Push(L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LString(host.Name))
			return 1
		}))
		return 1
	}

	v, ok := host.LValues[index]
	if v == nil || !ok {
		v = lua.LNil
	}

	L.Push(v)
	return 1
}

func hostNewindex(L *lua.LState) int {
	host := checkHost(L)
	index := L.CheckString(2)
	value := L.CheckAny(3)

	updateHost(L, host, index, value)

	return 0
}
