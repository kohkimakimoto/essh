package essh

import (
	"github.com/yuin/gopher-lua"
	"sort"
)

type Task struct {
	Name        string
	Description string
	Prepare     func() error
	Driver      string
	Pty         bool
	Script      []map[string]string
	File        string
	Backend     string
	Targets     []string
	Parallel    bool
	Privileged  bool
	Disabled    bool
	Hidden      bool
	Prefix      string
	UsePrefix  bool
	Registry    *Registry

	LValues map[string]lua.LValue
}

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
		Backend: TASK_BACKEND_LOCAL,
		Script:  []map[string]string{},
		LValues: map[string]lua.LValue{},
	}
}

func SortedTasks() []*Task {
	names := []string{}
	namesMap := map[string]bool{}
	tasks := []*Task{}

	for name, _ := range GlobalRegistry.Tasks {
		if namesMap[name] {
			// already registerd to names
			continue
		}

		names = append(names, name)
		namesMap[name] = true
	}

	for name, _ := range LocalRegistry.Tasks {
		if namesMap[name] {
			// already registerd to names
			continue
		}

		names = append(names, name)
		namesMap[name] = true
	}

	sort.Strings(names)

	for _, name := range names {
		if t, ok := GlobalRegistry.Tasks[name]; ok {
			tasks = append(tasks, t)
		}

		if t, ok := LocalRegistry.Tasks[name]; ok {
			tasks = append(tasks, t)
		}
	}

	return tasks
}

func GetEnabledTask(name string) *Task {
	for _, t := range SortedTasks() {
		if t.Name == name && !t.Disabled {
			return t
		}
	}

	return nil
}

func (t *Task) IsRemoteTask() bool {
	if t.Backend == TASK_BACKEND_REMOTE {
		return true
	} else {
		return false
	}
}

func (t *Task) Context() *Registry {
	return t.Registry
}

func (t *Task) TargetsSlice() []string {
	if len(t.Targets) >= 1 {
		return t.Targets
	}

	return []string{}
}

func (t *Task) DescriptionOrDefault() string {
	if t.Description == "" {
		return t.Name + " task"
	}

	return t.Description
}

type TaskContext struct {
	Task    *Task
	Payload string
}

func NewTaskContext(task *Task, payload string) *TaskContext {
	return &TaskContext{
		Task:    task,
		Payload: payload,
	}
}
