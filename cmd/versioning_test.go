package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestVersioningCmd(t *testing.T) {
	if versioningCmd == nil {
		t.Fatal("versioningCmd should not be nil")
	}

	if versioningCmd.Use != "versioning" {
		t.Errorf("versioningCmd.Use = %v, want versioning", versioningCmd.Use)
	}
}

func TestVersioningGetCmd(t *testing.T) {
	if versioningGetCmd == nil {
		t.Fatal("versioningGetCmd should not be nil")
	}

	if versioningGetCmd.Use != "get [bucket]" {
		t.Errorf("versioningGetCmd.Use = %v, want get [bucket]", versioningGetCmd.Use)
	}
}

func TestVersioningEnableCmd(t *testing.T) {
	if versioningEnableCmd == nil {
		t.Fatal("versioningEnableCmd should not be nil")
	}

	if versioningEnableCmd.Use != "enable [bucket]" {
		t.Errorf("versioningEnableCmd.Use = %v, want enable [bucket]", versioningEnableCmd.Use)
	}
}

func TestVersioningSuspendCmd(t *testing.T) {
	if versioningSuspendCmd == nil {
		t.Fatal("versioningSuspendCmd should not be nil")
	}

	if versioningSuspendCmd.Use != "suspend [bucket]" {
		t.Errorf("versioningSuspendCmd.Use = %v, want suspend [bucket]", versioningSuspendCmd.Use)
	}
}

func TestVersioningCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"get",
		"enable",
		"suspend",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range versioningCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in versioning command", subName)
		}
	}
}

func TestVersioningGetCmd_ArgsValidation(t *testing.T) {
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
			err := versioningGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVersioningEnableCmd_ArgsValidation(t *testing.T) {
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
			err := versioningEnableCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVersioningSuspendCmd_ArgsValidation(t *testing.T) {
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
			err := versioningSuspendCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
