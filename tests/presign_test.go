package tests

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestPresignGetObject(t *testing.T) {
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

	testKey := "presign-test.txt"
	testContent := "Content for presigned URL test"

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
		Body:   bytes.NewReader([]byte(testContent)),
	})
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	t.Run("generate presigned GET URL", func(t *testing.T) {
		presignClient := s3.NewPresignClient(tc.Client)

		presignedURL, err := presignClient.PresignGetObject(tc.Context, &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			t.Errorf("PresignGetObject() error = %v", err)
			return
		}

		t.Logf("Presigned GET URL: %s", presignedURL.URL)
		t.Logf("Expires at: %s", presignedURL.SignedHeader.Get("X-Amz-Expires"))

		if presignedURL.URL == "" {
			t.Error("Presigned URL is empty")
		}

		if !strings.Contains(presignedURL.URL, tc.BucketName) {
			t.Error("Presigned URL does not contain bucket name")
		}

		if !strings.Contains(presignedURL.URL, testKey) {
			t.Error("Presigned URL does not contain object key")
		}
	})

	t.Run("use presigned GET URL", func(t *testing.T) {
		presignClient := s3.NewPresignClient(tc.Client)

		presignedURL, err := presignClient.PresignGetObject(tc.Context, &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			t.Errorf("PresignGetObject() error = %v", err)
			return
		}

		resp, err := http.Get(presignedURL.URL)
		if err != nil {
			t.Errorf("HTTP GET error = %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Failed to read response body: %v", err)
			return
		}

		if string(body) != testContent {
			t.Errorf("Expected content %q, got %q", testContent, string(body))
		}

		t.Log("Presigned GET URL works correctly")
	})
}

func TestPresignPutObject(t *testing.T) {
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

	testKey := "presign-upload.txt"
	uploadContent := "Content uploaded via presigned URL"

	t.Run("generate presigned PUT URL", func(t *testing.T) {
		presignClient := s3.NewPresignClient(tc.Client)

		presignedURL, err := presignClient.PresignPutObject(tc.Context, &s3.PutObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			t.Errorf("PresignPutObject() error = %v", err)
			return
		}

		t.Logf("Presigned PUT URL: %s", presignedURL.URL)

		if presignedURL.URL == "" {
			t.Error("Presigned URL is empty")
		}
	})

	t.Run("use presigned PUT URL", func(t *testing.T) {
		presignClient := s3.NewPresignClient(tc.Client)

		presignedURL, err := presignClient.PresignPutObject(tc.Context, &s3.PutObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			t.Errorf("PresignPutObject() error = %v", err)
			return
		}

		req, err := http.NewRequest("PUT", presignedURL.URL, bytes.NewReader([]byte(uploadContent)))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "text/plain")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("HTTP PUT error = %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
			return
		}

		t.Log("Presigned PUT URL works correctly")
	})

	t.Run("verify uploaded content", func(t *testing.T) {
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
		if string(body) != uploadContent {
			t.Errorf("Expected content %q, got %q", uploadContent, string(body))
		}

		t.Log("Uploaded content verified")
	})
}

func TestPresignDeleteObject(t *testing.T) {
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

	testKey := "presign-delete.txt"

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
		Body:   bytes.NewReader([]byte("test content")),
	})
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	t.Run("generate presigned DELETE URL", func(t *testing.T) {
		presignClient := s3.NewPresignClient(tc.Client)

		presignedURL, err := presignClient.PresignDeleteObject(tc.Context, &s3.DeleteObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			t.Errorf("PresignDeleteObject() error = %v", err)
			return
		}

		t.Logf("Presigned DELETE URL: %s", presignedURL.URL)

		if presignedURL.URL == "" {
			t.Error("Presigned URL is empty")
		}
	})

	t.Run("use presigned DELETE URL", func(t *testing.T) {
		presignClient := s3.NewPresignClient(tc.Client)

		presignedURL, err := presignClient.PresignDeleteObject(tc.Context, &s3.DeleteObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			t.Errorf("PresignDeleteObject() error = %v", err)
			return
		}

		req, err := http.NewRequest("DELETE", presignedURL.URL, nil)
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("HTTP DELETE error = %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 204 or 200, got %d", resp.StatusCode)
			return
		}

		t.Log("Presigned DELETE URL works correctly")
	})

	t.Run("verify object deleted", func(t *testing.T) {
		_, err := tc.Client.HeadObject(tc.Context, &s3.HeadObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err == nil {
			t.Error("Object should have been deleted")
		} else {
			t.Logf("Object deleted successfully: %v", err)
		}
	})
}

func TestPresignURLExpiration(t *testing.T) {
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

	testKey := "expiration-test.txt"
	testContent := "test content"

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
		Body:   bytes.NewReader([]byte(testContent)),
	})
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	t.Run("verify expiration in URL", func(t *testing.T) {
		presignClient := s3.NewPresignClient(tc.Client)

		expiration := 3600 * time.Second
		presignedURL, err := presignClient.PresignGetObject(tc.Context, &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = expiration
		})
		if err != nil {
			t.Errorf("PresignGetObject() error = %v", err)
			return
		}

		parsedURL, err := url.Parse(presignedURL.URL)
		if err != nil {
			t.Errorf("Failed to parse URL: %v", err)
			return
		}

		expiresParam := parsedURL.Query().Get("X-Amz-Expires")
		if expiresParam == "" {
			t.Error("X-Amz-Expires parameter not found in URL")
		} else {
			t.Logf("Expiration parameter: %s seconds", expiresParam)
		}
	})

	t.Run("short expiration URL", func(t *testing.T) {
		presignClient := s3.NewPresignClient(tc.Client)

		presignedURL, err := presignClient.PresignGetObject(tc.Context, &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 60 * time.Second
		})
		if err != nil {
			t.Errorf("PresignGetObject() error = %v", err)
			return
		}

		t.Logf("Short expiration URL: %s", presignedURL.URL)
	})
}

func TestPresignNonExistentObject(t *testing.T) {
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

	t.Run("presign URL for non-existent object", func(t *testing.T) {
		presignClient := s3.NewPresignClient(tc.Client)

		presignedURL, err := presignClient.PresignGetObject(tc.Context, &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String("non-existent-file.txt"),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			t.Errorf("PresignGetObject() error = %v", err)
			return
		}

		resp, err := http.Get(presignedURL.URL)
		if err != nil {
			t.Errorf("HTTP GET error = %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Error("Expected error for non-existent object, got 200 OK")
		} else {
			t.Logf("Got expected error status: %d", resp.StatusCode)
		}
	})
}

func TestPresignWithSpecialKeyNames(t *testing.T) {
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
		{"with_spaces", "path with spaces/file.txt", "content"},
		{"with_unicode", "unicode/中文.txt", "中文内容"},
		{"with_special", "special/file-name_test.txt", "content"},
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

			presignClient := s3.NewPresignClient(tc.Client)

			presignedURL, err := presignClient.PresignGetObject(tc.Context, &s3.GetObjectInput{
				Bucket: aws.String(tc.BucketName),
				Key:    aws.String(testCase.key),
			}, func(opts *s3.PresignOptions) {
				opts.Expires = 15 * time.Minute
			})
			if err != nil {
				t.Errorf("PresignGetObject() error = %v", err)
				return
			}

			t.Logf("Presigned URL for key %q: %s", testCase.key, presignedURL.URL)

			resp, err := http.Get(presignedURL.URL)
			if err != nil {
				t.Errorf("HTTP GET error = %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
				return
			}

			body, _ := io.ReadAll(resp.Body)
			if string(body) != testCase.content {
				t.Errorf("Expected content %q, got %q", testCase.content, string(body))
			}

			t.Log("Presigned URL with special key works correctly")
		})
	}
}

func TestPresignMultipleOperations(t *testing.T) {
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

	testKey := "multi-op-test.txt"

	t.Run("generate multiple presigned URLs", func(t *testing.T) {
		presignClient := s3.NewPresignClient(tc.Client)

		putURL, err := presignClient.PresignPutObject(tc.Context, &s3.PutObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("PresignPutObject() error = %v", err)
			return
		}

		getURL, err := presignClient.PresignGetObject(tc.Context, &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("PresignGetObject() error = %v", err)
			return
		}

		deleteURL, err := presignClient.PresignDeleteObject(tc.Context, &s3.DeleteObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("PresignDeleteObject() error = %v", err)
			return
		}

		t.Logf("PUT URL: %s", putURL.URL)
		t.Logf("GET URL: %s", getURL.URL)
		t.Logf("DELETE URL: %s", deleteURL.URL)

		content := "Test content for multi-operation"
		req, _ := http.NewRequest("PUT", putURL.URL, bytes.NewReader([]byte(content)))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("PUT request error = %v", err)
			return
		}
		resp.Body.Close()
		t.Log("PUT via presigned URL successful")

		resp, err = http.Get(getURL.URL)
		if err != nil {
			t.Errorf("GET request error = %v", err)
			return
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if string(body) != content {
			t.Errorf("GET content mismatch")
		}
		t.Log("GET via presigned URL successful")

		req, _ = http.NewRequest("DELETE", deleteURL.URL, nil)
		resp, err = client.Do(req)
		if err != nil {
			t.Errorf("DELETE request error = %v", err)
			return
		}
		resp.Body.Close()
		t.Log("DELETE via presigned URL successful")
	})
}

func TestPresignWithContext(t *testing.T) {
	ctx := context.Background()
	client, err := config.GetS3Client(ctx)
	if err != nil {
		t.Fatalf("Failed to create S3 client: %v", err)
	}

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

	testKey := "context-test.txt"
	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
		Body:   bytes.NewReader([]byte("test")),
	})
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	t.Run("presign with context", func(t *testing.T) {
		presignClient := s3.NewPresignClient(client)

		presignedURL, err := presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("PresignGetObject() error = %v", err)
			return
		}

		if presignedURL.URL == "" {
			t.Error("Presigned URL is empty")
		}

		t.Logf("Presigned URL with context: %s", presignedURL.URL)
	})
}

func TestPresignURLFormat(t *testing.T) {
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

	testKey := "format-test.txt"
	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
		Body:   bytes.NewReader([]byte("test")),
	})
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	presignClient := s3.NewPresignClient(tc.Client)

	presignedURL, err := presignClient.PresignGetObject(tc.Context, &s3.GetObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(testKey),
	})
	if err != nil {
		t.Fatalf("PresignGetObject() error = %v", err)
	}

	t.Run("verify URL structure", func(t *testing.T) {
		parsedURL, err := url.Parse(presignedURL.URL)
		if err != nil {
			t.Errorf("Failed to parse URL: %v", err)
			return
		}

		t.Logf("Scheme: %s", parsedURL.Scheme)
		t.Logf("Host: %s", parsedURL.Host)
		t.Logf("Path: %s", parsedURL.Path)

		query := parsedURL.Query()
		requiredParams := []string{
			"X-Amz-Algorithm",
			"X-Amz-Credential",
			"X-Amz-Date",
			"X-Amz-Expires",
			"X-Amz-SignedHeaders",
			"X-Amz-Signature",
		}

		for _, param := range requiredParams {
			if query.Get(param) == "" {
				t.Errorf("Missing required parameter: %s", param)
			} else {
				t.Logf("Found parameter: %s", param)
			}
		}
	})
}

func TestPresignWithContentType(t *testing.T) {
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

	testKey := "content-type-test.txt"
	contentType := "application/json"

	t.Run("presign PUT with content type", func(t *testing.T) {
		presignClient := s3.NewPresignClient(tc.Client)

		presignedURL, err := presignClient.PresignPutObject(tc.Context, &s3.PutObjectInput{
			Bucket:      aws.String(tc.BucketName),
			Key:         aws.String(testKey),
			ContentType: aws.String(contentType),
		})
		if err != nil {
			t.Errorf("PresignPutObject() error = %v", err)
			return
		}

		t.Logf("Presigned PUT URL with content type: %s", presignedURL.URL)

		req, err := http.NewRequest("PUT", presignedURL.URL, bytes.NewReader([]byte(`{"test": true}`)))
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
			return
		}
		req.Header.Set("Content-Type", contentType)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("HTTP PUT error = %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("verify content type stored", func(t *testing.T) {
		output, err := tc.Client.HeadObject(tc.Context, &s3.HeadObjectInput{
			Bucket: aws.String(tc.BucketName),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Errorf("HeadObject() error = %v", err)
			return
		}

		if output.ContentType != nil && *output.ContentType != contentType {
			t.Errorf("Expected content type %s, got %s", contentType, *output.ContentType)
		}

		t.Logf("Stored content type: %s", *output.ContentType)
	})
}
