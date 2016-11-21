package essh

import "testing"

func TestUserHomeDir(t *testing.T) {
	dir := userHomeDir()
	if dir == "" {
		t.Error("home dir is empty")
	}
}

func TestEnvKeyEscape(t *testing.T) {
	if str := EnvKeyEscape("AAA-BBB"); str != "AAA_BBB" {
		t.Errorf("invalid result %s", str)
	}

	if str := EnvKeyEscape("AAA.BBB"); str != "AAA_BBB" {
		t.Errorf("invalid result %s", str)
	}
}
