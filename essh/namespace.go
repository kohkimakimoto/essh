package essh

import (
	"github.com/yuin/gopher-lua"
	"sort"
)

type Namespace struct {
	Name        string
	Hosts   map[string]*Host
	Tasks   map[string]*Task
	Drivers map[string]*Driver
	LValues map[string]lua.LValue
}

var Namespaces map[string]*Namespace

var DefaultNamespaceName = "default"

func NewNamespace() *Namespace {
	return &Namespace{
		Hosts: map[string]*Host{},
		Tasks: map[string]*Task{},
		Drivers: map[string]*Driver{
			DefaultDriverName: DefaultDriver,
		},
		LValues: map[string]lua.LValue{},
	}
}

func (namespace *Namespace) RegisterHost(host *Host) {
	namespace.Hosts[host.Name] = host
	host.Namespace = namespace
	removeHostInGlobalSpace(host)
}

func (namespace *Namespace) RegisterTask(task *Task) {
	namespace.Tasks[task.Name] = task
	task.Namespace = namespace
	removeTaskInGlobalSpace(task)
}

func (namespace *Namespace) RegisterDriver(driver *Driver) {
	namespace.Drivers[driver.Name] = driver
	driver.Namespace = namespace
	removeDriverInGlobalSpace(driver)
}

func (namespace *Namespace) SortedTasks() []*Task {
	names := []string{}
	namesMap := map[string]bool{}
	tasks := []*Task{}

	for name, _ := range namespace.Tasks {
		if namesMap[name] {
			// already registerd to names
			continue
		}

		names = append(names, name)
		namesMap[name] = true
	}

	sort.Strings(names)

	for _, name := range names {
		if t, ok := Tasks[name]; ok {
			tasks = append(tasks, t)
		}
	}

	return tasks
}

func SortedNamespaces() []*Namespace {
	names := []string{}
	namesMap := map[string]bool{}
	namespaces := []*Namespace{}

	for name, _ := range Namespaces {
		if namesMap[name] {
			// already registerd to names
			continue
		}

		names = append(names, name)
		namesMap[name] = true
	}

	sort.Strings(names)

	for _, name := range names {
		if j, ok := Namespaces[name]; ok {
			namespaces = append(namespaces, j)
		}
	}

	return namespaces
}
