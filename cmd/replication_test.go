package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestReplicationCmd(t *testing.T) {
	if replicationCmd == nil {
		t.Fatal("replicationCmd should not be nil")
	}

	if replicationCmd.Use != "replication" {
		t.Errorf("replicationCmd.Use = %v, want replication", replicationCmd.Use)
	}
}

func TestReplicationGetCmd(t *testing.T) {
	if replicationGetCmd == nil {
		t.Fatal("replicationGetCmd should not be nil")
	}

	if replicationGetCmd.Use != "get [bucket]" {
		t.Errorf("replicationGetCmd.Use = %v, want get [bucket]", replicationGetCmd.Use)
	}
}

func TestReplicationDeleteCmd(t *testing.T) {
	if replicationDeleteCmd == nil {
		t.Fatal("replicationDeleteCmd should not be nil")
	}

	if replicationDeleteCmd.Use != "delete [bucket]" {
		t.Errorf("replicationDeleteCmd.Use = %v, want delete [bucket]", replicationDeleteCmd.Use)
	}
}

func TestReplicationCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"get",
		"delete",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range replicationCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in replication command", subName)
		}
	}
}

func TestReplicationGetCmd_ArgsValidation(t *testing.T) {
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
			err := replicationGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReplicationDeleteCmd_ArgsValidation(t *testing.T) {
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
			err := replicationDeleteCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
