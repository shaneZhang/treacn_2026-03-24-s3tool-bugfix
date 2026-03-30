package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestMultipartCmd(t *testing.T) {
	if multipartCmd == nil {
		t.Fatal("multipartCmd should not be nil")
	}

	if multipartCmd.Use != "multipart" {
		t.Errorf("multipartCmd.Use = %v, want multipart", multipartCmd.Use)
	}
}

func TestMultipartInitCmd(t *testing.T) {
	if multipartInitCmd == nil {
		t.Fatal("multipartInitCmd should not be nil")
	}

	if multipartInitCmd.Use != "init [bucket] [key]" {
		t.Errorf("multipartInitCmd.Use = %v, want init [bucket] [key]", multipartInitCmd.Use)
	}
}

func TestMultipartUploadCmd(t *testing.T) {
	if multipartUploadCmd == nil {
		t.Fatal("multipartUploadCmd should not be nil")
	}

	if multipartUploadCmd.Use != "upload [bucket] [key] [upload-id] [part-number] [file]" {
		t.Errorf("multipartUploadCmd.Use = %v, want upload [bucket] [key] [upload-id] [part-number] [file]", multipartUploadCmd.Use)
	}

	partSizeFlag := multipartUploadCmd.Flags().Lookup("part-size")
	if partSizeFlag == nil {
		t.Error("part-size flag should exist")
	}
}

func TestMultipartListCmd(t *testing.T) {
	if multipartListCmd == nil {
		t.Fatal("multipartListCmd should not be nil")
	}

	expectedUse := "list [bucket] [key] [upload-id]"
	if multipartListCmd.Use != expectedUse {
		t.Errorf("multipartListCmd.Use = %v, want %s", multipartListCmd.Use, expectedUse)
	}
}

func TestMultipartCompleteCmd(t *testing.T) {
	if multipartCompleteCmd == nil {
		t.Fatal("multipartCompleteCmd should not be nil")
	}

	expectedUse := "complete [bucket] [key] [upload-id]"
	if multipartCompleteCmd.Use != expectedUse {
		t.Errorf("multipartCompleteCmd.Use = %v, want %s", multipartCompleteCmd.Use, expectedUse)
	}
}

func TestMultipartAbortCmd(t *testing.T) {
	if multipartAbortCmd == nil {
		t.Fatal("multipartAbortCmd should not be nil")
	}

	expectedUse := "abort [bucket] [key] [upload-id]"
	if multipartAbortCmd.Use != expectedUse {
		t.Errorf("multipartAbortCmd.Use = %v, want %s", multipartAbortCmd.Use, expectedUse)
	}
}

func TestMultipartCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"init",
		"upload",
		"list",
		"complete",
		"abort",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range multipartCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in multipart command", subName)
		}
	}
}

func TestMultipartInitCmd_ArgsValidation(t *testing.T) {
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
			err := multipartInitCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMultipartUploadCmd_ArgsValidation(t *testing.T) {
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
			name:    "四个参数",
			args:    []string{"bucket", "key", "upload-id", "1"},
			wantErr: true,
		},
		{
			name:    "五个参数",
			args:    []string{"bucket", "key", "upload-id", "1", "file.txt"},
			wantErr: false,
		},
		{
			name:    "六个参数",
			args:    []string{"bucket", "key", "upload-id", "1", "file.txt", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := multipartUploadCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMultipartListCmd_ArgsValidation(t *testing.T) {
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
			wantErr: true,
		},
		{
			name:    "三个参数",
			args:    []string{"bucket", "key", "upload-id"},
			wantErr: false,
		},
		{
			name:    "四个参数",
			args:    []string{"bucket", "key", "upload-id", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := multipartListCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMultipartCompleteCmd_ArgsValidation(t *testing.T) {
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
			wantErr: true,
		},
		{
			name:    "三个参数",
			args:    []string{"bucket", "key", "upload-id"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := multipartCompleteCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMultipartAbortCmd_ArgsValidation(t *testing.T) {
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
			wantErr: true,
		},
		{
			name:    "三个参数",
			args:    []string{"bucket", "key", "upload-id"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := multipartAbortCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
