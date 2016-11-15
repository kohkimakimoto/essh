package essh

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Songmu/wrapcommander"
	fatihColor "github.com/fatih/color"
	"github.com/kardianos/osext"
	"github.com/kohkimakimoto/essh/support/color"
	"github.com/kohkimakimoto/essh/support/helper"
	"github.com/yuin/gopher-lua"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"text/template"
)

// system configurations.
var (
	SystemWideConfigFile         string
	SystemWideOverrideConfigFile string
	UserConfigFile               string
	UserOverrideConfigFile       string
	UserDataDir                  string
	WorkingDirConfigFile         string
	WorkingDirOverrideConfigFile string
	WorkingDataDir               string
	WorkingDir                   string
	Executable                   string
)

// flags
var (
	versionFlag            bool
	helpFlag               bool
	printFlag              bool
	colorFlag              bool
	noColorFlag            bool
	debugFlag              bool
	hostsFlag              bool
	quietFlag              bool
	allFlag                bool
	tagsFlag               bool
	tasksFlag              bool
	genFlag                bool
	updateFlag             bool
	withGlobalFlag         bool
	cleanAllFlag           bool
	cleanModulesFlag       bool
	cleanTmpFlag           bool
	zshCompletionModeFlag  bool
	zshCompletionFlag      bool
	zshCompletionHostsFlag bool
	zshCompletionTagsFlag  bool
	zshCompletionTasksFlag bool
	bashCompletionFlag     bool
	aliasesFlag            bool
	execFlag               bool
	fileFlag               bool
	prefixFlag             bool
	parallelFlag           bool
	privilegedFlag         bool
	ptyFlag                bool
	workindDirVar          string
	configVar              string
	selectVar              []string = []string{}
	targetVar              []string = []string{}
	backendVar             string
	prefixStringVar        string
	driverVar              string
)

const (
	ExitErr = 1
)

func Start() (exitStatus int) {
	defer func() {
		if e := recover(); e != nil {
			exitStatus = ExitErr
			if zshCompletionModeFlag && !debugFlag {
				// suppress printing error in running completion code.
				return
			}

			printError(e)
		}
	}()

	if os.Getenv("ESSH_DEBUG") != "" {
		debugFlag = true
	}

	if len(os.Args) == 1 {
		printUsage()
		return
	}

	osArgs := os.Args[1:]
	args := []string{}
	doesNotParseOption := false

	for {
		if len(osArgs) == 0 {
			break
		}

		arg := osArgs[0]

		if doesNotParseOption {
			// restructure args to remove essh options.
			args = append(args, arg)
		} else if arg == "--print" {
			printFlag = true
		} else if arg == "--version" {
			versionFlag = true
		} else if arg == "--help" {
			helpFlag = true
		} else if arg == "--color" {
			colorFlag = true
		} else if arg == "--no-color" {
			noColorFlag = true
		} else if arg == "--debug" {
			debugFlag = true
		} else if arg == "--hosts" {
			hostsFlag = true
		} else if arg == "--quiet" {
			quietFlag = true
		} else if arg == "--all" {
			allFlag = true
		} else if arg == "--tasks" {
			tasksFlag = true
		} else if arg == "--select" {
			if len(osArgs) < 2 {
				printError("--select reguires an argument.")
				return ExitErr
			}
			selectVar = append(selectVar, osArgs[1])
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--select=") {
			selectVar = append(selectVar, strings.Split(arg, "=")[1])
		} else if arg == "--tags" {
			tagsFlag = true
		} else if arg == "--gen" {
			genFlag = true
		} else if arg == "--update" {
			updateFlag = true
		} else if arg == "--clean-modules" {
			cleanModulesFlag = true
		} else if arg == "--clean-tmp" {
			cleanTmpFlag = true
		} else if arg == "--clean-all" {
			cleanAllFlag = true
		} else if arg == "--with-global" {
			withGlobalFlag = true
		} else if arg == "--zsh-completion" {
			zshCompletionFlag = true
			zshCompletionModeFlag = true
		} else if arg == "--zsh-completion-hosts" {
			zshCompletionHostsFlag = true
			zshCompletionModeFlag = true
		} else if arg == "--zsh-completion-tags" {
			zshCompletionTagsFlag = true
			zshCompletionModeFlag = true
		} else if arg == "--zsh-completion-tasks" {
			zshCompletionTasksFlag = true
			zshCompletionModeFlag = true
		} else if arg == "--bash-completion" {
			// TODO
			bashCompletionFlag = true
		} else if arg == "--aliases" {
			aliasesFlag = true
		} else if arg == "--working-dir" {
			if len(osArgs) < 2 {
				printError("--working-dir reguires an argument.")
				return ExitErr
			}
			workindDirVar = osArgs[1]
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--working-dir=") {
			workindDirVar = strings.Split(arg, "=")[1]
		} else if arg == "--config" {
			if len(osArgs) < 2 {
				printError("--config reguires an argument.")
				return ExitErr
			}
			configVar = osArgs[1]
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--config=") {
			configVar = strings.Split(arg, "=")[1]
		} else if arg == "--exec" {
			execFlag = true
		} else if arg == "--privileged" {
			privilegedFlag = true
		} else if arg == "--parallel" {
			parallelFlag = true
		} else if arg == "--prefix" {
			prefixFlag = true
		} else if arg == "--prefix-string" {
			if len(osArgs) < 2 {
				printError("--prefix-string reguires an argument.")
				return ExitErr
			}
			prefixStringVar = osArgs[1]
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--prefix-string=") {
			prefixStringVar = strings.Split(arg, "=")[1]
		} else if arg == "--driver" {
			if len(osArgs) < 2 {
				printError("--driver reguires an argument.")
				return ExitErr
			}
			driverVar = osArgs[1]
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--driver=") {
			driverVar = strings.Split(arg, "=")[1]
		} else if arg == "--target" {
			if len(osArgs) < 2 {
				printError("--target reguires an argument.")
				return ExitErr
			}
			targetVar = append(targetVar, osArgs[1])
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--target=") {
			targetVar = append(targetVar, strings.Split(arg, "=")[1])
		} else if arg == "--backend" {
			if len(osArgs) < 2 {
				printError("--backend reguires an argument.")
				return ExitErr
			}
			backendVar = osArgs[1]
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--backend=") {
			backendVar = strings.Split(arg, "=")[1]
		} else if arg == "--file" {
			fileFlag = true
		} else if arg == "--pty" {
			ptyFlag = true
		} else if arg == "--" {
			doesNotParseOption = true
			// to behave same ssh. pass the `--` to the ssh.
			args = append(args, arg)
		} else {
			// restructure args to remove essh options.
			args = append(args, arg)
		}

		osArgs = osArgs[1:]
	}

	if colorFlag {
		fatihColor.NoColor = false
	}

	if noColorFlag {
		fatihColor.NoColor = true
	}

	if os.Getenv("ESSH_DEBUG") != "" {
		debugFlag = true
	}

	if workindDirVar != "" {
		err := os.Chdir(workindDirVar)
		if err != nil {
			printError(err)
			return ExitErr
		}
	}

	// decide the wokingDirConfigFile
	wd, err := os.Getwd()
	if err != nil {
		printError(fmt.Errorf("couldn't get working dir %v\n", err))
		return ExitErr
	}

	WorkingDir = wd
	WorkingDataDir = filepath.Join(wd, ".essh")
	WorkingDirConfigFile = filepath.Join(wd, "esshconfig.lua")

	workingDirConfigFileBasename := filepath.Base(WorkingDirConfigFile)
	workingDirConfigFileDir := filepath.Dir(WorkingDirConfigFile)
	workingDirConfigFileBasenameExtension := filepath.Ext(workingDirConfigFileBasename)
	workingDirConfigFileName := workingDirConfigFileBasename[0 : len(workingDirConfigFileBasename)-len(workingDirConfigFileBasenameExtension)]

	WorkingDirOverrideConfigFile = filepath.Join(workingDirConfigFileDir, workingDirConfigFileName+"_override"+workingDirConfigFileBasenameExtension)

	// overwrite config file path by --config option.
	if configVar != "" {
		if filepath.IsAbs(configVar) {
			WorkingDirConfigFile = configVar
		} else {
			WorkingDirConfigFile = filepath.Join(wd, configVar)
		}
	}

	if helpFlag {
		printUsage()
		return
	}

	if cleanAllFlag || cleanModulesFlag || cleanTmpFlag {
		err := removeRegistryData()
		if err != nil {
			printError(err)
			return ExitErr
		}
		return
	}

	if versionFlag {
		fmt.Printf("%s (%s)\n", Version, CommitHash)
		return
	}

	if zshCompletionFlag {
		s, err := sprintByTemplate(ZSH_COMPLETION)
		if err != nil {
			printError(err)
			return ExitErr
		}

		fmt.Print(s)
		return
	}

	if aliasesFlag {
		s, err := sprintByTemplate(ALIASES_CODE)
		if err != nil {
			printError(err)
			return ExitErr
		}

		fmt.Print(s)
		return
	}

	if bashCompletionFlag {
		fmt.Print(BASH_COMPLETION)
		return
	}

	// set up the lua state.
	L := lua.NewState()
	defer L.Close()
	InitLuaState(L)

	if debugFlag {
		fmt.Printf("[essh debug] init lua state\n")
	}

	// generate temporary ssh config file
	tmpFile, err := ioutil.TempFile("", "essh.ssh_config.")
	if err != nil {
		printError(err)
		return ExitErr
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())

		if debugFlag {
			fmt.Printf("[essh debug] deleted config file: %s \n", tmpFile.Name())
		}
	}()
	temporarySSHConfigFile := tmpFile.Name()

	if debugFlag {
		fmt.Printf("[essh debug] generated config file: %s \n", temporarySSHConfigFile)
	}

	lessh, ok := toLTable(L.GetGlobal("essh"))
	if !ok {
		printError(fmt.Errorf("essh must be a table"))
		return ExitErr
	}

	// set temporary ssh config file path
	lessh.RawSetString("ssh_config", lua.LString(temporarySSHConfigFile))

	// user context
	CurrentRegistry = NewRegistry(UserDataDir, RegistryTypeGlobal)
	GlobalRegistry = CurrentRegistry

	if err := CurrentRegistry.MkDirs(); err != nil {
		printError(err)
		return ExitErr
	}

	// load system wide config
	if _, err := os.Stat(SystemWideConfigFile); err == nil {
		if debugFlag {
			fmt.Printf("[essh debug] loading config file: %s\n", SystemWideConfigFile)
		}

		if err := CurrentRegistry.MkDirs(); err != nil {
			printError(err)
			return ExitErr

		}

		if err := L.DoFile(SystemWideConfigFile); err != nil {
			printError(err)
			return ExitErr
		}

		if debugFlag {
			fmt.Printf("[essh debug] loaded config file: %s\n", SystemWideConfigFile)
		}
	}

	// load per-user wide config
	if _, err := os.Stat(UserConfigFile); err == nil {
		if debugFlag {
			fmt.Printf("[essh debug] loading config file: %s\n", UserConfigFile)
		}

		if err := CurrentRegistry.MkDirs(); err != nil {
			printError(err)
			return ExitErr
		}

		if err := L.DoFile(UserConfigFile); err != nil {
			printError(err)
			return ExitErr
		}

		if debugFlag {
			fmt.Printf("[essh debug] loaded config file: %s\n", UserConfigFile)
		}
	}

	// load current dir config
	// change context to working dir context
	CurrentRegistry = NewRegistry(WorkingDataDir, RegistryTypeLocal)
	LocalRegistry = CurrentRegistry

	if _, err := os.Stat(WorkingDirConfigFile); err == nil {
		if debugFlag {
			fmt.Printf("[essh debug] loading config file: %s\n", WorkingDirConfigFile)
		}

		if err := CurrentRegistry.MkDirs(); err != nil {
			printError(err)
			return ExitErr
		}

		if err := L.DoFile(WorkingDirConfigFile); err != nil {
			printError(err)
			return ExitErr
		}

		if debugFlag {
			fmt.Printf("[essh debug] loaded config file: %s\n", WorkingDirConfigFile)
		}
	}

	// load override config
	if _, err := os.Stat(WorkingDirOverrideConfigFile); err == nil {
		if debugFlag {
			fmt.Printf("[essh debug] loading config file: %s\n", WorkingDirOverrideConfigFile)
		}

		if err := L.DoFile(WorkingDirOverrideConfigFile); err != nil {
			printError(err)
			return ExitErr
		}

		if debugFlag {
			fmt.Printf("[essh debug] loaded config file: %s\n", WorkingDirOverrideConfigFile)
		}
	}

	CurrentRegistry = GlobalRegistry
	// load override user config
	if _, err := os.Stat(UserOverrideConfigFile); err == nil {
		if debugFlag {
			fmt.Printf("[essh debug] loading config file: %s\n", UserOverrideConfigFile)
		}

		if err := CurrentRegistry.MkDirs(); err != nil {
			printError(err)
			return ExitErr
		}

		if err := L.DoFile(UserOverrideConfigFile); err != nil {
			printError(err)
			return ExitErr
		}

		if debugFlag {
			fmt.Printf("[essh debug] loaded config file: %s\n", UserOverrideConfigFile)
		}
	}

	// load override global config
	if _, err := os.Stat(SystemWideOverrideConfigFile); err == nil {
		if debugFlag {
			fmt.Printf("[essh debug] loading config file: %s\n", SystemWideOverrideConfigFile)
		}

		if err := CurrentRegistry.MkDirs(); err != nil {
			printError(err)
			return ExitErr
		}

		if err := L.DoFile(SystemWideOverrideConfigFile); err != nil {
			printError(err)
			return ExitErr
		}

		if debugFlag {
			fmt.Printf("[essh debug] loaded config file: %s\n", SystemWideOverrideConfigFile)
		}
	}

	// validate config
	if err := validateConfig(); err != nil {
		printError(err)
		return ExitErr
	}

	// show hosts for zsh completion
	if zshCompletionHostsFlag {
		for _, host := range SortedPublicHosts() {
			if !host.Hidden {
				fmt.Printf("%s\t%s\n", ColonEscape(host.Name), ColonEscape(host.DescriptionOrDefault()))
			}
		}

		return
	}

	// show tasks for zsh completion
	if zshCompletionTasksFlag {
		for _, task := range SortedTasks() {
			if !task.Disabled && !task.Hidden {
				fmt.Printf("%s\t%s\n", ColonEscape(task.Name), ColonEscape(task.DescriptionOrDefault()))
			}
		}
		return
	}

	if zshCompletionTagsFlag {
		for _, tag := range Tags() {
			fmt.Printf("%s\n", ColonEscape(tag))
		}
		return
	}

	// only print hosts list
	if hostsFlag {
		var hosts []*Host
		if len(selectVar) > 0 {
			hosts = HostsByNames(selectVar)
		} else {
			hosts = SortedHosts()
		}
		tb := helper.NewPlainTable(os.Stdout)
		if !quietFlag {
			if allFlag {
				tb.SetHeader([]string{"NAME", "DESCRIPTION", "TAGS", "REGISTRY", "SCOPE", "HIDDEN"})
			} else {
				tb.SetHeader([]string{"NAME", "DESCRIPTION", "TAGS", "REGISTRY"})
			}

		}
		for _, host := range hosts {
			if (!host.Hidden && !host.Private) || allFlag {
				if quietFlag {
					tb.Append([]string{host.Name})
				} else {
					if allFlag {
						hidden := "false"
						if host.Hidden {
							hidden = "true"
						}
						scope := "public"
						if host.Private {
							scope = "private"
						}
						tb.Append([]string{host.Name, host.Description, strings.Join(host.Tags, ","), host.Registry.TypeString(), scope, hidden})
					} else {
						tb.Append([]string{host.Name, host.Description, strings.Join(host.Tags, ","), host.Registry.TypeString()})
					}

				}
			}
		}
		tb.Render()

		return
	}

	// only print tags list
	if tagsFlag {
		tb := helper.NewPlainTable(os.Stdout)
		if !quietFlag {
			tb.SetHeader([]string{"NAME", "PUBLIC_HOSTS", "HOSTS"})
		}
		for _, tag := range Tags() {
			if quietFlag {
				tb.Append([]string{tag})
			} else {
				tb.Append([]string{tag, fmt.Sprintf("%d", len(HostsByTag(tag, true))), fmt.Sprintf("%d", len(HostsByTag(tag, false)))})
			}
		}
		tb.Render()

		return
	}

	// only print tasks list
	if tasksFlag {
		tb := helper.NewPlainTable(os.Stdout)
		if !quietFlag {
			if allFlag {
				tb.SetHeader([]string{"NAME", "DESCRIPTION", "REGISTRY", "DISABLED", "HIDDEN"})
			} else {
				tb.SetHeader([]string{"NAME", "DESCRIPTION", "REGISTRY"})
			}
		}

		for _, t := range SortedTasks() {
			if (!t.Hidden && !t.Disabled) || allFlag {
				if quietFlag {
					tb.Append([]string{t.Name})
				} else {
					if allFlag {
						tb.Append([]string{t.Name, t.Description, t.Registry.TypeString(), fmt.Sprintf("%v", t.Disabled), fmt.Sprintf("%v", t.Hidden)})
					} else {
						tb.Append([]string{t.Name, t.Description, t.Registry.TypeString()})
					}
				}
			}
		}
		tb.Render()

		return
	}

	outputConfig, ok := toString(lessh.RawGetString("ssh_config"))
	if !ok {
		printError(fmt.Errorf("invalid value %v in the 'ssh_config'", lessh.RawGetString("ssh_config")))
		return ExitErr
	}

	// generate ssh hosts config
	content, err := UpdateSSHConfig(outputConfig, SortedPublicHosts())
	if err != nil {
		printError(err)
		return ExitErr
	}

	// only print generated config
	if printFlag {
		fmt.Println(string(content))
		return
	}

	// only generating contents
	if genFlag {
		return
	}

	// select running mode and run it.
	if execFlag {
		if len(args) == 0 {
			printError("exec mode requires 1 parameter at latest.")
			return ExitErr
		}

		command := strings.Join(args, " ")

		// create temporary task
		task := NewTask()
		task.Name = "--exec"
		task.Pty = ptyFlag
		task.Parallel = parallelFlag
		task.Privileged = privilegedFlag
		task.Driver = driverVar
		if fileFlag {
			task.File = command
		} else {
			task.Script = []map[string]string{
				map[string]string{"code": command},
			}
		}
		if backendVar != "" {
			task.Backend = backendVar
		}
		task.Targets = targetVar

		if prefixStringVar == "" {
			if prefixFlag {
				if task.IsRemoteTask() {
					task.Prefix = DefaultPrefixRemote
				} else {
					task.Prefix = DefaultPrefixLocal
				}
			}
		} else {
			task.Prefix = prefixStringVar
		}

		err := runTask(outputConfig, task)
		if err != nil {
			printError(err)
			return ExitErr
		}

		return
	} else {
		// try to get a task.
		if len(args) > 0 {
			taskName := args[0]
			task := GetEnabledTask(taskName)
			if task != nil {
				if len(args) >= 2 {
					printError("too many arguments.")
					return ExitErr
				} else {
					err := runTask(outputConfig, task)
					if err != nil {
						printError(err)
						return ExitErr
					}
					return
				}
			}
		}

		if updateFlag && len(args) == 0 {
			// run just "essh --update"
			return
		}

		// no argument
		if len(args) == 0 {
			printUsage()
			return
		}

		// run ssh command
		err, ex := runSSH(L, outputConfig, args)
		if err != nil {
			printError(err)
			return ExitErr
		}

		exitStatus = ex
	}

	return
}

func UpdateSSHConfig(outputConfig string, enabledHosts []*Host) ([]byte, error) {
	if debugFlag {
		fmt.Printf("[essh debug] output ssh_config contents to the file: %s \n", outputConfig)
	}

	// generate ssh hosts config
	content, err := GenHostsConfig(enabledHosts)
	if err != nil {
		return nil, err
	}

	// update temporary ssh config file
	err = ioutil.WriteFile(outputConfig, content, 0644)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func runTask(config string, task *Task) error {
	if debugFlag {
		fmt.Printf("[essh debug] run task: %s\n", task.Name)
	}

	if task.Registry != nil {
		// change current registry
		CurrentRegistry = task.Registry
	}

	// re generate config (task).
	if task.Registry == nil {
		// this is "--exec" command mode. use only public hosts
		_, err := UpdateSSHConfig(config, SortedPublicHosts())
		if err != nil {
			return err
		}
	} else {
		_, err := UpdateSSHConfig(config, SameRegistryHosts(task.Registry.Type))
		if err != nil {
			return err
		}
	}

	if task.Prepare != nil {
		if debugFlag {
			fmt.Printf("[essh debug] run prepare function.\n")
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
		if task.Registry == nil {
			hosts = FindPublicHosts(task.TargetsSlice())
		} else {
			hosts = FindHostsInRegistry(task.TargetsSlice(), task.Registry.Type)
		}

		if len(hosts) == 0 {
			return fmt.Errorf("There are not hosts to run the command. you must specify the valid hosts.")
		}

		wg := &sync.WaitGroup{}
		m := new(sync.Mutex)
		for _, host := range hosts {
			if task.Parallel {
				wg.Add(1)
				go func(host *Host) {
					err := runRemoteTaskScript(config, task, host, hosts, m)
					if err != nil {
						fmt.Fprintf(os.Stderr, color.FgRB("essh error: %v\n", err))
						panic(err)
					}

					wg.Done()
				}(host)
			} else {
				err := runRemoteTaskScript(config, task, host, hosts, m)
				if err != nil {
					return err
				}
			}
		}
		wg.Wait()
	} else {
		// run locally.
		var hosts []*Host
		if task.Registry == nil {
			hosts = FindPublicHosts(task.TargetsSlice())
		} else {
			hosts = FindHostsInRegistry(task.TargetsSlice(), task.Registry.Type)
		}

		wg := &sync.WaitGroup{}
		m := new(sync.Mutex)

		if len(task.Targets) >= 1 && len(hosts) == 0 {
			return fmt.Errorf("There are not hosts to run the command. you must specify the valid hosts.")
		}

		if len(hosts) == 0 {
			err := runLocalTaskScript(config, task, nil, hosts, m)
			if err != nil {
				return err
			}
			return nil
		}

		for _, host := range hosts {
			if task.Parallel {
				wg.Add(1)
				go func(host *Host) {
					err := runLocalTaskScript(config, task, host, hosts, m)
					if err != nil {
						fmt.Fprintf(os.Stderr, color.FgRB("essh error: %v\n", err))
						panic(err)
					}

					wg.Done()
				}(host)
			} else {
				err := runLocalTaskScript(config, task, host, hosts, m)
				if err != nil {
					return err
				}
			}
		}
		wg.Wait()
	}

	return nil
}

func runRemoteTaskScript(sshConfigPath string, task *Task, host *Host, hosts []*Host, m *sync.Mutex) error {
	// setup ssh command args
	var sshCommandArgs []string
	if task.Pty {
		sshCommandArgs = []string{"-t", "-t", "-F", sshConfigPath, host.Name}
	} else {
		sshCommandArgs = []string{"-F", sshConfigPath, host.Name}
	}

	// generate commands by using driver
	driver := DefaultDriver
	if task.Driver != "" {
		driver = FindDriverInRegistry(task.Driver, task.Registry)
		if driver == nil {
			return fmt.Errorf("invalid driver name '%s'", task.Driver)
		}
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

	if task.Privileged {
		script = "sudo su - <<\\EOF-ESSH-PRIVILEGED\n" + script + "\n" + "EOF-ESSH-PRIVILEGED"
	}

	// inspired by https://github.com/laravel/envoy
	delimiter := "EOF-ESSH-SCRIPT"
	sshCommandArgs = append(sshCommandArgs, "bash", "-s", "<<\\"+delimiter+"\n"+script+"\n"+delimiter)

	cmd := exec.Command("ssh", sshCommandArgs[:]...)
	if debugFlag {
		fmt.Printf("[essh debug] real ssh command: %v \n", cmd.Args)
	}

	prefix := ""
	if task.Prefix != "" {
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
		tmpl, err := template.New("T").Funcs(funcMap).Parse(task.Prefix)
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

	cmd.Stdin = os.Stdin

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	// inspired by https://github.com/fujiwara/nssh/blob/master/nssh.go
	go scanLines(stdout, os.Stdout, prefix, m)
	go scanLines(stderr, os.Stderr, prefix, m)

	return cmd.Wait()
}

func runLocalTaskScript(sshConfigPath string, task *Task, host *Host, hosts []*Host, m *sync.Mutex) error {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "bash"
		flag = "-c"
	}

	// generate commands by using driver
	driver := DefaultDriver
	if task.Driver != "" {
		driver = FindDriverInRegistry(task.Driver, task.Registry)
		if driver == nil {
			return fmt.Errorf("invalid driver name '%s'", task.Driver)
		}
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

	if task.Privileged {
		script = "cd " + WorkingDir + "\n" + script
		script = "sudo su - <<\\EOF-ESSH-PRIVILEGED\n" + script + "\n" + "EOF-ESSH-PRIVILEGED"
	}

	cmd := exec.Command(shell, flag, script)
	if debugFlag {
		fmt.Printf("[essh debug] real local command: %v \n", cmd.Args)
	}

	prefix := ""
	if task.Prefix == DefaultPrefixLocal && host == nil {
		// simple local task (does not specify the hosts)
		// prevent to use invalid text template.
		// replace prefix string to the string that is not included "{{.Host}}"
		prefix = "[local] "
	} else if task.Prefix != "" {
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
		tmpl, err := template.New("T").Funcs(funcMap).Parse(task.Prefix)
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

	cmd.Stdin = os.Stdin

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	// inspired by https://github.com/fujiwara/nssh/blob/master/nssh.go
	go scanLines(stdout, os.Stdout, prefix, m)
	go scanLines(stderr, os.Stderr, prefix, m)

	return cmd.Wait()
}

func scanLines(src io.ReadCloser, dest io.Writer, prefix string, m *sync.Mutex) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		func(m *sync.Mutex) {
			m.Lock()
			defer m.Unlock()
			if prefix != "" {
				fmt.Fprintf(dest, "%s%s\n", color.FgCB(prefix), scanner.Text())
			} else {
				fmt.Fprintf(dest, "%s\n", scanner.Text())
			}
		}(m)
	}
}

func runSSH(L *lua.LState, config string, args []string) (error, int) {
	// hooks
	hooks := map[string][]interface{}{}

	// Limitation!
	// hooks fires only when the hostname is just specified.
	if len(args) == 1 {
		hostname := args[0]
		if host := GetPublicHost(hostname); host != nil {
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

func getHookScript(L *lua.LState, hooks []interface{}) (string, error) {
	hookScript := ""
	for _, hook := range hooks {
		code, err := convertHook(L, hook)
		if err != nil {
			return "", err
		}
		hookScript += code + "\n"
	}

	return hookScript, nil
}

func convertHook(L *lua.LState, hook interface{}) (string, error) {
	if hookFn, ok := hook.(*lua.LFunction); ok {
		err := L.CallByParam(lua.P{
			Fn:      hookFn,
			NRet:    1,
			Protect: false,
		})

		ret := L.Get(-1) // returned value
		L.Pop(1)

		if err != nil {
			return "", err
		}

		if ret == lua.LNil {
			return "", nil
		} else if retStr, ok := toString(ret); ok {
			return retStr, nil
		} else if retFn, ok := toLFunction(ret); ok {
			return convertHook(L, retFn)
		} else {
			return "", fmt.Errorf("hook function return value must be string or function.")
		}
	} else if hookStr, ok := hook.(string); ok {
		return hookStr, nil
	} else {
		return "", fmt.Errorf("invalid type hook: %v", hook)
	}
}

func runCommand(command string) error {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}
	cmd := exec.Command(shell, flag, command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func validateConfig() error {
	// check duplication of the host, task and tag names
	names := map[string]bool{}
	for _, host := range SortedPublicHosts() {
		if _, ok := names[host.Name]; ok {
			return fmt.Errorf("Host '%s' is duplicated", host.Name)
		}
		names[host.Name] = true
	}

	for _, task := range SortedTasks() {
		if _, ok := names[task.Name]; ok {
			return fmt.Errorf("Task '%s' is duplicated", task.Name)
		}
		names[task.Name] = true
	}

	for _, tag := range Tags() {
		if _, ok := names[tag]; ok {
			return fmt.Errorf("Tag '%s' is duplicated", tag)
		}
		names[tag] = true
	}

	return nil
}

type CallbackWriter struct {
	Func func(data []byte)
}

func (w *CallbackWriter) Write(data []byte) (int, error) {
	if w.Func != nil {
		w.Func(data)
	}
	return len(data), nil
}

func removeRegistryData() error {
	if withGlobalFlag {
		c := NewRegistry(UserDataDir, RegistryTypeGlobal)
		if cleanModulesFlag || cleanAllFlag {
			if _, err := os.Stat(c.ModulesDir()); err == nil {
				fmt.Fprintf(os.Stdout, "Deleting: '%s'\n", color.FgYB(c.ModulesDir()))
				err = os.RemoveAll(c.ModulesDir())
				if err != nil {
					return err
				}
			}
		}

		if cleanTmpFlag || cleanAllFlag {
			if _, err := os.Stat(c.TmpDir()); err == nil {
				fmt.Fprintf(os.Stdout, "Deleting: '%s'\n", color.FgYB(c.TmpDir()))
				err = os.RemoveAll(c.TmpDir())
				if err != nil {
					return err
				}
			}
		}
	}

	c := NewRegistry(WorkingDataDir, RegistryTypeLocal)
	if cleanModulesFlag || cleanAllFlag {
		if _, err := os.Stat(c.ModulesDir()); err == nil {
			fmt.Fprintf(os.Stdout, "Deleting: '%s'\n", color.FgYB(c.ModulesDir()))
			err = os.RemoveAll(c.ModulesDir())
			if err != nil {
				return err
			}
		}
	}

	if cleanTmpFlag || cleanAllFlag {
		if _, err := os.Stat(c.TmpDir()); err == nil {
			fmt.Fprintf(os.Stdout, "Deleting: '%s'\n", color.FgYB(c.TmpDir()))
			err = os.RemoveAll(c.TmpDir())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func printUsage() {
	fmt.Print(`Usage: essh [<options>] [<ssh options and args...>]

  Essh is an extended ssh command.
  version ` + Version + ` (` + CommitHash + `)

  Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
  The MIT License (MIT)

Options:
  (General Options)
  --print                       Print generated ssh config.
  --gen                         Only generate ssh config.
  --working-dir <dir>           Change working directory.
  --config <file>               Load per-project configuration from the file.
  --color                       Force ANSI output.
  --no-color                    Disable ANSI output.
  --debug                       Output debug log.

  (Manage Hosts, Tags And Tasks)
  --hosts                       List hosts.
  --tags                        List tags.
  --tasks                       List tasks.
  --select <tag|host>           (Using with --hosts option) Use only the hosts filtered with a tag or a host.
  --all                         (Using with --hosts, --tasks or --tags option) Show all that includs hidden objects.
  --quiet                       (Using with --hosts, --tasks or --tags option) Show only names.

  (Manage Modules)
  --update                      Update modules.
  --clean-modules               Clean downloaded modules.
  --clean-tmp                   Clean temporary data.
  --clean-all                   Clean all data.
  --with-global                 (Using with --update, --clean-modules, --clean-tmp or --clean-all option) Update or clean modules in the local and global both registry.

  (Execute Commands)
  --exec                        Execute commands with the hosts.
  --target <tag|host>           (Using with --exec option) Target hosts to run the commands.
  --backend <remote|local>      (Using with --exec option) Run the commands on local or remote hosts.
  --prefix                      (Using with --exec option) Enable outputing prefix.
  --prefix-string [<prefix>]    (Using with --exec option) Custom string of the prefix.
  --privileged                  (Using with --exec option) Run by the privileged user.
  --parallel                    (Using with --exec option) Run in parallel.
  --pty                         (Using with --exec option) Allocate pseudo-terminal. (add ssh option "-t -t" internally)
  --file                        (Using with --exec option) Load commands from a file.
  --driver                      (Using with --exec option) Specify a driver.

  (Completion)
  --zsh-completion              Output zsh completion code.
  --aliases                     Output aliases code.

  (Help)
  --version                     Print version.
  --help                        Print help.

See: https://github.com/kohkimakimoto/essh for updates, code and issues.

`)
}

func sprintByTemplate(tmplContent string) (string, error) {
	tmpl, err := template.New("T").Parse(tmplContent)
	if err != nil {
		return "", err
	}

	dict := map[string]interface{}{
		"Executable": Executable,
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, dict)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func printError(err interface{}) {
	fmt.Fprintf(os.Stderr, color.FgRB("essh error: %v\n", err))
}

func init() {
	// set SystemWideConfigFile
	SystemWideConfigFile = "/etc/essh/config.lua"
	SystemWideOverrideConfigFile = "/etc/essh/config_override.lua"

	// set UserDataDir
	home := userHomeDir()
	UserDataDir = filepath.Join(home, ".essh")

	// create UserDataDir, if it doesn't exist
	if _, err := os.Stat(UserDataDir); os.IsNotExist(err) {
		err = os.MkdirAll(UserDataDir, os.FileMode(0755))
		if err != nil {
			panic(err)
		}
	}

	UserConfigFile = filepath.Join(UserDataDir, "config.lua")
	UserOverrideConfigFile = filepath.Join(UserDataDir, "config_override.lua")

	if _bin, err := osext.Executable(); err == nil {
		Executable = _bin
	} else {
		Executable = "essh"
	}

}

var ZSH_COMPLETION = `# This is zsh completion code.
# If you want to use it. write the following code in your '.zshrc'
#   eval "$({{.Executable}} --zsh-completion)"
_essh_hosts() {
    local -a __essh_hosts
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_hosts=($({{.Executable}} --zsh-completion-hosts | awk -F'\t' '{print $1":"$2}'))
    IFS=$PRE_IFS
    _describe -t host "host" __essh_hosts
}

_essh_tasks() {
    local -a __essh_tasks
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_tasks=($({{.Executable}} --zsh-completion-tasks | awk -F'\t' '{print $1":"$2}'))
    IFS=$PRE_IFS
    _describe -t task "task" __essh_tasks
}

_essh_tags() {
    local -a __essh_tags
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_tags=($({{.Executable}} --zsh-completion-tags))
    IFS=$PRE_IFS
    _describe -t tag "tag" __essh_tags
}

_essh_global_options() {
}

_essh_options() {
    local -a __essh_options
    __essh_options=(
        '--version:Print version.'
        '--help:Print help.'
        '--print:Print generated ssh config.'
        '--color:Force ANSI output.'
        '--no-color:Disable ANSI output.'
        '--gen:Only generate ssh config.'
        '--update:Update modules.'
        '--clean-modules:Clean downloaded modules.'
        '--clean-tmp:Clean temporary data.'
        '--clean-all:Clean all data.'
        '--working-dir:Change working directory.'
        '--config:Load per-project configuration from the file.'
        '--hosts:List hosts.'
        '--tags:List tags.'
        '--tasks:List tasks.'
        '--debug:Output debug log.'
        '--exec:Execute commands with the hosts.'
        '--zsh-completion:Output zsh completion code.'
        '--aliases:Output aliases code.'
     )
    _describe -t option "option" __essh_options
}

_essh_hosts_options() {
    local -a __essh_options
    __essh_options=(
        '--debug:Output debug log.'
        '--quiet:Show only names.'
        '--all:Show all that includs hidden objects.'
        '--select:Use only the hosts filtered with a tag or a host.'
     )
    _describe -t option "option" __essh_options
}

_essh_tasks_options() {
    local -a __essh_options
    __essh_options=(
        '--debug:Output debug log.'
        '--quiet:Show only names.'
        '--all:Show all that includs hidden objects.'
     )
    _describe -t option "option" __essh_options
}

_essh_tags_options() {
    local -a __essh_options
    __essh_options=(
        '--debug:Output debug log.'
        '--quiet:Show only names.'
        '--all:Show all that includs hidden objects.'
     )
    _describe -t option "option" __essh_options
}

_essh_exec_options() {
    local -a __essh_options
    __essh_options=(
        '--debug:Output debug log.'
        '--backend:Run the commands on local or remote hosts.'
        '--target:Target hosts to run the commands.'
        '--prefix:Disable outputing prefix.'
        '--prefix-string:Custom string of the prefix.'
        '--privileged:Run by the privileged user.'
        '--parallel:Run in parallel.'
        '--pty:Allocate pseudo-terminal. (add ssh option "-t -t" internally)'
        '--file:Load commands from a file.'
        '--driver:Specify a driver.'
     )
    _describe -t option "option" __essh_options
}

_essh_registry_options() {
    local -a __essh_options
    __essh_options=(
        '--with-global:Update or clean modules in the local, global both registry.'
     )
    _describe -t option "option" __essh_options
}

_essh_backends() {
    local -a __essh_options
    __essh_options=(
        'local'
        'remote'
     )
    _describe -t option "option" __essh_options
}

_essh () {
    local curcontext="$curcontext" state line
    local last_arg arg execMode hostsMode tasksMode tagsMode

    typeset -A opt_args

    _arguments \
        '1: :->objects' \
        '*: :->args' \
        && ret=0

    case $state in
        objects)
            case $line[1] in
                -*)
                    _essh_options
                    ;;
                *)
                    _essh_tasks
                    _essh_hosts
                    ;;
            esac
            ;;
        args)
            last_arg="${line[${#line[@]}-1]}"

            for arg in ${line[@]}; do
                case $arg in
                    --exec)
                        execMode="on"
                        ;;
                    --hosts)
                        hostsMode="on"
                        ;;
                    --tasks)
                        tasksMode="on"
                        ;;
                    --tags)
                        tagsMode="on"
                        ;;
                    *)
                        ;;
                esac
            done

            case $last_arg in
                --print|--help|--version|--gen)
                    ;;
                --file|--config)
                    _files
                    ;;
                --select|--target)
                    _essh_hosts
                    _essh_tags
                    ;;
                --backend)
                    _essh_backends
                    ;;
                --clean-modules|--clean-tmp|--clean-all|--update)
                    _essh_registry_options
                    ;;
                *)
                    if [ "$execMode" = "on" ]; then
                        _essh_exec_options
                    elif [ "$hostsMode" = "on" ]; then
                        _essh_hosts_options
                    elif [ "$tasksMode" = "on" ]; then
                        _essh_tasks_options
                    elif [ "$tagsMode" = "on" ]; then
                        _essh_tags_options
                    else
                        _essh_options
                        _files
                    fi
                    ;;
            esac
            ;;
        *)
            _files
            ;;
    esac
}

compdef _essh essh
`

var ALIASES_CODE = `# This is aliases code.
# If you want to use it. write the following code in your '.zshrc'
#   eval "$({{.Executable}} --aliases)"
function escp() {
    {{.Executable}} --exec 'scp -F $ESSH_SSH_CONFIG' "$@"
}
function ersync() {
    {{.Executable}} --exec 'rsync -e "ssh -F $ESSH_SSH_CONFIG"' "$@"
}
`

var BASH_COMPLETION = ``
