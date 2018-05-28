package main_test

import (
	"os"
	"testing"
	vs "vsconfig-go"
)

func TestMain(t *testing.T) {
	err := vs.Vsgoinit()
	if err != nil {
		t.Error(err)
	}
	pwd, err := vs.GetAbsDir()
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(pwd + "/.vscode/tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			t.Failed()
		} else {
			t.Error(err)
		}
	}
}
