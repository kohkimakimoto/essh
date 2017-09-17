package essh

import (
	"github.com/yuin/gopher-lua"
	"sort"
)

type HostQuery struct {
	Datasource map[string]*Host
	Selections []string
	Filters    []string
}

func NewHostQuery() *HostQuery {
	return &HostQuery{
		Datasource: Hosts,
		Selections: []string{},
		Filters:    []string{},
	}
}

func (hostQuery *HostQuery) SetDatasource(datasource map[string]*Host) *HostQuery {
	hostQuery.Datasource = datasource
	return hostQuery
}

func (hostQuery *HostQuery) AppendSelection(selection string) *HostQuery {
	hostQuery.Selections = append(hostQuery.Selections, selection)
	return hostQuery
}

func (hostQuery *HostQuery) AppendSelections(selections []string) *HostQuery {
	hostQuery.Selections = append(hostQuery.Selections, selections...)
	return hostQuery
}

func (hostQuery *HostQuery) AppendFilter(filter string) *HostQuery {
	hostQuery.Filters = append(hostQuery.Filters, filter)
	return hostQuery
}

func (hostQuery *HostQuery) AppendFilters(filters []string) *HostQuery {
	hostQuery.Filters = append(hostQuery.Filters, filters...)
	return hostQuery
}

func (hostQuery *HostQuery) GetHosts() []*Host {
	hosts := hostQuery.getHostsList()

	if len(hostQuery.Selections) == 0 && len(hostQuery.Filters) == 0 {
		return hosts
	}

	hosts = hostQuery.selectHosts(hosts)

	for _, filter := range hostQuery.Filters {
		hosts = hostQuery.filterHosts(hosts, filter)
	}

	return hosts
}

type NameSortableHosts []*Host

func (h NameSortableHosts) Len() int {
	return len(h)
}

func (h NameSortableHosts) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h NameSortableHosts) Less(i, j int) bool {
	return h[i].Name < h[j].Name
}

func (hostQuery *HostQuery) GetHostsOrderByName() []*Host {
	hosts := hostQuery.GetHosts()

	sort.Sort(NameSortableHosts(hosts))

	return hosts
}

func (hostQuery *HostQuery) selectHosts(hosts []*Host) []*Host {
	if len(hostQuery.Selections) == 0 {
		return hosts
	}

	newHosts := []*Host{}
	selections := hostQuery.Selections

	for _, host := range hosts {
		selected := false

	B1:
		for _, selection := range selections {
			if host.Name == selection {
				newHosts = append(newHosts, host)
				selected = true
				break B1
			}
		}

		if selected {
			continue
		}

	B2:
		for _, tag := range host.Tags {
			for _, selection := range selections {
				if tag == selection {
					newHosts = append(newHosts, host)
					break B2
				}
			}
		}
	}

	return newHosts
}

func (hostQuery *HostQuery) filterHosts(hosts []*Host, filter string) []*Host {
	newHosts := []*Host{}
	for _, host := range hosts {
		if host.Name == filter {
			newHosts = append(newHosts, host)
			continue
		}

		for _, tag := range host.Tags {
			if tag == filter {
				newHosts = append(newHosts, host)
				break
			}
		}
	}

	return newHosts
}

func (hostQuery *HostQuery) getHostsList() []*Host {
	hostsSlice := []*Host{}
	for _, host := range hostQuery.Datasource {
		hostsSlice = append(hostsSlice, host)
	}
	return hostsSlice
}

func esshSelectHosts(L *lua.LState) int {
	hostQuery := NewHostQuery()

	if L.GetTop() > 1 {
		panic("select_hosts can receive max 1 argument.")
	}

	if L.GetTop() == 1 {
		value := L.CheckAny(1)
		selections := []string{}

		if selectionsStr, ok := toString(value); ok {
			selections = []string{selectionsStr}
		} else if selectionsSlice, ok := toSlice(value); ok {
			for _, selection := range selectionsSlice {
				if selectionStr, ok := selection.(string); ok {
					selections = append(selections, selectionStr)
				}
			}
		} else {
			panic("select_hosts can receive string or array table of strings.")
		}
		hostQuery.AppendSelections(selections)
	}

	L.Push(newLHostQuery(L, hostQuery))
	return 1
}

const LHostQueryClass = "HostQuery*"

func registerHostQueryClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LHostQueryClass)
	mt.RawSetString("__index", L.NewFunction(hostQueryIndex))
}

func newLHostQuery(L *lua.LState, hostQuery *HostQuery) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = hostQuery
	L.SetMetatable(ud, L.GetTypeMetatable(LHostQueryClass))
	return ud
}

func checkHostQuery(L *lua.LState) *HostQuery {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*HostQuery); ok {
		return v
	}
	L.ArgError(1, "HostQuery object expected")
	return nil
}

func hostQueryIndex(L *lua.LState) int {
	//_ := checkHostQuery(L)
	//_ := L.CheckUserData(1)
	index := L.CheckString(2)

	switch index {
	case "filter":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			hostQuery := checkHostQuery(L)
			ud := L.CheckUserData(1)
			if L.GetTop() != 2 {
				panic("filter must receive max 2 argument.")
			} else {
				filters := []string{}
				value := L.CheckAny(2)
				if filtersStr, ok := toString(value); ok {
					filters = []string{filtersStr}
				} else if filtersSlice, ok := toSlice(value); ok {
					for _, filter := range filtersSlice {
						if filterStr, ok := filter.(string); ok {
							filters = append(filters, filterStr)
						}
					}
				} else {
					panic("filter can receive string or array table of strings.")
				}

				hostQuery.AppendFilters(filters)
			}

			ud.Value = hostQuery
			L.Push(ud)
			return 1
		}))

		return 1
	case "get":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			hostQuery := checkHostQuery(L)

			lhosts := L.NewTable()
			for _, host := range hostQuery.GetHosts() {
				lhost := newLHost(L, host)
				lhosts.Append(lhost)
			}

			L.Push(lhosts)
			return 1
		}))

		return 1
	case "first":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			L.Push(L.NewFunction(func(L *lua.LState) int {
				hostQuery := checkHostQuery(L)

				hosts := hostQuery.GetHosts()
				if len(hosts) > 0 {
					L.Push(newLHost(L, hosts[0]))
					return 1
				}
				L.Push(lua.LNil)
				return 1
			}))
			return 1
		}))
		L.Push(lua.LNil)
		return 1
	default:
		L.Push(lua.LNil)
		return 1
	}
}
