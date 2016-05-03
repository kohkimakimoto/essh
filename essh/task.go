package essh

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
	Lock        bool
	Disabled    bool
	Hidden      bool
	Prefix      string
	Context     *Context
}

var Tasks []*Task = []*Task{}

var (
	DefaultPrefixRemote = "[{{.Host.Name}}] "
	DefaultPrefixLocal  = "[Local => {{.Host.Name}}] "
)

func NewTask() *Task {
	return &Task{
		On:      []string{},
		Foreach: []string{},
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
