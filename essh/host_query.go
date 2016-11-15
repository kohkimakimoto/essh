package essh

type HostQuery struct {
	Filters []string
}

func NewHostQuery() *HostQuery {
	return &HostQuery{
		Filters: []string{},
	}
}

func (hostQuery *HostQuery) AppendFilter(filter string) {
	hostQuery.Filters = append(hostQuery.Filters, filter)
}

func (hostQuery *HostQuery) GetHosts() []*Host {
	hosts := SortedHosts()
	if len(hostQuery.Filters) == 0 {
		return hosts
	}

	for _, filter := range hostQuery.Filters {
		hosts = hostQuery.filterHosts(hosts, filter)
	}

	return hosts
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