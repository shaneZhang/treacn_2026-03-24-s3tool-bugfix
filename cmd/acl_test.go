package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAclCmd(t *testing.T) {
	if aclCmd == nil {
		t.Fatal("aclCmd should not be nil")
	}

	if aclCmd.Use != "acl" {
		t.Errorf("aclCmd.Use = %v, want acl", aclCmd.Use)
	}
}

func TestAclBucketGetCmd(t *testing.T) {
	if aclBucketGetCmd == nil {
		t.Fatal("aclBucketGetCmd should not be nil")
	}

	if aclBucketGetCmd.Use != "bucket-get [bucket]" {
		t.Errorf("aclBucketGetCmd.Use = %v, want bucket-get [bucket]", aclBucketGetCmd.Use)
	}
}

func TestAclBucketSetCmd(t *testing.T) {
	if aclBucketSetCmd == nil {
		t.Fatal("aclBucketSetCmd should not be nil")
	}

	if aclBucketSetCmd.Use != "bucket-set [bucket] [acl]" {
		t.Errorf("aclBucketSetCmd.Use = %v, want bucket-set [bucket] [acl]", aclBucketSetCmd.Use)
	}
}

func TestAclObjectGetCmd(t *testing.T) {
	if aclObjectGetCmd == nil {
		t.Fatal("aclObjectGetCmd should not be nil")
	}

	if aclObjectGetCmd.Use != "object-get [bucket] [key]" {
		t.Errorf("aclObjectGetCmd.Use = %v, want object-get [bucket] [key]", aclObjectGetCmd.Use)
	}
}

func TestAclObjectSetCmd(t *testing.T) {
	if aclObjectSetCmd == nil {
		t.Fatal("aclObjectSetCmd should not be nil")
	}

	if aclObjectSetCmd.Use != "object-set [bucket] [key] [acl]" {
		t.Errorf("aclObjectSetCmd.Use = %v, want object-set [bucket] [key] [acl]", aclObjectSetCmd.Use)
	}
}

func TestAclCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"bucket-get",
		"bucket-set",
		"object-get",
		"object-set",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range aclCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in acl command", subName)
		}
	}
}

func TestAclBucketGetCmd_ArgsValidation(t *testing.T) {
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
			err := aclBucketGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAclBucketSetCmd_ArgsValidation(t *testing.T) {
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
			args:    []string{"bucket", "private"},
			wantErr: false,
		},
		{
			name:    "三个参数",
			args:    []string{"bucket", "private", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := aclBucketSetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAclObjectGetCmd_ArgsValidation(t *testing.T) {
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
			err := aclObjectGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAclObjectSetCmd_ArgsValidation(t *testing.T) {
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
			args:    []string{"bucket", "key", "private"},
			wantErr: false,
		},
		{
			name:    "四个参数",
			args:    []string{"bucket", "key", "private", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := aclObjectSetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
