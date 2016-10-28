package essh

type HostQuery struct {
	TagsFilter []string
}

func NewHostQuery() *HostQuery {
	return &HostQuery{
		TagsFilter: []string{},
	}
}

func (hostQuery *HostQuery) AppendTag(tag string) {
	hostQuery.TagsFilter = append(hostQuery.TagsFilter, tag)
}

func (hostQuery *HostQuery) GetHosts() []*Host {
	hosts := SortedHosts()
	if len(hostQuery.TagsFilter) == 0 {
		return hosts
	}

	for _, tag := range hostQuery.TagsFilter {
		hosts = hostQuery.filterHosts(hosts, tag)
	}

	return hosts
}

func (hostQuery *HostQuery) filterHosts(hosts []*Host, filterdTag string) []*Host {
	newHosts := []*Host{}
	for _, host := range hosts {
		for _, tag := range host.Tags {
			if tag == filterdTag {
				newHosts = append(newHosts, host)
			}
		}
	}

	return newHosts
}