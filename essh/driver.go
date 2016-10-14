package essh

import (
	"bytes"
	"fmt"
	"github.com/yuin/gopher-lua"
	"runtime"
	"strings"
	"text/template"
)

type Driver struct {
	Name    string
	Props   map[string]interface{}
	Engine  func(*Driver) (string, error)
	LValues map[string]lua.LValue
}

var DefaultDriver *Driver

var (
	BuiltinDefaultDriverName = "default"
)

func NewDriver() *Driver {
	return &Driver{
		Props:   map[string]interface{}{},
		LValues: map[string]lua.LValue{},
	}
}

func (driver *Driver) GenerateRunnableContent(sshConfigPath string, task *Task, host *Host) (string, error) {
	for key, value := range driver.LValues {
		driver.Props[key] = toGoValue(value)
	}

	if driver.Engine == nil {
		return "", fmt.Errorf("invalid driver '%s'. The engine was not defined.", driver.Name)
	}

	templateText, err := driver.Engine(driver)
	if err != nil {
		return "", err
	}

	scripts := []map[string]string{}
	if task.File != "" {
		tContent, err := GetContentFromPath(task.File)
		if err != nil {
			return "", err
		}
		scripts = append(scripts, map[string]string{"code": string(tContent)})
	} else {
		scripts = task.Script
	}

	funcMap := template.FuncMap{
		"ShellEscape":  ShellEscape,
		"ToUpper":      strings.ToUpper,
		"ToLower":      strings.ToLower,
		"EnvKeyEscape": EnvKeyEscape,
	}

	dict := map[string]interface{}{
		"GOARCH":        runtime.GOARCH,
		"GOOS":          runtime.GOOS,
		"Debug":         debugFlag,
		"Driver":        driver,
		"Task":          task,
		"Host":          host,
		"Scripts":       scripts,
		"SSHConfigPath": sshConfigPath,
	}

	baseTempl, err := template.New("base").Funcs(funcMap).Parse(templateText)
	if err != nil {
		return "", err
	}

	tmpl, err := baseTempl.Parse(EnvironmentTemplate)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, dict)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func FindDriverInRegistry(name string, registry *Registry) *Driver {
	if registry != nil {
		if driver, ok := registry.Drivers[name]; ok {
			return driver
		}
	} else {
		if driver, ok := LocalRegistry.Drivers[name]; ok {
			return driver
		}
		if driver, ok := GlobalRegistry.Drivers[name]; ok {
			return driver
		}
	}

	return nil
}

func init() {
	// set built-in drivers
	// default (just concatenate with new line code)
	driver := NewDriver()
	driver.Name = BuiltinDefaultDriverName
	driver.Engine = func(driver *Driver) (string, error) {
		return `
{{template "environment" .}}
{{range $i, $script := .Scripts}}{{$script.code}}
{{end}}`, nil
	}
	DefaultDriver = driver
}

const EnvironmentTemplate = `{{define "environment" -}}
export ESSH_TASK_NAME={{.Task.Name | ShellEscape}}
export ESSH_SSH_CONFIG={{.SSHConfigPath}}
export ESSH_DEBUG="{{if .Debug}}1{{end}}"

{{if .Host -}}
export ESSH_HOSTNAME={{.Host.Name | ShellEscape}}
export ESSH_HOST_HOSTNAME={{.Host.Name | ShellEscape}}
{{range $i, $kvpair := .Host.SortedSSHConfig -}}
{{range $key, $value := $kvpair -}}
export ESSH_HOST_SSH_{{$key | ToUpper}}={{$value | ShellEscape }}
{{end -}}
{{end -}}
{{range $key, $value := .Host.Props -}}
export ESSH_HOST_PROPS_{{$key | ToUpper | EnvKeyEscape}}={{$value | ShellEscape }}
{{end -}}
{{range $i, $value := .Host.Tags -}}
export ESSH_HOST_TAGS_{{$value | ToUpper | EnvKeyEscape}}=1
{{end -}}
{{end -}}
{{end}}
`
