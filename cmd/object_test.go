package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestObjectCmd(t *testing.T) {
	if objectCmd == nil {
		t.Fatal("objectCmd should not be nil")
	}

	if objectCmd.Use != "object" {
		t.Errorf("objectCmd.Use = %v, want object", objectCmd.Use)
	}
}

func TestObjectListCmd(t *testing.T) {
	if objectListCmd == nil {
		t.Fatal("objectListCmd should not be nil")
	}

	if objectListCmd.Use != "list [bucket]" {
		t.Errorf("objectListCmd.Use = %v, want list [bucket]", objectListCmd.Use)
	}

	prefixFlag := objectListCmd.Flags().Lookup("prefix")
	if prefixFlag == nil {
		t.Error("prefix flag should exist")
	}

	recursiveFlag := objectListCmd.Flags().Lookup("recursive")
	if recursiveFlag == nil {
		t.Error("recursive flag should exist")
	}

	maxKeysFlag := objectListCmd.Flags().Lookup("max-keys")
	if maxKeysFlag == nil {
		t.Error("max-keys flag should exist")
	}
}

func TestObjectPutCmd(t *testing.T) {
	if objectPutCmd == nil {
		t.Fatal("objectPutCmd should not be nil")
	}

	if objectPutCmd.Use != "put [bucket] [key] [file]" {
		t.Errorf("objectPutCmd.Use = %v, want put [bucket] [key] [file]", objectPutCmd.Use)
	}

	contentTypeFlag := objectPutCmd.Flags().Lookup("content-type")
	if contentTypeFlag == nil {
		t.Error("content-type flag should exist")
	}

	storageClassFlag := objectPutCmd.Flags().Lookup("storage-class")
	if storageClassFlag == nil {
		t.Error("storage-class flag should exist")
	}
}

func TestObjectGetCmd(t *testing.T) {
	if objectGetCmd == nil {
		t.Fatal("objectGetCmd should not be nil")
	}

	if objectGetCmd.Use != "get [bucket] [key] [local-file]" {
		t.Errorf("objectGetCmd.Use = %v, want get [bucket] [key] [local-file]", objectGetCmd.Use)
	}
}

func TestObjectDeleteCmd(t *testing.T) {
	if objectDeleteCmd == nil {
		t.Fatal("objectDeleteCmd should not be nil")
	}

	if objectDeleteCmd.Use != "delete [bucket] [key]" {
		t.Errorf("objectDeleteCmd.Use = %v, want delete [bucket] [key]", objectDeleteCmd.Use)
	}
}

func TestObjectCopyCmd(t *testing.T) {
	if objectCopyCmd == nil {
		t.Fatal("objectCopyCmd should not be nil")
	}

	if objectCopyCmd.Use != "copy [source-bucket] [source-key] [dest-bucket] [dest-key]" {
		t.Errorf("objectCopyCmd.Use = %v, want copy [source-bucket] [source-key] [dest-bucket] [dest-key]", objectCopyCmd.Use)
	}
}

func TestObjectInfoCmd(t *testing.T) {
	if objectInfoCmd == nil {
		t.Fatal("objectInfoCmd should not be nil")
	}

	if objectInfoCmd.Use != "info [bucket] [key]" {
		t.Errorf("objectInfoCmd.Use = %v, want info [bucket] [key]", objectInfoCmd.Use)
	}
}

func TestObjectUrlCmd(t *testing.T) {
	if objectUrlCmd == nil {
		t.Fatal("objectUrlCmd should not be nil")
	}

	if objectUrlCmd.Use != "url [bucket] [key]" {
		t.Errorf("objectUrlCmd.Use = %v, want url [bucket] [key]", objectUrlCmd.Use)
	}
}

func TestObjectMvCmd(t *testing.T) {
	if objectMvCmd == nil {
		t.Fatal("objectMvCmd should not be nil")
	}

	if objectMvCmd.Use != "mv [bucket] [source-key] [dest-key]" {
		t.Errorf("objectMvCmd.Use = %v, want mv [bucket] [source-key] [dest-key]", objectMvCmd.Use)
	}
}

func TestObjectCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"list",
		"put",
		"get",
		"delete",
		"copy",
		"info",
		"url",
		"mv",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range objectCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in object command", subName)
		}
	}
}

func TestObjectListCmd_ArgsValidation(t *testing.T) {
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
			args:    []string{"bucket1", "prefix"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := objectListCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestObjectPutCmd_ArgsValidation(t *testing.T) {
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
			wantErr: true,
		},
		{
			name:    "三个参数",
			args:    []string{"bucket", "key", "file.txt"},
			wantErr: false,
		},
		{
			name:    "四个参数",
			args:    []string{"bucket", "key", "file.txt", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := objectPutCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestObjectGetCmd_ArgsValidation(t *testing.T) {
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
			name:    "三个参数",
			args:    []string{"bucket", "key", "local-file"},
			wantErr: false,
		},
		{
			name:    "两个参数",
			args:    []string{"bucket", "key"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := objectGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestObjectDeleteCmd_ArgsValidation(t *testing.T) {
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
			err := objectDeleteCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestObjectCopyCmd_ArgsValidation(t *testing.T) {
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
			args:    []string{"src-bucket", "src-key", "dest-bucket", "dest-key"},
			wantErr: false,
		},
		{
			name:    "三个参数",
			args:    []string{"src-bucket", "src-key", "dest-bucket"},
			wantErr: true,
		},
		{
			name:    "五个参数",
			args:    []string{"src-bucket", "src-key", "dest-bucket", "dest-key", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := objectCopyCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestObjectMvCmd_ArgsValidation(t *testing.T) {
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
			name:    "三个参数",
			args:    []string{"bucket", "src-key", "dest-key"},
			wantErr: false,
		},
		{
			name:    "两个参数",
			args:    []string{"bucket", "src-key"},
			wantErr: true,
		},
		{
			name:    "四个参数",
			args:    []string{"bucket", "src-key", "dest-key", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := objectMvCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1024 * 1024, "1.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
		{1024 * 1024 * 1024 * 1024, "1.0 TB"},
	}

	for _, tt := range tests {
		result := formatBytes(tt.bytes)
		if result != tt.expected {
			t.Errorf("formatBytes(%d) = %v, want %v", tt.bytes, result, tt.expected)
		}
	}
}
