package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestLifecycleCmd(t *testing.T) {
	if lifecycleCmd == nil {
		t.Fatal("lifecycleCmd should not be nil")
	}

	if lifecycleCmd.Use != "lifecycle" {
		t.Errorf("lifecycleCmd.Use = %v, want lifecycle", lifecycleCmd.Use)
	}
}

func TestLifecycleGetCmd(t *testing.T) {
	if lifecycleGetCmd == nil {
		t.Fatal("lifecycleGetCmd should not be nil")
	}

	if lifecycleGetCmd.Use != "get [bucket]" {
		t.Errorf("lifecycleGetCmd.Use = %v, want get [bucket]", lifecycleGetCmd.Use)
	}
}

func TestLifecycleDeleteCmd(t *testing.T) {
	if lifecycleDeleteCmd == nil {
		t.Fatal("lifecycleDeleteCmd should not be nil")
	}

	if lifecycleDeleteCmd.Use != "delete [bucket]" {
		t.Errorf("lifecycleDeleteCmd.Use = %v, want delete [bucket]", lifecycleDeleteCmd.Use)
	}
}

func TestLifecycleCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"get",
		"delete",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range lifecycleCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in lifecycle command", subName)
		}
	}
}

func TestLifecycleGetCmd_ArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "无参数",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "一个参数",
			args:    []string{"my-bucket"},
			wantErr: false,
		},
		{
			name:    "两个参数",
			args:    []string{"bucket1", "bucket2"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lifecycleGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLifecycleDeleteCmd_ArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "无参数",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "一个参数",
			args:    []string{"my-bucket"},
			wantErr: false,
		},
		{
			name:    "两个参数",
			args:    []string{"bucket1", "bucket2"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lifecycleDeleteCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
