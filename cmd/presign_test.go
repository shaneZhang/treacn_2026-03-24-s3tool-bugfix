package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestPresignCmd(t *testing.T) {
	if presignCmd == nil {
		t.Fatal("presignCmd should not be nil")
	}

	if presignCmd.Use != "presign" {
		t.Errorf("presignCmd.Use = %v, want presign", presignCmd.Use)
	}
}

func TestPresignGetCmd(t *testing.T) {
	if presignGetCmd == nil {
		t.Fatal("presignGetCmd should not be nil")
	}

	if presignGetCmd.Use != "get [bucket] [key]" {
		t.Errorf("presignGetCmd.Use = %v, want get [bucket] [key]", presignGetCmd.Use)
	}

	expiresFlag := presignGetCmd.Flags().Lookup("expires")
	if expiresFlag == nil {
		t.Error("expires flag should exist")
	}
}

func TestPresignPutCmd(t *testing.T) {
	if presignPutCmd == nil {
		t.Fatal("presignPutCmd should not be nil")
	}

	if presignPutCmd.Use != "put [bucket] [key]" {
		t.Errorf("presignPutCmd.Use = %v, want put [bucket] [key]", presignPutCmd.Use)
	}

	expiresFlag := presignPutCmd.Flags().Lookup("expires")
	if expiresFlag == nil {
		t.Error("expires flag should exist")
	}
}

func TestPresignDeleteCmd(t *testing.T) {
	if presignDeleteCmd == nil {
		t.Fatal("presignDeleteCmd should not be nil")
	}

	if presignDeleteCmd.Use != "delete [bucket] [key]" {
		t.Errorf("presignDeleteCmd.Use = %v, want delete [bucket] [key]", presignDeleteCmd.Use)
	}

	expiresFlag := presignDeleteCmd.Flags().Lookup("expires")
	if expiresFlag == nil {
		t.Error("expires flag should exist")
	}
}

func TestPresignCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"get",
		"put",
		"delete",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range presignCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in presign command", subName)
		}
	}
}

func TestPresignGetCmd_ArgsValidation(t *testing.T) {
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
			args:    []string{"bucket"},
			wantErr: true,
		},
		{
			name:    "两个参数",
			args:    []string{"bucket", "key"},
			wantErr: false,
		},
		{
			name:    "三个参数",
			args:    []string{"bucket", "key", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := presignGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPresignPutCmd_ArgsValidation(t *testing.T) {
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
			name:    "两个参数",
			args:    []string{"bucket", "key"},
			wantErr: false,
		},
		{
			name:    "一个参数",
			args:    []string{"bucket"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := presignPutCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPresignDeleteCmd_ArgsValidation(t *testing.T) {
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
			name:    "两个参数",
			args:    []string{"bucket", "key"},
			wantErr: false,
		},
		{
			name:    "一个参数",
			args:    []string{"bucket"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := presignDeleteCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
