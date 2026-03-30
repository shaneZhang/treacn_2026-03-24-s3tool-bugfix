package testutil

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// ExecuteCommand 执行cobra命令并返回输出
func ExecuteCommand(root *cobra.Command, args ...string) (string, string, error) {
	_, output, err := ExecuteCommandC(root, args...)
	return output.Stdout, output.Stderr, err
}

// ExecuteCommandC 执行cobra命令并返回详细输出
func ExecuteCommandC(root *cobra.Command, args ...string) (*cobra.Command, *Output, error) {
	bufOut := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)
	root.SetOut(bufOut)
	root.SetErr(bufErr)
	root.SetArgs(args)

	cmd, err := root.ExecuteC()

	return cmd, &Output{Stdout: bufOut.String(), Stderr: bufErr.String()}, err
}

// Output 命令输出
type Output struct {
	Stdout string
	Stderr string
}

// CreateTempFile 创建临时文件
func CreateTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "test-*")
	require.NoError(t, err)

	_, err = tmpfile.Write([]byte(content))
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	return tmpfile.Name()
}

// CreateTestBucketName 创建测试用的桶名（带随机后缀）
func CreateTestBucketName(prefix string) string {
	return prefix + "-test-" + randomString(6)
}

// randomString 生成随机字符串
func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}

// SetupTestConfig 设置测试配置
func SetupTestConfig(t *testing.T) {
	t.Helper()
	// 创建一个测试用的配置文件
	content := `
region: us-east-1
access_key: test-key
secret_key: test-secret
`
	tmpfile := CreateTempFile(t, content)
	t.Setenv("S3TOOL_CONFIG", tmpfile)
}

// AssertContains 断言输出包含指定字符串
func AssertContains(t *testing.T, s, substr string, msg ...interface{}) {
	t.Helper()
	require.Contains(t, s, substr, msg...)
}

// AssertNotContains 断言输出不包含指定字符串
func AssertNotContains(t *testing.T, s, substr string, msg ...interface{}) {
	t.Helper()
	require.NotContains(t, s, substr, msg...)
}
