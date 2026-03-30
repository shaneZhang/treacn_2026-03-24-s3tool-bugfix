package cmd

import (
	"os"
	"testing"

	"s3tool/config"
	"s3tool/pkg/testutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试presign get命令参数验证
func TestPresignGetCommand(t *testing.T) {
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
			err := presignGetCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试presign put命令参数验证
func TestPresignPutCommand(t *testing.T) {
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
			err := presignPutCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试presign delete命令参数验证
func TestPresignDeleteCommand(t *testing.T) {
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
			err := presignDeleteCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试presign子命令存在
func TestPresignSubCommands(t *testing.T) {
	expectedSubCommands := []string{"get", "put", "delete"}

	for _, subCmdName := range expectedSubCommands {
		t.Run("子命令存在:"+subCmdName, func(t *testing.T) {
			cmd, _, err := presignCmd.Traverse([]string{subCmdName})
			require.NoError(t, err)
			assert.Equal(t, subCmdName, cmd.Name())
		})
	}
}

// 测试presign命令标志
func TestPresignCommandFlags(t *testing.T) {
	testCases := []struct {
		cmd         *cobra.Command
		expectFlags []string
	}{
		{cmd: presignGetCmd, expectFlags: []string{"expires"}},
		{cmd: presignPutCmd, expectFlags: []string{"expires"}},
		{cmd: presignDeleteCmd, expectFlags: []string{"expires"}},
	}

	for _, tc := range testCases {
		for _, flagName := range tc.expectFlags {
			t.Run(flagName, func(t *testing.T) {
				flags := tc.cmd.Flags()
				assert.NotNil(t, flags.Lookup(flagName))
			})
		}
	}
}

// 测试expires参数默认值
func TestPresignExpiresFlag(t *testing.T) {
	// 测试get命令的expires默认值
	flags := presignGetCmd.Flags()
	expiresFlag := flags.Lookup("expires")
	assert.NotNil(t, expiresFlag)
	assert.Equal(t, "3600", expiresFlag.DefValue)

	// 测试put命令
	flags = presignPutCmd.Flags()
	expiresFlag = flags.Lookup("expires")
	assert.NotNil(t, expiresFlag)
	assert.Equal(t, "3600", expiresFlag.DefValue)

	// 测试delete命令
	flags = presignDeleteCmd.Flags()
	expiresFlag = flags.Lookup("expires")
	assert.NotNil(t, expiresFlag)
	assert.Equal(t, "3600", expiresFlag.DefValue)
}

// 集成测试：测试帮助信息
func TestPresignCommandHelp(t *testing.T) {
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
			args:           []string{"--config", tmpfile, "presign", "get", "--help"},
			expectedOutput: "生成用于下载对象的预签名URL",
		},
		{
			args:           []string{"--config", tmpfile, "presign", "put", "--help"},
			expectedOutput: "生成用于上传对象的预签名URL",
		},
		{
			args:           []string{"--config", tmpfile, "presign", "delete", "--help"},
			expectedOutput: "生成用于删除对象的预签名URL",
		},
		{
			args:           []string{"--config", tmpfile, "presign", "--help"},
			expectedOutput: "预签名URL",
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

// 测试expires参数边界值
func TestPresignExpiresBoundary(t *testing.T) {
	// 测试命令是否接受expires参数
	flags := presignGetCmd.Flags()

	// 设置一个非常小的值
	err := flags.Set("expires", "60") // 1分钟
	require.NoError(t, err)
	expires, err := flags.GetInt("expires")
	require.NoError(t, err)
	assert.Equal(t, 60, expires)

	// 设置一个大的值
	err = flags.Set("expires", "604800") // 7天
	require.NoError(t, err)
	expires, err = flags.GetInt("expires")
	require.NoError(t, err)
	assert.Equal(t, 604800, expires)

	// 无效值应该报错
	err = flags.Set("expires", "invalid")
	assert.Error(t, err)
}

// 测试presign命令描述
func TestPresignCommandDescriptions(t *testing.T) {
	assert.Contains(t, presignCmd.Short, "预签名URL")
	assert.Contains(t, presignGetCmd.Short, "GET")
	assert.Contains(t, presignPutCmd.Short, "PUT")
	assert.Contains(t, presignDeleteCmd.Short, "DELETE")
}
