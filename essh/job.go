package essh

import (
	"github.com/yuin/gopher-lua"
	"sort"
)

type Job struct {
	Name        string
	Description string
	Env         map[string]string
	Config      *lua.LTable
	Prepare     func() error
	Hosts       map[string]*Host
	Tasks       map[string]*Task
	Drivers     map[string]*Driver
	LValues     map[string]lua.LValue
}

var Jobs map[string]*Job

var DefaultJobName = "default"

func NewJob() *Job {
	return &Job{
		Hosts: map[string]*Host{},
		Tasks: map[string]*Task{},
		Drivers: map[string]*Driver{
			DefaultDriverName: DefaultDriver,
		},
		LValues: map[string]lua.LValue{},
	}
}

func (job *Job) DescriptionOrDefault() string {
	if job.Description == "" {
		return job.Name + " job"
	}

	return job.Description
}

func (job *Job) RegisterHost(host *Host) {
	job.Hosts[host.Name] = host
	host.Job = job
	removeHostInGlobalSpace(host)
}

func (job *Job) RegisterTask(task *Task) {
	job.Tasks[task.Name] = task
	task.Job = job
	removeTaskInGlobalSpace(task)
}

func (job *Job) RegisterDriver(driver *Driver) {
	job.Drivers[driver.Name] = driver
	driver.Job = job
	removeDriverInGlobalSpace(driver)
}

func (job *Job) SortedTasks() []*Task {
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

func SortedJobs() []*Job {
	names := []string{}
	namesMap := map[string]bool{}
	jobs := []*Job{}

	for name, _ := range Jobs {
		if namesMap[name] {
			// already registerd to names
			continue
		}

		names = append(names, name)
		namesMap[name] = true
	}

	sort.Strings(names)

	for _, name := range names {
		if j, ok := Jobs[name]; ok {
			jobs = append(jobs, j)
		}
	}

	return jobs
}
