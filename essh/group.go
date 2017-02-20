package essh

import (
	"github.com/yuin/gopher-lua"
)

type GroupType int

const (
	GroupTypeUndefined GroupType = iota
	GroupTypeHosts
	GroupTypeTasks
	GroupTypeDrivers
)

type Group struct {
	Type    GroupType
	Hosts   map[string]*Host
	Tasks   map[string]*Task
	Drivers map[string]*Driver
	LValues map[string]lua.LValue
}

func NewGroup() *Group {
	return &Group{
		Type:    GroupTypeUndefined,
		Hosts:   map[string]*Host{},
		Tasks:   map[string]*Task{},
		Drivers: map[string]*Driver{},
		LValues: map[string]lua.LValue{},
	}
}

func (group *Group) RegisterHost(host *Host) {
	group.Type = GroupTypeHosts

	group.Hosts[host.Name] = host
	host.Group = group
}

func (group *Group) RegisterTask(task *Task) {
	group.Type = GroupTypeTasks

	group.Tasks[task.Name] = task
	task.Group = group
}

func (group *Group) RegisterDriver(driver *Driver) {
	group.Type = GroupTypeDrivers

	group.Drivers[driver.Name] = driver
	driver.Group = group
}
