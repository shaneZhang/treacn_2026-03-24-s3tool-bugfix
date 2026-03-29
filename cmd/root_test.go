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

func TestRootCommand(t *testing.T) {
	// 重置配置
	viper.Reset()
	config.GlobalConfig = config.Config{}

	// 创建临时配置文件
	tmpfile := testutil.CreateTempFile(t, `
region: us-east-1
access_key: test-key
secret_key: test-secret
`)
	defer os.Remove(tmpfile)

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "显示帮助信息",
			args:           []string{"--help"},
			expectedOutput: "S3 命令行管理工具",
			expectedError:  false,
		},
		{
			name:           "显示版本/帮助",
			args:           []string{"-h"},
			expectedOutput: "Usage:",
			expectedError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 重置命令状态
			rootCmd.SetArgs(nil)

			args := append([]string{"--config", tmpfile}, tc.args...)
			stdout, stderr, err := testutil.ExecuteCommand(rootCmd, args...)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tc.expectedOutput != "" {
				assert.Contains(t, stdout+stderr, tc.expectedOutput)
			}
		})
	}
}

func TestRootCommand_ConfigFile(t *testing.T) {
	// 重置配置
	viper.Reset()
	config.GlobalConfig = config.Config{}

	testCases := []struct {
		name          string
		configContent string
		expected      config.Config
		expectedError bool
	}{
		{
			name: "加载有效配置",
			configContent: `
region: cn-north-1
access_key: AKIA123456
secret_key: secret123
endpoint: http://localhost:4566
force_path_style: true
`,
			expected: config.Config{
				Region:         "cn-north-1",
				AccessKey:      "AKIA123456",
				SecretKey:      "secret123",
				Endpoint:       "http://localhost:4566",
				ForcePathStyle: true,
				UseAccelerate:  false,
			},
			expectedError: false,
		},
		{
			name:          "空配置使用默认值",
			configContent: "",
			expected: config.Config{
				Region:         "us-east-1",
				ForcePathStyle: false,
				UseAccelerate:  false,
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 重置配置
			viper.Reset()
			config.GlobalConfig = config.Config{}

			tmpfile := testutil.CreateTempFile(t, tc.configContent)
			defer os.Remove(tmpfile)

			// 使用自定义的配置加载逻辑
			err := config.LoadConfig(tmpfile)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// 验证配置
			assert.Equal(t, tc.expected.Region, config.GlobalConfig.Region)
			assert.Equal(t, tc.expected.AccessKey, config.GlobalConfig.AccessKey)
			assert.Equal(t, tc.expected.SecretKey, config.GlobalConfig.SecretKey)
			assert.Equal(t, tc.expected.Endpoint, config.GlobalConfig.Endpoint)
			assert.Equal(t, tc.expected.ForcePathStyle, config.GlobalConfig.ForcePathStyle)
			assert.Equal(t, tc.expected.UseAccelerate, config.GlobalConfig.UseAccelerate)
		})
	}
}

func TestRootCommand_SubCommands(t *testing.T) {
	// 验证所有子命令都已注册
	expectedCommands := []string{
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
		"init",
	}

	for _, cmdName := range expectedCommands {
		t.Run("子命令存在:"+cmdName, func(t *testing.T) {
			cmd, _, err := rootCmd.Traverse([]string{cmdName})
			require.NoError(t, err)
			assert.Equal(t, cmdName, cmd.Name())
		})
	}
}

func TestRootCommand_GlobalFlags(t *testing.T) {
	// 测试全局标志
	flag := rootCmd.Flag("config")
	assert.NotNil(t, flag)
	assert.Equal(t, "config", flag.Name)
	assert.Equal(t, "c", flag.Shorthand)

	flag = rootCmd.Flag("verbose")
	assert.NotNil(t, flag)
	assert.Equal(t, "verbose", flag.Name)
	assert.Equal(t, "v", flag.Shorthand)
}

// 测试命令行参数解析
func TestCommandArgumentParsing(t *testing.T) {
	testCases := []struct {
		command     *cobra.Command
		name        string
		args        []string
		expectError bool
	}{
		{
			command:     bucketCreateCmd,
			name:        "bucket create 无参数",
			args:        []string{},
			expectError: true, // 需要桶名参数
		},
		{
			command:     bucketCreateCmd,
			name:        "bucket create 有参数",
			args:        []string{"my-bucket"},
			expectError: false,
		},
		{
			command:     bucketDeleteCmd,
			name:        "bucket delete 无参数",
			args:        []string{},
			expectError: true, // 需要桶名参数
		},
		{
			command:     objectListCmd,
			name:        "object list 无参数",
			args:        []string{},
			expectError: true, // 需要桶参数
		},
		{
			command:     objectListCmd,
			name:        "object list 有参数",
			args:        []string{"my-bucket"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 检查参数验证
			err := tc.command.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
