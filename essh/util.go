package essh

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

func EnvKeyEscape(s string) string {
	return strings.Replace(strings.Replace(s, "-", "_", -1), ".", "_", -1)
}

func ColonEscape(s string) string {
	return strings.Replace(s, ":", "\\:", -1)
}

func GetContentFromPath(shellPath string) ([]byte, error) {
	var scriptContent []byte
	if strings.HasPrefix(shellPath, "http://") || strings.HasPrefix(shellPath, "https://") {
		// get script from remote using http.
		if debugFlag {
			fmt.Printf("[essh debug] get script using http from '%s'\n", shellPath)
		}

		var httpClient *http.Client = &http.Client{}
		//if strings.HasPrefix(shellPath, "https://") {
		//	tr := &http.Transport{
		//		// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//	}
		//	httpClient = &http.Client{Transport: tr}
		//} else {
		//	httpClient = &http.Client{}
		//}

		resp, err := httpClient.Get(shellPath)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		scriptContent = b
	} else {
		// get script from the file system.
		b, err := ioutil.ReadFile(shellPath)
		if err != nil {
			return nil, err
		}
		scriptContent = b
	}

	return scriptContent, nil
}
