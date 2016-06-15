package essh

import "sort"

type Task struct {
	Name        string
	Description string
	Configure   func() error
	Prepare     func(task *TaskContext) error
	Driver      string
	Pty         bool
	Script      []map[string]string
	File        string
	On          []string
	Foreach     []string
	Parallel    bool
	Privileged  bool
	// Lock is deprecated. use "bash.lock" in `modules/bash/index.lua`
	Lock     bool
	Disabled bool
	Hidden   bool
	Backend  string
	Prefix   string
	Context  *Context
}

var Tasks map[string]*Task = map[string]*Task{}

var (
	DefaultPrefixRemote = "[{{.Host.Name}}] "
	DefaultPrefixLocal  = "[Local => {{.Host.Name}}] "
)

const (
	TASK_BACKEND_LOCAL  = "local"
	TASK_BACKEND_REMOTE = "remote"
)

func NewTask() *Task {
	return &Task{
		On:      []string{},
		Foreach: []string{},
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
	if len(t.On) >= 1 {
		return true
	} else {
		return false
	}
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
