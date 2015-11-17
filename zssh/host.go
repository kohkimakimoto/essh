package zssh

import (
	"bytes"
	"github.com/yuin/gopher-lua"
	"sort"
	"text/template"
)

type Host struct {
	Name        string
	Config      *lua.LTable
	Hooks       map[string]func() error
	Description string
	Hidden      bool
	Tags        map[string][]string
}

const LHostClass = "ZsshHost*"

var Hosts []*Host = []*Host{}

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
