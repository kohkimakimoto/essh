package essh

import (
	"bytes"
	"github.com/yuin/gopher-lua"
	"runtime"
	"strings"
	"text/template"
)

type Driver struct {
	Name   string
	Config *lua.LTable
	Props  map[string]interface{}
	Engine func(*Driver) (string, error)
}

var Drivers map[string]*Driver = map[string]*Driver{}

var (
	BuiltinDefaultDriverName = "default"
)

func NewDriver() *Driver {
	return &Driver{
		Props: map[string]interface{}{},
	}
}

func (driver *Driver) GenerateRunnableContent(task *Task, host *Host) (string, error) {
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
		"GOARCH":  runtime.GOARCH,
		"GOOS":    runtime.GOOS,
		"Debug":   debugFlag,
		"Driver":  driver,
		"Task":    task,
		"Host":    host,
		"Scripts": scripts,
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

func init() {
	// set built-in drivers
	ResetDrivers()
}

const EnvironmentTemplate = `{{define "environment" -}}
export ESSH_TASK_NAME={{.Task.Name | ShellEscape}}
{{if .Host -}}
export ESSH_HOSTNAME={{.Host.Name | ShellEscape}}
export ESSH_HOST_HOSTNAME={{.Host.Name | ShellEscape}}
{{range $i, $kvpair := .Host.SSHConfig -}}
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

func ResetDrivers() {
	Drivers = map[string]*Driver{}

	// default (just concatenate with new line code)
	driver := NewDriver()
	driver.Name = BuiltinDefaultDriverName
	driver.Engine = func(driver *Driver) (string, error) {
		return `
{{template "environment" .}}
{{range $i, $script := .Scripts}}{{$script.code}}
{{end}}`, nil
	}
	Drivers[driver.Name] = driver
}
