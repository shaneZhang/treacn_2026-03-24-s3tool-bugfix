package tests

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

func TestMain(m *testing.M) {
	if err := config.LoadConfig("/Users/zhangyuqing/.s3tool.yaml"); err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestBucketList(t *testing.T) {
	ctx := context.Background()
	client, err := config.GetS3Client(ctx)
	if err != nil {
		t.Fatalf("Failed to create S3 client: %v", err)
	}

	output, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		t.Errorf("ListBuckets() error = %v", err)
		return
	}

	if output == nil {
		t.Error("ListBuckets() returned nil output")
		return
	}

	t.Logf("Found %d buckets", len(output.Buckets))
	for _, bucket := range output.Buckets {
		t.Logf("  - %s (created: %s)", *bucket.Name, bucket.CreationDate.Format("2006-01-02"))
	}
}

func TestBucketCreateAndDelete(t *testing.T) {
	tc, err := NewTestContext()
	if err != nil {
		t.Fatalf("Failed to create test context: %v", err)
	}
	defer tc.Cleanup()

	t.Run("create bucket", func(t *testing.T) {
		err := tc.CreateTestBucket()
		if err != nil {
			t.Errorf("CreateTestBucket() error = %v", err)
			return
		}

		exists, err := BucketExists(tc.Client, tc.Context, tc.BucketName)
		if err != nil {
			t.Errorf("BucketExists() error = %v", err)
			return
		}
		if !exists {
			t.Error("Bucket was not created")
		}
		t.Logf("Bucket %s created successfully", tc.BucketName)
	})

	t.Run("delete bucket", func(t *testing.T) {
		err := tc.DeleteTestBucket()
		if err != nil {
			t.Errorf("DeleteTestBucket() error = %v", err)
			return
		}

		time.Sleep(1 * time.Second)

		exists, _ := BucketExists(tc.Client, tc.Context, tc.BucketName)
		if exists {
			t.Error("Bucket was not deleted")
		}
		t.Logf("Bucket %s deleted successfully", tc.BucketName)
	})
}

func TestBucketLocation(t *testing.T) {
	tc, err := NewTestContext()
	if err != nil {
		t.Fatalf("Failed to create test context: %v", err)
	}
	defer tc.Cleanup()

	err = tc.CreateTestBucket()
	if err != nil {
		t.Fatalf("Failed to create test bucket: %v", err)
	}
	defer tc.DeleteTestBucket()

	output, err := tc.Client.GetBucketLocation(tc.Context, &s3.GetBucketLocationInput{
		Bucket: aws.String(tc.BucketName),
	})
	if err != nil {
		t.Errorf("GetBucketLocation() error = %v", err)
		return
	}

	location := string(output.LocationConstraint)
	if location == "" {
		location = "us-east-1"
	}
	t.Logf("Bucket %s is in region: %s", tc.BucketName, location)
}

func TestBucketEmpty(t *testing.T) {
	tc, err := NewTestContext()
	if err != nil {
		t.Fatalf("Failed to create test context: %v", err)
	}
	defer tc.Cleanup()

	err = tc.CreateTestBucket()
	if err != nil {
		t.Fatalf("Failed to create test bucket: %v", err)
	}
	defer tc.DeleteTestBucket()

	testObjects := []string{"file1.txt", "file2.txt", "dir/file3.txt"}
	for _, key := range testObjects {
		_, err := tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(key),
			Body:   bytes.NewReader([]byte("test content")),
		})
		if err != nil {
			t.Fatalf("Failed to upload test object %s: %v", key, err)
		}
	}

	t.Run("verify objects exist", func(t *testing.T) {
		listOutput, err := tc.Client.ListObjectsV2(tc.Context, &s3.ListObjectsV2Input{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			t.Errorf("ListObjectsV2() error = %v", err)
			return
		}
		if len(listOutput.Contents) != 3 {
			t.Errorf("Expected 3 objects, got %d", len(listOutput.Contents))
		}
	})

	t.Run("empty bucket", func(t *testing.T) {
		deletedCount := 0
		for {
			listOutput, err := tc.Client.ListObjectsV2(tc.Context, &s3.ListObjectsV2Input{
				Bucket: aws.String(tc.BucketName),
			})
			if err != nil {
				t.Errorf("ListObjectsV2() error = %v", err)
				return
			}

			if len(listOutput.Contents) == 0 {
				break
			}

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
				t.Errorf("DeleteObjects() error = %v", err)
				return
			}

			deletedCount += len(listOutput.Contents)

			if listOutput.IsTruncated == nil || !*listOutput.IsTruncated {
				break
			}
		}

		t.Logf("Deleted %d objects", deletedCount)

		listOutput, _ := tc.Client.ListObjectsV2(tc.Context, &s3.ListObjectsV2Input{
			Bucket: aws.String(tc.BucketName),
		})
		if len(listOutput.Contents) != 0 {
			t.Errorf("Bucket not empty, still has %d objects", len(listOutput.Contents))
		}
	})
}

func TestBucketCreateWithInvalidName(t *testing.T) {
	ctx := context.Background()
	client, err := config.GetS3Client(ctx)
	if err != nil {
		t.Fatalf("Failed to create S3 client: %v", err)
	}

	invalidNames := []string{
		"InvalidBucketName",
		"invalid.bucket.name",
		"invalid_bucket_name_",
		"invalid",
	}

	for _, name := range invalidNames {
		t.Run(name, func(t *testing.T) {
			_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
				Bucket: aws.String(name),
			})
			if err == nil {
				t.Errorf("Expected error for invalid bucket name %s, got nil", name)
				client.DeleteBucket(ctx, &s3.DeleteBucketInput{
					Bucket: aws.String(name),
				})
			} else {
				t.Logf("Got expected error for %s: %v", name, err)
			}
		})
	}
}

func TestBucketDeleteNonEmpty(t *testing.T) {
	tc, err := NewTestContext()
	if err != nil {
		t.Fatalf("Failed to create test context: %v", err)
	}
	defer tc.Cleanup()

	err = tc.CreateTestBucket()
	if err != nil {
		t.Fatalf("Failed to create test bucket: %v", err)
	}

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String("test.txt"),
		Body:   bytes.NewReader([]byte("test")),
	})
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	_, err = tc.Client.DeleteBucket(tc.Context, &s3.DeleteBucketInput{
		Bucket: aws.String(tc.BucketName),
	})
	if err == nil {
		t.Error("Expected error when deleting non-empty bucket, got nil")
	} else {
		t.Logf("Got expected error: %v", err)
	}

	tc.DeleteTestBucket()
}

func TestBucketTagging(t *testing.T) {
	tc, err := NewTestContext()
	if err != nil {
		t.Fatalf("Failed to create test context: %v", err)
	}
	defer tc.Cleanup()

	err = tc.CreateTestBucket()
	if err != nil {
		t.Fatalf("Failed to create test bucket: %v", err)
	}
	defer tc.DeleteTestBucket()

	t.Run("put bucket tags", func(t *testing.T) {
		tags := []types.Tag{
			{Key: aws.String("Environment"), Value: aws.String("Test")},
			{Key: aws.String("Project"), Value: aws.String("S3Tool")},
		}

		_, err := tc.Client.PutBucketTagging(tc.Context, &s3.PutBucketTaggingInput{
			Bucket: aws.String(tc.BucketName),
			Tagging: &types.Tagging{
				TagSet: tags,
			},
		})
		if err != nil {
			t.Errorf("PutBucketTagging() error = %v", err)
			return
		}
		t.Log("Bucket tags set successfully")
	})

	t.Run("get bucket tags", func(t *testing.T) {
		output, err := tc.Client.GetBucketTagging(tc.Context, &s3.GetBucketTaggingInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			t.Errorf("GetBucketTagging() error = %v", err)
			return
		}

		if len(output.TagSet) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(output.TagSet))
		}

		for _, tag := range output.TagSet {
			t.Logf("Tag: %s = %s", *tag.Key, *tag.Value)
		}
	})

	t.Run("delete bucket tags", func(t *testing.T) {
		_, err := tc.Client.DeleteBucketTagging(tc.Context, &s3.DeleteBucketTaggingInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			t.Errorf("DeleteBucketTagging() error = %v", err)
			return
		}
		t.Log("Bucket tags deleted successfully")
	})
}

func TestBucketVersioning(t *testing.T) {
	tc, err := NewTestContext()
	if err != nil {
		t.Fatalf("Failed to create test context: %v", err)
	}
	defer tc.Cleanup()

	err = tc.CreateTestBucket()
	if err != nil {
		t.Fatalf("Failed to create test bucket: %v", err)
	}
	defer tc.DeleteTestBucket()

	t.Run("get initial versioning status", func(t *testing.T) {
		output, err := tc.Client.GetBucketVersioning(tc.Context, &s3.GetBucketVersioningInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			t.Errorf("GetBucketVersioning() error = %v", err)
			return
		}
		t.Logf("Initial versioning status: %s", output.Status)
	})

	t.Run("enable versioning", func(t *testing.T) {
		_, err := tc.Client.PutBucketVersioning(tc.Context, &s3.PutBucketVersioningInput{
			Bucket: aws.String(tc.BucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: "Enabled",
			},
		})
		if err != nil {
			t.Errorf("PutBucketVersioning() error = %v", err)
			return
		}
		t.Log("Versioning enabled successfully")
	})

	t.Run("verify versioning enabled", func(t *testing.T) {
		output, err := tc.Client.GetBucketVersioning(tc.Context, &s3.GetBucketVersioningInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			t.Errorf("GetBucketVersioning() error = %v", err)
			return
		}
		if output.Status != types.BucketVersioningStatusEnabled {
			t.Errorf("Expected versioning status Enabled, got %s", output.Status)
		}
	})

	t.Run("suspend versioning", func(t *testing.T) {
		_, err := tc.Client.PutBucketVersioning(tc.Context, &s3.PutBucketVersioningInput{
			Bucket: aws.String(tc.BucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: "Suspended",
			},
		})
		if err != nil {
			t.Errorf("PutBucketVersioning() error = %v", err)
			return
		}
		t.Log("Versioning suspended successfully")
	})
}

func TestBucketEncryption(t *testing.T) {
	tc, err := NewTestContext()
	if err != nil {
		t.Fatalf("Failed to create test context: %v", err)
	}
	defer tc.Cleanup()

	err = tc.CreateTestBucket()
	if err != nil {
		t.Fatalf("Failed to create test bucket: %v", err)
	}
	defer tc.DeleteTestBucket()

	t.Run("enable encryption", func(t *testing.T) {
		_, err := tc.Client.PutBucketEncryption(tc.Context, &s3.PutBucketEncryptionInput{
			Bucket: aws.String(tc.BucketName),
			ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
				Rules: []types.ServerSideEncryptionRule{
					{
						ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
							SSEAlgorithm: types.ServerSideEncryptionAes256,
						},
					},
				},
			},
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Encryption not supported by this S3 implementation")
			}
			t.Errorf("PutBucketEncryption() error = %v", err)
			return
		}
		t.Log("Encryption enabled successfully")
	})

	t.Run("get encryption config", func(t *testing.T) {
		output, err := tc.Client.GetBucketEncryption(tc.Context, &s3.GetBucketEncryptionInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") || strings.Contains(err.Error(), "ServerSideEncryptionConfigurationNotFoundError") {
				t.Skip("Encryption not supported or not configured")
			}
			t.Errorf("GetBucketEncryption() error = %v", err)
			return
		}

		if output.ServerSideEncryptionConfiguration != nil {
			for _, rule := range output.ServerSideEncryptionConfiguration.Rules {
				t.Logf("Encryption algorithm: %s", rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
			}
		}
	})

	t.Run("disable encryption", func(t *testing.T) {
		_, err := tc.Client.DeleteBucketEncryption(tc.Context, &s3.DeleteBucketEncryptionInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Encryption not supported by this S3 implementation")
			}
			t.Errorf("DeleteBucketEncryption() error = %v", err)
			return
		}
		t.Log("Encryption disabled successfully")
	})
}
