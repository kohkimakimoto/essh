package essh

import (
	"sort"
	"strings"
)

type TaskQuery struct {
	Datasource map[string]*Task
}

func NewTaskQuery() *TaskQuery {
	// merge tasks and namespace's tasks
	datasource := Tasks

	if len(Namespaces) > 0 {
		for _, namespace := range Namespaces {

			for _, task := range namespace.Tasks {
				datasource[task.PublicName()] = task
			}
		}
	}

	return &TaskQuery{
		Datasource: datasource,
	}
}

func (taskQuery *TaskQuery) SetDatasource(datasource map[string]*Task) *TaskQuery {
	taskQuery.Datasource = datasource
	return taskQuery
}

func (taskQuery *TaskQuery) GetTasks() []*Task {
	tasks := taskQuery.getTasksList()

	return tasks
}

type NameSortableTasks []*Task

func (t NameSortableTasks) Len() int {
	return len(t)
}

func (t NameSortableTasks) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t NameSortableTasks) Less(i, j int) bool {
	return t[i].PublicName() < t[j].PublicName()
}

func (taskQuery *TaskQuery) GetTasksOrderByName() []*Task {
	tasks := taskQuery.GetTasks()

	sort.Sort(NameSortableTasks(tasks))

	return tasks
}

func (taskQuery *TaskQuery) getTasksList() []*Task {
	tasksSlice := []*Task{}
	for _, task := range taskQuery.Datasource {
		tasksSlice = append(tasksSlice, task)
	}
	return tasksSlice
}

func GetEnabledTask(name string, namespaceName string) *Task {
	if namespaceName != "" && !strings.Contains(name, ":") && namespaceName != DefaultNamespaceName {
		name = namespaceName + ":" + name
	}

	for _, t := range NewTaskQuery().GetTasksOrderByName() {
		if t.PublicName() == name && !t.Disabled {
			return t
		}
	}

	return nil
}
