package essh

import (
	"bytes"
	"github.com/yuin/gopher-lua"
	"runtime"
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

func (driver *Driver) GenerateRunnableContent(task *Task) (string, error) {
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
		"ShellEscape": func(str string) string {
			return ShellEscape(str)
		},
	}

	dict := map[string]interface{}{
		"GOARCH":  runtime.GOARCH,
		"GOOS":    runtime.GOOS,
		"Debug":   debugFlag,
		"Driver":  driver,
		"Task":    task,
		"Scripts": scripts,
	}

	tmpl, err := template.New("T").Funcs(funcMap).Parse(templateText)
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

func ResetDrivers() {
	Drivers = map[string]*Driver{}

	// default (just concatenate with new line code)
	driver := NewDriver()
	driver.Name = BuiltinDefaultDriverName
	driver.Engine = func(driver *Driver) (string, error) {
		return `{{range $i, $script := .Scripts}}{{$script.code}}
{{end}}`, nil
	}
	Drivers[driver.Name] = driver
}
