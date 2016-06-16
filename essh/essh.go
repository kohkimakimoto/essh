package essh

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kohkimakimoto/essh/color"
	"github.com/kohkimakimoto/essh/helper"
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
	SystemWideConfigFile string
	UserConfigFile       string
	UserDataDir          string
	WorkingDirConfigFile string
	WorkingDataDir       string
	WorkingDir           string
)

// flags
var (
	versionFlag            bool
	helpFlag               bool
	printFlag              bool
	configFlag             bool
	userConfigFlag         bool
	systemConfigFlag       bool
	debugFlag              bool
	hostsFlag              bool
	quietFlag              bool
	allFlag                bool
	tagsFlag               bool
	tasksFlag              bool
	genFlag                bool
	updateFlag             bool
	noGlobalFlag           bool
	cleanFlag              bool
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
	rsyncFlag              bool
	scpFlag                bool

	workindDirVar   string
	filtersVar      []string = []string{}
	onVar           []string = []string{}
	foreachVar      []string = []string{}
	prefixStringVar string
	driverVar       string
	// beta implementation
	formatVar string
)

func Start() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}

		if err != nil && zshCompletionModeFlag && !debugFlag {
			// suppress the error in running completion code.
			err = nil
		}
	}()

	err = start()

	return err
}

func start() error {
	if len(os.Args) == 1 {
		printUsage()
		return nil
	}

	osArgs := os.Args[1:]
	args := []string{}
	osArgsIndex := 0

	for {
		if len(osArgs) == 0 {
			break
		}

		arg := osArgs[0]
		if arg == "--print" {
			printFlag = true
		} else if arg == "--version" {
			versionFlag = true
		} else if arg == "--help" {
			helpFlag = true
		} else if arg == "--config" {
			configFlag = true
		} else if arg == "--user-config" {
			userConfigFlag = true
		} else if arg == "--system-config" {
			systemConfigFlag = true
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
		} else if arg == "--filter" {
			if len(osArgs) < 2 {
				return fmt.Errorf("--filter reguires an argument.")
			}
			filtersVar = append(filtersVar, osArgs[1])
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--filter=") {
			filtersVar = append(filtersVar, strings.Split(arg, "=")[1])
		} else if arg == "--format" {
			if len(osArgs) < 2 {
				return fmt.Errorf("--format reguires an argument.")
			}
			formatVar = osArgs[1]
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--format=") {
			formatVar = strings.Split(arg, "=")[1]
		} else if arg == "--tags" {
			tagsFlag = true
		} else if arg == "--gen" {
			genFlag = true
		} else if arg == "--update" {
			updateFlag = true
		} else if arg == "--clean" {
			cleanFlag = true
		} else if arg == "--no-global" {
			noGlobalFlag = true
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
			bashCompletionFlag = true
		} else if arg == "--aliases" {
			aliasesFlag = true
		} else if arg == "--working-dir" {
			if len(osArgs) < 2 {
				return fmt.Errorf("--working-dir reguires an argument.")
			}
			workindDirVar = osArgs[1]
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--working-dir=") {
			workindDirVar = strings.Split(arg, "=")[1]
		} else if arg == "--exec" {
			execFlag = true
		} else if arg == "--on" {
			if len(osArgs) < 2 {
				return fmt.Errorf("--on reguires an argument.")
			}
			onVar = append(onVar, osArgs[1])
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--on=") {
			onVar = append(onVar, strings.Split(arg, "=")[1])
		} else if arg == "--foreach" {
			if len(osArgs) < 2 {
				return fmt.Errorf("--foreach reguires an argument.")
			}
			foreachVar = append(foreachVar, osArgs[1])
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--foreach=") {
			foreachVar = append(foreachVar, strings.Split(arg, "=")[1])
		} else if arg == "--privileged" {
			privilegedFlag = true
		} else if arg == "--parallel" {
			parallelFlag = true
		} else if arg == "--prefix" {
			prefixFlag = true
		} else if arg == "--prefix-string" {
			if len(osArgs) < 2 {
				return fmt.Errorf("--prefix-string reguires an argument.")
			}
			prefixStringVar = osArgs[1]
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--prefix-string=") {
			prefixStringVar = strings.Split(arg, "=")[1]
		} else if arg == "--driver" {
			if len(osArgs) < 2 {
				return fmt.Errorf("--driver reguires an argument.")
			}
			driverVar = osArgs[1]
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--driver=") {
			driverVar = strings.Split(arg, "=")[1]
		} else if arg == "--file" {
			fileFlag = true
		} else if arg == "--pty" {
			ptyFlag = true
		} else if arg == "--rsync" {
			if osArgsIndex != 0 {
				return fmt.Errorf("--rsync must be the first option.")
			}
			rsyncFlag = true
		} else if arg == "--scp" {
			if osArgsIndex != 0 {
				return fmt.Errorf("--scp must be the first option.")
			}
			scpFlag = true
		} else if !rsyncFlag && !scpFlag && strings.HasPrefix(arg, "--") {
			// rsync can have long options
			return fmt.Errorf("invalid option '%s'.", arg)
		} else {
			// restructure args to remove essh options.
			args = append(args, arg)
		}

		osArgsIndex++
		osArgs = osArgs[1:]
	}

	if os.Getenv("ESSH_DEBUG") != "" {
		debugFlag = true
	}

	if workindDirVar != "" {
		err := os.Chdir(workindDirVar)
		if err != nil {
			return err
		}
	}

	// decide the wokingDirConfigFile
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("couldn't get working dir %v\n", err)
	}
	WorkingDir = wd
	WorkingDataDir = filepath.Join(wd, ".essh")
	WorkingDirConfigFile = filepath.Join(wd, "essh.lua")

	if helpFlag {
		printHelp()
		return nil
	}

	if cleanFlag {
		err := removeModules()
		if err != nil {
			return err
		}
		return nil
	}

	if versionFlag {
		fmt.Printf("%s (%s)\n", Version, CommitHash)
		return nil
	}

	if zshCompletionFlag {
		fmt.Print(ZSH_COMPLETION)
		return nil
	}

	if aliasesFlag {
		fmt.Print(ALIASES_CODE)
		return nil
	}

	if bashCompletionFlag {
		fmt.Print(BASH_COMPLETION)
		return nil
	}

	if configFlag {
		runCommand(getEditor() + " " + WorkingDirConfigFile)
		return nil
	}

	if userConfigFlag {
		runCommand(getEditor() + " " + UserConfigFile)
		return nil
	}

	if systemConfigFlag {
		runCommand(getEditor() + " " + SystemWideConfigFile)
		return nil
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
		return err
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
		return fmt.Errorf("essh must be a table")
	}

	// set temporary ssh config file path
	lessh.RawSetString("ssh_config", lua.LString(temporarySSHConfigFile))

	// user context
	CurrentContext = NewContext(UserDataDir, ContextTypeGlobal)
	ContextMap[CurrentContext.Key] = CurrentContext

	if err := CurrentContext.MkDirs(); err != nil {
		return err
	}

	// load system wide config
	if _, err := os.Stat(SystemWideConfigFile); err == nil {
		if debugFlag {
			fmt.Printf("[essh debug] loading config file: %s\n", SystemWideConfigFile)
		}

		if err := L.DoFile(SystemWideConfigFile); err != nil {
			return err
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

		if err := L.DoFile(UserConfigFile); err != nil {
			return err
		}

		if debugFlag {
			fmt.Printf("[essh debug] loaded config file: %s\n", UserConfigFile)
		}
	}

	// load current dir config
	if WorkingDirConfigFile != "" {
		if _, err := os.Stat(WorkingDirConfigFile); err == nil {

			if debugFlag {
				fmt.Printf("[essh debug] loading config file: %s\n", WorkingDirConfigFile)
			}

			// change context to working dir context
			CurrentContext = NewContext(WorkingDataDir, ContextTypeLocal)
			ContextMap[CurrentContext.Key] = CurrentContext

			if err := CurrentContext.MkDirs(); err != nil {
				return err
			}

			if err := L.DoFile(WorkingDirConfigFile); err != nil {
				return err
			}

			if debugFlag {
				fmt.Printf("[essh debug] loaded config file: %s\n", WorkingDirConfigFile)
			}
		}

		// load additional config files.
		files, err := filepath.Glob(filepath.Join(WorkingDir, "essh.*.lua"))
		if err != nil {
			return err
		}
		for _, file := range files {
			if debugFlag {
				fmt.Printf("[essh debug] loading config file: %s\n", file)
			}

			if err := L.DoFile(file); err != nil {
				return err
			}

			if debugFlag {
				fmt.Printf("[essh debug] loaded config file: %s\n", file)
			}
		}
	}

	// basic configuration loading is completed.

	// override config using task configuration?
	taskConfigureContextKey := os.Getenv("ESSH_TASK_CONFIGURE_CONTEXT_KEY")
	if taskConfigureContextKey != "" {
		// check context
		if ctx, ok := ContextMap[taskConfigureContextKey]; ok {
			if debugFlag {
				fmt.Printf("[essh debug] got a context for configuring '%s' (%s) \n", ctx.Key, ctx.DataDir)
			}

			taskConfigureTask := os.Getenv("ESSH_TASK_CONFIGURE_TASK")
			if taskConfigureTask != "" {
				task := GetEnabledTask(taskConfigureTask)
				if task == nil {
					return fmt.Errorf("load configuration by using ESSH_TASK_CONFIGURE_TASK. but used unknown task '%s'", taskConfigureTask)
				}
				if err := processTaskConfigure(task); err != nil {
					return err
				}
			}
		} else {
			if debugFlag {
				fmt.Printf("[essh debug] a context for configuring is '%s'. but is not included in now context map.\n", taskConfigureContextKey)
			}
		}
	}

	// validate config
	if err := validateConfig(); err != nil {
		return err
	}

	// show hosts for zsh completion
	if zshCompletionHostsFlag {
		for _, host := range SortedPublicHosts() {
			if !host.Hidden {
				fmt.Printf("%s\t%s\n", ColonEscape(host.Name), ColonEscape(host.DescriptionOrDefault()))
			}
		}

		return nil
	}

	// show tasks for zsh completion
	if zshCompletionTasksFlag {
		for _, task := range Tasks {
			if !task.Disabled && !task.Hidden {
				fmt.Printf("%s\t%s\n", ColonEscape(task.Name), ColonEscape(task.DescriptionOrDefault()))
			}
		}
		return nil
	}

	if zshCompletionTagsFlag {
		for _, tag := range Tags() {
			fmt.Printf("%s\n", ColonEscape(tag))
		}
		return nil
	}

	// only print hosts list
	if hostsFlag {
		var hosts []*Host
		if len(filtersVar) > 0 {
			hosts = HostsByNames(filtersVar)
		} else {
			hosts = SortedHosts()
		}
		tb := helper.NewPlainTable(color.StdoutWriter)
		if !quietFlag {
			tb.SetHeader([]string{"NAME", "DESCRIPTION", "TAGS", "REGISTRY", "HIDDEN", "SCOPE"})
		}
		for _, host := range hosts {
			if (!host.Hidden && !host.Private) || allFlag {
				if quietFlag {
					tb.Append([]string{host.Name})
				} else {
					//
					hidden := ""
					if host.Hidden {
						hidden = "true"
					}
					scope := ""
					if host.Private {
						scope = "private"
					}
					tb.Append([]string{host.Name, host.Description, strings.Join(host.Tags, ","), host.Context.TypeString(), hidden, scope})

					//if host.Private {
					//	green := color.FgG
					//	tb.Append([]string{green(host.Name), green(host.Description), green(strings.Join(host.Tags, ",")), green(host.Context.TypeString()), green(hidden), green(scope)})
					//} else if host.Hidden {
					//	yellow := color.FgY
					//	tb.Append([]string{yellow(host.Name), yellow(host.Description), yellow(strings.Join(host.Tags, ",")), yellow(host.Context.TypeString()), yellow(hidden), yellow(scope)})
					//} else {
					//	tb.Append([]string{host.Name, host.Description, strings.Join(host.Tags, ","), host.Context.TypeString(), hidden, scope})
					//}
				}
			}
		}
		tb.Render()

		return nil
	}

	// only print tags list
	if tagsFlag {
		tb := helper.NewPlainTable(os.Stdout)
		if !quietFlag {
			tb.SetHeader([]string{"NAME"})
		}
		for _, tag := range Tags() {
			tb.Append([]string{tag})
		}
		tb.Render()

		return nil
	}

	// only print tasks list
	if tasksFlag {
		tb := helper.NewPlainTable(os.Stdout)
		if !quietFlag {
			tb.SetHeader([]string{"NAME", "DESCRIPTION", "DISABLED", "HIDDEN"})
		}
		for _, t := range SortedTasks() {
			if (!t.Hidden && !t.Disabled) || allFlag {
				if quietFlag {
					tb.Append([]string{t.Name})
				} else {
					//
					if t.Disabled {
						red := color.FgR
						tb.Append([]string{red(t.Name), red(t.Description), red(fmt.Sprintf("%v", t.Disabled)), red(fmt.Sprintf("%v", t.Hidden))})
					} else if t.Hidden {
						yellow := color.FgY
						tb.Append([]string{yellow(t.Name), yellow(t.Description), yellow(fmt.Sprintf("%v", t.Disabled)), yellow(fmt.Sprintf("%v", t.Hidden))})
					} else {
						tb.Append([]string{t.Name, t.Description, fmt.Sprintf("%v", t.Disabled), fmt.Sprintf("%v", t.Hidden)})
					}
				}
			}
		}
		tb.Render()

		return nil
	}

	outputConfig, ok := toString(lessh.RawGetString("ssh_config"))
	if !ok {
		return fmt.Errorf("invalid value %v in the 'ssh_config'", lessh.RawGetString("ssh_config"))
	}

	// generate ssh hosts config
	content, err := UpdateSSHConfig(outputConfig, SortedPublicHosts())
	if err != nil {
		return err
	}

	// only print generated config
	if printFlag {
		fmt.Println(string(content))
		return nil
	}

	// only generating contents
	if genFlag {
		return nil
	}

	// select running mode and run it.
	if execFlag {
		if len(args) == 0 {
			return fmt.Errorf("exec mode requires 1 parameter at latest.")
		}

		command := args[0]
		payload := ""
		if len(args) == 2 {
			payload = args[1]
		}

		// create temporary task
		task := NewTask()
		task.Name = "--exec"
		task.Pty = ptyFlag
		task.Lock = false
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
		task.On = onVar
		task.Foreach = foreachVar

		if len(task.Foreach) >= 1 && len(task.On) >= 1 {
			return fmt.Errorf("invalid options: can't use '--foreach' and '--on' at the same time.")
		}

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

		return runTask(outputConfig, task, payload)
	} else if rsyncFlag {
		err = runRsync(outputConfig, args)
	} else if scpFlag {
		err = runSCP(outputConfig, args)
	} else {
		// try to get a task.
		if len(args) > 0 {
			taskName := args[0]
			task := GetEnabledTask(taskName)
			if task != nil {
				if len(args) > 2 {
					return fmt.Errorf("too many arguments.")
				} else if len(args) == 2 {
					return runTask(outputConfig, task, args[1])
				} else {
					return runTask(outputConfig, task, "")
				}
			}
		}

		if updateFlag && len(args) == 0 {
			// run just "essh --update"
			return nil
		}

		// no argument
		if len(args) == 0 {
			printUsage()
			return nil
		}
		// run ssh command
		err = runSSH(L, outputConfig, args)
	}

	return err
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

func printJson(hosts []*Host, indent string) {
	convHosts := []map[string]map[string]interface{}{}

	for _, host := range hosts {
		h := map[string]map[string]interface{}{}

		hv := map[string]interface{}{}
		for _, pair := range host.SSHConfig() {
			for k, v := range pair {
				hv[k] = v
			}
		}
		h[host.Name] = hv

		hv["description"] = host.Description
		hv["Hidden"] = host.Hidden
		hv["tags"] = host.Tags

		convHosts = append(convHosts, h)
	}

	if indent == "" {
		b, err := json.Marshal(convHosts)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	} else {
		b, err := json.MarshalIndent(convHosts, "", indent)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}
}

func processTaskConfigure(task *Task) error {
	// configure function cleans global config and uses custom config that is defined in a task.
	if task.Configure == nil {
		return nil
	}

	if debugFlag {
		fmt.Printf("[essh debug] run configure function.\n")
	}

	// clean hosts.
	ResetHosts()

	err := os.Setenv("ESSH_TASK_CONFIGURE_TASK", task.Name)
	if err != nil {
		return err
	}

	err = os.Setenv("ESSH_TASK_CONFIGURE_CONTEXT_KEY", task.Context.Key)
	if err != nil {
		return err
	}

	err = task.Configure()
	if err != nil {
		return err
	}

	return nil
}

func runTask(config string, task *Task, payload string) error {
	if debugFlag {
		fmt.Printf("[essh debug] run task: %s\n", task.Name)
	}

	// re generate config (task u).
	_, err := UpdateSSHConfig(config, SameContextHosts(task.Context.Type))
	if err != nil {
		return err
	}

	if err := processTaskConfigure(task); err != nil {
		return err
	}

	if task.Configure != nil {
		// re generate config.
		_, err := UpdateSSHConfig(config, SortedHosts())
		if err != nil {
			return err
		}
	}

	if task.Prepare != nil {
		if debugFlag {
			fmt.Printf("[essh debug] run prepare function.\n")
		}

		ctx := NewTaskContext(task, payload)
		err := task.Prepare(ctx)
		if err != nil {
			return err
		}

		payload = ctx.Payload
	}

	// get target hosts.
	if task.IsRemoteTask() {
		// run remotely.
		hosts := FindHosts(task.TargetsSlice(), task.Context.Type)
		wg := &sync.WaitGroup{}
		m := new(sync.Mutex)
		for _, host := range hosts {
			if task.Parallel {
				wg.Add(1)
				go func(host *Host) {
					err := runRemoteTaskScript(config, task, payload, host, m)
					if err != nil {
						fmt.Fprintf(color.StderrWriter, color.FgRB("[essh error] %v\n", err))
						panic(err)
					}

					wg.Done()
				}(host)
			} else {
				err := runRemoteTaskScript(config, task, payload, host, m)
				if err != nil {
					return err
				}
			}
		}
		wg.Wait()
	} else {
		// run locally.
		hosts := FindHosts(task.TargetsSlice(), task.Context.Type)
		wg := &sync.WaitGroup{}
		m := new(sync.Mutex)

		if len(hosts) == 0 {
			err := runLocalTaskScript(task, payload, nil, m)
			if err != nil {
				return err
			}
			return nil
		}

		for _, host := range hosts {
			if task.Parallel {
				wg.Add(1)
				go func(host *Host) {
					err := runLocalTaskScript(task, payload, host, m)
					if err != nil {
						fmt.Fprintf(color.StderrWriter, color.FgRB("[essh error] %v\n", err))
						panic(err)
					}

					wg.Done()
				}(host)
			} else {
				err := runLocalTaskScript(task, payload, host, m)
				if err != nil {
					return err
				}
			}
		}
		wg.Wait()
	}

	return nil
}

func runRemoteTaskScript(config string, task *Task, payload string, host *Host, m *sync.Mutex) error {
	// setup ssh command args
	var sshCommandArgs []string
	if task.Pty {
		sshCommandArgs = []string{"-t", "-t", "-F", config, host.Name}
	} else {
		sshCommandArgs = []string{"-F", config, host.Name}
	}

	// generate commands by using driver
	driver := Drivers[BuiltinDefaultDriverName]
	if task.Driver != "" {
		driver = Drivers[task.Driver]
		if driver == nil {
			return fmt.Errorf("invalid driver name '%s'", task.Driver)
		}
	}
	if debugFlag {
		fmt.Printf("[essh debug] driver: %s \n", driver.Name)
	}

	var script string
	content, err := driver.GenerateRunnableContent(task, host)
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
		dict := map[string]interface{}{
			"Host": host,
			"Task": task,
		}
		tmpl, err := template.New("T").Parse(task.Prefix)
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

	cmd.Stdin = bytes.NewBufferString(payload)

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
	go scanLines(stdout, color.StdoutWriter, prefix, m)
	go scanLines(stderr, color.StderrWriter, prefix, m)

	return cmd.Wait()
}

func runLocalTaskScript(task *Task, payload string, host *Host, m *sync.Mutex) error {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "bash"
		flag = "-c"
	}

	// generate commands by using driver
	driver := Drivers[BuiltinDefaultDriverName]
	if task.Driver != "" {
		driver = Drivers[task.Driver]
		if driver == nil {
			return fmt.Errorf("invalid driver name '%s'", task.Driver)
		}
	}
	if debugFlag {
		fmt.Printf("[essh debug] driver: %s \n", driver.Name)
	}

	var script string
	content, err := driver.GenerateRunnableContent(task, host)
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
		prefix = "[Local] "
	} else if task.Prefix != "" {
		dict := map[string]interface{}{
			"Host": host,
			"Task": task,
		}
		tmpl, err := template.New("T").Parse(task.Prefix)
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

	cmd.Stdin = bytes.NewBufferString(payload)

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
	go scanLines(stdout, color.StdoutWriter, prefix, m)
	go scanLines(stderr, color.StderrWriter, prefix, m)

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

func runSSH(L *lua.LState, config string, args []string) error {
	// hooks
	var hooks map[string][]interface{}

	// Limitation!
	// hooks fires only when the hostname is just specified.
	if len(args) == 1 {
		hostname := args[0]
		if host := GetPublicHost(hostname); host != nil {
			hooks = host.Hooks
		}
	}

	// run before_connect hook
	if before := hooks["before_connect"]; before != nil {
		if debugFlag {
			fmt.Printf("[essh debug] run before_connect hook\n")
		}
		hookScript, err := getHookScript(L, before)
		if err != nil {
			return err
		}
		if debugFlag {
			fmt.Printf("[essh debug] before_connect hook script: %s\n", hookScript)
		}
		if err := runCommand(hookScript); err != nil {
			return err
		}
	}

	// register after_disconnect hook
	defer func() {
		// after hook
		if after := hooks["after_disconnect"]; after != nil {
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
	if afterConnect := hooks["after_connect"]; afterConnect != nil {
		hookScript, err := getHookScript(L, afterConnect)
		if err != nil {
			return err
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

	return cmd.Run()
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

func runSCP(config string, args []string) error {
	if debugFlag {
		fmt.Printf("[essh debug] use scp mode.\n")
	}

	if len(args) < 2 {
		return fmt.Errorf("scp mode requires 2 parameters at least.")
	}

	// In the scp mode.
	// the arguments must be scp command options and args.
	sshCommandArgs := []string{"-F", config}
	sshCommandArgs = append(sshCommandArgs, args[:]...)

	// execute ssh commmand
	cmd := exec.Command("scp", sshCommandArgs[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if debugFlag {
		fmt.Printf("[essh debug] real ssh command: %v \n", cmd.Args)
	}

	return cmd.Run()
}

func runRsync(config string, args []string) error {
	if debugFlag {
		fmt.Printf("[essh debug] use rsync mode.\n")
	}

	if len(args) < 1 {
		return fmt.Errorf("rsync mode requires 1 parameters at least.")
	}

	// In the rsync mode.
	// the arguments must be rsync command options and args.
	sshCommandArgs := []string{"-F", config}
	rsyncSSHOption := `-e "ssh ` + strings.Join(sshCommandArgs, " ") + `"`

	rsyncCommand := "rsync " + rsyncSSHOption + " " + strings.Join(args, " ")

	if debugFlag {
		fmt.Printf("[essh debug] real rsync command: %v\n", rsyncCommand)
	}

	return runCommand(rsyncCommand)
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
			return fmt.Errorf("'%s' is duplicated", host.Name)
		}
		names[host.Name] = true
	}

	for _, task := range Tasks {
		if _, ok := names[task.Name]; ok {
			return fmt.Errorf("'%s' is duplicated", task.Name)
		}
		names[task.Name] = true
	}

	for _, tag := range Tags() {
		if _, ok := names[tag]; ok {
			return fmt.Errorf("'%s' is duplicated", tag)
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

func removeModules() error {
	if !noGlobalFlag {
		c := NewContext(UserDataDir, ContextTypeGlobal)
		if _, err := os.Stat(c.ModulesDir()); err == nil {
			err = os.RemoveAll(c.ModulesDir())
			if err != nil {
				return err
			}
		}

		if _, err := os.Stat(c.TmpDir()); err == nil {
			err = os.RemoveAll(c.TmpDir())
			if err != nil {
				return err
			}
		}
	}

	c := NewContext(WorkingDataDir, ContextTypeLocal)
	if _, err := os.Stat(c.ModulesDir()); err == nil {
		err = os.RemoveAll(c.ModulesDir())
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(c.TmpDir()); err == nil {
		err = os.RemoveAll(c.TmpDir())
		if err != nil {
			return err
		}
	}

	return nil
}

func printHelp() {
	fmt.Print(`Usage: essh [<options>] [<ssh options and args...>]

  essh is an extended ssh command.
  version ` + Version + ` (` + CommitHash + `)

  Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
  The MIT License (MIT)

`)
	printOptionsInfo()
}

func printUsage() {
	fmt.Print(`Usage: essh [<options>] [<ssh options and args...>]

  Essh is an extended ssh command.
  version ` + Version + ` (` + CommitHash + `)

  Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
  The MIT License (MIT)

See also:
  essh --help

`)
}

func printOptionsInfo() {
	fmt.Print(`Options:
general options.
  --version                     Print version.
  --help                        Print help.
  --print                       Print generated ssh config.
  --gen                         Only generate ssh config.
  --config                      Edit config file in the current directory.
  --user-config                 Edit per-user config file.
  --system-config               Edit system wide config file.
  --working-dir <dir>           Change working directory.
  --debug                       Output debug log.

manage hosts, tags and tasks.
  --hosts                       List hosts.
  --tags                        List tags.
  --tasks                       List tasks.
  --quiet                       (Using with --hosts, --tasks or --tags option) Show only names.
  --filter <tag|host>           (Using with --hosts option) Use only the hosts filtered with a tag or a host.
  --all                         (Using with --hosts or --tasks option) Show all that includs hidden objects.

manage modules.
  --update                      Update modules.
  --clean                       Clean the downloaded modules.
  --no-global                   (Using with --update or --clean option) Update or clean only the modules about per-project config.

execute commands using hosts configuration.
  --exec                        Execute commands with the hosts.
  --on <tag|host>               (Using with --exec option) Run commands on remote hosts.
  --foreach <tag|host>          (Using with --exec option) Run commands locally for each hosts.
  --prefix                      (Using with --exec option) Enable outputing prefix.
  --prefix-string [<prefix>]    (Using with --exec option) Custom string of the prefix.
  --privileged                  (Using with --exec option) Run by the privileged user.
  --parallel                    (Using with --exec option) Run in parallel.
  --pty                         (Using with --exec option) Allocate pseudo-terminal. (add ssh option "-t -t" internally)
  --file                        (Using with --exec option) Load commands from a file.
  --driver                      (Using with --exec option) Specify a driver.

integrate other ssh related commands.
  --rsync                       Run rsync with essh configuration.
  --scp                         Run scp with essh configuration.

utility for zsh.
  --zsh-completion              Output zsh completion code.
  --aliases                     Output aliases code.

Github:
  https://github.com/kohkimakimoto/essh

`)
}

func getEditor() string {
	editor := os.Getenv("ESSH_EDITOR")
	if editor != "" {
		return editor
	}

	return os.Getenv("EDITOR")
}

func init() {
	// set SystemWideConfigFile
	SystemWideConfigFile = "/etc/essh/config.lua"

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
}

var ZSH_COMPLETION = `# This is zsh completion code.
# If you want to use it. write the following code in your '.zshrc'
#   eval "$(essh --zsh-completion)"
_essh_hosts() {
    local -a __essh_hosts
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_hosts=($(essh --zsh-completion-hosts | awk -F'\t' '{print $1":"$2}'))
    IFS=$PRE_IFS
    _describe -t host "host" __essh_hosts
}

_essh_tasks() {
    local -a __essh_tasks
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_tasks=($(essh --zsh-completion-tasks | awk -F'\t' '{print $1":"$2}'))
    IFS=$PRE_IFS
    _describe -t task "task" __essh_tasks
}

_essh_tags() {
    local -a __essh_tags
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_tags=($(essh --zsh-completion-tags))
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
        '--gen:Only generate ssh config.'
        '--update:Update modules.'
        '--clean:Clean the downloaded modules.'
        '--no-global:Update or clean only the modules about per-project config.'
        '--config:Edit config file in the current directory.'
        '--user-config:Edit per-user config file.'
        '--system-config:Edit system wide config file.'
        '--working-dir:Change working directory.'
        '--hosts:List hosts.'
        '--tags:List tags.'
        '--tasks:List tasks.'
        '--debug:Output debug log.'
        '--exec:Execute commands with the hosts.'
        '--rsync:Run rsync with essh configuration.'
        '--scp:Run scp with essh configuration.'
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
        '--filter:Use only the hosts filtered with a tag or a host'
     )
    _describe -t option "option" __essh_options
}

_essh_tags_options() {
    local -a __essh_options
    __essh_options=(
        '--debug:Output debug log.'
        '--quiet:Show only names.'
     )
    _describe -t option "option" __essh_options
}

_essh_exec_options() {
    local -a __essh_options
    __essh_options=(
        '--debug:Output debug log.'
        '--on:Run commands on remote hosts.'
        '--foreach:Run commands locally for each hosts.'
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
                --print|--help|--version|--gen|--config|--system-config)
                    ;;
                --file|--config-file)
                    _files
                    ;;
                --filter|--on|--foreach)
                    _essh_hosts
                    _essh_tags
                    ;;
                *)
                    if [ "$execMode" = "on" ]; then
                        _essh_exec_options
                    elif [ "$hostsMode" = "on" ]; then
                        _essh_hosts_options
                    elif [ "$tasksMode" = "on" ]; then
                        _essh_hosts_options
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

var ALIASES_CODE = `# This is aliaes code.
# If you want to use it. write the following code in your '.zshrc'
#   eval "$(essh --aliases)"
alias escp='essh --scp'
alias ersync='essh --rsync'
`

var BASH_COMPLETION = `
_essh_targets() {

}

_essh () {

}

complete -F _essh essh

`
