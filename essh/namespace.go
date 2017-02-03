package essh

import (
	"github.com/yuin/gopher-lua"
	"sort"
)

type Namespace struct {
	Name        string
	Description string
	// Props       map[string]string
	Hidden  bool
	Prepare func() error
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

func (job *Namespace) DescriptionOrDefault() string {
	if job.Description == "" {
		return job.Name + " job"
	}

	return job.Description
}

func (job *Namespace) RegisterHost(host *Host) {
	job.Hosts[host.Name] = host
	host.Job = job
	removeHostInGlobalSpace(host)
}

func (job *Namespace) RegisterTask(task *Task) {
	job.Tasks[task.Name] = task
	task.Namespace = job
	removeTaskInGlobalSpace(task)
}

func (job *Namespace) RegisterDriver(driver *Driver) {
	job.Drivers[driver.Name] = driver
	driver.Namespace = job
	removeDriverInGlobalSpace(driver)
}

func (job *Namespace) SortedTasks() []*Task {
	names := []string{}
	namesMap := map[string]bool{}
	tasks := []*Task{}

	for name, _ := range job.Tasks {
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

func SortedJobs() []*Namespace {
	names := []string{}
	namesMap := map[string]bool{}
	jobs := []*Namespace{}

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
			jobs = append(jobs, j)
		}
	}

	return jobs
}
