package essh

import (
	"fmt"
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

func esshGroup(L *lua.LState) int {
	first := L.CheckTable(1)
	j := registerGroup(L)
	setupGroup(L, j, first)
	L.Push(newLGroup(L, j))
	return 1
}

func registerGroup(L *lua.LState) *Group {
	j := NewGroup()
	return j
}

func setupGroup(L *lua.LState, group *Group, config *lua.LTable) {
	// guarantee evaluating a key/value dictionary at first.
	config.ForEach(func(k, v lua.LValue) {
		if kstr, ok := toString(k); ok {
			updateGroup(L, group, kstr, v)
		}
	})

	config.ForEach(func(k, v lua.LValue) {
		if _, ok := toString(k); ok {
			return
		} else if _, ok := toFloat64(k); ok {
			// set a host, task or driver
			lv, ok := v.(*lua.LUserData)
			if !ok {
				panic(fmt.Sprintf("expected userdata (host, task or driver) but got '%v'\n", v))
			}

			switch resource := lv.Value.(type) {
			case *Host:
				if group.Type != GroupTypeUndefined && group.Type != GroupTypeHosts {
					panic("group can use only one type of resources. \n")
				}

				// set host table data
				if group.LValues["hosts"] == nil {
					group.LValues["hosts"] = L.NewTable()
				}
				hosts, ok := toLTable(group.LValues["hosts"])
				if !ok {
					panic("broken 'hosts' table")
				}
				host := L.NewTable()
				resource.MapLValuesToLTable(host)
				hosts.RawSetString(resource.Name, host)

				// register host object
				group.RegisterHost(resource)
			case *Task:
				if group.Type != GroupTypeUndefined && group.Type != GroupTypeTasks {
					panic("group can use only one type of resources. \n")
				}

				// set task table data
				if group.LValues["tasks"] == nil {
					group.LValues["tasks"] = L.NewTable()
				}
				tasks, ok := toLTable(group.LValues["tasks"])
				if !ok {
					panic("broken 'tasks' table")
				}
				task := L.NewTable()
				resource.MapLValuesToLTable(task)
				tasks.RawSetString(resource.Name, task)

				// register task object
				group.RegisterTask(resource)
			case *Driver:
				if group.Type != GroupTypeUndefined && group.Type != GroupTypeDrivers {
					panic("group can use only one type of resources. \n")
				}

				// set task table data
				if group.LValues["drivers"] == nil {
					group.LValues["drivers"] = L.NewTable()
				}
				drivers, ok := toLTable(group.LValues["drivers"])
				if !ok {
					panic("broken 'drivers' table")
				}
				driver := L.NewTable()
				resource.MapLValuesToLTable(driver)
				drivers.RawSetString(resource.Name, driver)

				// register task object
				group.RegisterDriver(resource)
			default:
				panic(fmt.Sprintf("expected host, task or driver but got '%v'\n", resource))
			}
		} else {
			panic("invalid operation\n")
		}
	})

	applyGroupDefaultValues(L, group)
}

func updateGroup(L *lua.LState, group *Group, key string, value lua.LValue) {
	group.LValues[key] = value

	switch key {
	case "hosts":
		if tb, ok := toLTable(value); ok {
			if group.Type != GroupTypeUndefined && group.Type != GroupTypeHosts {
				panic("group can use only one type of resources. \n")
			}

			// initialize
			group.Hosts = map[string]*Host{}

			tb.ForEach(func(k, v lua.LValue) {
				name, ok := toString(k)
				if !ok {
					panic(fmt.Sprintf("expected string of host's name but got '%v'\n", k))
				}

				config, ok := toLTable(v)
				if !ok {
					panic(fmt.Sprintf("expected table of host's config but got '%v'\n", v))
				}

				h := registerHost(L, name)
				setupHost(L, h, config)
				group.RegisterHost(h)
			})
		} else {
			panic(fmt.Sprintf("expected table but got '%v'\n", value))
		}
	case "tasks":
		if tb, ok := toLTable(value); ok {
			if group.Type != GroupTypeUndefined && group.Type != GroupTypeTasks {
				panic("group can use only one type of resources. \n")
			}

			// initialize
			group.Tasks = map[string]*Task{}

			tb.ForEach(func(k, v lua.LValue) {
				name, ok := toString(k)
				if !ok {
					panic(fmt.Sprintf("expected string of task's name but got '%v'\n", k))
				}

				config, ok := toLTable(v)
				if !ok {
					panic(fmt.Sprintf("expected table of task's config but got '%v'\n", v))
				}

				t := registerTask(L, name)
				setupTask(L, t, config)
				group.RegisterTask(t)
			})
		} else {
			panic(fmt.Sprintf("expected table but got '%v'\n", value))
		}
	case "drivers":
		if tb, ok := toLTable(value); ok {
			if group.Type != GroupTypeUndefined && group.Type != GroupTypeDrivers {
				panic("group can use only one type of resources. \n")
			}

			// initialize
			group.Drivers = map[string]*Driver{
				DefaultDriverName: DefaultDriver,
			}

			tb.ForEach(func(k, v lua.LValue) {
				name, ok := toString(k)
				if !ok {
					panic(fmt.Sprintf("expected string of driver's name but got '%v'\n", k))
				}

				config, ok := toLTable(v)
				if !ok {
					panic(fmt.Sprintf("expected table of driver's config but got '%v'\n", v))
				}

				d := registerDriver(L, name)
				setupDriver(L, d, config)
				group.RegisterDriver(d)
			})
		} else {
			panic(fmt.Sprintf("expected table but got '%v'\n", value))
		}
	}
}

func applyGroupDefaultValues(L *lua.LState, group *Group) {

	isSkipKey := func(k string) bool {
		if k == "hosts" || k == "tasks" || k == "drivers" {
			return true
		} else {
			return false
		}
	}

	switch group.Type {
	case GroupTypeHosts:
		for _, h := range group.Hosts {
			for k, v := range group.LValues {
				if !isSkipKey(k) {
					if h.LValues[k] == nil {
						updateHost(L, h, k, v)
					}
				}
			}
		}
	case GroupTypeTasks:
		for _, t := range group.Tasks {
			for k, v := range group.LValues {
				if !isSkipKey(k) {
					if t.LValues[k] == nil {
						updateTask(L, t, k, v)
					}
				}
			}
		}
	case GroupTypeDrivers:
		for _, d := range group.Drivers {
			for k, v := range group.LValues {
				if !isSkipKey(k) {
					if d.LValues[k] == nil {
						updateDriver(L, d, k, v)
					}
				}
			}
		}
	}
}

const LGroupClass = "Group*"

func registerGroupClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LGroupClass)
	mt.RawSetString("__call", L.NewFunction(groupCall))
	mt.RawSetString("__index", L.NewFunction(groupIndex))
	mt.RawSetString("__newindex", L.NewFunction(groupNewindex))
}

func newLGroup(L *lua.LState, group *Group) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = group
	L.SetMetatable(ud, L.GetTypeMetatable(LGroupClass))
	return ud
}

func checkGroup(L *lua.LState) *Group {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Group); ok {
		return v
	}
	L.ArgError(1, "Group object expected")
	return nil
}

func groupCall(L *lua.LState) int {
	group := checkGroup(L)
	tb := L.CheckTable(2)

	setupGroup(L, group, tb)

	return 0
}

func groupIndex(L *lua.LState) int {
	group := checkGroup(L)
	index := L.CheckString(2)

	switch index {
	default:
		v, ok := group.LValues[index]
		if v == nil || !ok {
			v = lua.LNil
		}
		L.Push(v)
	}

	return 1
}

func groupNewindex(L *lua.LState) int {
	panic("unsupport to override group's properties")

	return 0
}
