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
  --update-only	Only update ssh config. doesn't run ssh command.
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

	print := false
	updateOnly := false
	for _, arg := range args {
		if arg == "--print" {
			print = true
		}
		if arg == "--update-only" {
			updateOnly = true
		}

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

	if print {
		fmt.Println(string(content))
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

	if updateOnly {
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