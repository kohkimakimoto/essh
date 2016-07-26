package essh

import "sort"

type Task struct {
	Name        string
	Description string
	// Configure deprecated.
	Prepare func(task *TaskContext) error
	Driver  string
	Pty     bool
	Script  []map[string]string
	File    string

	Backend string
	Targets []string

	Parallel   bool
	Privileged bool
	// Lock is deprecated. use "bash.lock" in `modules/bash/index.lua`
	Lock     bool
	Disabled bool
	Hidden   bool

	Prefix   string
	Registry *Registry
}

var Tasks map[string]*Task = map[string]*Task{}

var (
	DefaultPrefixLocal  = `local - {{.Host.Name}}:{{HostnameAlignString " "}}`
	DefaultPrefixRemote = `remote - {{.Host.Name}}:{{HostnameAlignString " "}}`
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
	}
}

func SortedTasks() []*Task {
	names := []string{}
	tasks := []*Task{}

	for name, _ := range Tasks {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		tasks = append(tasks, Tasks[name])
	}

	return tasks
}

func GetEnabledTask(name string) *Task {
	if task, ok := Tasks[name]; ok {
		if !task.Disabled {
			return task
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
