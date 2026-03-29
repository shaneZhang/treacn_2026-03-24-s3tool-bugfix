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

// 测试object list命令参数验证
func TestObjectListCommand(t *testing.T) {
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
			name:        "有参数",
			args:        []string{"my-bucket"},
			expectError: false,
		},
		{
			name:        "参数过多",
			args:        []string{"bucket1", "bucket2"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := objectListCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试object put命令参数验证
func TestObjectPutCommand(t *testing.T) {
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
			name:        "缺少文件参数",
			args:        []string{"my-bucket", "my-key"},
			expectError: true,
		},
		{
			name:        "完整参数",
			args:        []string{"my-bucket", "my-key", "/path/to/file"},
			expectError: false,
		},
		{
			name:        "参数过多",
			args:        []string{"bucket", "key", "file", "extra"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := objectPutCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试object get命令参数验证
func TestObjectGetCommand(t *testing.T) {
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
			name:        "缺少本地文件参数",
			args:        []string{"my-bucket", "my-key"},
			expectError: true,
		},
		{
			name:        "完整参数",
			args:        []string{"my-bucket", "my-key", "/path/to/local/file"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := objectGetCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试object delete命令参数验证
func TestObjectDeleteCommand(t *testing.T) {
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
			err := objectDeleteCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试object copy命令参数验证
func TestObjectCopyCommand(t *testing.T) {
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
			name:        "缺少目标参数",
			args:        []string{"src-bucket", "src-key", "dest-bucket"},
			expectError: true,
		},
		{
			name:        "完整参数",
			args:        []string{"src-bucket", "src-key", "dest-bucket", "dest-key"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := objectCopyCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试object info命令参数验证
func TestObjectInfoCommand(t *testing.T) {
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
			err := objectInfoCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试object url命令参数验证
func TestObjectUrlCommand(t *testing.T) {
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
			err := objectUrlCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试object mv命令参数验证
func TestObjectMvCommand(t *testing.T) {
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
			name:        "缺少目标key参数",
			args:        []string{"my-bucket", "source-key"},
			expectError: true,
		},
		{
			name:        "完整参数",
			args:        []string{"my-bucket", "source-key", "dest-key"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := objectMvCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试object子命令存在
func TestObjectSubCommands(t *testing.T) {
	expectedSubCommands := []string{"list", "put", "get", "delete", "copy", "info", "url", "mv"}

	for _, subCmdName := range expectedSubCommands {
		t.Run("子命令存在:"+subCmdName, func(t *testing.T) {
			cmd, _, err := objectCmd.Traverse([]string{subCmdName})
			require.NoError(t, err)
			assert.Equal(t, subCmdName, cmd.Name())
		})
	}
}

// 测试object命令标志
func TestObjectCommandFlags(t *testing.T) {
	t.Run("object list flags", func(t *testing.T) {
		flags := objectListCmd.Flags()
		assert.NotNil(t, flags.Lookup("prefix"))
		assert.NotNil(t, flags.Lookup("recursive"))
		assert.NotNil(t, flags.Lookup("max-keys"))
	})

	t.Run("object put flags", func(t *testing.T) {
		flags := objectPutCmd.Flags()
		assert.NotNil(t, flags.Lookup("content-type"))
		assert.NotNil(t, flags.Lookup("storage-class"))
	})
}

// 测试formatBytes辅助函数
func TestFormatBytes(t *testing.T) {
	testCases := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1024 * 1024, "1.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
		{1536, "1.5 KB"},
		{1024 * 1024 * 5, "5.0 MB"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := formatBytes(tc.bytes)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// 集成测试：测试帮助信息
func TestObjectCommandHelp(t *testing.T) {
	// 创建临时配置文件
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
			args:           []string{"--config", tmpfile, "object", "list", "--help"},
			expectedOutput: "列出存储桶中的对象",
		},
		{
			args:           []string{"--config", tmpfile, "object", "put", "--help"},
			expectedOutput: "上传本地文件到 S3 存储桶",
		},
		{
			args:           []string{"--config", tmpfile, "object", "--help"},
			expectedOutput: "对象相关操作",
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

// 测试边界条件：大文件上传（参数验证）
func TestObjectPutLargeFile(t *testing.T) {
	// 仅测试参数验证
	err := objectPutCmd.ValidateArgs([]string{"my-bucket", "large-file-key", "/path/to/large/file"})
	require.NoError(t, err)

	// 验证flag是否存在
	flags := objectPutCmd.Flags()
	assert.NotNil(t, flags.Lookup("storage-class"))
}
