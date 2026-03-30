package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// 测试配置
const (
	testBucketPrefix = "s3tool-test-"
	testRegion       = "a1283"
)

// 获取测试用的 S3 客户端
func getTestClient(t *testing.T) *s3.Client {
	// 加载测试配置
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("无法获取用户主目录: %v", err)
	}
	configPath := homeDir + "/.s3tool.yaml"

	if err := config.LoadConfig(configPath); err != nil {
		t.Fatalf("加载配置文件失败: %v", err)
	}

	ctx := context.Background()
	client, err := config.GetS3Client(ctx)
	if err != nil {
		t.Fatalf("初始化 S3 客户端失败: %v", err)
	}

	return client
}

// 生成唯一的测试桶名
func generateTestBucketName() string {
	return fmt.Sprintf("%s%d", testBucketPrefix, time.Now().UnixNano())
}

// 清理测试桶
func cleanupTestBucket(t *testing.T, client *s3.Client, bucketName string) {
	ctx := context.Background()

	// 先清空桶
	listOutput, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err == nil && len(listOutput.Contents) > 0 {
		objectsToDelete := make([]types.ObjectIdentifier, len(listOutput.Contents))
		for i, obj := range listOutput.Contents {
			objectsToDelete[i] = types.ObjectIdentifier{Key: obj.Key}
		}

		_, _ = client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &types.Delete{
				Objects: objectsToDelete,
			},
		})
	}

	// 删除桶
	_, _ = client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
}

// TestIntegration_BucketOperations 测试存储桶操作
func TestIntegration_BucketOperations(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()
	bucketName := generateTestBucketName()

	// 测试结束后清理
	defer cleanupTestBucket(t, client, bucketName)

	t.Run("创建存储桶", func(t *testing.T) {
		input := &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		}

		// 如果区域不是 us-east-1，需要设置 LocationConstraint
		if config.GlobalConfig.Region != "us-east-1" {
			input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(config.GlobalConfig.Region),
			}
		}

		_, err := client.CreateBucket(ctx, input)
		if err != nil {
			t.Fatalf("创建存储桶失败: %v", err)
		}
		t.Logf("存储桶 %s 创建成功", bucketName)
	})

	t.Run("列出存储桶", func(t *testing.T) {
		output, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			t.Fatalf("列出存储桶失败: %v", err)
		}

		found := false
		for _, bucket := range output.Buckets {
			if *bucket.Name == bucketName {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("新创建的存储桶 %s 不在列表中", bucketName)
		}
		t.Logf("找到 %d 个存储桶", len(output.Buckets))
	})

	t.Run("获取存储桶位置", func(t *testing.T) {
		output, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("获取存储桶位置失败: %v", err)
		}

		location := string(output.LocationConstraint)
		if location == "" {
			location = "us-east-1"
		}
		t.Logf("存储桶位置: %s", location)
	})
}

// TestIntegration_ObjectOperations 测试对象操作
func TestIntegration_ObjectOperations(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()
	bucketName := generateTestBucketName()
	objectKey := "test-object.txt"
	objectContent := "Hello, S3Tool Integration Test!"

	// 测试结束后清理
	defer cleanupTestBucket(t, client, bucketName)

	// 先创建桶
	t.Run("创建测试桶", func(t *testing.T) {
		input := &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		}
		if config.GlobalConfig.Region != "us-east-1" {
			input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(config.GlobalConfig.Region),
			}
		}
		_, err := client.CreateBucket(ctx, input)
		if err != nil {
			t.Fatalf("创建存储桶失败: %v", err)
		}
	})

	t.Run("上传对象", func(t *testing.T) {
		input := &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(objectKey),
			Body:        strings.NewReader(objectContent),
			ContentType: aws.String("text/plain"),
		}

		_, err := client.PutObject(ctx, input)
		if err != nil {
			t.Fatalf("上传对象失败: %v", err)
		}
		t.Logf("对象 %s 上传成功", objectKey)
	})

	t.Run("列出对象", func(t *testing.T) {
		output, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("列出对象失败: %v", err)
		}

		if len(output.Contents) != 1 {
			t.Errorf("期望 1 个对象，实际 %d 个", len(output.Contents))
		}

		if *output.Contents[0].Key != objectKey {
			t.Errorf("期望对象键 %s，实际 %s", objectKey, *output.Contents[0].Key)
		}
		t.Logf("找到 %d 个对象", len(output.Contents))
	})

	t.Run("获取对象信息", func(t *testing.T) {
		output, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			t.Fatalf("获取对象信息失败: %v", err)
		}

		if output.ContentLength == nil || *output.ContentLength != int64(len(objectContent)) {
			t.Errorf("对象大小不匹配，期望 %d，实际 %d", len(objectContent), *output.ContentLength)
		}
		t.Logf("对象大小: %d bytes", *output.ContentLength)
	})

	t.Run("下载对象", func(t *testing.T) {
		output, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			t.Fatalf("下载对象失败: %v", err)
		}
		defer output.Body.Close()

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(output.Body)
		if err != nil {
			t.Fatalf("读取对象内容失败: %v", err)
		}

		if buf.String() != objectContent {
			t.Errorf("对象内容不匹配，期望 %s，实际 %s", objectContent, buf.String())
		}
		t.Logf("对象内容: %s", buf.String())
	})

	t.Run("复制对象", func(t *testing.T) {
		destKey := "test-object-copy.txt"
		copySource := bucketName + "/" + objectKey

		_, err := client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String(destKey),
			CopySource: aws.String(copySource),
		})
		if err != nil {
			t.Fatalf("复制对象失败: %v", err)
		}

		// 验证复制成功
		output, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("列出对象失败: %v", err)
		}

		if len(output.Contents) != 2 {
			t.Errorf("期望 2 个对象，实际 %d 个", len(output.Contents))
		}
		t.Logf("复制对象 %s 成功", destKey)
	})

	t.Run("删除对象", func(t *testing.T) {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			t.Fatalf("删除对象失败: %v", err)
		}

		// 验证删除成功
		output, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("列出对象失败: %v", err)
		}

		if len(output.Contents) != 1 {
			t.Errorf("期望 1 个对象，实际 %d 个", len(output.Contents))
		}
		t.Logf("对象 %s 删除成功", objectKey)
	})
}

// TestIntegration_BucketEmpty 测试清空存储桶
func TestIntegration_BucketEmpty(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()
	bucketName := generateTestBucketName()

	// 测试结束后清理
	defer cleanupTestBucket(t, client, bucketName)

	// 创建桶
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	if config.GlobalConfig.Region != "us-east-1" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(config.GlobalConfig.Region),
		}
	}
	_, err := client.CreateBucket(ctx, input)
	if err != nil {
		t.Fatalf("创建存储桶失败: %v", err)
	}

	// 上传多个对象
	for i := 1; i <= 5; i++ {
		key := fmt.Sprintf("object-%d.txt", i)
		content := fmt.Sprintf("Content %d", i)
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   strings.NewReader(content),
		})
		if err != nil {
			t.Fatalf("上传对象 %s 失败: %v", key, err)
		}
	}

	// 验证对象数量
	listOutput, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		t.Fatalf("列出对象失败: %v", err)
	}
	if len(listOutput.Contents) != 5 {
		t.Errorf("期望 5 个对象，实际 %d 个", len(listOutput.Contents))
	}

	// 清空桶
	for {
		listOutput, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("列出对象失败: %v", err)
		}

		if len(listOutput.Contents) == 0 {
			break
		}

		objectsToDelete := make([]types.ObjectIdentifier, len(listOutput.Contents))
		for i, obj := range listOutput.Contents {
			objectsToDelete[i] = types.ObjectIdentifier{Key: obj.Key}
		}

		_, err = client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &types.Delete{
				Objects: objectsToDelete,
			},
		})
		if err != nil {
			t.Fatalf("删除对象失败: %v", err)
		}

		if listOutput.IsTruncated == nil || !*listOutput.IsTruncated {
			break
		}
	}

	// 验证桶已空
	listOutput, err = client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		t.Fatalf("列出对象失败: %v", err)
	}
	if len(listOutput.Contents) != 0 {
		t.Errorf("期望 0 个对象，实际 %d 个", len(listOutput.Contents))
	}
	t.Log("存储桶清空成功")
}

// TestIntegration_Versioning 测试版本控制
func TestIntegration_Versioning(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()
	bucketName := generateTestBucketName()

	// 测试结束后清理
	defer cleanupTestBucket(t, client, bucketName)

	// 创建桶
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	if config.GlobalConfig.Region != "us-east-1" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(config.GlobalConfig.Region),
		}
	}
	_, err := client.CreateBucket(ctx, input)
	if err != nil {
		t.Fatalf("创建存储桶失败: %v", err)
	}

	t.Run("获取版本控制状态", func(t *testing.T) {
		output, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("获取版本控制状态失败: %v", err)
		}
		t.Logf("版本控制状态: %s", output.Status)
	})

	t.Run("启用版本控制", func(t *testing.T) {
		_, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusEnabled,
			},
		})
		if err != nil {
			t.Fatalf("启用版本控制失败: %v", err)
		}

		// 验证启用成功
		output, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("获取版本控制状态失败: %v", err)
		}
		if output.Status != types.BucketVersioningStatusEnabled {
			t.Errorf("版本控制状态应为 Enabled，实际为 %s", output.Status)
		}
		t.Log("版本控制已启用")
	})

	t.Run("暂停版本控制", func(t *testing.T) {
		_, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusSuspended,
			},
		})
		if err != nil {
			t.Fatalf("暂停版本控制失败: %v", err)
		}

		// 验证暂停成功
		output, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("获取版本控制状态失败: %v", err)
		}
		if output.Status != types.BucketVersioningStatusSuspended {
			t.Errorf("版本控制状态应为 Suspended，实际为 %s", output.Status)
		}
		t.Log("版本控制已暂停")
	})
}

// TestIntegration_Tags 测试标签管理
func TestIntegration_Tags(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()
	bucketName := generateTestBucketName()

	// 测试结束后清理
	defer cleanupTestBucket(t, client, bucketName)

	// 创建桶
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	if config.GlobalConfig.Region != "us-east-1" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(config.GlobalConfig.Region),
		}
	}
	_, err := client.CreateBucket(ctx, input)
	if err != nil {
		t.Fatalf("创建存储桶失败: %v", err)
	}

	t.Run("设置存储桶标签", func(t *testing.T) {
		tags := []types.Tag{
			{Key: aws.String("Environment"), Value: aws.String("Test")},
			{Key: aws.String("Project"), Value: aws.String("S3Tool")},
		}

		_, err := client.PutBucketTagging(ctx, &s3.PutBucketTaggingInput{
			Bucket: aws.String(bucketName),
			Tagging: &types.Tagging{
				TagSet: tags,
			},
		})
		if err != nil {
			t.Fatalf("设置存储桶标签失败: %v", err)
		}
		t.Log("存储桶标签设置成功")
	})

	t.Run("获取存储桶标签", func(t *testing.T) {
		output, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("获取存储桶标签失败: %v", err)
		}

		if len(output.TagSet) != 2 {
			t.Errorf("期望 2 个标签，实际 %d 个", len(output.TagSet))
		}

		for _, tag := range output.TagSet {
			t.Logf("标签: %s = %s", *tag.Key, *tag.Value)
		}
	})

	t.Run("删除存储桶标签", func(t *testing.T) {
		_, err := client.DeleteBucketTagging(ctx, &s3.DeleteBucketTaggingInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("删除存储桶标签失败: %v", err)
		}

		// 验证删除成功
		_, err = client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
			Bucket: aws.String(bucketName),
		})
		// 应该返回错误，因为没有标签了
		if err == nil {
			t.Error("期望获取标签时返回错误，但没有")
		}
		t.Log("存储桶标签已删除")
	})
}

// TestIntegration_CORS 测试 CORS 配置
func TestIntegration_CORS(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()
	bucketName := generateTestBucketName()

	// 测试结束后清理
	defer cleanupTestBucket(t, client, bucketName)

	// 创建桶
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	if config.GlobalConfig.Region != "us-east-1" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(config.GlobalConfig.Region),
		}
	}
	_, err := client.CreateBucket(ctx, input)
	if err != nil {
		t.Fatalf("创建存储桶失败: %v", err)
	}

	t.Run("设置 CORS 配置", func(t *testing.T) {
		corsRules := []types.CORSRule{
			{
				AllowedHeaders: []string{"*"},
				AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
				AllowedOrigins: []string{"*"},
				MaxAgeSeconds:  aws.Int32(3000),
			},
		}

		_, err := client.PutBucketCors(ctx, &s3.PutBucketCorsInput{
			Bucket: aws.String(bucketName),
			CORSConfiguration: &types.CORSConfiguration{
				CORSRules: corsRules,
			},
		})
		if err != nil {
			t.Fatalf("设置 CORS 配置失败: %v", err)
		}
		t.Log("CORS 配置设置成功")
	})

	t.Run("获取 CORS 配置", func(t *testing.T) {
		output, err := client.GetBucketCors(ctx, &s3.GetBucketCorsInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("获取 CORS 配置失败: %v", err)
		}

		if len(output.CORSRules) != 1 {
			t.Errorf("期望 1 个 CORS 规则，实际 %d 个", len(output.CORSRules))
		}
		t.Logf("CORS 规则数量: %d", len(output.CORSRules))
	})

	t.Run("删除 CORS 配置", func(t *testing.T) {
		_, err := client.DeleteBucketCors(ctx, &s3.DeleteBucketCorsInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("删除 CORS 配置失败: %v", err)
		}

		// 验证删除成功
		_, err = client.GetBucketCors(ctx, &s3.GetBucketCorsInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			t.Error("期望获取 CORS 时返回错误，但没有")
		}
		t.Log("CORS 配置已删除")
	})
}

// TestIntegration_Connection 测试连接
func TestIntegration_Connection(t *testing.T) {
	t.Run("测试 S3 连接", func(t *testing.T) {
		client := getTestClient(t)
		ctx := context.Background()

		// 尝试列出存储桶来验证连接
		_, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			t.Fatalf("无法连接到 S3: %v", err)
		}
		t.Log("S3 连接成功")
	})

	t.Run("验证配置", func(t *testing.T) {
		if config.GlobalConfig.Region == "" {
			t.Error("Region 未配置")
		}
		if config.GlobalConfig.AccessKey == "" {
			t.Error("AccessKey 未配置")
		}
		if config.GlobalConfig.SecretKey == "" {
			t.Error("SecretKey 未配置")
		}
		if config.GlobalConfig.Endpoint == "" {
			t.Error("Endpoint 未配置")
		}

		t.Logf("Region: %s", config.GlobalConfig.Region)
		t.Logf("Endpoint: %s", config.GlobalConfig.Endpoint)
		t.Logf("ForcePathStyle: %v", config.GlobalConfig.ForcePathStyle)
	})
}
