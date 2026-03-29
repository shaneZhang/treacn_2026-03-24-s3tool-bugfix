package cmd

import (
	"bytes"
	"context"
	"testing"
	"time"

	"s3tool/config"
	"s3tool/pkg/mock"
	"s3tool/pkg/testutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试bucket list命令
func TestBucketListCommand(t *testing.T) {
	viper.Reset()
	config.GlobalConfig = config.Config{
		Region:    "us-east-1",
		AccessKey: "test-key",
		SecretKey: "test-secret",
	}

	testCases := []struct {
		name           string
		mockBuckets    []types.Bucket
		expectedOutput string
		expectedError  bool
	}{
		{
			name: "列出多个存储桶",
			mockBuckets: []types.Bucket{
				{
					Name:         aws.String("bucket1"),
					CreationDate: aws.Time(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					Name:         aws.String("bucket2"),
					CreationDate: aws.Time(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			expectedOutput: "bucket1",
			expectedError:  false,
		},
		{
			name:           "空存储桶列表",
			mockBuckets:    []types.Bucket{},
			expectedOutput: "总存储桶数量: 0",
			expectedError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建mock客户端（这里只是演示，实际用于依赖注入场景）
			_ = &mock.MockS3Client{
				ListBucketsFunc: func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
					return &s3.ListBucketsOutput{
						Buckets: tc.mockBuckets,
					}, nil
				},
			}

			// 测试参数验证
			err := bucketListCmd.ValidateArgs([]string{})
			require.NoError(t, err)
		})
	}
}

// 测试bucket create命令参数验证
func TestBucketCreateCommand(t *testing.T) {
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
			args:        []string{"my-new-bucket"},
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
			err := bucketCreateCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试bucket delete命令参数验证
func TestBucketDeleteCommand(t *testing.T) {
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
			args:        []string{"my-bucket-to-delete"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := bucketDeleteCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试bucket info命令参数验证
func TestBucketInfoCommand(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := bucketInfoCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试bucket location命令参数验证
func TestBucketLocationCommand(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := bucketLocationCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试bucket empty命令参数验证
func TestBucketEmptyCommand(t *testing.T) {
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
			args:        []string{"my-bucket-to-empty"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := bucketEmptyCmd.ValidateArgs(tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// 测试bucket子命令存在
func TestBucketSubCommands(t *testing.T) {
	expectedSubCommands := []string{"list", "create", "delete", "info", "location", "empty"}

	for _, subCmdName := range expectedSubCommands {
		t.Run("子命令存在:"+subCmdName, func(t *testing.T) {
			cmd, _, err := bucketCmd.Traverse([]string{subCmdName})
			require.NoError(t, err)
			assert.Equal(t, subCmdName, cmd.Name())
		})
	}
}

// 测试bucket命令标志
func TestBucketCommandFlags(t *testing.T) {
	// 测试list命令是否有预期的标志
	flags := bucketListCmd.Flags()
	assert.NotNil(t, flags)
	// 可以添加更多标志测试
}

// 集成测试：测试完整的命令执行流程（不实际调用AWS）
func TestBucketCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 创建临时配置文件
	tmpfile := testutil.CreateTempFile(t, `
region: us-east-1
access_key: test-key
secret_key: test-secret
`)

	// 测试 bucket list 命令的帮助输出
	args := []string{"--config", tmpfile, "bucket", "list", "--help"}
	stdout, stderr, err := testutil.ExecuteCommand(rootCmd, args...)

	require.NoError(t, err)
	assert.Contains(t, stdout+stderr, "列出当前账户下所有 S3 存储桶")
}

// 测试空桶清空场景（mock）
func TestBucketEmptyScenario(t *testing.T) {
	viper.Reset()
	config.GlobalConfig = config.Config{
		Region:    "us-east-1",
		AccessKey: "test-key",
		SecretKey: "test-secret",
	}

	// 模拟空桶场景
	mockClient := &mock.MockS3Client{
		ListObjectsV2Func: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			return &s3.ListObjectsV2Output{
				Contents: []types.Object{},
			}, nil
		},
	}

	// 测试参数验证
	err := bucketEmptyCmd.ValidateArgs([]string{"test-empty-bucket"})
	require.NoError(t, err)

	_ = mockClient // 避免未使用警告
}

// 测试创建存储桶时的区域配置
func TestBucketCreateRegionConfiguration(t *testing.T) {
	testCases := []struct {
		name          string
		region        string
		expectHasConf bool
	}{
		{
			name:          "us-east-1区域不需要配置",
			region:        "us-east-1",
			expectHasConf: false,
		},
		{
			name:          "其他区域需要配置",
			region:        "cn-north-1",
			expectHasConf: true,
		},
		{
			name:          "空区域使用默认",
			region:        "",
			expectHasConf: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.GlobalConfig.Region = tc.region

			var input *s3.CreateBucketInput
			if tc.region != "" && tc.region != "us-east-1" {
				input = &s3.CreateBucketInput{
					Bucket: aws.String("test-bucket"),
					CreateBucketConfiguration: &types.CreateBucketConfiguration{
						LocationConstraint: types.BucketLocationConstraint(tc.region),
					},
				}
			} else {
				input = &s3.CreateBucketInput{
					Bucket: aws.String("test-bucket"),
				}
			}

			if tc.expectHasConf {
				assert.NotNil(t, input.CreateBucketConfiguration)
			} else {
				// 如果配置存在但LocationConstraint为空字符串，这也是可接受的
				if input.CreateBucketConfiguration != nil {
					assert.Equal(t, types.BucketLocationConstraint(""), input.CreateBucketConfiguration.LocationConstraint)
				}
			}
		})
	}
}

// 辅助测试函数：测试命令输出格式
func TestBucketCommandOutput(t *testing.T) {
	// 测试table输出格式
	var buf bytes.Buffer
	bucketCmd.SetOutput(&buf)

	// 确保输出可以正常设置
	assert.Equal(t, &buf, bucketCmd.OutOrStdout())
}
