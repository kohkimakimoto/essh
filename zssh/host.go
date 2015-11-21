package zssh

import (
	"bytes"
	"github.com/yuin/gopher-lua"
	"sort"
	"text/template"
	"os/exec"
	"runtime"
	"os"
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

func (h *Host) Run(command string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	realWriter := &RealWriter{
		NewLine: true,
		Prefix:  "[" + h.Name + "]: ",
	}
	outWriter := &writer{
		realWriter: realWriter,
		Type:       1,
	}
	errWriter := &writer{
		realWriter: realWriter,
		Type:       2,
	}

	cmd.Stdout = outWriter
	cmd.Stderr = errWriter
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
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

func HostsByTags(tags map[string][]string) []*Host {
	var hosts = []*Host{}

	for _, host := range Hosts {
		alreadyExists := false
		for tag, values := range tags {
			hostTagValues := host.Tags[tag]
			for _, serverTagValue := range hostTagValues {
				for _, value := range values {
					if serverTagValue == value && !alreadyExists {
						hosts = append(hosts, host)
						alreadyExists = true
					}
				}
			}
		}
	}

	return hosts
}