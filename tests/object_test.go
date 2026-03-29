package tests

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestObjectPutAndGet(t *testing.T) {
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

	testKey := "test-object.txt"
	testContent := "Hello, S3Tool Test!"

	t.Run("put object", func(t *testing.T) {
		_, err := tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
			Bucket:      aws.String(tc.BucketName),
			Key:         aws.String(testKey),
			Body:        bytes.NewReader([]byte(testContent)),
			ContentType: aws.String("text/plain"),
		})
		if err != nil {
			t.Errorf("PutObject() error = %v", err)
			return
		}
		t.Logf("Object %s uploaded successfully", testKey)
	})

	t.Run("get object", func(t *testing.T) {
		output, err := tc.Client.GetObject(tc.Context, &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("GetObject() error = %v", err)
			return
		}
		defer output.Body.Close()

		body, err := io.ReadAll(output.Body)
		if err != nil {
			t.Errorf("Failed to read object body: %v", err)
			return
		}

		if string(body) != testContent {
			t.Errorf("Expected content %q, got %q", testContent, string(body))
		}

		if output.ContentType != nil && *output.ContentType != "text/plain" {
			t.Errorf("Expected content type text/plain, got %s", *output.ContentType)
		}

		t.Logf("Object content: %s", string(body))
	})
}

func TestObjectList(t *testing.T) {
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

	testObjects := []struct {
		key     string
		content string
	}{
		{"file1.txt", "content1"},
		{"file2.txt", "content2"},
		{"dir/file3.txt", "content3"},
		{"dir/subdir/file4.txt", "content4"},
	}

	for _, obj := range testObjects {
		_, err := tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(obj.key),
			Body:   bytes.NewReader([]byte(obj.content)),
		})
		if err != nil {
			t.Fatalf("Failed to upload test object %s: %v", obj.key, err)
		}
	}

	t.Run("list all objects", func(t *testing.T) {
		output, err := tc.Client.ListObjectsV2(tc.Context, &s3.ListObjectsV2Input{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			t.Errorf("ListObjectsV2() error = %v", err)
			return
		}

		if len(output.Contents) != 4 {
			t.Errorf("Expected 4 objects, got %d", len(output.Contents))
		}

		for _, obj := range output.Contents {
			t.Logf("  - %s (%d bytes)", *obj.Key, obj.Size)
		}
	})

	t.Run("list with prefix", func(t *testing.T) {
		output, err := tc.Client.ListObjectsV2(tc.Context, &s3.ListObjectsV2Input{
			Bucket: aws.String(tc.BucketName),
			Prefix: aws.String("dir/"),
		})
		if err != nil {
			t.Errorf("ListObjectsV2() error = %v", err)
			return
		}

		if len(output.Contents) != 2 {
			t.Errorf("Expected 2 objects with prefix 'dir/', got %d", len(output.Contents))
		}

		for _, obj := range output.Contents {
			t.Logf("  - %s", *obj.Key)
		}
	})

	t.Run("list with delimiter", func(t *testing.T) {
		output, err := tc.Client.ListObjectsV2(tc.Context, &s3.ListObjectsV2Input{
			Bucket:    aws.String(tc.BucketName),
			Delimiter: aws.String("/"),
		})
		if err != nil {
			t.Errorf("ListObjectsV2() error = %v", err)
			return
		}

		t.Logf("Objects at root level:")
		for _, obj := range output.Contents {
			t.Logf("  - %s", *obj.Key)
		}

		t.Logf("Common prefixes:")
		for _, prefix := range output.CommonPrefixes {
			t.Logf("  - %s", *prefix.Prefix)
		}
	})
}

func TestObjectDelete(t *testing.T) {
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

	testKey := "delete-test.txt"

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
		Body:   bytes.NewReader([]byte("test content")),
	})
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	t.Run("delete object", func(t *testing.T) {
		_, err := tc.Client.DeleteObject(tc.Context, &s3.DeleteObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("DeleteObject() error = %v", err)
			return
		}
		t.Logf("Object %s deleted successfully", testKey)
	})

	t.Run("verify deletion", func(t *testing.T) {
		_, err := tc.Client.HeadObject(tc.Context, &s3.HeadObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err == nil {
			t.Error("Expected error for deleted object, got nil")
		} else {
			t.Logf("Got expected error: %v", err)
		}
	})
}

func TestObjectCopy(t *testing.T) {
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

	srcKey := "source.txt"
	destKey := "destination.txt"
	testContent := "content to copy"

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(srcKey),
		Body:   bytes.NewReader([]byte(testContent)),
	})
	if err != nil {
		t.Fatalf("Failed to upload source object: %v", err)
	}

	t.Run("copy object", func(t *testing.T) {
		copySource := tc.BucketName + "/" + srcKey
		_, err := tc.Client.CopyObject(tc.Context, &s3.CopyObjectInput{
			Bucket:     aws.String(tc.BucketName),
			Key:        aws.String(destKey),
			CopySource: aws.String(copySource),
		})
		if err != nil {
			t.Errorf("CopyObject() error = %v", err)
			return
		}
		t.Logf("Object copied from %s to %s", srcKey, destKey)
	})

	t.Run("verify copy", func(t *testing.T) {
		output, err := tc.Client.GetObject(tc.Context, &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(destKey),
		})
		if err != nil {
			t.Errorf("GetObject() error = %v", err)
			return
		}
		defer output.Body.Close()

		body, _ := io.ReadAll(output.Body)
		if string(body) != testContent {
			t.Errorf("Expected content %q, got %q", testContent, string(body))
		}
		t.Log("Copy verified successfully")
	})
}

func TestObjectMove(t *testing.T) {
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

	srcKey := "original.txt"
	destKey := "moved.txt"
	testContent := "content to move"

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(srcKey),
		Body:   bytes.NewReader([]byte(testContent)),
	})
	if err != nil {
		t.Fatalf("Failed to upload source object: %v", err)
	}

	t.Run("move object", func(t *testing.T) {
		copySource := tc.BucketName + "/" + srcKey

		_, err := tc.Client.CopyObject(tc.Context, &s3.CopyObjectInput{
			Bucket:     aws.String(tc.BucketName),
			Key:        aws.String(destKey),
			CopySource: aws.String(copySource),
		})
		if err != nil {
			t.Errorf("CopyObject() error = %v", err)
			return
		}

		_, err = tc.Client.DeleteObject(tc.Context, &s3.DeleteObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(srcKey),
		})
		if err != nil {
			t.Errorf("DeleteObject() error = %v", err)
			return
		}

		t.Logf("Object moved from %s to %s", srcKey, destKey)
	})

	t.Run("verify move", func(t *testing.T) {
		_, err := tc.Client.HeadObject(tc.Context, &s3.HeadObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(srcKey),
		})
		if err == nil {
			t.Error("Source object should not exist after move")
		}

		output, err := tc.Client.GetObject(tc.Context, &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(destKey),
		})
		if err != nil {
			t.Errorf("GetObject() error = %v", err)
			return
		}
		defer output.Body.Close()

		body, _ := io.ReadAll(output.Body)
		if string(body) != testContent {
			t.Errorf("Expected content %q, got %q", testContent, string(body))
		}
		t.Log("Move verified successfully")
	})
}

func TestObjectHead(t *testing.T) {
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

	testKey := "head-test.txt"
	testContent := "test content for head"

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket:      aws.String(tc.BucketName),
		Key:         aws.String(testKey),
		Body:        bytes.NewReader([]byte(testContent)),
		ContentType: aws.String("text/plain"),
	})
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	t.Run("head object", func(t *testing.T) {
		output, err := tc.Client.HeadObject(tc.Context, &s3.HeadObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("HeadObject() error = %v", err)
			return
		}

		t.Logf("Content Length: %d", output.ContentLength)
		t.Logf("Content Type: %s", *output.ContentType)
		t.Logf("Last Modified: %s", output.LastModified.Format("2006-01-02 15:04:05"))
		t.Logf("ETag: %s", *output.ETag)

		if output.ContentLength == nil || *output.ContentLength != int64(len(testContent)) {
			t.Errorf("Expected content length %d, got %d", len(testContent), output.ContentLength)
		}
	})

	t.Run("head non-existent object", func(t *testing.T) {
		_, err := tc.Client.HeadObject(tc.Context, &s3.HeadObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String("non-existent.txt"),
		})
		if err == nil {
			t.Error("Expected error for non-existent object, got nil")
		} else {
			t.Logf("Got expected error: %v", err)
		}
	})
}

func TestObjectTagging(t *testing.T) {
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

	testKey := "tagged-object.txt"

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
		Body:   bytes.NewReader([]byte("test content")),
	})
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	t.Run("put object tags", func(t *testing.T) {
		tags := []types.Tag{
			{Key: aws.String("Type"), Value: aws.String("Test")},
			{Key: aws.String("Owner"), Value: aws.String("S3Tool")},
		}

		_, err := tc.Client.PutObjectTagging(tc.Context, &s3.PutObjectTaggingInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
			Tagging: &types.Tagging{
				TagSet: tags,
			},
		})
		if err != nil {
			t.Errorf("PutObjectTagging() error = %v", err)
			return
		}
		t.Log("Object tags set successfully")
	})

	t.Run("get object tags", func(t *testing.T) {
		output, err := tc.Client.GetObjectTagging(tc.Context, &s3.GetObjectTaggingInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("GetObjectTagging() error = %v", err)
			return
		}

		if len(output.TagSet) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(output.TagSet))
		}

		for _, tag := range output.TagSet {
			t.Logf("Tag: %s = %s", *tag.Key, *tag.Value)
		}
	})

	t.Run("delete object tags", func(t *testing.T) {
		_, err := tc.Client.DeleteObjectTagging(tc.Context, &s3.DeleteObjectTaggingInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("DeleteObjectTagging() error = %v", err)
			return
		}
		t.Log("Object tags deleted successfully")
	})
}

func TestObjectWithMetadata(t *testing.T) {
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

	testKey := "metadata-test.txt"
	metadata := map[string]string{
		"author":      "s3tool-test",
		"project":     "s3tool",
		"description": "test file with metadata",
	}

	t.Run("put object with metadata", func(t *testing.T) {
		_, err := tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
			Bucket:   aws.String(tc.BucketName),
			Key:      aws.String(testKey),
			Body:     bytes.NewReader([]byte("test content")),
			Metadata: metadata,
		})
		if err != nil {
			t.Errorf("PutObject() error = %v", err)
			return
		}
		t.Log("Object with metadata uploaded successfully")
	})

	t.Run("get object metadata", func(t *testing.T) {
		output, err := tc.Client.HeadObject(tc.Context, &s3.HeadObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("HeadObject() error = %v", err)
			return
		}

		for key, expectedValue := range metadata {
			if output.Metadata[key] != expectedValue {
				t.Errorf("Expected metadata %s=%s, got %s", key, expectedValue, output.Metadata[key])
			}
		}

		t.Logf("Metadata: %v", output.Metadata)
	})
}

func TestObjectWithStorageClass(t *testing.T) {
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

	testKey := "storage-class-test.txt"

	t.Run("put object with storage class", func(t *testing.T) {
		_, err := tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
			Bucket:       aws.String(tc.BucketName),
			Key:          aws.String(testKey),
			Body:         bytes.NewReader([]byte("test content")),
			StorageClass: types.StorageClassStandard,
		})
		if err != nil {
			t.Errorf("PutObject() error = %v", err)
			return
		}
		t.Log("Object with storage class uploaded successfully")
	})

	t.Run("verify storage class", func(t *testing.T) {
		output, err := tc.Client.HeadObject(tc.Context, &s3.HeadObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("HeadObject() error = %v", err)
			return
		}

		t.Logf("Storage class: %s", output.StorageClass)
	})
}

func TestLargeObjectUpload(t *testing.T) {
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

	testKey := "large-object.bin"
	size := 5 * 1024 * 1024

	largeContent := make([]byte, size)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	t.Run("upload large object", func(t *testing.T) {
		start := time.Now()
		_, err := tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
			Body:   bytes.NewReader(largeContent),
		})
		if err != nil {
			t.Errorf("PutObject() error = %v", err)
			return
		}
		t.Logf("Large object (%d bytes) uploaded in %v", size, time.Since(start))
	})

	t.Run("download and verify large object", func(t *testing.T) {
		output, err := tc.Client.GetObject(tc.Context, &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("GetObject() error = %v", err)
			return
		}
		defer output.Body.Close()

		downloaded, err := io.ReadAll(output.Body)
		if err != nil {
			t.Errorf("Failed to read object body: %v", err)
			return
		}

		if len(downloaded) != size {
			t.Errorf("Expected size %d, got %d", size, len(downloaded))
		}

		for i := 0; i < size; i++ {
			if downloaded[i] != largeContent[i] {
				t.Errorf("Content mismatch at position %d", i)
				break
			}
		}

		t.Log("Large object verified successfully")
	})
}

func TestObjectWithSpecialCharacters(t *testing.T) {
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

	testCases := []struct {
		name    string
		key     string
		content string
	}{
		{"spaces", "path with spaces/file name.txt", "content with spaces"},
		{"unicode", "unicode/中文文件名.txt", "中文内容"},
		{"special_chars", "special/file-name_test.txt", "content"},
		{"nested_path", "a/b/c/d/e/deeply-nested.txt", "nested content"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
				Bucket: aws.String(tc.BucketName),
				Key:    aws.String(testCase.key),
				Body:   bytes.NewReader([]byte(testCase.content)),
			})
			if err != nil {
				t.Errorf("PutObject() error = %v", err)
				return
			}

			output, err := tc.Client.GetObject(tc.Context, &s3.GetObjectInput{
				Bucket: aws.String(tc.BucketName),
				Key:    aws.String(testCase.key),
			})
			if err != nil {
				t.Errorf("GetObject() error = %v", err)
				return
			}
			defer output.Body.Close()

			body, _ := io.ReadAll(output.Body)
			if string(body) != testCase.content {
				t.Errorf("Expected content %q, got %q", testCase.content, string(body))
			}

			t.Logf("Successfully handled key: %s", testCase.key)
		})
	}
}

func TestObjectACL(t *testing.T) {
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

	testKey := "acl-test.txt"

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
		Body:   bytes.NewReader([]byte("test content")),
	})
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	t.Run("get object ACL", func(t *testing.T) {
		output, err := tc.Client.GetObjectAcl(tc.Context, &s3.GetObjectAclInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("ACL not supported by this S3 implementation")
			}
			t.Errorf("GetObjectAcl() error = %v", err)
			return
		}

		if output.Owner != nil && output.Owner.DisplayName != nil {
			t.Logf("Owner: %s", *output.Owner.DisplayName)
		} else {
			t.Log("Owner information not available")
		}
		for _, grant := range output.Grants {
			if grant.Grantee != nil && grant.Grantee.DisplayName != nil {
				t.Logf("Grant: %s -> %s", grant.Permission, *grant.Grantee.DisplayName)
			} else {
				t.Logf("Grant: %s", grant.Permission)
			}
		}
	})

	t.Run("set object ACL", func(t *testing.T) {
		_, err := tc.Client.PutObjectAcl(tc.Context, &s3.PutObjectAclInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
			ACL:    types.ObjectCannedACLPrivate,
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("ACL not supported by this S3 implementation")
			}
			t.Errorf("PutObjectAcl() error = %v", err)
			return
		}
		t.Log("Object ACL set successfully")
	})
}

func TestFileUploadDownload(t *testing.T) {
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

	testContent := "This is a test file for upload/download"
	localFile, err := tc.CreateTestFile(testContent)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	testKey := "uploaded-file.txt"

	t.Run("upload from file", func(t *testing.T) {
		file, err := os.Open(localFile)
		if err != nil {
			t.Fatalf("Failed to open local file: %v", err)
		}
		defer file.Close()

		_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
			Body:   file,
		})
		if err != nil {
			t.Errorf("PutObject() error = %v", err)
			return
		}
		t.Log("File uploaded successfully")
	})

	t.Run("download to file", func(t *testing.T) {
		output, err := tc.Client.GetObject(tc.Context, &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("GetObject() error = %v", err)
			return
		}
		defer output.Body.Close()

		downloadPath := fmt.Sprintf("%s/downloaded.txt", tc.TestDir)
		file, err := os.Create(downloadPath)
		if err != nil {
			t.Fatalf("Failed to create download file: %v", err)
		}
		defer file.Close()

		_, err = io.Copy(file, output.Body)
		if err != nil {
			t.Errorf("Failed to write to file: %v", err)
			return
		}

		content, err := os.ReadFile(downloadPath)
		if err != nil {
			t.Fatalf("Failed to read downloaded file: %v", err)
		}

		if string(content) != testContent {
			t.Errorf("Expected content %q, got %q", testContent, string(content))
		}

		t.Log("File downloaded and verified successfully")
	})
}
