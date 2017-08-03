package essh

import (
	"github.com/yuin/gopher-lua"
)

type Task struct {
	Name        string
	Description string
	Props       map[string]string
	Prepare     func() error
	Driver      string
	Pty         bool
	Script      []map[string]string
	File        string
	Backend     string
	Targets     []string
	Filters     []string
	Parallel    bool
	Privileged  bool
	User        string
	// deprecated? use only hidden?
	Disabled  bool
	Hidden    bool
	Prefix    string
	UsePrefix bool
	Registry  *Registry
	Namespace *Namespace
	Group     *Group
	Args      []string
	LValues   map[string]lua.LValue
	Parent    *Task
	Child     *Task
}

var Tasks map[string]*Task

var DefaultTaskName = "default"

var (
	DefaultPrefixLocal  = `[local:{{.Host.Name}}]{{HostnameAlignString " "}}`
	DefaultPrefixRemote = `[remote:{{.Host.Name}}]{{HostnameAlignString " "}}`
)

const (
	TASK_BACKEND_LOCAL  = "local"
	TASK_BACKEND_REMOTE = "remote"
)

func NewTask() *Task {
	return &Task{
		Targets: []string{},
		Filters: []string{},
		Backend: TASK_BACKEND_LOCAL,
		Script:  []map[string]string{},
		Args:    []string{},
		LValues: map[string]lua.LValue{},
	}
}

func (t *Task) MapLValuesToLTable(tb *lua.LTable) {
	for key, value := range t.LValues {
		tb.RawSetString(key, value)
	}
}

func (t *Task) PublicName() string {
	if t.Namespace != nil && t.Namespace.Name != DefaultNamespaceName {
		if t.Name == DefaultTaskName {
			return t.Namespace.Name
		}

		return t.Namespace.Name + ":" + t.Name
	}

	return t.Name
}

func (t *Task) IsRemoteTask() bool {
	if t.Backend == TASK_BACKEND_REMOTE {
		return true
	} else {
		return false
	}
}

func (t *Task) TargetsSlice() []string {
	if len(t.Targets) >= 1 {
		return t.Targets
	}

	return []string{}
}

func (t *Task) FiltersSlice() []string {
	if len(t.Filters) >= 1 {
		return t.Filters
	}

	return []string{}
}

func (t *Task) DescriptionOrDefault() string {
	if t.Description == "" {
		return t.Name + " task"
	}

	return t.Description
}

func removeTaskInGlobalSpace(task *Task) {
	t := Tasks[task.Name]
	if t == task {
		if t.Child != nil {
			newTask := t.Child
			Tasks[newTask.Name] = newTask
			newTask.Parent = nil
		} else {
			delete(Tasks, t.Name)
		}
	}
}
