package essh

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"text/template"

	fatihColor "github.com/fatih/color"
	"github.com/kardianos/osext"
	"github.com/sevir/essh/support/color"
	"github.com/sevir/essh/support/helper"
	lua "github.com/yuin/gopher-lua"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// system configurations.
var (
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
	versionFlag bool
	helpFlag    bool
	printFlag   bool
	colorFlag   bool
	noColorFlag bool
	debugFlag   bool
	hostsFlag   bool
	quietFlag   bool
	allFlag     bool
	tagsFlag    bool
	tasksFlag   bool
	evalFlag    bool
	evalFileVar string
	menuFlag    bool
	genFlag     bool
	globalFlag  bool

	zshCompletionModeFlag       bool
	zshCompletionFlag           bool
	zshCompletionHostsFlag      bool
	zshCompletionTagsFlag       bool
	zshCompletionTasksFlag      bool
	zshCompletionNamespacesFlag bool

	bashCompletionModeFlag       bool
	bashCompletionFlag           bool
	bashCompletionHostsFlag      bool
	bashCompletionTagsFlag       bool
	bashCompletionTasksFlag      bool
	bashCompletionNamespacesFlag bool

	aliasesFlag     bool
	execFlag        bool
	fileFlag        bool
	prefixFlag      bool
	parallelFlag    bool
	privilegedFlag  bool
	userVar         string
	ptyFlag         bool
	SSHConfigFlag   bool
	workindDirVar   string
	configVar       string
	selectVar       []string
	targetVar       []string
	filterVar       []string
	backendVar      string
	prefixStringVar string
	driverVar       string
)

const (
	ExitErr = 1
)

func initResources() {
	// Flags
	helpFlag = false
	printFlag = false
	colorFlag = false
	noColorFlag = false
	debugFlag = false
	hostsFlag = false
	quietFlag = false
	allFlag = false
	tagsFlag = false
	tasksFlag = false
	menuFlag = false
	evalFlag = false
	evalFileVar = ""
	genFlag = false
	globalFlag = false
	zshCompletionModeFlag = false
	zshCompletionFlag = false
	zshCompletionHostsFlag = false
	zshCompletionTagsFlag = false
	zshCompletionTasksFlag = false
	zshCompletionNamespacesFlag = false
	bashCompletionModeFlag = false
	bashCompletionFlag = false
	bashCompletionHostsFlag = false
	bashCompletionTagsFlag = false
	bashCompletionTasksFlag = false
	bashCompletionNamespacesFlag = false
	aliasesFlag = false
	execFlag = false
	fileFlag = false
	prefixFlag = false
	parallelFlag = false
	privilegedFlag = false
	userVar = ""
	ptyFlag = false
	SSHConfigFlag = false
	workindDirVar = ""
	configVar = ""
	selectVar = []string{}
	targetVar = []string{}
	filterVar = []string{}
	backendVar = ""
	prefixStringVar = ""
	driverVar = ""

	// Registry
	CurrentRegistry = nil
	GlobalRegistry = nil
	LocalRegistry = nil

	// Hosts, Tasks, Drivers,
	Hosts = map[string]*Host{}
	Tasks = map[string]*Task{}
	Drivers = map[string]*Driver{}

	// set built-in drivers
	driver := NewDriver()
	driver.Name = DefaultDriverName
	driver.Engine = func(driver *Driver) (string, error) {
		return `
{{template "environment" .}}
{{template "functions" .}}
{{range $i, $script := .Scripts}}{{$script.code}}
{{end}}`, nil
	}
	Drivers[DefaultDriverName] = driver
	DefaultDriver = driver
}

// Add new types for the TUI
type item struct {
	name, desc string
	isHost     bool
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.name }

type model struct {
	list     list.Model
	selected item
	quitting bool
}

// New styles for the TUI
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0"))

	docStyle = lipgloss.NewStyle().Margin(1, 2)

	hostStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("25")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)

	taskStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("208")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected = i
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return "\n" + m.list.View()
}

// Custom delegate embedding the default delegate
type delegate struct {
	list.DefaultDelegate
}

func (d delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	itm, ok := listItem.(item)
	if !ok {
		return
	}

	var tagType string
	var tagStyle lipgloss.Style
	switch itm.isHost {
	case true:
		tagStyle = hostStyle
		tagType = "H"
	case false:
		tagStyle = taskStyle
		tagType = "T"
	default:
		tagStyle = hostStyle
	}

	tag := tagStyle.Render(tagType)
	desc := descStyle.Render(itm.desc)

	// If the item is selected, render the name with underline
	var name string
	if index == m.Index() {
		name = titleStyle.Underline(true).Render(itm.name)
		// Reset the title style to prevent affecting other items
		titleStyle = titleStyle.Underline(false)
	} else {
		name = titleStyle.Render(itm.name)
	}

	var symbol string
	if index == m.Index() {
		symbol = "»"
	} else {
		symbol = " "
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Left, symbol, tag, " ", name),
		desc,
	)

	fmt.Fprint(w, content)
}

// ---- end of new types for the TUI

func Run(osArgs []string) (exitStatus int) {
	defer func() {
		if e := recover(); e != nil {
			exitStatus = ExitErr
			if zshCompletionModeFlag && !debugFlag {
				// suppress printing error in running completion code.
				return
			}

			if bashCompletionModeFlag && !debugFlag {
				// suppress printing error in running completion code.
				return
			}

			printError(e)
		}
	}()

	initResources()

	if os.Getenv("ESSH_DEBUG") != "" {
		debugFlag = true
	}

	if len(osArgs) == 0 {
		printUsage()
		return

	}

	args := []string{}
	doesNotParseOption := false

	// parsing options
	// Essh uses only double dash options like `--print`,
	// because of preventing conflict to the ssh options.
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
		} else if arg == "--ssh-config" {
			SSHConfigFlag = true
		} else if arg == "--quiet" {
			quietFlag = true
		} else if arg == "--all" {
			allFlag = true
		} else if arg == "--tasks" {
			tasksFlag = true
		} else if arg == "--eval" {
			evalFlag = true
		} else if arg == "--eval-file" {
			if len(osArgs) < 2 {
				printError("--eval-file requires an argument.")
				return ExitErr
			}
			evalFileVar = osArgs[1]
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--eval-file=") {
			evalFileVar = strings.Split(arg, "=")[1]
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
		} else if arg == "--global" {
			globalFlag = true
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
			bashCompletionModeFlag = true
		} else if arg == "--bash-completion-hosts" {
			bashCompletionHostsFlag = true
			bashCompletionModeFlag = true
		} else if arg == "--bash-completion-tags" {
			bashCompletionTagsFlag = true
			bashCompletionModeFlag = true
		} else if arg == "--bash-completion-tasks" {
			bashCompletionTasksFlag = true
			bashCompletionModeFlag = true
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
		} else if arg == "--user" {
			if len(osArgs) < 2 {
				printError("--user reguires an argument.")
				return ExitErr
			}
			userVar = osArgs[1]
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--user=") {
			userVar = strings.Split(arg, "=")[1]
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
		} else if arg == "--filter" {
			if len(osArgs) < 2 {
				printError("--filter reguires an argument.")
				return ExitErr
			}
			filterVar = append(filterVar, osArgs[1])
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--filter=") {
			filterVar = append(filterVar, strings.Split(arg, "=")[1])
		} else if arg == "--backend" {
			if len(osArgs) < 2 {
				printError("--backend reguires an argument.")
				return ExitErr
			}
			backendVar = osArgs[1]
			osArgs = osArgs[1:]
		} else if strings.HasPrefix(arg, "--backend=") {
			backendVar = strings.Split(arg, "=")[1]
		} else if arg == "--script-file" {
			fileFlag = true
		} else if arg == "--pty" {
			ptyFlag = true
		} else if arg == "--menu" {
			menuFlag = true
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

	// If exists hidden file is preferred
	if _, err := os.Stat(filepath.Join(wd, ".esshconfig.lua")); err == nil {
		WorkingDirConfigFile = filepath.Join(wd, ".esshconfig.lua")
	}

	// use config file path from environment variable if it set.
	if configVar == "" && os.Getenv("ESSH_CONFIG") != "" {
		configVar = os.Getenv("ESSH_CONFIG")
	}

	// overwrite config file path by --config option.
	if configVar != "" {
		if filepath.IsAbs(configVar) {
			WorkingDirConfigFile = configVar
			WorkingDataDir = filepath.Join(filepath.Dir(WorkingDirConfigFile), ".essh")
		} else {
			WorkingDirConfigFile = filepath.Join(wd, configVar)
			WorkingDataDir = filepath.Join(filepath.Dir(WorkingDirConfigFile), ".essh")
		}

		if _, err := os.Stat(WorkingDirConfigFile); err != nil {
			printError(err)
			return ExitErr
		}
	}

	workingDirConfigFileBasename := filepath.Base(WorkingDirConfigFile)
	workingDirConfigFileDir := filepath.Dir(WorkingDirConfigFile)
	workingDirConfigFileBasenameExtension := filepath.Ext(workingDirConfigFileBasename)
	workingDirConfigFileName := workingDirConfigFileBasename[0 : len(workingDirConfigFileBasename)-len(workingDirConfigFileBasenameExtension)]

	WorkingDirOverrideConfigFile = filepath.Join(workingDirConfigFileDir, workingDirConfigFileName+"_override"+workingDirConfigFileBasenameExtension)

	if helpFlag {
		printHelp()
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

	if bashCompletionFlag {
		s, err := sprintByTemplate(BASH_COMPLETION)
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

	// extend lua package path.
	libdir := filepath.Join(UserDataDir, "lib")
	libdir2 := filepath.Join(WorkingDataDir, "lib")
	if os.PathSeparator == '/' { // unix-like
		lua.LuaPathDefault = libdir2 + "/?.lua;" + libdir + "/?.lua;" + "/usr/local/share/essh/lib/?.lua;" + lua.LuaPathDefault
	} else {
		lua.LuaPathDefault = libdir2 + "\\?.lua;" + libdir + "\\?.lua;" + lua.LuaPathDefault
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
		os.Remove(tmpFile.Name())

		if debugFlag {
			fmt.Printf("[essh debug] deleted config file: %s \n", tmpFile.Name())
		}
	}()

	temporarySSHConfigFile := tmpFile.Name()
	tmpFile.Close()

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
	GlobalRegistry = NewRegistry(UserDataDir, RegistryTypeGlobal)
	LocalRegistry = NewRegistry(WorkingDataDir, RegistryTypeLocal)

	CurrentRegistry = GlobalRegistry

	if _, err := os.Stat(WorkingDirConfigFile); err == nil && !globalFlag {
		// has working directroy config file

		// change context to working dir context
		CurrentRegistry = LocalRegistry

		// load working directory config
		if _, err := os.Stat(WorkingDirConfigFile); err == nil {
			if debugFlag {
				fmt.Printf("[essh debug] loading config file: %s\n", WorkingDirConfigFile)
			}

			if err := L.DoFile(WorkingDirConfigFile); err != nil {
				printError(err)
				return ExitErr
			}

			if debugFlag {
				fmt.Printf("[essh debug] loaded config file: %s\n", WorkingDirConfigFile)
			}
		}
	} else {
		// does not have working directory config file

		// load per-user configuration file.
		if _, err := os.Stat(UserConfigFile); err == nil {
			if debugFlag {
				fmt.Printf("[essh debug] loading config file: %s\n", UserConfigFile)
			}

			if err := L.DoFile(UserConfigFile); err != nil {
				printError(err)
				return ExitErr
			}

			if debugFlag {
				fmt.Printf("[essh debug] loaded config file: %s\n", UserConfigFile)
			}
		}
	}

	// change context to working dir context
	CurrentRegistry = LocalRegistry

	// load working directory override config
	if _, err := os.Stat(WorkingDirOverrideConfigFile); err == nil && !globalFlag {
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

	// change context to global
	CurrentRegistry = GlobalRegistry

	// load override global config
	if _, err := os.Stat(UserOverrideConfigFile); err == nil {
		if debugFlag {
			fmt.Printf("[essh debug] loading config file: %s\n", UserOverrideConfigFile)
		}

		if err := L.DoFile(UserOverrideConfigFile); err != nil {
			printError(err)
			return ExitErr
		}

		if debugFlag {
			fmt.Printf("[essh debug] loaded config file: %s\n", UserOverrideConfigFile)
		}
	}

	// validate config
	if err := validateResources(NewTaskQuery().Datasource, NewHostQuery().Datasource); err != nil {
		printError(err)
		return ExitErr
	}

	// show hosts for zsh completion
	if zshCompletionHostsFlag {
		for _, host := range NewHostQuery().GetHostsOrderByName() {
			if !host.Hidden {
				fmt.Printf("%s\t%s\n", ColonEscape(host.Name), ColonEscape(host.DescriptionOrDefault()))
			}
		}

		return
	}

	if bashCompletionHostsFlag {
		for _, host := range NewHostQuery().GetHostsOrderByName() {
			if !host.Hidden {
				fmt.Printf("%s\n", ColonEscape(host.Name))
			}
		}

		return
	}

	// show tasks for zsh completion
	if zshCompletionTasksFlag {
		for _, t := range NewTaskQuery().GetTasksOrderByName() {
			hidden := t.Hidden
			if !t.Disabled && !hidden {
				fmt.Printf("%s\t%s\n", ColonEscape(t.PublicName()), ColonEscape(t.DescriptionOrDefault()))
			}
		}
		return
	}

	if bashCompletionTasksFlag {
		for _, t := range NewTaskQuery().GetTasksOrderByName() {
			hidden := t.Hidden
			if !t.Disabled && !hidden {
				fmt.Printf("%s\n", ColonEscape(t.PublicName()))
			}
		}
		return
	}

	if zshCompletionTagsFlag || bashCompletionTagsFlag {
		for _, tag := range GetTags(Hosts) {
			fmt.Printf("%s\n", ColonEscape(tag))
		}
		return
	}

	// only print hosts list
	if hostsFlag {
		if len(selectVar) == 0 && len(filterVar) > 0 {
			printError("--filter must be used with --select option.")
			return ExitErr
		}

		query := NewHostQuery().AppendSelections(selectVar).AppendFilters(filterVar)
		if !allFlag {
			query = query.isVisible()
		}
		filteredHosts := query.GetHostsOrderByName()

		if SSHConfigFlag {
			outputConfig, ok := toString(lessh.RawGetString("ssh_config"))
			if !ok {
				printError(fmt.Errorf("invalid value %v in the 'ssh_config'", lessh.RawGetString("ssh_config")))
				return ExitErr
			}

			// generate ssh hosts config
			content, err := UpdateSSHConfig(outputConfig, filteredHosts)
			if err != nil {
				printError(err)
				return ExitErr
			}

			// print generated config
			fmt.Println(string(content))
		} else {
			tb := helper.NewPlainTable(os.Stdout)
			if !quietFlag {
				tb.SetHeader([]string{"NAME", "DESCRIPTION", "TAGS", "HIDDEN"})
			}

			for _, host := range filteredHosts {
				if quietFlag {
					tb.Append([]string{host.Name})
				} else {
					hidden := "false"
					if host.Hidden {
						hidden = "true"
					}
					tb.Append([]string{host.Name, host.Description, strings.Join(host.Tags, ","), hidden})
				}
			}

			tb.Render()
		}

		return
	}

	// only print tags list
	if tagsFlag {
		tb := helper.NewPlainTable(os.Stdout)
		if !quietFlag {
			tb.SetHeader([]string{"NAME"})
		}
		for _, tag := range GetTags(Hosts) {
			tb.Append([]string{tag})
		}
		tb.Render()

		return
	}

	// only print tasks list
	if tasksFlag {
		tb := helper.NewPlainTable(os.Stdout)
		if !quietFlag {
			tb.SetHeader([]string{"NAME", "DESCRIPTION", "HIDDEN"})
		}
		for _, t := range NewTaskQuery().GetTasksOrderByName() {
			hidden := t.Hidden
			if (!hidden && !t.Disabled) || allFlag {
				if quietFlag {
					tb.Append([]string{t.PublicName()})
				} else {
					tb.Append([]string{t.PublicName(), t.Description, fmt.Sprintf("%v", t.Hidden)})
				}
			}
		}
		tb.Render()

		return
	}

	// only eval lua code
	if evalFlag {

		var code string = "print('No Lua code is passed')"
		//lua code can be passed as arguments, standard input, or env variable ESSH_EVAL
		if os.Getenv("ESSH_EVAL") != "" {
			code = os.Getenv("ESSH_EVAL")
			// check if code is passed as standard input including all code
		} else if len(args) == 0 {
			stdin := bufio.NewReader(os.Stdin)
			code_buffer, err := io.ReadAll(stdin)

			if err != nil {
				printError(err)
				return ExitErr
			}
			code = string(code_buffer)
			// check if code is passed as arguments
		} else if len(args) > 0 {
			code = strings.Join(args, "\n")
		}

		//print code before eval
		if debugFlag {
			fmt.Println("[essh debug] lua code to be executed:")
			fmt.Println(code)
			fmt.Println("[essh debug] end of lua code")
			fmt.Println("[essh debug] executing lua code:")
		}

		if err := L.DoString(code); err != nil {
			printError(err)
			return ExitErr
		}

		return
	}

	// only eval lua code from file
	if evalFileVar != "" {
		if debugFlag {
			fmt.Printf("[essh debug] evaluating lua file: %s\n", evalFileVar)
		}

		if _, err := os.Stat(evalFileVar); os.IsNotExist(err) {
			printError(fmt.Errorf("file not found: %s", evalFileVar))
			return ExitErr
		}

		// detect arguments and pass them to lua script
		argTable := L.NewTable()
		scriptPath := os.Args[2]
		L.SetTable(argTable, lua.LNumber(0), lua.LString(scriptPath))
		luaArgs := os.Args[3:]
		for i, arg := range luaArgs {
			L.SetTable(argTable, lua.LNumber(i+1), lua.LString(arg))
		}

		L.SetGlobal("arg", argTable)

		if err := L.DoFile(evalFileVar); err != nil {
			//printError(err)
			print(fmt.Sprintf("error in file %s: %v", evalFileVar, err))
			return ExitErr
		}

		return
	}

	outputConfig, ok := toString(lessh.RawGetString("ssh_config"))
	if !ok {
		printError(fmt.Errorf("invalid value %v in the 'ssh_config'", lessh.RawGetString("ssh_config")))
		return ExitErr
	}

	// generate ssh hosts config
	content, err := UpdateSSHConfig(outputConfig, NewHostQuery().GetHostsOrderByName())
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

	if menuFlag {
		// Create list of items combining hosts and tasks
		items := []list.Item{}

		// Add hosts first
		hosts := NewHostQuery().GetHostsOrderByName()
		for _, host := range hosts {
			if !host.Hidden {
				items = append(items, item{
					name:   host.Name,
					desc:   host.DescriptionOrDefault(),
					isHost: true,
				})
			}
		}

		// Add tasks
		tasks := NewTaskQuery().GetTasksOrderByName()
		for _, task := range tasks {
			if !task.Hidden && !task.Disabled {
				items = append(items, item{
					name:   task.PublicName(),
					desc:   task.DescriptionOrDefault(),
					isHost: false,
				})
			}
		}

		// Setup list

		// delegate := list.NewDefaultDelegate()
		// delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		// 	Background(lipgloss.Color("57")).
		// 	Foreground(lipgloss.Color("230"))

		l := list.New(items, list.NewDefaultDelegate(), 80, 20) // Set width and height to ensure visibility
		l.Title = "ESSH Hosts and Tasks"
		l.Styles.Title = titleStyle
		l.SetShowStatusBar(true)
		l.SetFilteringEnabled(true)
		l.Styles.Title = titleStyle

		// Custom delegate to render items with tags
		l.SetDelegate(delegate{list.NewDefaultDelegate()})
		l.Styles.NoItems = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

		m := model{list: l}

		p := tea.NewProgram(m, tea.WithAltScreen())
		finalModel, err := p.Run()
		if err != nil {
			printError(err)
			return ExitErr
		}

		if m, ok := finalModel.(model); ok && !m.quitting {
			selected := m.selected
			if selected.isHost {
				// Run SSH for host
				err, ex := runSSH(L, outputConfig, []string{selected.name})
				if err != nil {
					printError(err)
					return ExitErr
				}
				return ex
			} else {
				// Run task
				task := GetEnabledTask(selected.name)
				if task != nil {
					err := runTask(outputConfig, task, []string{}, L)
					if err != nil {
						printError(err)
						return ExitErr
					}
				}
			}
		}
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
		task.User = userVar
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

		if len(targetVar) == 0 && len(filterVar) > 0 {
			printError("--filter must be used with --target option.")
			return ExitErr
		}

		task.Targets = targetVar
		task.Filters = filterVar

		if prefixFlag || prefixStringVar != "" {
			task.UsePrefix = true
		}

		if prefixStringVar != "" {
			task.Prefix = prefixStringVar
		}

		err := runTask(outputConfig, task, []string{}, L)
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
				var taskargs []string
				if len(args) >= 2 {
					taskargs = args[1:]
				} else {
					taskargs = []string{}
				}

				err := runTask(outputConfig, task, taskargs, L)
				if err != nil {
					printError(err)
					return ExitErr
				}
				return
			}
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

// this code is borrowed from https://github.com/fujiwara/nssh/blob/master/nssh.go
func processStdin(chs []chan []byte) {
	buf := make([]byte, 1024)
	for {
		n, err := io.ReadAtLeast(os.Stdin, buf, 1)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, color.FgRB("essh error in reading stdin: %v\n", err))
			}
			break
		}
		for _, ch := range chs {
			ch <- buf[0:n]
		}
	}

	// STDIN is EOF. close channels
	for _, ch := range chs {
		close(ch)
	}
}

// this code is borrowed from https://github.com/fujiwara/nssh/blob/master/nssh.go
func handleInput(stdinCh chan []byte, dest io.WriteCloser) {
	for {
		b, more := <-stdinCh
		if more {
			_, err := dest.Write(b)
			if err != nil {
				if e, ok := err.(*os.PathError); ok && e.Err == syscall.EPIPE {
					// broken pipe. suppress and ignore this error.
					dest.Close()
					break
				} else {
					fmt.Fprintf(os.Stderr, color.FgRB("essh error in writing stdin: %v (data: %v)\n", err, b))
					dest.Close()
					break
				}
			}
		} else {
			dest.Close()
			break
		}
	}
}

// this code is borrowed from https://github.com/fujiwara/nssh/blob/master/nssh.go
func scanLines(src io.ReadCloser, dest io.Writer, prefix string, m *sync.Mutex) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		// prevent mixing data in a line.
		m.Lock()
		if prefix != "" {
			fmt.Fprintf(dest, "%s%s\n", color.FgCB(prefix), scanner.Text())
		} else {
			fmt.Fprintf(dest, "%s\n", scanner.Text())
		}
		m.Unlock()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, color.FgRB("essh error: scanner.Scan() returns error: %v\n", err))
	}
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
			return "", fmt.Errorf("hook function return value must be string or function")
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
		shell = "bash"
		flag = "-c"
	}
	cmd := exec.Command(shell, flag, command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func validateResources(tasks map[string]*Task, hosts map[string]*Host) error {
	// check duplication of the host, task and tag names
	for _, task := range tasks {
		taskName := task.PublicName()
		if _, ok := hosts[taskName]; ok {
			return fmt.Errorf("Task '%s' is duplicated with hostname.", taskName)
		}
	}

	tags := GetTags(hosts)
	for _, tag := range tags {
		if _, ok := hosts[tag]; ok {
			return fmt.Errorf("Tag '%s' is duplicated with hostname.", tag)
		}
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

func printUsage() {
	fmt.Print(`Usage: essh [<options>] [<ssh options and args...>]

Essh is an extended ssh command.
version ` + Version + ` (` + CommitHash + `)

Original from Kohki Makimoto <kohki.makimoto@gmail.com>
Forked and extended by José F. Rives <jose@sevir.org>

The MIT License (MIT)

See more detail, use '--help'.

`)
}

//go:embed help_template.txt
var helpTextTemplate string

func printHelp() {
	fmt.Printf(helpTextTemplate, Version, CommitHash)
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

//go:embed zsh_completion.sh
var ZSH_COMPLETION string

//go:embed bash_completion.sh
var BASH_COMPLETION string

var ALIASES_CODE = `# This is aliases code.
# If you want to use it. write the following code in your '.zshrc'
#   eval "$(essh --aliases)"
function escp() {
    {{.Executable}} --exec 'scp -F $ESSH_SSH_CONFIG' "$@"
}
function ersync() {
    {{.Executable}} --exec 'rsync -e "ssh -F $ESSH_SSH_CONFIG"' "$@"
}
`
