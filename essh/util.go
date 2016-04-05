package essh

import (
	"os"
	"runtime"
	"strings"
)

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

func ShellEscape(s string) string {
	return "'" + strings.Replace(s, "'", "'\"'\"'", -1) + "'"
}
