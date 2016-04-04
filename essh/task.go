package essh

type Task struct {
	Name        string
	Description string
	Prepare     func(task *Task, payload string) error
	Tty         bool
	Script      string
	On          []string
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
