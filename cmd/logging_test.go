package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestLoggingCmd(t *testing.T) {
	if loggingCmd == nil {
		t.Fatal("loggingCmd should not be nil")
	}

	if loggingCmd.Use != "logging" {
		t.Errorf("loggingCmd.Use = %v, want logging", loggingCmd.Use)
	}
}

func TestLoggingGetCmd(t *testing.T) {
	if loggingGetCmd == nil {
		t.Fatal("loggingGetCmd should not be nil")
	}

	if loggingGetCmd.Use != "get [bucket]" {
		t.Errorf("loggingGetCmd.Use = %v, want get [bucket]", loggingGetCmd.Use)
	}
}

func TestLoggingDisableCmd(t *testing.T) {
	if loggingDisableCmd == nil {
		t.Fatal("loggingDisableCmd should not be nil")
	}

	if loggingDisableCmd.Use != "disable [bucket]" {
		t.Errorf("loggingDisableCmd.Use = %v, want disable [bucket]", loggingDisableCmd.Use)
	}
}

func TestLoggingCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"get",
		"disable",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range loggingCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in logging command", subName)
		}
	}
}

func TestLoggingGetCmd_ArgsValidation(t *testing.T) {
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
			err := loggingGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoggingDisableCmd_ArgsValidation(t *testing.T) {
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
			err := loggingDisableCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
