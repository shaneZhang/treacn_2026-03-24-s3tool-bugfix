package cmd

import (
	"testing"
)

func TestInitCmd(t *testing.T) {
	if initCmd == nil {
		t.Fatal("initCmd should not be nil")
	}

	if initCmd.Use != "init" {
		t.Errorf("initCmd.Use = %v, want init", initCmd.Use)
	}
}

func TestInitCmd_Args(t *testing.T) {
	if initCmd.Args != nil {
		err := initCmd.Args(nil, []string{"extra"})
		if err == nil {
			t.Error("initCmd should not accept extra arguments")
		}
	}
}

func TestInitCmd_RunE(t *testing.T) {
	if initCmd.RunE == nil {
		t.Error("initCmd.RunE should not be nil")
	}
}
