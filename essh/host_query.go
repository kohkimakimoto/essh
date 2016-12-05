package essh

import (
	"sort"
)

type HostQuery struct {
	Datasource map[string]*Host
	Selections []string
	Filters    []string
}

func NewHostQuery() *HostQuery {
	return &HostQuery{
		Datasource:  Hosts,
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

