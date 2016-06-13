package essh

type Task struct {
	Name        string
	Description string
	// deprecated...
	Configure   func() error
	Prepare     func(task *TaskContext) error
	Driver      string
	Pty         bool
	Script      []map[string]string
	File        string
	Tags        []string
	Parallel    bool
	Privileged  bool
	// Lock is deprecated. use "bash.lock" in `modules/bash/index.lua`
	Lock       bool
	Disabled   bool
	Hidden     bool
	Backend    string
	Registries []string
	Prefix     string
	Context    *Context
}

var Tasks []*Task = []*Task{}

var (
	DefaultPrefixRemote = "[{{.Host.Name}}] "
	DefaultPrefixLocal  = "[Local => {{.Host.Name}}] "
)

const (
	TASK_BACKEND_LOCAL = "local"
	TASK_BACKEND_REMOTE = "remote"
)

func NewTask() *Task {
	return &Task{
		Tags: []string{},
		Registries: []string{},
		Backend: TASK_BACKEND_LOCAL,
		Script:  []map[string]string{},
	}
}

func GetTask(name string) *Task {
	for _, task := range Tasks {
		if task.Name == name && !task.Disabled {
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
