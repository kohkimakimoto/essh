package essh

import (
	"sort"
)

type HostQuery struct {
	Selections []string
	Filters []string
}

func NewHostQuery() *HostQuery {
	return &HostQuery{
		Selections: []string{},
		Filters: []string{},
	}
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
	hosts := getAllHosts()

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

	sort.Sort(ScopeSortableHosts(hosts))
	sort.Sort(NameSortableHosts(hosts))

	return hosts
}

func (hostQuery *HostQuery) GetPublicHostsOrderByName() []*Host {
	hosts := hostQuery.GetHosts()

	sort.Sort(ScopeSortableHosts(hosts))
	sort.Sort(NameSortableHosts(hosts))

	filteredHosts := []*Host{}
	for _, h := range hosts {
		if !h.Private {
			filteredHosts = append(filteredHosts, h)
		}
	}

	return filteredHosts
}

func (hostQuery *HostQuery) GetSameRegistryHostsOrderByName(registryType int) []*Host {
	hosts := hostQuery.GetHosts()

	sort.Sort(ScopeSortableHosts(hosts))
	sort.Sort(NameSortableHosts(hosts))

	filteredHosts := []*Host{}
	for _, h := range hosts {
		if h.Registry.Type == registryType {
			filteredHosts = append(filteredHosts, h)
		}
	}

	return filteredHosts
}

type ScopeSortableHosts []*Host

func (h ScopeSortableHosts) Len() int {
	return len(h)
}

func (h ScopeSortableHosts) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h ScopeSortableHosts) Less(i, j int) bool {
	return !h[i].Private && h[j].Private
}

func (hostQuery *HostQuery) GetHostsOrderByScopeAndName() []*Host {
	hosts := hostQuery.GetHosts()

	sort.Sort(NameSortableHosts(hosts))
	sort.Sort(ScopeSortableHosts(hosts))

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

func getAllHosts() []*Host {
	hosts := []*Host{}

	for _, host := range GlobalRegistry.Hosts {
		hosts = append(hosts, host)
	}

	for _, host := range LocalRegistry.Hosts {
		hosts = append(hosts, host)
	}

	return hosts
}

