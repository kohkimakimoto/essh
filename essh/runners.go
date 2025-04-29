package essh

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"text/template"

	"github.com/Songmu/wrapcommander"
	"github.com/sevir/essh/support/color"
	lua "github.com/yuin/gopher-lua"
)

func runTask(config string, task *Task, args []string, L *lua.LState) error {
	if debugFlag {
		fmt.Printf("[essh debug] run task: %s\n", task.Name)
		fmt.Printf("[essh debug] task's args: %v\n", args)
	}

	if task.Registry != nil {
		// change current registry
		CurrentRegistry = task.Registry
	}

	// compose args
	argstb := L.NewTable()
	for i := 0; i < len(args); i++ {
		L.RawSet(argstb, lua.LNumber(i+1), lua.LString(args[i]))
	}
	updateTask(L, task, "args", argstb)

	if task.Prepare != nil {
		if debugFlag {
			fmt.Printf("[essh debug] run task's prepare function.\n")
		}

		err := task.Prepare()
		if err != nil {
			return err
		}
	}

	// get target hosts.
	if task.IsRemoteTask() {
		// run remotely.
		var hosts []*Host
		if len(task.TargetsSlice()) == 0 {
			hosts = []*Host{}
		} else {
			hosts = NewHostQuery().
				AppendSelections(task.TargetsSlice()).
				AppendFilters(task.FiltersSlice()).
				GetHostsOrderByName()
		}

		if len(hosts) == 0 {
			return fmt.Errorf("There are not hosts to run the command. you must specify the valid hosts.")
		}

		// see https://github.com/kohkimakimoto/essh/issues/38
		//// handle stdin
		stdinChs := make([]chan ([]byte), len(hosts))
		for i, _ := range hosts {
			stdinChs[i] = make(chan []byte, 256)
		}
		go func() {
			processStdin(stdinChs)
		}()

		wg := &sync.WaitGroup{}
		m := new(sync.Mutex)
		for i, host := range hosts {
			if task.Parallel {
				wg.Add(1)
				go func(host *Host) {
					err := runRemoteTaskScript(config, task, host, hosts, stdinChs[i], m)
					if err != nil {
						fmt.Fprintf(os.Stderr, color.FgRB("essh error: %v\n", err))
						panic(err)
					}

					wg.Done()
				}(host)
			} else {
				err := runRemoteTaskScript(config, task, host, hosts, stdinChs[i], m)
				if err != nil {
					return err
				}
			}
		}
		wg.Wait()
	} else {
		// run locally.
		var hosts []*Host
		if len(task.TargetsSlice()) == 0 {
			hosts = []*Host{}
		} else {
			hosts = NewHostQuery().
				AppendSelections(task.TargetsSlice()).
				AppendFilters(task.FiltersSlice()).
				GetHostsOrderByName()
		}

		if len(task.Targets) >= 1 && len(hosts) == 0 {
			return fmt.Errorf("There are not hosts to run the command. you must specify the valid hosts.")
		}

		wg := &sync.WaitGroup{}
		m := new(sync.Mutex)

		if len(hosts) == 0 {
			// local no host task
			// This pattern should run just exec. should not use magic to pipe stdin to multi targets.
			err := runLocalTaskScript(config, task, nil, hosts, nil, m)
			if err != nil {
				return err
			}
			return nil
		}

		// see https://github.com/kohkimakimoto/essh/issues/38
		// handle stdin
		stdinChs := make([]chan ([]byte), len(hosts))
		for i, _ := range hosts {
			stdinChs[i] = make(chan []byte, 256)
		}
		go func() {
			processStdin(stdinChs)
		}()

		for i, host := range hosts {
			if task.Parallel {
				wg.Add(1)
				go func(host *Host) {
					err := runLocalTaskScript(config, task, host, hosts, stdinChs[i], m)
					if err != nil {
						fmt.Fprintf(os.Stderr, color.FgRB("essh error: %v\n", err))
						panic(err)
					}

					wg.Done()
				}(host)
			} else {
				err := runLocalTaskScript(config, task, host, hosts, stdinChs[i], m)
				if err != nil {
					return err
				}
			}
		}
		wg.Wait()
	}

	return nil
}

func runRemoteTaskScript(sshConfigPath string, task *Task, host *Host, hosts []*Host, stdinCh chan []byte, m *sync.Mutex) error {
	// setup ssh command args
	var sshCommandArgs []string
	if task.Pty {
		sshCommandArgs = []string{"-t", "-t", "-F", sshConfigPath, host.Name}
	} else {
		sshCommandArgs = []string{"-F", sshConfigPath, host.Name}
	}

	// generate commands by using driver
	if task.Driver == "" {
		task.Driver = DefaultDriverName
	}

	driver := Drivers[task.Driver]
	if driver == nil {
		return fmt.Errorf("invalid driver name '%s'", task.Driver)
	}

	if debugFlag {
		fmt.Printf("[essh debug] driver: %s \n", driver.Name)
	}

	var script string
	content, err := driver.GenerateRunnableContent(sshConfigPath, task, host)
	if err != nil {
		return err
	}
	script += content

	if task.User != "" {
		script = "sudo -u " + ShellEscape(task.User) + " bash -l -c " + ShellEscape(script)
	} else if task.Privileged {
		script = "sudo bash -l -c " + ShellEscape(script)
	}

	sshCommandArgs = append(sshCommandArgs, "bash", "-c", ShellEscape(script))

	if task.SSHOptions != nil {
		sshCommandArgs = append(task.SSHOptions, sshCommandArgs[:]...)
	}

	cmd := exec.Command("ssh", sshCommandArgs[:]...)
	if debugFlag {
		fmt.Printf("[essh debug] real ssh command: %v \n", cmd.Args)
	}

	prefix := ""
	if task.UsePrefix {
		prefixTmp := task.Prefix
		if prefixTmp == "" {
			if task.IsRemoteTask() {
				prefixTmp = DefaultPrefixRemote
			} else {
				prefixTmp = DefaultPrefixLocal
			}
		}

		funcMap := template.FuncMap{
			"ShellEscape":         ShellEscape,
			"ToUpper":             strings.ToUpper,
			"ToLower":             strings.ToLower,
			"EnvKeyEscape":        EnvKeyEscape,
			"HostnameAlignString": HostnameAlignString(host, hosts),
		}

		dict := map[string]interface{}{
			"Host": host,
			"Task": task,
		}
		tmpl, err := template.New("T").Funcs(funcMap).Parse(prefixTmp)
		if err != nil {
			return err
		}
		var b bytes.Buffer
		err = tmpl.Execute(&b, dict)
		if err != nil {
			return err
		}

		prefix = b.String()
	}

	// cmd.Stdin = os.Stdin

	// see https://github.com/kohkimakimoto/essh/issues/38
	if stdinCh == nil {
		cmd.Stdin = os.Stdin
	} else {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		go handleInput(stdinCh, stdin)
	}

	wg := &sync.WaitGroup{}
	if len(hosts) <= 1 && prefix == "" {
		cmd.Stdout = os.Stdout
	} else {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			scanLines(stdout, os.Stdout, prefix, m)
			wg.Done()
		}()
	}

	if len(hosts) <= 1 && prefix == "" {
		cmd.Stderr = os.Stderr
	} else {
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			scanLines(stderr, os.Stderr, prefix, m)
			wg.Done()
		}()
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	wg.Wait()

	return cmd.Wait()
}

func runLocalTaskScript(sshConfigPath string, task *Task, host *Host, hosts []*Host, stdinCh chan []byte, m *sync.Mutex) error {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "bash"
		flag = "-c"
	}

	// generate commands by using driver
	if task.Driver == "" {
		task.Driver = DefaultDriverName
	}

	driver := Drivers[task.Driver]
	if driver == nil {
		return fmt.Errorf("invalid driver name '%s'", task.Driver)
	}

	if debugFlag {
		fmt.Printf("[essh debug] driver: %s \n", driver.Name)
	}

	var script string
	content, err := driver.GenerateRunnableContent(sshConfigPath, task, host)
	if err != nil {
		return err
	}
	script += content

	if task.User != "" {
		script = "cd " + WorkingDir + "\n" + script
		script = "sudo -u " + ShellEscape(task.User) + " bash -l -c " + ShellEscape(script)
	} else if task.Privileged {
		script = "cd " + WorkingDir + "\n" + script
		script = "sudo bash -l -c " + ShellEscape(script)
	}

	cmd := exec.Command(shell, flag, script)
	if debugFlag {
		fmt.Printf("[essh debug] real local command: %v \n", cmd.Args)
	}

	prefix := ""
	if host == nil && task.UsePrefix {
		// simple local task (does not specify the hosts)
		// prevent to use invalid text template.
		// replace prefix string to the string that is not included "{{.Host}}"
		prefix = "[local] "
	} else if task.UsePrefix {
		prefixTmp := task.Prefix
		if prefixTmp == "" {
			if task.IsRemoteTask() {
				prefixTmp = DefaultPrefixRemote
			} else {
				prefixTmp = DefaultPrefixLocal
			}
		}

		funcMap := template.FuncMap{
			"ShellEscape":         ShellEscape,
			"ToUpper":             strings.ToUpper,
			"ToLower":             strings.ToLower,
			"EnvKeyEscape":        EnvKeyEscape,
			"HostnameAlignString": HostnameAlignString(host, hosts),
		}

		dict := map[string]interface{}{
			"Host": host,
			"Task": task,
		}
		tmpl, err := template.New("T").Funcs(funcMap).Parse(prefixTmp)
		if err != nil {
			return err
		}
		var b bytes.Buffer
		err = tmpl.Execute(&b, dict)
		if err != nil {
			return err
		}

		prefix = b.String()
	}

	// cmd.Stdin = os.Stdin

	// see https://github.com/kohkimakimoto/essh/issues/38
	if stdinCh == nil {
		cmd.Stdin = os.Stdin
	} else {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		go handleInput(stdinCh, stdin)
	}

	wg := &sync.WaitGroup{}
	if len(hosts) <= 1 && prefix == "" {
		cmd.Stdout = os.Stdout
	} else {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			scanLines(stdout, os.Stdout, prefix, m)
			wg.Done()
		}()
	}

	if len(hosts) <= 1 && prefix == "" {
		cmd.Stderr = os.Stderr
	} else {
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			scanLines(stderr, os.Stderr, prefix, m)
			wg.Done()
		}()
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	wg.Wait()

	return cmd.Wait()
}

func runSSH(L *lua.LState, config string, args []string) (error, int) {
	// hooks
	hooks := map[string][]interface{}{}

	// Limitation!
	// hooks fires only when the hostname is just specified.
	if len(args) == 1 {
		hostname := args[0]
		if host := Hosts[hostname]; host != nil {
			hooks["before_connect"] = host.HooksBeforeConnect
			hooks["after_disconnect"] = host.HooksAfterDisconnect
			hooks["after_connect"] = host.HooksAfterConnect
		}
	}

	// run before_connect hook
	if before := hooks["before_connect"]; before != nil && len(before) > 0 {
		if debugFlag {
			fmt.Printf("[essh debug] run before_connect hook\n")
		}
		hookScript, err := getHookScript(L, before)
		if err != nil {
			return err, ExitErr
		}
		if debugFlag {
			fmt.Printf("[essh debug] before_connect hook script: %s\n", hookScript)
		}
		if err := runCommand(hookScript); err != nil {
			return err, ExitErr
		}
	}

	// register after_disconnect hook
	defer func() {
		// after hook
		if after := hooks["after_disconnect"]; after != nil && len(after) > 0 {
			if debugFlag {
				fmt.Printf("[essh debug] run after_disconnect hook\n")
			}
			hookScript, err := getHookScript(L, after)
			if err != nil {
				panic(err)
			}
			if debugFlag {
				fmt.Printf("[essh debug] after_disconnect hook script: %s\n", hookScript)
			}
			if err := runCommand(hookScript); err != nil {
				panic(err)
			}
		}
	}()

	// setup ssh command args
	var sshCommandArgs []string

	// run after_connect hook
	if afterConnect := hooks["after_connect"]; afterConnect != nil && len(afterConnect) > 0 {
		hookScript, err := getHookScript(L, afterConnect)
		if err != nil {
			return err, ExitErr
		}

		script := hookScript
		script += "\nexec $SHELL\n"

		hasTOption := false
		for _, arg := range args {
			if arg == "-t" {
				hasTOption = true
			}
		}

		if hasTOption {
			sshCommandArgs = []string{"-F", config}
		} else {
			sshCommandArgs = []string{"-t", "-F", config}
		}

		sshCommandArgs = append(sshCommandArgs, args[:]...)
		sshCommandArgs = append(sshCommandArgs, script)
	} else {
		sshCommandArgs = []string{"-F", config}
		sshCommandArgs = append(sshCommandArgs, args[:]...)
	}

	// execute ssh commmand
	cmd := exec.Command("ssh", sshCommandArgs[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if debugFlag {
		fmt.Printf("[essh debug] real ssh command: %v \n", cmd.Args)
	}

	err := cmd.Run()
	ex := wrapcommander.ResolveExitCode(err)

	// Running as a wrapper of ssh command suppress printing error.
	// Printing error is essh's behavior. ssh does not have it.
	return nil, ex
}
