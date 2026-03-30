package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestPolicyCmd(t *testing.T) {
	if policyCmd == nil {
		t.Fatal("policyCmd should not be nil")
	}

	if policyCmd.Use != "policy" {
		t.Errorf("policyCmd.Use = %v, want policy", policyCmd.Use)
	}
}

func TestPolicyGetCmd(t *testing.T) {
	if policyGetCmd == nil {
		t.Fatal("policyGetCmd should not be nil")
	}

	if policyGetCmd.Use != "get [bucket]" {
		t.Errorf("policyGetCmd.Use = %v, want get [bucket]", policyGetCmd.Use)
	}
}

func TestPolicySetCmd(t *testing.T) {
	if policySetCmd == nil {
		t.Fatal("policySetCmd should not be nil")
	}

	if policySetCmd.Use != "set [bucket] [policy-file]" {
		t.Errorf("policySetCmd.Use = %v, want set [bucket] [policy-file]", policySetCmd.Use)
	}
}

func TestPolicyDeleteCmd(t *testing.T) {
	if policyDeleteCmd == nil {
		t.Fatal("policyDeleteCmd should not be nil")
	}

	if policyDeleteCmd.Use != "delete [bucket]" {
		t.Errorf("policyDeleteCmd.Use = %v, want delete [bucket]", policyDeleteCmd.Use)
	}
}

func TestPolicyCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"get",
		"set",
		"delete",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range policyCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in policy command", subName)
		}
	}
}

func TestPolicyGetCmd_ArgsValidation(t *testing.T) {
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
			err := policyGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPolicySetCmd_ArgsValidation(t *testing.T) {
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
			args:    []string{"bucket", "policy.json"},
			wantErr: false,
		},
		{
			name:    "三个参数",
			args:    []string{"bucket", "policy.json", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := policySetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPolicyDeleteCmd_ArgsValidation(t *testing.T) {
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
			err := policyDeleteCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
