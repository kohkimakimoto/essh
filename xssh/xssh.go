package xssh

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"log"
	"path/filepath"
	"github.com/yuin/gopher-lua"
	"fmt"
	"io/ioutil"
)

var ConfigFile string
var SSHConfigFile string

func Main() int {
	if len(os.Args) == 1 {
		fmt.Println(`xssh: extended ssh command.

Custom options:
  --print:	Print generated ssh config.
  --list	List hosts.
  --update	Only update ssh config file. doesn't run ssh command.
		`)
		// show ssh help
		run("ssh")
		return 0
	}

	var args []string
	if len(os.Args) >= 2 {
		// remove the command name
		args = os.Args[1:]
	}

	printFlag := false
	updateFlag := false
	listFlag := false
	zshCompletinFlag := false

	for _, arg := range args {
		if arg == "--print" {
			printFlag = true
		}
		if arg == "--update" {
			updateFlag = true
		}
		if arg == "--list" {
			listFlag = true
		}
		if arg == "--zsh-completion" {
			zshCompletinFlag = true
		}
	}

	if zshCompletinFlag {
		fmt.Print(ZSH_COMPLETION)
		return 0
	}

	lstate := lua.NewState()
	defer lstate.Close()

	LoadFunctions(lstate)

	if _, err := os.Stat(ConfigFile); err == nil {
		if err := lstate.DoFile(ConfigFile); err != nil {
			log.Printf("Error: %s", err)
			return 1
		}
	}

	content, err := GenHostsConfig()
	if err != nil {
		log.Printf("Error: %s", err)
		return 1
	}

	if printFlag {
		fmt.Println(string(content))
		if !updateFlag {
			return 0
		}
	}

	if !printFlag && listFlag {
		for _, host := range Hosts {
			if host.Description != "" {
				fmt.Printf("%s\t%s\n", host.Name, host.Description)
			} else {
				fmt.Printf("%s\n", host.Name)
			}
		}

		return 0
	}

	// check modification.
	isModified := true
	if _, err := os.Stat(SSHConfigFile); err == nil {
		b, err := ioutil.ReadFile(SSHConfigFile)
		if err != nil {
			log.Printf("Error: %s", err)
			return 1
		}

		if string(b) == string(content) {
			isModified = false
		}
	}

	// update .ssh/config
	if isModified {
		err = ioutil.WriteFile(SSHConfigFile, content, 0644)
		if err != nil {
			log.Printf("Error: %s", err)
			return 1
		}
	}

	if updateFlag {
		return 0
	}

	// setup ssh command
	cmdline := "ssh " + strings.Join(args, " ")

	// got hooks
	var hooks map[string]func() error
	if len(args) >= 1 {
		hostname := args[0]
		if host := GetHost(hostname); host != nil {
			hooks = host.Hooks
		}
	}

	// before hook
	if before := hooks["before"]; before != nil {
		err := before()
		if err != nil {
			log.Printf("Error: %s", err)
			return 1
		}
	}

	// run ssh
	err = run(cmdline)

	// after hook
	if after := hooks["after"]; after != nil {
		err := after()
		if err != nil {
			log.Printf("Error: %s", err)
			return 1
		}
	}

	if err != nil {
		return 1
	}

	return 0
}


func run(command string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func init() {
	if ConfigFile == "" {
		home := userHomeDir()
		ConfigFile = filepath.Join(home, ".ssh/xssh.lua")
	}

	if SSHConfigFile == "" {
		home := userHomeDir()
		SSHConfigFile = filepath.Join(home, ".ssh/config")
	}

}

var ZSH_COMPLETION = `
_xssh_hosts() {
    local -a __xssh_hosts
    PRE_IFS=$IFS
    IFS=$'\n'
    __xssh_hosts=($(xssh --list | awk -F'\t' '{print $1":"$2}'))
    IFS=$PRE_IFS
    _describe -t commands "xssh_hosts" __xssh_hosts
}

_xssh () {
    local curcontext="$curcontext" state line
    typeset -A opt_args

    _arguments \
        '1: :->command'

    case $state in
        command)
            _xssh_hosts
            ;;
        *)
            _files
            ;;
    esac
}

compdef _xssh xssh

`