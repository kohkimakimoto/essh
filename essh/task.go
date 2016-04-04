package essh

import (
	"errors"
	"fmt"
)

type Task struct {
	Name        string
	Description string
	Prepare     func(task *Task, payload string) error
	Tty         bool
	Script      string
	On          string
}

func (m *Task) Run() error {

	return nil
}

var Tasks []*Task = []*Task{}

func GetTask(name string) (*Task, error) {
	for _, task := range Tasks {
		if task.Name == name {
			return task, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("not found '%s' task.", name))
}
