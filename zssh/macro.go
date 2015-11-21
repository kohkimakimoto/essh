package zssh

import(
	"bufio"
	"fmt"
	"errors"
	"os"
	"sync"
	"strings"
)

type Macro struct {
	Name        string
	Description string
	Parallel    bool
	Tty         bool
	Command     string                                               `gluamapper:"-"`
	CommandFunc func(host *Host) (string, error)     `gluamapper:"-"`
	Confirm     bool                                                 `gluamapper:"-"`
	ConfirmText string                                               `gluamapper:"-"`
	OnServers   []string                                             `gluamapper:"-"`
	OnTags      map[string][]string                                  `gluamapper:"-"`
	RunLocally  bool                                                 `gluamapper:"-"`
}

var Macros []*Macro = []*Macro{}

func GetMacro(name string) (*Macro, error) {
	for _, macro := range Macros {
		if macro.Name == name {
			return macro, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("not found '%s' task.", name))
}

func (m *Macro) Run() error {
	if m.Confirm {
		if !m.AskYesOrNo() {
			return nil
		}
	}

	if m.RunLocally {
		// run locally
		script, err := m.Script(nil)
		if err != nil {
			return err
		}
		if script == "" {
			return nil
		}

		err = Run(script)
		if err != nil {
			return err
		}

		return nil
	}

	hosts, err := m.TargetHosts()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for _, host := range hosts {
		script, err := m.Script(host)
		if err != nil {
			return err
		}

		var cmd string
		if m.Tty {
			cmd = "echo '" + script + "' | ssh -t -t " + host.Name+ " bash -se"
		} else {
			cmd = "echo '" + script + "' | ssh " + host.Name+ " bash -se"
		}

		if m.Parallel {
			wg.Add(1)
			go func(host *Host, cmd string) {
				host.Run(cmd)
				wg.Done()
			}(host, cmd)
		} else {
			// ignore err
			host.Run(cmd)
		}
	}
	wg.Wait()


	return nil
}

func (m *Macro) Script(host *Host) (string, error) {
	script := m.Command
	if m.CommandFunc != nil {
		ret, err := m.CommandFunc(host)
		if err != nil {
			return "", err
		}

		script = ret
	}

	return script, nil
}

func (m *Macro) TargetHosts() ([]*Host, error) {

	var targets = []*Host{}

	hosts := HostsByTags(m.OnTags)
	targets = append(targets, hosts...)

	for _, v := range m.OnServers {
		host := GetHost(v)
		if host != nil {
			targets = append(targets, host)
		}
	}

	// remove duplication
	var ret = []*Host{}
	var checkDup = map[string]bool{}
	for _, v := range targets {
		if _, ok := checkDup[v.Name]; !ok {
			checkDup[v.Name] = true
			ret = append(ret, v)
		}
	}

	return ret, nil
}

func (m *Macro) AskYesOrNo() bool {
	var msg string
	if m.ConfirmText != "" {
		msg = fmt.Sprintf("%s %s: ", m.ConfirmText, FgYB("[y/N]"))
	} else {
		msg = fmt.Sprintf("Are you sure you want to run the [%s] task? %s: ", FgY(m.Name), FgYB("[y/N]"))
	}

	fmt.Printf(msg)
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
	str = strings.Trim(str, "\r\n")

	if str == "y" {
		return true
	} else {
		return false
	}
}

func Escape(s string) string {
	return "'" + strings.Replace(s, "'", "'\"'\"'", -1) + "'"
}
