package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	if rootCmd.Use != "s3tool" {
		t.Errorf("rootCmd.Use = %v, want s3tool", rootCmd.Use)
	}
}

func TestExecute(t *testing.T) {
	// Skip this test as it requires a valid config file
	t.Skip("Skipping Execute test - requires valid config file")
}

func TestRootCmd_Flags(t *testing.T) {
	configFlag := rootCmd.PersistentFlags().Lookup("config")
	if configFlag == nil {
		t.Error("config flag should exist")
	} else {
		if configFlag.Shorthand != "c" {
			t.Errorf("config shorthand = %v, want c", configFlag.Shorthand)
		}
	}

	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("verbose flag should exist")
	} else {
		if verboseFlag.Shorthand != "v" {
			t.Errorf("verbose shorthand = %v, want v", verboseFlag.Shorthand)
		}
	}
}

func TestRootCmd_Subcommands(t *testing.T) {
	expectedCommands := []string{
		"init",
		"bucket",
		"object",
		"presign",
		"multipart",
		"policy",
		"lifecycle",
		"versioning",
		"tags",
		"website",
		"cors",
		"acl",
		"logging",
		"encryption",
		"replication",
		"notification",
	}

	for _, cmdName := range expectedCommands {
		cmd, _, err := rootCmd.Find([]string{cmdName})
		if err != nil {
			t.Errorf("Command %s not found: %v", cmdName, err)
			continue
		}
		if cmd.Use != cmdName {
			t.Errorf("Command use = %v, want %v", cmd.Use, cmdName)
		}
	}
}

func TestRootCmd_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("Help output should not be empty")
	}

	expectedStrings := []string{
		"s3tool",
		"Usage:",
		"Available Commands:",
		"Flags:",
	}

	for _, str := range expectedStrings {
		if !bytes.Contains([]byte(output), []byte(str)) {
			t.Errorf("Help output should contain %q", str)
		}
	}
}

func TestRootCmd_Version(t *testing.T) {
	if rootCmd.Version != "" {
		t.Logf("Version: %s", rootCmd.Version)
	}
}

func TestRootCmd_ShortAndLong(t *testing.T) {
	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}
	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}
}
