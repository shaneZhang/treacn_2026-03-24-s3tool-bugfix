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

// 测试acl bucket-get命令参数验证
func TestAclBucketGetCommand(t *testing.T) {
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
			name:        "完整参数",
			args:        []string{"my-bucket"},
			expectError: false,
		},
		{
			name:        "参数过多",
			args:        []string{"my-bucket", "extra"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := aclBucketGetCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试acl bucket-set命令参数验证
func TestAclBucketSetCommand(t *testing.T) {
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
			name:        "缺少acl参数",
			args:        []string{"my-bucket"},
			expectError: true,
		},
		{
			name:        "完整参数",
			args:        []string{"my-bucket", "public-read"},
			expectError: false,
		},
		{
			name:        "参数过多",
			args:        []string{"my-bucket", "public-read", "extra"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := aclBucketSetCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试acl object-get命令参数验证
func TestAclObjectGetCommand(t *testing.T) {
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
		{
			name:        "参数过多",
			args:        []string{"my-bucket", "my-key", "extra"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := aclObjectGetCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试acl object-set命令参数验证
func TestAclObjectSetCommand(t *testing.T) {
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
			name:        "缺少acl参数",
			args:        []string{"my-bucket", "my-key"},
			expectError: true,
		},
		{
			name:        "完整参数",
			args:        []string{"my-bucket", "my-key", "public-read"},
			expectError: false,
		},
		{
			name:        "参数过多",
			args:        []string{"my-bucket", "my-key", "public-read", "extra"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := aclObjectSetCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试acl子命令存在
func TestAclSubCommands(t *testing.T) {
	expectedSubCommands := []string{"bucket-get", "bucket-set", "object-get", "object-set"}

	for _, subCmdName := range expectedSubCommands {
		t.Run("子命令存在:"+subCmdName, func(t *testing.T) {
			cmd, _, err := aclCmd.Traverse([]string{subCmdName})
			require.NoError(t, err)
			assert.Equal(t, subCmdName, cmd.Name())
		})
	}
}

// 集成测试：测试帮助信息
func TestAclCommandHelp(t *testing.T) {
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
			args:           []string{"--config", tmpfile, "acl", "bucket-get", "--help"},
			expectedOutput: "获取指定存储桶的ACL",
		},
		{
			args:           []string{"--config", tmpfile, "acl", "bucket-set", "--help"},
			expectedOutput: "设置存储桶的ACL",
		},
		{
			args:           []string{"--config", tmpfile, "acl", "object-get", "--help"},
			expectedOutput: "获取指定对象的ACL",
		},
		{
			args:           []string{"--config", tmpfile, "acl", "object-set", "--help"},
			expectedOutput: "设置对象的ACL",
		},
		{
			args:           []string{"--config", tmpfile, "acl", "--help"},
			expectedOutput: "管理存储桶和对象的访问控制列表",
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

// 测试ACL可选值文档
func TestAclCommandDocumentation(t *testing.T) {
	// 验证bucket-set命令的Long描述包含可选值信息
	assert.Contains(t, aclBucketSetCmd.Long, "private | public-read | public-read-write | authenticated-read")

	// 验证object-set命令的Long描述包含可选值信息
	assert.Contains(t, aclObjectSetCmd.Long, "private | public-read | public-read-write | authenticated-read | bucket-owner-read | bucket-owner-full-control")
}

// 测试ACL命令短描述
func TestAclCommandShortDescriptions(t *testing.T) {
	assert.Contains(t, aclCmd.Short, "访问控制列表管理")
	assert.Contains(t, aclBucketGetCmd.Short, "获取存储桶ACL")
	assert.Contains(t, aclBucketSetCmd.Short, "设置存储桶ACL")
	assert.Contains(t, aclObjectGetCmd.Short, "获取对象ACL")
	assert.Contains(t, aclObjectSetCmd.Short, "设置对象ACL")
}
