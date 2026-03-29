package cmd

import (
	"os"
	"testing"

	"s3tool/config"
	"s3tool/pkg/testutil"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试multipart init命令参数验证
func TestMultipartInitCommand(t *testing.T) {
	viper.Reset()
	config.GlobalConfig = config.Config{
		Region:    "us-east-1",
		AccessKey: "test-key",
		SecretKey: "test-secret",
	}

	testCases := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "无参数",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "缺少key参数",
			args:        []string{"my-bucket"},
			expectError: true,
		},
		{
			name:        "完整参数",
			args:        []string{"my-bucket", "my-key"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := multipartInitCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试multipart upload命令参数验证
func TestMultipartUploadCommand(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "无参数",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "缺少upload-id参数",
			args:        []string{"my-bucket", "my-key"},
			expectError: true,
		},
		{
			name:        "缺少part-number参数",
			args:        []string{"my-bucket", "my-key", "upload-id-123"},
			expectError: true,
		},
		{
			name:        "缺少file参数",
			args:        []string{"my-bucket", "my-key", "upload-id-123", "1"},
			expectError: true,
		},
		{
			name:        "完整参数",
			args:        []string{"my-bucket", "my-key", "upload-id-123", "1", "/path/to/part/file"},
			expectError: false,
		},
		{
			name:        "无效的分片编号",
			args:        []string{"my-bucket", "my-key", "upload-id-123", "not-a-number", "/path/to/part/file"},
			expectError: false, // 参数验证不检查类型，只检查数量
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := multipartUploadCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试multipart list命令参数验证
func TestMultipartListCommand(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "无参数",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "缺少key参数",
			args:        []string{"my-bucket"},
			expectError: true,
		},
		{
			name:        "缺少upload-id参数",
			args:        []string{"my-bucket", "my-key"},
			expectError: true,
		},
		{
			name:        "完整参数",
			args:        []string{"my-bucket", "my-key", "upload-id-123"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := multipartListCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试multipart complete命令参数验证
func TestMultipartCompleteCommand(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "无参数",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "缺少key参数",
			args:        []string{"my-bucket"},
			expectError: true,
		},
		{
			name:        "缺少upload-id参数",
			args:        []string{"my-bucket", "my-key"},
			expectError: true,
		},
		{
			name:        "完整参数",
			args:        []string{"my-bucket", "my-key", "upload-id-123"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := multipartCompleteCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试multipart abort命令参数验证
func TestMultipartAbortCommand(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "无参数",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "缺少key参数",
			args:        []string{"my-bucket"},
			expectError: true,
		},
		{
			name:        "缺少upload-id参数",
			args:        []string{"my-bucket", "my-key"},
			expectError: true,
		},
		{
			name:        "完整参数",
			args:        []string{"my-bucket", "my-key", "upload-id-123"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := multipartAbortCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试multipart子命令存在
func TestMultipartSubCommands(t *testing.T) {
	expectedSubCommands := []string{"init", "upload", "list", "complete", "abort"}

	for _, subCmdName := range expectedSubCommands {
		t.Run("子命令存在:"+subCmdName, func(t *testing.T) {
			cmd, _, err := multipartCmd.Traverse([]string{subCmdName})
			require.NoError(t, err)
			assert.Equal(t, subCmdName, cmd.Name())
		})
	}
}

// 测试multipart命令标志
func TestMultipartCommandFlags(t *testing.T) {
	t.Run("multipart upload flags", func(t *testing.T) {
		flags := multipartUploadCmd.Flags()
		assert.NotNil(t, flags.Lookup("part-size"))
	})
}

// 测试part-size默认值
func TestMultipartPartSizeFlag(t *testing.T) {
	flags := multipartUploadCmd.Flags()
	partSizeFlag := flags.Lookup("part-size")
	assert.NotNil(t, partSizeFlag)
	assert.Equal(t, "5242880", partSizeFlag.DefValue) // 5MB
}

// 集成测试：测试帮助信息
func TestMultipartCommandHelp(t *testing.T) {
	tmpfile := testutil.CreateTempFile(t, `
region: us-east-1
access_key: test-key
secret_key: test-secret
`)
	defer os.Remove(tmpfile)

	testCases := []struct {
		args           []string
		expectedOutput string
	}{
		{
			args:           []string{"--config", tmpfile, "multipart", "init", "--help"},
			expectedOutput: "初始化一个大文件的多部分上传",
		},
		{
			args:           []string{"--config", tmpfile, "multipart", "upload", "--help"},
			expectedOutput: "上传多部分上传的一个分片",
		},
		{
			args:           []string{"--config", tmpfile, "multipart", "list", "--help"},
			expectedOutput: "列出多部分上传的所有已上传分片",
		},
		{
			args:           []string{"--config", tmpfile, "multipart", "complete", "--help"},
			expectedOutput: "完成多部分上传",
		},
		{
			args:           []string{"--config", tmpfile, "multipart", "abort", "--help"},
			expectedOutput: "中止一个进行中的多部分上传",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedOutput, func(t *testing.T) {
			stdout, stderr, err := testutil.ExecuteCommand(rootCmd, tc.args...)
			require.NoError(t, err)
			assert.Contains(t, stdout+stderr, tc.expectedOutput)
		})
	}
}

// 测试part-size参数配置
func TestMultipartPartSizeConfiguration(t *testing.T) {
	flags := multipartUploadCmd.Flags()

	// 测试默认值
	partSize, err := flags.GetInt("part-size")
	require.NoError(t, err)
	assert.Equal(t, 5*1024*1024, partSize)

	// 设置自定义值
	err = flags.Set("part-size", "10485760") // 10MB
	require.NoError(t, err)
	partSize, err = flags.GetInt("part-size")
	require.NoError(t, err)
	assert.Equal(t, 10*1024*1024, partSize)

	// 无效值测试
	err = flags.Set("part-size", "invalid")
	assert.Error(t, err)
}
