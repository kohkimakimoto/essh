package essh

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGetAllHosts(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer func() {
		os.Remove(tmpDir)
	}()

	tmpDir2, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer func() {
		os.Remove(tmpDir2)
	}()

	GlobalRegistry = NewRegistry(tmpDir, RegistryTypeGlobal)
	LocalRegistry = NewRegistry(tmpDir2, RegistryTypeLocal)

	h := NewHost()
	h.Name = "essh-host-global-web01"
	h.Registry = GlobalRegistry
	h.Tags = []string{
		"web",
	}
	GlobalRegistry.Hosts[h.Name] = h

	h = NewHost()
	h.Name = "essh-host-global-web02"
	h.Registry = GlobalRegistry
	h.Tags = []string{
		"web",
	}
	GlobalRegistry.Hosts[h.Name] = h

	h = NewHost()
	h.Name = "essh-host-global-db01"
	h.Registry = GlobalRegistry
	h.Tags = []string{
		"db",
	}
	GlobalRegistry.Hosts[h.Name] = h

	h = NewHost()
	h.Name = "essh-host-local-web01"
	h.Registry = LocalRegistry
	h.Tags = []string{
		"web",
	}
	LocalRegistry.Hosts[h.Name] = h

	hosts := getAllHosts()
	if len(hosts) != 4 {
		t.Errorf("host number should be 4 but %d", len(hosts))
	}
}

func TestHostQueryGetHostsForSelectedHosts(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer func() {
		os.Remove(tmpDir)
	}()

	tmpDir2, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer func() {
		os.Remove(tmpDir2)
	}()

	GlobalRegistry = NewRegistry(tmpDir, RegistryTypeGlobal)
	LocalRegistry = NewRegistry(tmpDir2, RegistryTypeLocal)

	h := NewHost()
	h.Name = "essh-host-global-web01"
	h.Registry = GlobalRegistry
	h.Tags = []string{
		"web",
	}
	GlobalRegistry.Hosts[h.Name] = h

	h = NewHost()
	h.Name = "essh-host-global-web02"
	h.Registry = GlobalRegistry
	h.Tags = []string{
		"web",
	}
	GlobalRegistry.Hosts[h.Name] = h

	h = NewHost()
	h.Name = "essh-host-global-db01"
	h.Registry = GlobalRegistry
	h.Tags = []string{
		"db",
	}
	GlobalRegistry.Hosts[h.Name] = h

	h = NewHost()
	h.Name = "essh-host-local-web01"
	h.Registry = LocalRegistry
	h.Tags = []string{
		"web",
	}
	LocalRegistry.Hosts[h.Name] = h

	// pattern1
	q := NewHostQuery()
	q.AppendSelection("essh-host-global-web01")
	q.AppendSelection("web")
	hosts := q.GetHosts()
	t.Logf("selected hosts\n")
	for _, host := range hosts {
		t.Logf("host = %v\n", host.Name)
	}

	if len(hosts) != 3 {
		t.Errorf("host number should be 3 but %d", len(hosts))
	}

	// pattern2
	q = NewHostQuery()
	q.AppendSelection("essh-host-global-web01")
	q.AppendSelection("db")
	hosts = q.GetHosts()

	t.Logf("selected hosts\n")
	for _, host := range hosts {
		t.Logf("host = %v\n", host.Name)
	}

	if len(hosts) != 2 {
		t.Errorf("host number should be 2 but %d", len(hosts))
	}

	// pattern3
	q = NewHostQuery()
	q.AppendFilter("web")
	hosts = q.GetHosts()

	t.Logf("selected hosts\n")
	for _, host := range hosts {
		t.Logf("host = %v\n", host.Name)
	}

	if len(hosts) != 3 {
		t.Errorf("host number should be 3 but %d", len(hosts))
	}

	// pattern4
	q = NewHostQuery()
	q.AppendSelection("web")
	q.AppendFilter("db")
	hosts = q.GetHosts()

	t.Logf("selected hosts\n")
	for _, host := range hosts {
		t.Logf("host = %v\n", host.Name)
	}

	if len(hosts) != 0 {
		t.Errorf("host number should be 0 but %d", len(hosts))
	}

	// pattern5
	q = NewHostQuery()
	q.AppendSelection("web")
	q.AppendFilter("essh-host-global-web01")
	hosts = q.GetHosts()

	t.Logf("selected hosts\n")
	for _, host := range hosts {
		t.Logf("host = %v\n", host.Name)
	}

	if len(hosts) != 1 {
		t.Errorf("host number should be 0 but %d", len(hosts))
	}

	// pattern6
	q = NewHostQuery()
	q.AppendSelection("web")
	q.AppendSelection("db")
	hosts = q.GetHostsOrderByName()

	t.Logf("selected hosts\n")
	for _, host := range hosts {
		t.Logf("host = %v\n", host.Name)
	}

	if len(hosts) != 4 {
		t.Errorf("host number should be 0 but %d", len(hosts))
	}

	if hosts[0].Name != "essh-host-global-db01" {
		t.Errorf("the first host should be essh-host-global-db01 but %d", hosts[0].Name)
	}
}
