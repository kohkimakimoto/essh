package docker

import (
	"bytes"
	"github.com/kohkimakimoto/essh/essh"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLocalExecTaskInDocker(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer func() {
		os.RemoveAll(tmpDir)
	}()
	configFile := filepath.Join(tmpDir, "esshconfig.lua")
	if err = ioutil.WriteFile(configFile, []byte(testConfig), 0644); err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	os.Chdir(tmpDir)

	// capture stdout
	// borrowed from http://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// test just running
	exitCode := essh.Run([]string{"example"})
	if exitCode != 0 {
		t.Error("exited with non zero")
	}

	// back to normal state
	w.Close()
	os.Stdout = old
	out := <-outC

	t.Log(out)
}

var testConfig string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	modDir := filepath.Dir(filename)

	testConfig = `
	local docker = import "` + modDir + `"

	driver "docker-centos7" {
		engine = docker.driver,
		image = "centos:centos7",
		remove_terminated_container = true,
	}

	task "example" {
		backend = "local",
		driver = "docker-centos7",
		script = {
			"cat /etc/redhat-release",
			"echo foo",
		},
	}
`

}
