package zssh

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"fmt"
)

func Run(command string) error {
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

type writer struct {
	Type       int
	realWriter *RealWriter
}

type RealWriter struct {
	Prefix  string
	NewLine bool
}

func (w *RealWriter) Write(dataType int, data []byte) {
	dataStr := string(data)
	dataStr = strings.Replace(dataStr, "\r\n", "\n", -1)

	if w.NewLine {
		w.NewLine = false
		if dataType == 1 {
			fmt.Fprintf(os.Stdout, "%s", w.Prefix)
		} else {
			fmt.Fprintf(os.Stderr, "%s", w.Prefix)
		}
	}

	if strings.Contains(dataStr, "\n") {
		lineCount := strings.Count(dataStr, "\n")

		if dataStr[len(dataStr)-1:] == "\n" {
			w.NewLine = true
		}

		if w.NewLine {
			dataStr = strings.Replace(dataStr, "\n", "\n"+w.Prefix, lineCount-1)
		} else {
			dataStr = strings.Replace(dataStr, "\n", "\n"+w.Prefix, -1)
		}

		if dataType == 1 {
			fmt.Fprintf(os.Stdout, "%s", dataStr)
		} else {
			fmt.Fprintf(os.Stderr, "%s", dataStr)
		}

	} else {
		if dataType == 1 {
			fmt.Fprintf(os.Stdout, "%s", dataStr)
		} else {
			fmt.Fprintf(os.Stderr, "%s", dataStr)
		}
	}
}

func (w *writer) Write(data []byte) (int, error) {
	w.realWriter.Write(w.Type, data)

	return len(data), nil
}
