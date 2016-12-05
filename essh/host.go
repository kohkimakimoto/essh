package essh

import (
	"bytes"
	"github.com/yuin/gopher-lua"
	"sort"
	"strings"
	"text/template"
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
	Registry    *Registry
	Job         *Job
	LValues              map[string]lua.LValue
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

func SortedTags() []string {
	tagsMap := map[string]string{}
	tags := []string{}

	for _, host := range NewHostQuery().GetHostsOrderByName() {
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

func removeHostInGlobalSpace(host *Host) {
	h := Hosts[host.Name]
	if h == host {
		delete(Hosts, host.Name)
	}
}