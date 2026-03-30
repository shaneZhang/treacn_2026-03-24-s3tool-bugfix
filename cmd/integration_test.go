package cmd

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 集成测试前缀，避免冲突
const (
	testBucketPrefix = "s3tool-test-"
	testObjectPrefix = "integration-test/"
)

// TestIntegration_Setup 测试初始化
func TestIntegration_Setup(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试，使用 -short 标志跳过")
	}

	// 加载配置
	err := config.LoadConfig("")
	require.NoError(t, err, "加载配置失败")

	// 创建S3客户端
	client, err := config.GetS3Client(context.Background())
	require.NoError(t, err, "创建S3客户端失败")
	assert.NotNil(t, client)
}

// TestIntegration_BucketList 测试列出存储桶
func TestIntegration_BucketList(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := getTestClient(t)

	// 列出存储桶
	output, err := client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	require.NoError(t, err, "列出存储桶失败")
	assert.NotNil(t, output.Buckets)
	t.Logf("找到 %d 个存储桶", len(output.Buckets))
}

// TestIntegration_BucketCRUD 存储桶CRUD测试
func TestIntegration_BucketCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := getTestClient(t)
	bucketName := generateTestBucketName()

	// 1. 创建存储桶
	t.Run("创建存储桶", func(t *testing.T) {
		_, err := client.CreateBucket(context.Background(), &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		require.NoError(t, err, "创建存储桶失败")
		t.Logf("存储桶 %s 创建成功", bucketName)

		// 等待桶创建完成（最终一致性）
		time.Sleep(2 * time.Second)
	})

	// 2. 验证存储桶存在
	t.Run("验证存储桶存在", func(t *testing.T) {
		// 尝试HEAD桶
		_, err := client.HeadBucket(context.Background(), &s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})
		assert.NoError(t, err, "存储桶不存在")
	})

	// 3. 删除存储桶
	t.Run("删除存储桶", func(t *testing.T) {
		_, err := client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})
		assert.NoError(t, err, "删除存储桶失败")
		t.Logf("存储桶 %s 删除成功", bucketName)
	})
}

// TestIntegration_ObjectCRUD 对象操作测试
func TestIntegration_ObjectCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := getTestClient(t)
	testBucket := getExistingBucket(t, client)
	testKey := testObjectPrefix + "test-object.txt"
	testContent := "Hello, S3Tool Integration Test!"

	// 1. 上传对象
	t.Run("上传对象", func(t *testing.T) {
		_, err := client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(testKey),
			Body:   strings.NewReader(testContent),
		})
		require.NoError(t, err, "上传对象失败")
		t.Logf("对象 %s 上传成功", testKey)
	})

	// 2. 下载对象
	t.Run("下载对象", func(t *testing.T) {
		output, err := client.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(testKey),
		})
		require.NoError(t, err, "下载对象失败")
		defer output.Body.Close()

		content, err := io.ReadAll(output.Body)
		require.NoError(t, err)
		assert.Equal(t, testContent, string(content))
		t.Logf("下载内容验证通过: %s", string(content))
	})

	// 3. 列出对象
	t.Run("列出对象", func(t *testing.T) {
		output, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
			Bucket: aws.String(testBucket),
			Prefix: aws.String(testObjectPrefix),
		})
		require.NoError(t, err, "列出对象失败")
		assert.True(t, len(output.Contents) > 0, "应该至少有一个对象")
	})

	// 4. 复制对象
	t.Run("复制对象", func(t *testing.T) {
		destKey := testObjectPrefix + "copied-object.txt"
		copySource := fmt.Sprintf("%s/%s", testBucket, testKey)

		_, err := client.CopyObject(context.Background(), &s3.CopyObjectInput{
			Bucket:     aws.String(testBucket),
			Key:        aws.String(destKey),
			CopySource: aws.String(copySource),
		})
		require.NoError(t, err, "复制对象失败")
		t.Logf("对象复制成功: %s -> %s", testKey, destKey)

		// 清理复制的对象
		_, _ = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(destKey),
		})
	})

	// 5. 删除对象
	t.Run("删除对象", func(t *testing.T) {
		_, err := client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(testKey),
		})
		assert.NoError(t, err, "删除对象失败")
		t.Logf("对象 %s 删除成功", testKey)
	})
}

// TestIntegration_PresignURL 预签名URL测试
func TestIntegration_PresignURL(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := getTestClient(t)
	testBucket := getExistingBucket(t, client)
	testKey := testObjectPrefix + "presign-test-object.txt"
	testContent := "Presign URL Test Content"

	// 上传测试对象
	_, err := client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(testBucket),
		Key:    aws.String(testKey),
		Body:   strings.NewReader(testContent),
	})
	require.NoError(t, err)
	defer func() {
		_, _ = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(testKey),
		})
	}()

	// 1. 生成GET预签名URL
	t.Run("生成GET预签名URL", func(t *testing.T) {
		// 创建预签名客户端
		presignClient := s3.NewPresignClient(client)
		presignResult, err := presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(testKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Hour // 1小时过期
		})

		require.NoError(t, err, "生成GET预签名URL失败")
		assert.NotEmpty(t, presignResult.URL)
		t.Logf("GET预签名URL: %s", presignResult.URL)

		// 验证URL可访问
		resp, err := http.Get(presignResult.URL)
		require.NoError(t, err, "访问预签名URL失败")
		defer resp.Body.Close()

		content, _ := io.ReadAll(resp.Body)
		assert.Equal(t, testContent, string(content))
		t.Log("预签名URL访问验证通过")
	})

	// 2. 生成PUT预签名URL
	t.Run("生成PUT预签名URL", func(t *testing.T) {
		presignClient := s3.NewPresignClient(client)
		putKey := testObjectPrefix + "presign-put-test.txt"

		presignResult, err := presignClient.PresignPutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(putKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Hour
		})

		require.NoError(t, err, "生成PUT预签名URL失败")
		assert.NotEmpty(t, presignResult.URL)
		t.Logf("PUT预签名URL: %s", presignResult.URL)

		// 使用PUT URL上传
		putContent := "Content uploaded via presigned URL"
		req, _ := http.NewRequest("PUT", presignResult.URL, strings.NewReader(putContent))
		req.Header.Set("Content-Type", "text/plain")

		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		require.NoError(t, err, "PUT上传失败")
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// 验证上传内容
		getOutput, _ := client.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(putKey),
		})
		content, _ := io.ReadAll(getOutput.Body)
		getOutput.Body.Close()
		assert.Equal(t, putContent, string(content))

		// 清理
		_, _ = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(putKey),
		})
	})
}

// TestIntegration_FileOperations 文件操作测试
func TestIntegration_FileOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := getTestClient(t)
	testBucket := getExistingBucket(t, client)

	// 创建临时测试文件
	t.Run("文件上传下载", func(t *testing.T) {
		// 创建1MB测试文件
		testFile, err := os.CreateTemp("", "s3test-*.bin")
		require.NoError(t, err)
		defer os.Remove(testFile.Name())
		defer testFile.Close()

		// 写入随机数据
		data := make([]byte, 1024*1024) // 1MB
		for i := range data {
			data[i] = byte(i % 256)
		}
		_, err = testFile.Write(data)
		require.NoError(t, err)
		testFile.Close()

		// 计算MD5
		hasher := md5.New()
		hasher.Write(data)
		expectedMD5 := hex.EncodeToString(hasher.Sum(nil))

		// 上传文件
		fileKey := testObjectPrefix + "large-test-file.bin"
		file, err := os.Open(testFile.Name())
		require.NoError(t, err)
		defer file.Close()

		_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(fileKey),
			Body:   file,
		})
		require.NoError(t, err, "文件上传失败")
		t.Logf("1MB文件上传成功: %s", fileKey)

		// 下载验证
		output, err := client.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(fileKey),
		})
		require.NoError(t, err)
		defer output.Body.Close()

		downloadData, err := io.ReadAll(output.Body)
		require.NoError(t, err)

		hasher.Reset()
		hasher.Write(downloadData)
		actualMD5 := hex.EncodeToString(hasher.Sum(nil))
		assert.Equal(t, expectedMD5, actualMD5, "文件MD5校验失败")
		t.Log("文件下载MD5校验通过")

		// 清理
		_, _ = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(fileKey),
		})
	})
}

// TestIntegration_AclOperations ACL操作测试
func TestIntegration_AclOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := getTestClient(t)
	testBucket := getExistingBucket(t, client)
	testKey := testObjectPrefix + "acl-test-object.txt"

	// 上传测试对象
	_, err := client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(testBucket),
		Key:    aws.String(testKey),
		Body:   strings.NewReader("ACL Test"),
	})
	require.NoError(t, err)
	defer func() {
		_, _ = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(testKey),
		})
	}()

	t.Run("获取存储桶ACL", func(t *testing.T) {
		output, err := client.GetBucketAcl(context.Background(), &s3.GetBucketAclInput{
			Bucket: aws.String(testBucket),
		})
		require.NoError(t, err, "获取存储桶ACL失败")
		assert.NotNil(t, output.Owner)
		t.Logf("存储桶ACL Owner: %s", aws.ToString(output.Owner.DisplayName))
	})

	t.Run("获取对象ACL", func(t *testing.T) {
		output, err := client.GetObjectAcl(context.Background(), &s3.GetObjectAclInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(testKey),
		})
		require.NoError(t, err, "获取对象ACL失败")
		assert.NotNil(t, output.Owner)
		t.Logf("对象ACL获取成功，Grants数量: %d", len(output.Grants))
	})

	t.Run("设置对象ACL", func(t *testing.T) {
		_, err := client.PutObjectAcl(context.Background(), &s3.PutObjectAclInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(testKey),
			ACL:    types.ObjectCannedACLPublicRead,
		})
		require.NoError(t, err, "设置对象ACL失败")
		t.Log("对象ACL设置为public-read成功")

		// 验证
		output, _ := client.GetObjectAcl(context.Background(), &s3.GetObjectAclInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(testKey),
		})
		assert.True(t, len(output.Grants) >= 1, "应该有公开读权限")
	})
}

// TestIntegration_ErrorScenarios 错误场景测试
func TestIntegration_ErrorScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := getTestClient(t)
	testBucket := getExistingBucket(t, client)

	t.Run("访问不存在的对象", func(t *testing.T) {
		_, err := client.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(testBucket),
			Key:    aws.String(testObjectPrefix + "nonexistent-key"),
		})
		assert.Error(t, err)
		var nske *types.NoSuchKey
		assert.ErrorAs(t, err, &nske, "应该返回NoSuchKey错误")
		t.Log("不存在的对象错误验证通过")
	})

	t.Run("访问不存在的存储桶", func(t *testing.T) {
		_, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
			Bucket: aws.String("non-existent-bucket-" + time.Now().Format("20060102")),
		})
		assert.Error(t, err)
		t.Log("不存在的存储桶错误验证通过")
	})
}

// 辅助函数

func getTestClient(t *testing.T) *s3.Client {
	// 使用用户主目录下的配置文件
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err, "获取用户主目录失败")
	configPath := homeDir + "/.s3tool.yaml"

	t.Logf("使用配置文件: %s", configPath)
	err = config.LoadConfig(configPath)
	require.NoError(t, err, "加载配置文件失败")

	client, err := config.GetS3Client(context.Background())
	require.NoError(t, err, "创建S3客户端失败")
	return client
}

func getExistingBucket(t *testing.T, client *s3.Client) string {
	output, err := client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	require.NoError(t, err)
	require.True(t, len(output.Buckets) > 0, "需要至少一个现有存储桶进行测试")

	// 返回第一个存在的桶
	bucketName := aws.ToString(output.Buckets[0].Name)
	t.Logf("使用测试存储桶: %s", bucketName)
	return bucketName
}

func generateTestBucketName() string {
	return fmt.Sprintf("%s%s", testBucketPrefix, time.Now().Format("20060102-150405"))
}

// TestCleanup 清理测试对象
func TestCleanup_TestObjects(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过清理")
	}

	client := getTestClient(t)
	testBucket := getExistingBucket(t, client)

	t.Logf("开始清理测试前缀对象: %s", testObjectPrefix)

	// 列出所有测试前缀的对象并删除
	output, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(testBucket),
		Prefix: aws.String(testObjectPrefix),
	})
	if err != nil {
		t.Logf("列出测试对象失败: %v", err)
		return
	}

	if len(output.Contents) == 0 {
		t.Log("没有需要清理的测试对象")
		return
	}

	// 批量删除
	var objects []types.ObjectIdentifier
	for _, obj := range output.Contents {
		objects = append(objects, types.ObjectIdentifier{
			Key: obj.Key,
		})
	}

	_, err = client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
		Bucket: aws.String(testBucket),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	})
	if err != nil {
		t.Logf("清理测试对象失败: %v", err)
	} else {
		t.Logf("成功清理 %d 个测试对象", len(objects))
	}
}
