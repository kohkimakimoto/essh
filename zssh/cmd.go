package zssh

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
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

type CallbackFunc func(stdout string, stderr string)

type Writer struct {
	CallbackFunc CallbackFunc
	Type         int
}

func (w *Writer) Write(data []byte) (int, error) {
	if w.CallbackFunc != nil {
		for _, s := range strings.Split(string(data), "\n") {
			if w.Type == 1 {
				// stdout
				w.CallbackFunc(s, "")
			} else {
				// stderr
				w.CallbackFunc("", s)
			}
		}
	}
	return len(data), nil
}

func RunWithCallback(command string, callback CallbackFunc) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	outWriter := &Writer{
		CallbackFunc: callback,
		Type:         1,
	}

	errWriter := &Writer{
		CallbackFunc: callback,
		Type:         2,
	}

	cmd.Stdout = outWriter
	cmd.Stderr = errWriter
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
