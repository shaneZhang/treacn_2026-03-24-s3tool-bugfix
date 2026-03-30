package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestTagsCmd(t *testing.T) {
	if tagsCmd == nil {
		t.Fatal("tagsCmd should not be nil")
	}

	if tagsCmd.Use != "tags" {
		t.Errorf("tagsCmd.Use = %v, want tags", tagsCmd.Use)
	}
}

func TestTagsBucketGetCmd(t *testing.T) {
	if tagsBucketGetCmd == nil {
		t.Fatal("tagsBucketGetCmd should not be nil")
	}

	if tagsBucketGetCmd.Use != "bucket-get [bucket]" {
		t.Errorf("tagsBucketGetCmd.Use = %v, want bucket-get [bucket]", tagsBucketGetCmd.Use)
	}
}

func TestTagsBucketPutCmd(t *testing.T) {
	if tagsBucketPutCmd == nil {
		t.Fatal("tagsBucketPutCmd should not be nil")
	}

	if tagsBucketPutCmd.Use != "bucket-put [bucket] [key1=value1] [key2=value2]..." {
		t.Errorf("tagsBucketPutCmd.Use = %v, want bucket-put [bucket] [key1=value1] [key2=value2]...", tagsBucketPutCmd.Use)
	}
}

func TestTagsObjectGetCmd(t *testing.T) {
	if tagsObjectGetCmd == nil {
		t.Fatal("tagsObjectGetCmd should not be nil")
	}

	if tagsObjectGetCmd.Use != "object-get [bucket] [key]" {
		t.Errorf("tagsObjectGetCmd.Use = %v, want object-get [bucket] [key]", tagsObjectGetCmd.Use)
	}
}

func TestTagsCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"bucket-get",
		"bucket-put",
		"bucket-delete",
		"object-get",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range tagsCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in tags command", subName)
		}
	}
}

func TestTagsBucketGetCmd_ArgsValidation(t *testing.T) {
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
			err := tagsBucketGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTagsBucketPutCmd_ArgsValidation(t *testing.T) {
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
			args:    []string{"bucket", "env=prod"},
			wantErr: false,
		},
		{
			name:    "三个参数",
			args:    []string{"bucket", "env=prod", "team=dev"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tagsBucketPutCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTagsObjectGetCmd_ArgsValidation(t *testing.T) {
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
			err := tagsObjectGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
