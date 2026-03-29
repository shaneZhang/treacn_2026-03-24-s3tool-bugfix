package tests

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestMultipartInitAndAbort(t *testing.T) {
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

	testKey := "multipart-test.txt"
	var uploadID string

	t.Run("initiate multipart upload", func(t *testing.T) {
		output, err := tc.Client.CreateMultipartUpload(tc.Context, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("CreateMultipartUpload() error = %v", err)
			return
		}

		uploadID = *output.UploadId
		t.Logf("Multipart upload initiated, UploadId: %s", uploadID)

		if uploadID == "" {
			t.Error("UploadId is empty")
		}
	})

	t.Run("list multipart uploads", func(t *testing.T) {
		output, err := tc.Client.ListMultipartUploads(tc.Context, &s3.ListMultipartUploadsInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			t.Errorf("ListMultipartUploads() error = %v", err)
			return
		}

		found := false
		for _, upload := range output.Uploads {
			if *upload.UploadId == uploadID {
				found = true
				break
			}
		}

		if !found {
			t.Error("Upload not found in list")
		}

		t.Logf("Found %d multipart uploads", len(output.Uploads))
	})

	t.Run("abort multipart upload", func(t *testing.T) {
		_, err := tc.Client.AbortMultipartUpload(tc.Context, &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(tc.BucketName),
			Key:      aws.String(testKey),
			UploadId: aws.String(uploadID),
		})
		if err != nil {
			t.Errorf("AbortMultipartUpload() error = %v", err)
			return
		}

		t.Log("Multipart upload aborted successfully")
	})

	t.Run("verify upload aborted", func(t *testing.T) {
		output, err := tc.Client.ListMultipartUploads(tc.Context, &s3.ListMultipartUploadsInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			t.Errorf("ListMultipartUploads() error = %v", err)
			return
		}

		for _, upload := range output.Uploads {
			if *upload.UploadId == uploadID {
				t.Error("Upload should have been aborted")
			}
		}

		t.Log("Upload successfully aborted and removed from list")
	})
}

func TestMultipartUploadComplete(t *testing.T) {
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

	testKey := "multipart-complete.txt"
	partSize := 5 * 1024 * 1024

	part1 := make([]byte, partSize)
	part2 := make([]byte, partSize)
	for i := range part1 {
		part1[i] = byte(i % 256)
	}
	for i := range part2 {
		part2[i] = byte((i + 128) % 256)
	}

	var uploadID string
	var completedParts []types.CompletedPart

	t.Run("initiate multipart upload", func(t *testing.T) {
		output, err := tc.Client.CreateMultipartUpload(tc.Context, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("CreateMultipartUpload() error = %v", err)
			return
		}
		uploadID = *output.UploadId
		t.Logf("UploadId: %s", uploadID)
	})

	t.Run("upload part 1", func(t *testing.T) {
		output, err := tc.Client.UploadPart(tc.Context, &s3.UploadPartInput{
			Bucket:     aws.String(tc.BucketName),
			Key:        aws.String(testKey),
			UploadId:   aws.String(uploadID),
			PartNumber: aws.Int32(1),
			Body:       bytes.NewReader(part1),
		})
		if err != nil {
			t.Errorf("UploadPart() error = %v", err)
			return
		}

		completedParts = append(completedParts, types.CompletedPart{
			ETag:       output.ETag,
			PartNumber: aws.Int32(1),
		})
		t.Logf("Part 1 uploaded, ETag: %s", *output.ETag)
	})

	t.Run("upload part 2", func(t *testing.T) {
		output, err := tc.Client.UploadPart(tc.Context, &s3.UploadPartInput{
			Bucket:     aws.String(tc.BucketName),
			Key:        aws.String(testKey),
			UploadId:   aws.String(uploadID),
			PartNumber: aws.Int32(2),
			Body:       bytes.NewReader(part2),
		})
		if err != nil {
			t.Errorf("UploadPart() error = %v", err)
			return
		}

		completedParts = append(completedParts, types.CompletedPart{
			ETag:       output.ETag,
			PartNumber: aws.Int32(2),
		})
		t.Logf("Part 2 uploaded, ETag: %s", *output.ETag)
	})

	t.Run("list parts", func(t *testing.T) {
		output, err := tc.Client.ListParts(tc.Context, &s3.ListPartsInput{
			Bucket:   aws.String(tc.BucketName),
			Key:      aws.String(testKey),
			UploadId: aws.String(uploadID),
		})
		if err != nil {
			t.Errorf("ListParts() error = %v", err)
			return
		}

		if len(output.Parts) != 2 {
			t.Errorf("Expected 2 parts, got %d", len(output.Parts))
		}

		for _, part := range output.Parts {
			t.Logf("Part %d: Size=%d, ETag=%s", *part.PartNumber, part.Size, *part.ETag)
		}
	})

	t.Run("complete multipart upload", func(t *testing.T) {
		output, err := tc.Client.CompleteMultipartUpload(tc.Context, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(tc.BucketName),
			Key:      aws.String(testKey),
			UploadId: aws.String(uploadID),
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: completedParts,
			},
		})
		if err != nil {
			t.Errorf("CompleteMultipartUpload() error = %v", err)
			return
		}

		t.Logf("Multipart upload completed, Location: %s", *output.Location)
	})

	t.Run("verify completed object", func(t *testing.T) {
		output, err := tc.Client.GetObject(tc.Context, &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("GetObject() error = %v", err)
			return
		}
		defer output.Body.Close()

		body, _ := io.ReadAll(output.Body)
		expectedSize := partSize * 2
		if len(body) != expectedSize {
			t.Errorf("Expected size %d, got %d", expectedSize, len(body))
		}

		t.Logf("Object size: %d bytes", len(body))
	})
}

func TestMultipartUploadWithMetadata(t *testing.T) {
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

	testKey := "multipart-metadata.txt"
	metadata := map[string]string{
		"author":  "s3tool-test",
		"project": "multipart-test",
	}

	t.Run("initiate with metadata", func(t *testing.T) {
		output, err := tc.Client.CreateMultipartUpload(tc.Context, &s3.CreateMultipartUploadInput{
			Bucket:   aws.String(tc.BucketName),
			Key:      aws.String(testKey),
			Metadata: metadata,
		})
		if err != nil {
			t.Errorf("CreateMultipartUpload() error = %v", err)
			return
		}

		uploadID := *output.UploadId
		t.Logf("UploadId: %s", uploadID)

		partContent := []byte("test part content")
		partOutput, err := tc.Client.UploadPart(tc.Context, &s3.UploadPartInput{
			Bucket:     aws.String(tc.BucketName),
			Key:        aws.String(testKey),
			UploadId:   aws.String(uploadID),
			PartNumber: aws.Int32(1),
			Body:       bytes.NewReader(partContent),
		})
		if err != nil {
			t.Errorf("UploadPart() error = %v", err)
			return
		}

		_, err = tc.Client.CompleteMultipartUpload(tc.Context, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(tc.BucketName),
			Key:      aws.String(testKey),
			UploadId: aws.String(uploadID),
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: []types.CompletedPart{
					{
						ETag:       partOutput.ETag,
						PartNumber: aws.Int32(1),
					},
				},
			},
		})
		if err != nil {
			t.Errorf("CompleteMultipartUpload() error = %v", err)
			return
		}
	})

	t.Run("verify metadata", func(t *testing.T) {
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

func TestMultipartUploadWithContentType(t *testing.T) {
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

	testKey := "multipart-content-type.bin"
	contentType := "application/octet-stream"

	output, err := tc.Client.CreateMultipartUpload(tc.Context, &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(tc.BucketName),
		Key:         aws.String(testKey),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		t.Fatalf("CreateMultipartUpload() error = %v", err)
	}

	uploadID := *output.UploadId
	t.Logf("UploadId: %s", uploadID)

	partContent := []byte("test content")
	partOutput, err := tc.Client.UploadPart(tc.Context, &s3.UploadPartInput{
		Bucket:     aws.String(tc.BucketName),
		Key:        aws.String(testKey),
		UploadId:   aws.String(uploadID),
		PartNumber: aws.Int32(1),
		Body:       bytes.NewReader(partContent),
	})
	if err != nil {
		t.Fatalf("UploadPart() error = %v", err)
	}

	_, err = tc.Client.CompleteMultipartUpload(tc.Context, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(tc.BucketName),
		Key:      aws.String(testKey),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: []types.CompletedPart{
				{
					ETag:       partOutput.ETag,
					PartNumber: aws.Int32(1),
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("CompleteMultipartUpload() error = %v", err)
	}

	headOutput, err := tc.Client.HeadObject(tc.Context, &s3.HeadObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
	})
	if err != nil {
		t.Fatalf("HeadObject() error = %v", err)
	}

	if headOutput.ContentType != nil && *headOutput.ContentType != contentType {
		t.Errorf("Expected content type %s, got %s", contentType, *headOutput.ContentType)
	}

	t.Logf("Content-Type: %s", *headOutput.ContentType)
}

func TestMultipartUploadWrongPartOrder(t *testing.T) {
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

	testKey := "multipart-wrong-order.txt"

	output, err := tc.Client.CreateMultipartUpload(tc.Context, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
	})
	if err != nil {
		t.Fatalf("CreateMultipartUpload() error = %v", err)
	}

	uploadID := *output.UploadId
	defer tc.Client.AbortMultipartUpload(tc.Context, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(tc.BucketName),
		Key:      aws.String(testKey),
		UploadId: aws.String(uploadID),
	})

	partContent := []byte("test content")
	part1Output, _ := tc.Client.UploadPart(tc.Context, &s3.UploadPartInput{
		Bucket:     aws.String(tc.BucketName),
		Key:        aws.String(testKey),
		UploadId:   aws.String(uploadID),
		PartNumber: aws.Int32(1),
		Body:       bytes.NewReader(partContent),
	})

	part2Output, _ := tc.Client.UploadPart(tc.Context, &s3.UploadPartInput{
		Bucket:     aws.String(tc.BucketName),
		Key:        aws.String(testKey),
		UploadId:   aws.String(uploadID),
		PartNumber: aws.Int32(2),
		Body:       bytes.NewReader(partContent),
	})

	_, err = tc.Client.CompleteMultipartUpload(tc.Context, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(tc.BucketName),
		Key:      aws.String(testKey),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: []types.CompletedPart{
				{
					ETag:       part2Output.ETag,
					PartNumber: aws.Int32(2),
				},
				{
					ETag:       part1Output.ETag,
					PartNumber: aws.Int32(1),
				},
			},
		},
	})
	if err != nil {
		t.Logf("Complete with wrong order returned error: %v", err)
	} else {
		t.Log("Complete with wrong order succeeded (S3 may reorder parts)")
	}
}

func TestMultipartUploadInvalidUploadId(t *testing.T) {
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

	testKey := "multipart-invalid.txt"

	t.Run("upload part with invalid upload ID", func(t *testing.T) {
		_, err := tc.Client.UploadPart(tc.Context, &s3.UploadPartInput{
			Bucket:     aws.String(tc.BucketName),
			Key:        aws.String(testKey),
			UploadId:   aws.String("invalid-upload-id"),
			PartNumber: aws.Int32(1),
			Body:       bytes.NewReader([]byte("test")),
		})
		if err == nil {
			t.Error("Expected error for invalid upload ID, got nil")
		} else {
			t.Logf("Got expected error: %v", err)
		}
	})

	t.Run("complete with invalid upload ID", func(t *testing.T) {
		_, err := tc.Client.CompleteMultipartUpload(tc.Context, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(tc.BucketName),
			Key:      aws.String(testKey),
			UploadId: aws.String("invalid-upload-id"),
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: []types.CompletedPart{},
			},
		})
		if err == nil {
			t.Error("Expected error for invalid upload ID, got nil")
		} else {
			t.Logf("Got expected error: %v", err)
		}
	})

	t.Run("abort with invalid upload ID", func(t *testing.T) {
		_, err := tc.Client.AbortMultipartUpload(tc.Context, &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(tc.BucketName),
			Key:      aws.String(testKey),
			UploadId: aws.String("invalid-upload-id"),
		})
		if err == nil {
			t.Error("Expected error for invalid upload ID, got nil")
		} else {
			t.Logf("Got expected error: %v", err)
		}
	})
}

func TestMultipartUploadMultipleParts(t *testing.T) {
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

	testKey := "multipart-multiple.txt"
	numParts := 5
	partSize := 5 * 1024 * 1024

	output, err := tc.Client.CreateMultipartUpload(tc.Context, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
	})
	if err != nil {
		t.Fatalf("CreateMultipartUpload() error = %v", err)
	}

	uploadID := *output.UploadId
	t.Logf("UploadId: %s", uploadID)

	var completedParts []types.CompletedPart

	for i := 1; i <= numParts; i++ {
		partContent := make([]byte, partSize)
		for j := range partContent {
			partContent[j] = byte((i*50 + j) % 256)
		}

		partOutput, err := tc.Client.UploadPart(tc.Context, &s3.UploadPartInput{
			Bucket:     aws.String(tc.BucketName),
			Key:        aws.String(testKey),
			UploadId:   aws.String(uploadID),
			PartNumber: aws.Int32(int32(i)),
			Body:       bytes.NewReader(partContent),
		})
		if err != nil {
			t.Errorf("UploadPart(%d) error = %v", i, err)
			tc.Client.AbortMultipartUpload(tc.Context, &s3.AbortMultipartUploadInput{
				Bucket:   aws.String(tc.BucketName),
				Key:      aws.String(testKey),
				UploadId: aws.String(uploadID),
			})
			return
		}

		completedParts = append(completedParts, types.CompletedPart{
			ETag:       partOutput.ETag,
			PartNumber: aws.Int32(int32(i)),
		})
		t.Logf("Part %d uploaded", i)
	}

	_, err = tc.Client.CompleteMultipartUpload(tc.Context, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(tc.BucketName),
		Key:      aws.String(testKey),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		t.Errorf("CompleteMultipartUpload() error = %v", err)
		return
	}

	headOutput, err := tc.Client.HeadObject(tc.Context, &s3.HeadObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
	})
	if err != nil {
		t.Errorf("HeadObject() error = %v", err)
		return
	}

	expectedSize := int64(numParts * partSize)
	if *headOutput.ContentLength != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, *headOutput.ContentLength)
	}

	t.Logf("Multi-part upload completed: %d bytes", *headOutput.ContentLength)
}

func TestMultipartUploadCopyPart(t *testing.T) {
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

	sourceKey := "source-object.txt"
	destKey := "multipart-copy.txt"
	sourceContent := "This is the source object content for copy part"

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(sourceKey),
		Body:   bytes.NewReader([]byte(sourceContent)),
	})
	if err != nil {
		t.Fatalf("PutObject() error = %v", err)
	}

	t.Run("initiate multipart upload", func(t *testing.T) {
		output, err := tc.Client.CreateMultipartUpload(tc.Context, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(destKey),
		})
		if err != nil {
			t.Errorf("CreateMultipartUpload() error = %v", err)
			return
		}
		uploadID := *output.UploadId
		t.Logf("UploadId: %s", uploadID)

		copySource := tc.BucketName + "/" + sourceKey
		partOutput, err := tc.Client.UploadPartCopy(tc.Context, &s3.UploadPartCopyInput{
			Bucket:     aws.String(tc.BucketName),
			Key:        aws.String(destKey),
			UploadId:   aws.String(uploadID),
			PartNumber: aws.Int32(1),
			CopySource: aws.String(copySource),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("UploadPartCopy not supported by this S3 implementation")
			}
			t.Errorf("UploadPartCopy() error = %v", err)
			return
		}

		t.Logf("Part copied, ETag: %s", *partOutput.CopyPartResult.ETag)

		_, err = tc.Client.CompleteMultipartUpload(tc.Context, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(tc.BucketName),
			Key:      aws.String(destKey),
			UploadId: aws.String(uploadID),
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: []types.CompletedPart{
					{
						ETag:       partOutput.CopyPartResult.ETag,
						PartNumber: aws.Int32(1),
					},
				},
			},
		})
		if err != nil {
			t.Errorf("CompleteMultipartUpload() error = %v", err)
			return
		}
	})

	t.Run("verify copied content", func(t *testing.T) {
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
		if string(body) != sourceContent {
			t.Errorf("Content mismatch")
		}
		t.Log("Content verified successfully")
	})
}

func TestListMultipartUploads(t *testing.T) {
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

	var uploadIDs []string

	for i := 0; i < 3; i++ {
		output, err := tc.Client.CreateMultipartUpload(tc.Context, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(fmt.Sprintf("multipart-list-%d.txt", i)),
		})
		if err != nil {
			t.Fatalf("CreateMultipartUpload() error = %v", err)
		}
		uploadIDs = append(uploadIDs, *output.UploadId)
	}

	t.Run("list all multipart uploads", func(t *testing.T) {
		output, err := tc.Client.ListMultipartUploads(tc.Context, &s3.ListMultipartUploadsInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			t.Errorf("ListMultipartUploads() error = %v", err)
			return
		}

		if len(output.Uploads) < 3 {
			t.Errorf("Expected at least 3 uploads, got %d", len(output.Uploads))
		}

		for _, upload := range output.Uploads {
			t.Logf("Upload: Key=%s, UploadId=%s", *upload.Key, *upload.UploadId)
		}
	})

	for i, uploadID := range uploadIDs {
		_, err := tc.Client.AbortMultipartUpload(tc.Context, &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(tc.BucketName),
			Key:      aws.String(fmt.Sprintf("multipart-list-%d.txt", i)),
			UploadId: aws.String(uploadID),
		})
		if err != nil {
			t.Logf("Warning: failed to abort upload %s: %v", uploadID, err)
		}
	}
}
