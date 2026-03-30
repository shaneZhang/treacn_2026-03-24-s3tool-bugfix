package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestEncryptionCmd(t *testing.T) {
	if encryptionCmd == nil {
		t.Fatal("encryptionCmd should not be nil")
	}

	if encryptionCmd.Use != "encryption" {
		t.Errorf("encryptionCmd.Use = %v, want encryption", encryptionCmd.Use)
	}
}

func TestEncryptionGetCmd(t *testing.T) {
	if encryptionGetCmd == nil {
		t.Fatal("encryptionGetCmd should not be nil")
	}

	if encryptionGetCmd.Use != "get [bucket]" {
		t.Errorf("encryptionGetCmd.Use = %v, want get [bucket]", encryptionGetCmd.Use)
	}
}

func TestEncryptionEnableCmd(t *testing.T) {
	if encryptionEnableCmd == nil {
		t.Fatal("encryptionEnableCmd should not be nil")
	}

	if encryptionEnableCmd.Use != "enable [bucket]" {
		t.Errorf("encryptionEnableCmd.Use = %v, want enable [bucket]", encryptionEnableCmd.Use)
	}
}

func TestEncryptionDisableCmd(t *testing.T) {
	if encryptionDisableCmd == nil {
		t.Fatal("encryptionDisableCmd should not be nil")
	}

	if encryptionDisableCmd.Use != "disable [bucket]" {
		t.Errorf("encryptionDisableCmd.Use = %v, want disable [bucket]", encryptionDisableCmd.Use)
	}
}

func TestEncryptionCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"get",
		"enable",
		"disable",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range encryptionCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in encryption command", subName)
		}
	}
}

func TestEncryptionGetCmd_ArgsValidation(t *testing.T) {
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
			err := encryptionGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptionEnableCmd_ArgsValidation(t *testing.T) {
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
			err := encryptionEnableCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptionDisableCmd_ArgsValidation(t *testing.T) {
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
			err := encryptionDisableCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
