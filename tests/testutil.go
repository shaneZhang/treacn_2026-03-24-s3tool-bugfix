package tests

import (
	"context"
	"fmt"
	"os"
	"time"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	TestBucketPrefix = "s3tool-test-"
	TestConfigPath   = "/Users/zhangyuqing/.s3tool.yaml"
)

type TestContext struct {
	Client     *s3.Client
	Context    context.Context
	BucketName string
	TestDir    string
}

func NewTestContext() (*TestContext, error) {
	ctx := context.Background()

	if err := config.LoadConfig(TestConfigPath); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client, err := config.GetS3Client(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	bucketName := fmt.Sprintf("%s%d", TestBucketPrefix, time.Now().UnixNano())

	testDir, err := os.MkdirTemp("", "s3tool-test-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	return &TestContext{
		Client:     client,
		Context:    ctx,
		BucketName: bucketName,
		TestDir:    testDir,
	}, nil
}

func (tc *TestContext) CreateTestBucket() error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(tc.BucketName),
	}

	if config.GlobalConfig.Region != "" && config.GlobalConfig.Region != "us-east-1" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(config.GlobalConfig.Region),
		}
	}

	_, err := tc.Client.CreateBucket(tc.Context, input)
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", tc.BucketName, err)
	}

	return nil
}

func (tc *TestContext) DeleteTestBucket() error {
	listOutput, err := tc.Client.ListObjectsV2(tc.Context, &s3.ListObjectsV2Input{
		Bucket: aws.String(tc.BucketName),
	})
	if err == nil && len(listOutput.Contents) > 0 {
		objectsToDelete := make([]types.ObjectIdentifier, len(listOutput.Contents))
		for i, obj := range listOutput.Contents {
			objectsToDelete[i] = types.ObjectIdentifier{Key: obj.Key}
		}

		_, err = tc.Client.DeleteObjects(tc.Context, &s3.DeleteObjectsInput{
			Bucket: aws.String(tc.BucketName),
			Delete: &types.Delete{
				Objects: objectsToDelete,
			},
		})
		if err != nil {
			fmt.Printf("Warning: failed to delete objects: %v\n", err)
		}
	}

	_, err = tc.Client.DeleteBucket(tc.Context, &s3.DeleteBucketInput{
		Bucket: aws.String(tc.BucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete bucket %s: %w", tc.BucketName, err)
	}

	return nil
}

func (tc *TestContext) CreateTestFile(content string) (string, error) {
	filePath := fmt.Sprintf("%s/test-%d.txt", tc.TestDir, time.Now().UnixNano())
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create test file: %w", err)
	}
	return filePath, nil
}

func (tc *TestContext) Cleanup() {
	if tc.TestDir != "" {
		os.RemoveAll(tc.TestDir)
	}
}

func BucketExists(client *s3.Client, ctx context.Context, bucketName string) (bool, error) {
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func ObjectExists(client *s3.Client, ctx context.Context, bucket, key string) (bool, error) {
	_, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func UploadTestObject(client *s3.Client, ctx context.Context, bucket, key, content string) error {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   nil,
	})
	return err
}
