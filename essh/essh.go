package essh

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/yuin/gopher-lua"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"encoding/json"
	"runtime"
)

// system configurations.
var (
	SystemWideConfigFile string
	PerUserConfigFile    string
	CurrentDirConfigFile    string
)

// flags
var (
	versionFlag bool
	helpFlag bool
	printFlag bool
	configFlag bool
	systemConfigFlag bool
	debugFlag bool
	hostsFlag bool
	verboseFlag bool
	tagsFlag bool
	genFlag bool
	zshCompletinFlag bool
	bashCompletinFlag bool
	shellFlag bool
	rsyncFlag bool
	scpFlag bool

	configFile string
	filters []string = []string{}
	format string
)

func Start() error {
	if len(os.Args) == 1 {
		printUsage()
		return nil
	}

	args := os.Args[1:]
	for {
		if len(args) == 0 {
			break
		}

		arg := args[0]
		if arg == "--print" {
			printFlag = true
		} else if arg == "--version" {
			versionFlag = true
		} else if arg == "--help" {
			helpFlag = true
		} else if arg == "--config" {
			configFlag = true
		} else if arg == "--system-config" {
			systemConfigFlag = true
		} else if arg == "--debug" {
			debugFlag = true
		} else if arg == "--hosts" {
			hostsFlag = true
		} else if arg == "--verbose" {
			verboseFlag = true
		} else if arg == "--filter" {
			if len(args) < 2 {
				return fmt.Errorf("--filter reguires an argument.")
			}
			filters = append(filters, args[1])
			args = args[1:]
		} else if strings.HasPrefix(arg, "--filter=") {
			filters = append(filters, strings.Split(arg, "=")[1])
		} else if arg == "--format" {
			if len(args) < 2 {
				return fmt.Errorf("--format reguires an argument.")
			}
			format = args[1]
			args = args[1:]
		} else if strings.HasPrefix(arg, "--format=") {
			format = strings.Split(arg, "=")[1]
		} else if arg == "--tags" {
			tagsFlag = true
		} else if arg == "--gen" {
			genFlag = true
		} else if arg == "--zsh-completion" {
			zshCompletinFlag = true
		} else if arg == "--bash-completion" {
			bashCompletinFlag = true
		} else if arg == "--config-file" {
			if len(args) < 2 {
				return fmt.Errorf("--config-file reguires an argument.")
			}
			configFile = args[1]
			args = args[1:]
		} else if strings.HasPrefix(arg, "--config-file=") {
			configFile = strings.Split(arg, "=")[1]
		} else if arg == "--shell" {
			shellFlag = true
		} else if arg == "--rsync" {
			rsyncFlag = true
		} else if arg == "--scp" {
			scpFlag = true
		} else {
			break
		}

		args = args[1:]
	}

	if helpFlag {
		printUsage()
		return nil
	}
	
	if versionFlag {
		fmt.Printf("%s (%s)\n", Version, CommitHash)
		return nil
	}

	if zshCompletinFlag {
		fmt.Print(ZSH_COMPLETION)
		return nil
	}

	if bashCompletinFlag {
		fmt.Print(BASH_COMPLETION)
		return nil
	}

	if configFlag {
		shellExec("$EDITOR " + PerUserConfigFile)
		return nil
	}

	if systemConfigFlag {
		shellExec("$EDITOR " + SystemWideConfigFile)
		return nil
	}

	// set up the lua state.
	L := lua.NewState()
	defer L.Close()

	// init lua state
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

	// set temporary ssh config file path
	lessh.RawSetString("ssh_config", lua.LString(temporarySSHConfigFile))

	// load specific config file
	if configFile != "" {
		_, err := os.Stat(configFile)
		if err != nil {
			return err
		}

		if err := L.DoFile(configFile); err != nil {
			return err
		}

		if debugFlag {
			fmt.Printf("[essh debug] loaded config file: %s \n", configFile)
		}

	} else {
		// load system wide config
		if _, err := os.Stat(SystemWideConfigFile); err == nil {
			if err := L.DoFile(SystemWideConfigFile); err != nil {
				return err
			}

			if debugFlag {
				fmt.Printf("[essh debug] loaded config file: %s \n", SystemWideConfigFile)
			}
		}

		// load per-user wide config
		if _, err := os.Stat(PerUserConfigFile); err == nil {
			if err := L.DoFile(PerUserConfigFile); err != nil {
				return err
			}

			if debugFlag {
				fmt.Printf("[essh debug] loaded config file: %s \n", PerUserConfigFile)
			}
		}

		// load current dir config
		if CurrentDirConfigFile != "" {
			if _, err := os.Stat(CurrentDirConfigFile); err == nil {
				if err := L.DoFile(CurrentDirConfigFile); err != nil {
					return err
				}

				if debugFlag {
					fmt.Printf("[essh debug] loaded config file: %s \n", CurrentDirConfigFile)
				}
			}
		}
	}


	// generate ssh hosts config
	content, err := GenHostsConfig()
	if err != nil {
		return err
	}

	// only print generated config
	if printFlag {
		fmt.Println(string(content))
		return nil
	}

	// only print hosts list
	if hostsFlag {
		var hosts []*Host
		if len(filters) > 0 {
			hosts = HostsByTags(filters)
		} else {
			hosts = Hosts
		}

		if format == "json" {
			printJson(hosts, "")
		} else if format == "prettyjson" {
			printJson(hosts, "    ")
		} else {
			for _, host := range hosts {
				if !host.Hidden {
					if verboseFlag {
						fmt.Printf("%s\t%s\n", host.Name, host.Description)
					} else {
						fmt.Printf("%s\n", host.Name)
					}
				}
			}
		}

		return nil
	}

	// only print tags list
	if tagsFlag {
		for _, tag := range Tags() {
			fmt.Printf("%s\n", tag)
		}

		return nil
	}

	outputConfig, ok := toString(lessh.RawGetString("ssh_config"))
	if !ok {
		return fmt.Errorf("invalid value %v in the 'ssh_config'", lessh.RawGetString("ssh_config"))
	}

	if debugFlag {
		fmt.Printf("[essh debug] output ssh_config contents to the file: %s \n", outputConfig)
	}

	// update temporary ssh config file
	err = ioutil.WriteFile(outputConfig, content, 0644)
	if err != nil {
		return err
	}

	// only generating contents
	if genFlag {
		return nil
	}

	// select running mode and run it.
	if shellFlag {
		err = runShellScript(outputConfig, args)
	} else if rsyncFlag {
		err = runRsync(outputConfig, args)
	} else if scpFlag {
		err = runSCP(outputConfig, args)
	} else {
		err = runSSH(outputConfig, args)
	}

	return err
}

func printJson(hosts []*Host, indent string) {
	convHosts := []map[string]map[string]interface{}{}

	for _, host :=range hosts {
		h := map[string]map[string]interface{}{}

		hv := map[string]interface{}{}
		for _, pair := range host.Values() {
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


func runSSH(config string, args []string) error {
	// hooks
	var hooks map[string]interface{}

	// Limitation!
	// hooks fires only when the hostname is specified by the last argument.
	if len(args) > 0 {
		hostname := args[len(args)-1]
		if host := GetHost(hostname); host != nil {
			hooks = host.Hooks
		}
	}

	// run before hook
	if before := hooks["before_connect"]; before != nil {
		if debugFlag {
			fmt.Printf("[essh debug] run before_connect hook\n")
		}
		err := runHook(before)
		if err != nil {
			return err
		}
	} else if before := hooks["before"]; before != nil {
		// for backward compatibility
		if debugFlag {
			fmt.Printf("[essh debug] run before hook\n")
		}
		err := runHook(before)
		if err != nil {
			return err
		}
	}

	// register after hook
	defer func() {
		// after hook
		if after := hooks["after_disconnect"]; after != nil {
			if debugFlag {
				fmt.Printf("[essh debug] run after_disconnect hook\n")
			}
			err := runHook(after)
			if err != nil {
				panic(err)
			}
		} else if after := hooks["after"]; after != nil {
			// for backward compatibility
			if debugFlag {
				fmt.Printf("[essh debug] run after hook\n")
			}
			err := runHook(after)
			if err != nil {
				panic(err)
			}
		}
	}()

	// setup ssh command args
	var sshComandArgs []string

	if afterConnect := hooks["after_connect"]; afterConnect != nil {
		sshComandArgs = []string{"-t", "-F", config}
		sshComandArgs = append(sshComandArgs, args[:]...)

		script := afterConnect.(string)
		script += `
exec $SHELL
`
		sshComandArgs = append(sshComandArgs, script)

	} else {
		sshComandArgs = []string{"-F", config}
		sshComandArgs = append(sshComandArgs, args[:]...)
	}

	// execute ssh commmand
	cmd := exec.Command("ssh", sshComandArgs[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if debugFlag {
		fmt.Printf("[essh debug] real ssh command: %v \n", cmd.Args)
	}

	return cmd.Run()
}

func runHook(hook interface{}) error {
	if hookFunc, ok := hook.(func() error); ok {
		err := hookFunc()
		if err != nil {
			return err
		}
	} else if hookString, ok := hook.(string); ok {
		err := runCommand(hookString)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid type hook: %v", hook)
	}
	return nil
}


func runShellScript(config string, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("shell script mode requires 2 parameters at least.")
	}

	// In the shell script mode.
	// the last argument must be a script file path.
	shellPath := args[len(args)-1]
	// remove it
	args = args[:len(args)-1]

	var scriptContent []byte
	if strings.HasPrefix(shellPath, "http://") || strings.HasPrefix(shellPath, "https://") {
		// get script from remote using http.
		if debugFlag {
			fmt.Printf("[essh debug] get script using http from '%s'\n", shellPath)
		}

		var httpClient *http.Client
		if strings.HasPrefix(shellPath, "https://") {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			httpClient = &http.Client{Transport: tr}
		} else {
			httpClient = &http.Client{}
		}

		resp, err := httpClient.Get(shellPath)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		scriptContent = b
	} else {
		// get script from the file system.
		b, err := ioutil.ReadFile(shellPath)
		if err != nil {
			return err
		}
		scriptContent = b
	}

	if debugFlag {
		fmt.Printf("[essh debug] script:\n%s\n", string(scriptContent))
	}

	// setup ssh command args
	sshComandArgs := []string{"-F", config}
	sshComandArgs = append(sshComandArgs, args[:]...)
	sshComandArgs = append(sshComandArgs, "bash", "-se")

	cmd := exec.Command("ssh", sshComandArgs[:]...)
	cmd.Stdin = bytes.NewBuffer(scriptContent)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if debugFlag {
		fmt.Printf("[essh debug] real ssh command: %v \n", cmd.Args)
	}

	return cmd.Run()
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
	sshComandArgs := []string{"-F", config}
	sshComandArgs = append(sshComandArgs, args[:]...)

	// execute ssh commmand
	cmd := exec.Command("scp", sshComandArgs[:]...)
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
	sshComandArgs := []string{"-F", config}
	rsyncSSHOption := `-e "ssh ` + strings.Join(sshComandArgs, " ") + `"`

	rsyncCommand := "rsync " + rsyncSSHOption + " " + strings.Join(args, " ")

	if debugFlag {
		fmt.Printf("[essh debug] real rsync command: %v\n", rsyncCommand)
	}

	return shellExec(rsyncCommand)
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
	cmd.Stdin =  os.Stdin

	return cmd.Run()
}


func printUsage() {
	// print usage.
	fmt.Println(`Usage: essh [<options>] [<ssh options and args...>]

essh is an extended ssh command.
version ` + Version + ` (` + CommitHash + `)

Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
The MIT License (MIT)

Options:
  --version               Print version.
  --help                  Print help.
  --print                 Print generated ssh config.
  --gen                   Only generating ssh config.
  --config                Edit per-user config file.
  --system-config         Edit system wide config file.
  --config-file <file>    Load configuration from the specific file.
                          If you use this option, it does not use other default config files like a "/etc/essh/config.lua".

  --hosts                 List hosts. This option can use with additional options.
  --filter <tag>          (Using with --hosts option) Show only the hosts configured with a tag.
  --verbose               (Using with --hosts option) List hosts with description.

  --tags                  List tags.
  --format <format>       (Using with --hosts or --tags option) Output specified format (json|prettyjson)

  --zsh-completion        Output zsh completion code.
  --debug                 Output debug log

  --shell     Change behavior to execute a shell script on the remote host.
              Take a look "Running shell script" section.
  --rsync     Change behavior to execute rsync.
              Take a look "Running rsync" section.
  --scp       Change behavior to execute scp.
              Take a look "Running scp" section.

Running shell script:
  ESSH supports easily running a bash script on the remote server.
  Syntax:

    essh --shell [<ssh options and args...> <script path|script url>

  Examples:

    essh --shell web01.localhost /path/to/script.sh
    essh --shell web01.localhost https://example/script.sh

Running rsyc:
  You can use essh config for rsync using --rsync option.
  Syntax:

    essh --rsync <rsync options and args...>

  Examples:

    essh --rsync -avz /local/dir/ web01.localhost:/path/to/remote/dir

Running scp:
  You can use essh config for scp using --scp option.
  Syntax:

    essh --scp <scp options and args...>

  Examples:

    essh --scp web01.localhost:/path/to/file ./local/file

See also:
  ssh, rsync, scp
`)
}

func init() {
	if SystemWideConfigFile == "" {
		SystemWideConfigFile = "/etc/essh/config.lua"
	}
	if PerUserConfigFile == "" {
		home := userHomeDir()
		PerUserConfigFile = filepath.Join(home, ".essh/config.lua")
	}

	if CurrentDirConfigFile == "" {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("couldn't get working dir %v\n", err)
		} else {
			CurrentDirConfigFile = filepath.Join(wd, ".essh.lua")
		}
	}
}

var ZSH_COMPLETION = `
_essh_hosts() {
    local -a __essh_hosts
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_hosts=($(essh --hosts --verbose | awk -F'\t' '{print $1":"$2}'))
    IFS=$PRE_IFS
    _describe -t host "host" __essh_hosts
}

_essh () {
    local curcontext="$curcontext" state line
    typeset -A opt_args

    _arguments \
        '1: :->command'

    case $state in
        command)
            _essh_hosts
            ;;
        *)
            _files
            ;;
    esac
}

compdef _essh essh

`

var BASH_COMPLETION = `
_essh_hosts() {

}

_essh () {

}

complete -F _essh essh

`

