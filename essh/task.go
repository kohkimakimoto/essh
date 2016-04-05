package essh

type Task struct {
	Name        string
	Description string
	Prepare     func(task *TaskContext) error
	Tty         bool
	Script      string
	On          []string
	Parallel    bool
	Privileged  bool
	Prefix      string
}

var Tasks []*Task = []*Task{}

func GetTask(name string) *Task {
	for _, task := range Tasks {
		if task.Name == name {
			return task
		}
	}
	return nil
}

type TaskContext struct {
	Task *Task
	Payload string
}

func NewTaskContext(task *Task, payload string) *TaskContext {
	return &TaskContext{
		Task: task,
		Payload: payload,
	}
}