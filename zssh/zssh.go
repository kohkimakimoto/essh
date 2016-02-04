package zssh

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
)

// system configurations.
var (
	SystemWideConfigFile string
	PerUserConfigFile    string
)

// flags
var (
	versionFlag bool
	printFlag bool
	configFlag bool
	systemConfigFlag bool
	debugFlag bool
	hostsFlag bool
	verboseFlag bool
	tagsFlag bool
	zshCompletinFlag bool
	bashCompletinFlag bool
	shellFlag bool
	rsyncFlag bool
	scpFlag bool

	configFile string
	filters []string = []string{}
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
		} else if arg == "--tags" {
			tagsFlag = true
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

	// load lua custom functions
	LoadFunctions(L)

	if debugFlag {
		fmt.Printf("[zssh debug] loaded lua functions\n")
	}

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
			fmt.Printf("[zssh debug] loaded config file: %s \n", configFile)
		}

	} else {
		// load system wide config
		if _, err := os.Stat(SystemWideConfigFile); err == nil {
			if err := L.DoFile(SystemWideConfigFile); err != nil {
				return err
			}

			if debugFlag {
				fmt.Printf("[zssh debug] loaded config file: %s \n", SystemWideConfigFile)
			}

		}

		// load per-user wide config
		if _, err := os.Stat(PerUserConfigFile); err == nil {
			if err := L.DoFile(PerUserConfigFile); err != nil {
				return err
			}

			if debugFlag {
				fmt.Printf("[zssh debug] loaded config file: %s \n", PerUserConfigFile)
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

		for _, host := range hosts {
			if !host.Hidden {
				if verboseFlag {
					fmt.Printf("%s\t%s\n", host.Name, host.Description)
				} else {
					fmt.Printf("%s\n", host.Name)
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

	// generate temporary ssh config file
	tmpFile, err := ioutil.TempFile("", "zssh.ssh_config.")
	if err != nil {
		return err
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())

		if debugFlag {
			fmt.Printf("[zssh debug] deleted config file: %s \n", tmpFile.Name())
		}

	}()
	generatedSSHConfigFile := tmpFile.Name()

	if debugFlag {
		fmt.Printf("[zssh debug] generated config file: %s \n", generatedSSHConfigFile)
	}

	// update temporary sss config file
	err = ioutil.WriteFile(generatedSSHConfigFile, content, 0644)
	if err != nil {
		return err
	}


	// select running mode and run it.

	if shellFlag {
		err = runShellScript(generatedSSHConfigFile, args)
	} else if rsyncFlag {
		err = runRsync(generatedSSHConfigFile, args)
	} else if scpFlag {
		err = runSCP(generatedSSHConfigFile, args)
	} else {
		err = runSSH(generatedSSHConfigFile, args)
	}

	return err
}

func runSSH(config string, args []string) error {
	// hooks
	var hooks map[string]func() error

	// Limitation!
	// hooks fires only when the hostname is specified by the first argument.
	if len(args) > 0 {
		hostname := args[0]
		if host := GetHost(hostname); host != nil {
			hooks = host.Hooks
		}
	}

	// run before hook
	if before := hooks["before"]; before != nil {
		if debugFlag {
			fmt.Printf("[zssh debug] run before hook\n")
		}
		err := before()
		if err != nil {
			return err
		}
	}

	// register after hook
	defer func() {
		// after hook
		if after := hooks["after"]; after != nil {
			if debugFlag {
				fmt.Printf("[zssh debug] run after hook\n")
			}
			err := after()
			if err != nil {
				panic(err)
			}
		}
	}()

	// setup ssh command args
	sshComandArgs := []string{"-F", config}
	sshComandArgs = append(sshComandArgs, args[:]...)

	// execute ssh commmand
	cmd := exec.Command("ssh", sshComandArgs[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if debugFlag {
		fmt.Printf("[zssh debug] real ssh command: %v \n", cmd.Args)
	}

	return cmd.Run()
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
			fmt.Printf("[zssh debug] get script using http from '%s'\n", shellPath)
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
		fmt.Printf("[zssh debug] script:\n%s\n", string(scriptContent))
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
		fmt.Printf("[zssh debug] real ssh command: %v \n", cmd.Args)
	}

	return cmd.Run()
}

func runSCP(config string, args []string) error {
	if debugFlag {
		fmt.Printf("[zssh debug] use scp mode.\n")
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
		fmt.Printf("[zssh debug] real ssh command: %v \n", cmd.Args)
	}

	return cmd.Run()
}


func runRsync(config string, args []string) error {
	if debugFlag {
		fmt.Printf("[zssh debug] use rsync mode.\n")
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
		fmt.Printf("[zssh debug] real rsync command: %v\n", rsyncCommand)
	}

	return shellExec(rsyncCommand)
}

func printUsage() {
	// print usage.
	fmt.Println(`Usage: zssh [<options>] [<ssh options and args...>]

ZSSH is an extended ssh command.
version ` + Version + ` (` + CommitHash + `)

Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
The MIT License (MIT)

Options:
  --version               Print version.
  --print                 Print generated ssh config.
  --config                Edit per-user config file.
  --system-config         Edit system wide config file.
  --config-file <file>    Load configuration from the specific file.
                          If you use this option, it does not use other default config files like a "/etc/zssh/config.lua".

  --hosts                 List hosts. This option can use with additional options.
  --filter <tag>          (Using with --hosts option) Show only the hosts configured with a tag.
  --verbose               (Using with --hosts option) List hosts with description.

  --tags                  List tags.

  --zsh-completion        Output zsh completion code.
  --debug                 Output debug log

  --shell     Change behavior to execute a shell script on the remote host.
              Take a look "Running shell script" section.
  --rsync     Change behavior to execute rsync.
              Take a look "Running rsync" section.
  --scp       Change behavior to execute scp.
              Take a look "Running scp" section.

Running shell script:
  ZSSH supports easily running a bash script on the remote server.
  Syntax:

    zssh --shell [<ssh options and args...> <script path|script url>

  Examples:

    zssh --shell web01.localhost /path/to/script.sh
    zssh --shell web01.localhost https://example/script.sh

Running rsyc:
  You can use zssh config for rsync using --rsync option.
  Syntax:

    zssh --rsync <rsync options and args...>

  Examples:

    zssh --rsync -avz /local/dir/ web01.localhost:/path/to/remote/dir

Running scp:
  You can use zssh config for scp using --scp option.
  Syntax:

    zssh --scp <scp options and args...>

  Examples:

    zssh --scp web01.localhost:/path/to/file ./local/file

See also:
  ssh, rsync, scp
`)
}

func init() {
	if SystemWideConfigFile == "" {
		SystemWideConfigFile = "/etc/zssh/config.lua"
	}
	if PerUserConfigFile == "" {
		home := userHomeDir()
		PerUserConfigFile = filepath.Join(home, ".zssh/config.lua")
	}
}

var ZSH_COMPLETION = `
_zssh_hosts() {
    local -a __zssh_hosts
    PRE_IFS=$IFS
    IFS=$'\n'
    __zssh_hosts=($(zssh --hosts --verbose | awk -F'\t' '{print $1":"$2}'))
    IFS=$PRE_IFS
    _describe -t host "host" __zssh_hosts
}

_zssh () {
    local curcontext="$curcontext" state line
    typeset -A opt_args

    _arguments \
        '1: :->command'

    case $state in
        command)
            _zssh_hosts
            ;;
        *)
            _files
            ;;
    esac
}

compdef _zssh zssh

`

var BASH_COMPLETION = `
_zssh_hosts() {

}

_zssh () {

}

complete -F _zssh zssh

`

