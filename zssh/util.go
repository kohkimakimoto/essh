package zssh

import (
	"os"

	"runtime"
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

//
//func getStringFlagValue(args []string, flags ...string) string {
//	if len(args) != 0 {
//		for i, s := range args {
//			if s == "--" {
//				break
//			}
//			for _, flag := range flags {
//				if len(s) >= len(flag) && s[0:len(flag)] == flag {
//					if len(s) >= len(flag + "=") && s[0:len(flag + "=")] == flag + "=" {
//						return strings.Split(s, "=")[1], i
//					} else if s == flag {
//						if len(args) <= (i + 1) {
//							return "", -1
//						}
//						return args[i + 1], i + 1
//					}
//				}
//			}
//
//		}
//	}
//
//	return "", -1
//}
//
//func getBoolFlagValue(args []string, flags ...string) bool {
//	if len(args) != 0 {
//		for i, s := range args {
//			if s == "--" {
//				break
//			}
//			for _, flag := range flags {
//				if s == flag {
//					return true, i
//				}
//			}
//		}
//	}
//
//	return false, -1
//}
