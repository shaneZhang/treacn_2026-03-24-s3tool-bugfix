package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCorsCmd(t *testing.T) {
	if corsCmd == nil {
		t.Fatal("corsCmd should not be nil")
	}

	if corsCmd.Use != "cors" {
		t.Errorf("corsCmd.Use = %v, want cors", corsCmd.Use)
	}
}

func TestCorsGetCmd(t *testing.T) {
	if corsGetCmd == nil {
		t.Fatal("corsGetCmd should not be nil")
	}

	if corsGetCmd.Use != "get [bucket]" {
		t.Errorf("corsGetCmd.Use = %v, want get [bucket]", corsGetCmd.Use)
	}
}

func TestCorsDeleteCmd(t *testing.T) {
	if corsDeleteCmd == nil {
		t.Fatal("corsDeleteCmd should not be nil")
	}

	if corsDeleteCmd.Use != "delete [bucket]" {
		t.Errorf("corsDeleteCmd.Use = %v, want delete [bucket]", corsDeleteCmd.Use)
	}
}

func TestCorsCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"get",
		"delete",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range corsCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in cors command", subName)
		}
	}
}

func TestCorsGetCmd_ArgsValidation(t *testing.T) {
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
			err := corsGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCorsDeleteCmd_ArgsValidation(t *testing.T) {
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
			err := corsDeleteCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
