package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNotificationCmd(t *testing.T) {
	if notificationCmd == nil {
		t.Fatal("notificationCmd should not be nil")
	}

	if notificationCmd.Use != "notification" {
		t.Errorf("notificationCmd.Use = %v, want notification", notificationCmd.Use)
	}
}

func TestNotificationGetCmd(t *testing.T) {
	if notificationGetCmd == nil {
		t.Fatal("notificationGetCmd should not be nil")
	}

	if notificationGetCmd.Use != "get [bucket]" {
		t.Errorf("notificationGetCmd.Use = %v, want get [bucket]", notificationGetCmd.Use)
	}
}

func TestNotificationCmd_Subcommands(t *testing.T) {
	expectedSubcommands := []string{
		"get",
	}

	for _, subName := range expectedSubcommands {
		found := false
		for _, cmd := range notificationCmd.Commands() {
			if cmd.Use == subName || (len(subName) < len(cmd.Use) && cmd.Use[:len(subName)] == subName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %s not found in notification command", subName)
		}
	}
}

func TestNotificationGetCmd_ArgsValidation(t *testing.T) {
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
			err := notificationGetCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
