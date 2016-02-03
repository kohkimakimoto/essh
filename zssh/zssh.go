package zssh

import (
	"bytes"
	"crypto/tls"
	"flag"
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
	Version              = "0.5.0"
)

var IgnoreError flag.ErrorHandling = 9999

func Start() error {
	var printFlag, configFlag, systemConfigFlag, debugFlag, hostsFlag, verboseFlag, tagsFlag, zshCompletinFlag, bashCompletinFlag bool
	var configFile, shellPath string
	var rsyncArg string
	filters := []string{}

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
			if len(args) < 2 {
				return fmt.Errorf("--shell reguires an argument.")
			}
			shellPath = args[1]
			args = args[1:]
		} else if strings.HasPrefix(arg, "--shell=") {
			shellPath = strings.Split(arg, "=")[1]
		} else if arg == "--rsync" {
			if len(args) < 2 {
				return fmt.Errorf("--rsync reguires an argument.")
			}
			rsyncArg = args[1]
			args = args[1:]
		} else if strings.HasPrefix(arg, "--rsync=") {
			rsyncArg = strings.Split(arg, "=")[1]
		} else {
			break
		}

		args = args[1:]
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
		ShellExec("$EDITOR " + PerUserConfigFile)
		return nil
	}

	if systemConfigFlag {
		ShellExec("$EDITOR " + SystemWideConfigFile)
		return nil
	}

	var shellContent []byte

	if shellPath != "" {
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

			shellContent = b
		} else {
			// get script from the file system.
			b, err := ioutil.ReadFile(shellPath)
			if err != nil {
				return err
			}
			shellContent = b
		}
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

	if shellPath == "" && rsyncArg == "" {
		// get hooks
		var hooks map[string]func() error

		// Limitation!: hooks fires only when the hostname is specified by the first argument.
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
	}

	// setup ssh command args
	sshComandArgs := []string{"-F", generatedSSHConfigFile}
	sshComandArgs = append(sshComandArgs, args[:]...)
	if shellPath != "" {
		sshComandArgs = append(sshComandArgs, "bash", "-se")
	}

	if rsyncArg != "" {
		// rsync mode.
		if debugFlag {
			fmt.Printf("[zssh debug] use rsync mode.\n")
		}

		rsyncSSHOption := `-e "ssh ` + strings.Join(sshComandArgs, " ") + `"`
		rsyncCommand := "rsync " + rsyncSSHOption + " " + rsyncArg

		if debugFlag {
			fmt.Printf("[zssh debug] real rsync command: %v\n", rsyncCommand)
		}

		err := ShellExec(rsyncCommand)
		if err != nil {
			return err
		}
	} else {
		// normal ssh mode.

		// execute ssh commmand
		cmd := exec.Command("ssh", sshComandArgs[:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if shellPath != "" {
			// input shell script from stdin.
			cmd.Stdin = bytes.NewBuffer(shellContent)

			if debugFlag {
				fmt.Printf("[zssh debug] shell content:\n%s\n", string(shellContent))
			}

		} else {
			cmd.Stdin = os.Stdin
		}

		if debugFlag {
			fmt.Printf("[zssh debug] real ssh command: %v \n", cmd.Args)
		}

		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func printUsage() {
	// print usage.
	fmt.Println(`Usage: zssh [<options>] [<ssh options and args...>]

zssh is an extended ssh command.
version ` + Version + `

zssh options:
  --print                 Print generated ssh config.
  --config                Edit per-user config file.
  --system-config         Edit system wide config file.
  --config-file <FILE>    Load configuration from the specific file.
  --hosts                 List hosts. This option can use with additional options.
  --filter <TAG>          (Using with --hosts option) Show only the hosts configured with a tag.
  --verbose               (Using with --hosts option) List hosts with description.
  --tags                  List tags.
  --zsh-completion        Output zsh completion code.
  --debug                 Output debug log

zssh options for convenient functionalities:
  --shell <PATH> <HOSTNAME>              Execute a shell script of the path on the remote host.
  --rsync <rsync options and args...>    Execute rsync using zssh config.


And the following is original ssh command usage...
`)
	// show ssh help
	cmd := exec.Command("ssh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
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

