package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestWebsiteCmd(t *testing.T) {
	if websiteCmd == nil {
		t.Fatal("websiteCmd should not be nil")
	}

	if websiteCmd.Use != "website" {
		t.Errorf("websiteCmd.Use = %v, want website", websiteCmd.Use)
	}
}

func TestWebsiteGetCmd(t *testing.T) {
	if websiteGetCmd == nil {
		t.Fatal("websiteGetCmd should not be nil")
	}

	if websiteGetCmd.Use != "get [bucket]" {
		t.Errorf("websiteGetCmd.Use = %v, want get [bucket]", websiteGetCmd.Use)
	}
}

func TestWebsiteEnableCmd(t *testing.T) {
	if websiteEnableCmd == nil {
		t.Fatal("websiteEnableCmd should not be nil")
	}

	if websiteEnableCmd.Use != "enable [bucket] [index-document] [error-document]" {
		t.Errorf("websiteEnableCmd.Use = %v, want enable [bucket] [index-document] [error-document]", websiteEnableCmd.Use)
	}
}

func TestWebsiteDisableCmd(t *testing.T) {
	if websiteDisableCmd == nil {
		t.Fatal("websiteDisableCmd should not be nil")
	}

	if websiteDisableCmd.Use != "disable [bucket]" {
		t.Errorf("websiteDisableCmd.Use = %v, want disable [bucket]", websiteDisableCmd.Use)
	}
}

func TestWebsiteCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"get",
		"enable",
		"disable",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range websiteCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in website command", subName)
		}
	}
}

func TestWebsiteGetCmd_ArgsValidation(t *testing.T) {
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
			err := websiteGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebsiteEnableCmd_ArgsValidation(t *testing.T) {
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
			args:    []string{"bucket", "index.html"},
			wantErr: false,
		},
		{
			name:    "三个参数",
			args:    []string{"bucket", "index.html", "error.html"},
			wantErr: false,
		},
		{
			name:    "四个参数",
			args:    []string{"bucket", "index.html", "error.html", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := websiteEnableCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebsiteDisableCmd_ArgsValidation(t *testing.T) {
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
			err := websiteDisableCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
