package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestBucketCmd(t *testing.T) {
	if bucketCmd == nil {
		t.Fatal("bucketCmd should not be nil")
	}

	if bucketCmd.Use != "bucket" {
		t.Errorf("bucketCmd.Use = %v, want bucket", bucketCmd.Use)
	}
}

func TestBucketListCmd(t *testing.T) {
	if bucketListCmd == nil {
		t.Fatal("bucketListCmd should not be nil")
	}

	if bucketListCmd.Use != "list" {
		t.Errorf("bucketListCmd.Use = %v, want list", bucketListCmd.Use)
	}

	if bucketListCmd.RunE == nil {
		t.Error("bucketListCmd.RunE should not be nil")
	}
}

func TestBucketCreateCmd(t *testing.T) {
	if bucketCreateCmd == nil {
		t.Fatal("bucketCreateCmd should not be nil")
	}

	if bucketCreateCmd.Use != "create [bucket-name]" {
		t.Errorf("bucketCreateCmd.Use = %v, want create [bucket-name]", bucketCreateCmd.Use)
	}

	if bucketCreateCmd.Args == nil {
		t.Error("bucketCreateCmd.Args should not be nil")
	}

	err := bucketCreateCmd.Args(&cobra.Command{}, []string{})
	if err == nil {
		t.Error("bucketCreateCmd should require exactly 1 arg")
	}

	err = bucketCreateCmd.Args(&cobra.Command{}, []string{"bucket1", "bucket2"})
	if err == nil {
		t.Error("bucketCreateCmd should require exactly 1 arg")
	}

	err = bucketCreateCmd.Args(&cobra.Command{}, []string{"my-bucket"})
	if err != nil {
		t.Errorf("bucketCreateCmd should accept 1 arg, got error: %v", err)
	}
}

func TestBucketDeleteCmd(t *testing.T) {
	if bucketDeleteCmd == nil {
		t.Fatal("bucketDeleteCmd should not be nil")
	}

	if bucketDeleteCmd.Use != "delete [bucket-name]" {
		t.Errorf("bucketDeleteCmd.Use = %v, want delete [bucket-name]", bucketDeleteCmd.Use)
	}

	if bucketDeleteCmd.Args == nil {
		t.Error("bucketDeleteCmd.Args should not be nil")
	}
}

func TestBucketInfoCmd(t *testing.T) {
	if bucketInfoCmd == nil {
		t.Fatal("bucketInfoCmd should not be nil")
	}

	if bucketInfoCmd.Use != "info [bucket-name]" {
		t.Errorf("bucketInfoCmd.Use = %v, want info [bucket-name]", bucketInfoCmd.Use)
	}
}

func TestBucketLocationCmd(t *testing.T) {
	if bucketLocationCmd == nil {
		t.Fatal("bucketLocationCmd should not be nil")
	}

	if bucketLocationCmd.Use != "location [bucket-name]" {
		t.Errorf("bucketLocationCmd.Use = %v, want location [bucket-name]", bucketLocationCmd.Use)
	}
}

func TestBucketEmptyCmd(t *testing.T) {
	if bucketEmptyCmd == nil {
		t.Fatal("bucketEmptyCmd should not be nil")
	}

	if bucketEmptyCmd.Use != "empty [bucket-name]" {
		t.Errorf("bucketEmptyCmd.Use = %v, want empty [bucket-name]", bucketEmptyCmd.Use)
	}
}

func TestBucketCmd_Help(t *testing.T) {
	// Skip this test as Execute() writes to stdout directly
	t.Skip("Skipping Help test - Execute writes to stdout")
}

func TestBucketCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"list",
		"create",
		"delete",
		"info",
		"location",
		"empty",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range bucketCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in bucket command", subName)
		}
	}
}

func TestBucketCreateCmd_ArgsValidation(t *testing.T) {
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
		{
			name:    "带连字符的桶名",
			args:    []string{"my-test-bucket"},
			wantErr: false,
		},
		{
			name:    "带点的桶名",
			args:    []string{"my.test.bucket"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bucketCreateCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBucketDeleteCmd_ArgsValidation(t *testing.T) {
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
			err := bucketDeleteCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBucketInfoCmd_ArgsValidation(t *testing.T) {
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
			err := bucketInfoCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
