package zssh

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var ConfigFile string
var SSHConfigFile string
var Version = "0.2.1"


func Main() int {
	log.SetFlags(0)

	if len(os.Args) == 1 {
		fmt.Println(`zssh: extended ssh command.

version ` + Version + `

zssh custom options:
  --print	Print generated ssh config.
  --config	Edit config file.
  --hosts	List hosts.
  --macros	List macros.
  --update	Only update ssh config file. doesn't run ssh command.
  --zsh-completion	Output zsh completion code.
`)
		// show ssh help
		Run("ssh")
		return 0
	}

	var args []string
	if len(os.Args) >= 2 {
		// remove the command name
		args = os.Args[1:]
	}

	firstArg := args[0]

	printFlag := false
	updateFlag := false
	hostsFlag := false
	macrosFlag := false
	configFlag := false
	zshCompletinFlag := false

	for _, arg := range args {
		if arg == "--print" {
			printFlag = true
		}
		if arg == "--update" {
			updateFlag = true
		}
		if arg == "--hosts" {
			hostsFlag = true
		}
		if arg == "--macros" {
			macrosFlag = true
		}
		if arg == "--config" {
			configFlag = true
		}
		if arg == "--zsh-completion" {
			zshCompletinFlag = true
		}
	}

	if zshCompletinFlag {
		fmt.Print(ZSH_COMPLETION)
		return 0
	}

	if configFlag {
		Run("$EDITOR " + ConfigFile)
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

	if hostsFlag {
		for _, host := range Hosts {
			if !host.Hidden {
				if host.Description != "" {
					fmt.Printf("%s\t%s\n", host.Name, host.Description)
				} else {
					fmt.Printf("%s\n", host.Name)
				}
			}
		}

		return 0
	}

	if macrosFlag {
		for _, macro := range Macros {
			if macro.Description != "" {
				fmt.Printf("%s\t%s\n", macro.Name, macro.Description)
			} else {
				fmt.Printf("%s\n", macro.Name)
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

	if macro, err := GetMacro(firstArg); err == nil {
		// there is a macro
		err := macro.Run()
		if err != nil {
			log.Printf("Error: %s", err)
			return 1
		}
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
	err = Run(cmdline)

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
		ConfigFile = filepath.Join(home, ".ssh/zssh.lua")
	}

	if SSHConfigFile == "" {
		home := userHomeDir()
		SSHConfigFile = filepath.Join(home, ".ssh/config")
	}

}

var ZSH_COMPLETION = `
_zssh_hosts() {
    local -a __zssh_hosts
    local -a __zssh_macros
    PRE_IFS=$IFS
    IFS=$'\n'
    __zssh_hosts=($(zssh --hosts | awk -F'\t' '{print $1":"$2}'))
    __zssh_macros=($(zssh --macros | awk -F'\t' '{print $1":"$2}'))
    IFS=$PRE_IFS
    _describe -t host "host" __zssh_hosts
    _describe -t macro "macro" __zssh_macros
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
