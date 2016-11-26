package essh

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestShowSSHConfig(t *testing.T) {
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

	// test just running
	exitCode := Run([]string{"--hosts", "--select=essh_test_ssh_server", "--ssh-config"})
	if exitCode != 0 {
		t.Error("exited with non zero")
	}
}

func TestRunSSH(t *testing.T) {
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

	// test
	exitCode := Run([]string{"essh_test_ssh_server1", "bash -c 'echo -n connected'"})
	if exitCode != 0 {
		t.Error("exited with non zero")
	}

	// back to normal state
	w.Close()
	os.Stdout = old
	out := <-outC

	if out != "connected" {
		t.Errorf("invalid output: %v", out)
	}
}

var testConfig string

func init() {
	port := os.Getenv("TEST_SSH_SERVER_PORT")
	testConfig = `
	host "essh_test_ssh_server1" {
		HostName = "127.0.0.1",
		Port = "` + port + `",
		User = "root",
		StrictHostKeyChecking = "no",
		UserKnownHostsFile = "/dev/null",
		LogLevel = "ERROR",
	}

	host "essh_test_ssh_server2" {
		HostName = "127.0.0.1",
		Port = "` + port + `",
		User = "root",
		StrictHostKeyChecking = "no",
		UserKnownHostsFile = "/dev/null",
		LogLevel = "ERROR",
	}
	
	host "essh_test_ssh_server3" {
		HostName = "127.0.0.1",
		Port = "` + port + `",
		User = "root",
		StrictHostKeyChecking = "no",
		UserKnownHostsFile = "/dev/null",
		LogLevel = "ERROR",
	}
`

}
