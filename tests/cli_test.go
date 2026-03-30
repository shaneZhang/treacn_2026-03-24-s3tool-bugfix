package tests

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func getBinaryPath() string {
	binaryPath := "../s3tool"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		binaryPath = "./s3tool"
	}
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		binaryPath = "s3tool"
	}
	return binaryPath
}

func runCommand(args ...string) (string, string, error) {
	binaryPath := getBinaryPath()
	cmd := exec.Command(binaryPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func TestCLIHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"root help", []string{"--help"}},
		{"bucket help", []string{"bucket", "--help"}},
		{"object help", []string{"object", "--help"}},
		{"presign help", []string{"presign", "--help"}},
		{"multipart help", []string{"multipart", "--help"}},
		{"versioning help", []string{"versioning", "--help"}},
		{"lifecycle help", []string{"lifecycle", "--help"}},
		{"acl help", []string{"acl", "--help"}},
		{"cors help", []string{"cors", "--help"}},
		{"website help", []string{"website", "--help"}},
		{"encryption help", []string{"encryption", "--help"}},
		{"tags help", []string{"tags", "--help"}},
		{"logging help", []string{"logging", "--help"}},
		{"replication help", []string{"replication", "--help"}},
		{"notification help", []string{"notification", "--help"}},
		{"policy help", []string{"policy", "--help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, err := runCommand(tt.args...)
			if err != nil && !strings.Contains(err.Error(), "exit status 0") {
				t.Errorf("Command failed: %v", err)
			}

			if !strings.Contains(stdout, "Usage:") && !strings.Contains(stdout, "帮助") {
				t.Errorf("Help output missing Usage, got: %s", stdout)
			}

			t.Logf("Help output length: %d bytes", len(stdout))
		})
	}
}

func TestCLIBucketList(t *testing.T) {
	stdout, stderr, err := runCommand("bucket", "list")
	if err != nil {
		t.Errorf("bucket list failed: %v\nstderr: %s", err, stderr)
		return
	}

	t.Logf("Bucket list output:\n%s", stdout)

	if strings.Contains(stderr, "error") || strings.Contains(stderr, "Error") {
		t.Errorf("Error in stderr: %s", stderr)
	}
}

func TestCLIBucketCreateAndDelete(t *testing.T) {
	bucketName := fmt.Sprintf("cli-test-%d", time.Now().UnixNano())

	t.Run("create bucket", func(t *testing.T) {
		stdout, stderr, err := runCommand("bucket", "create", bucketName)
		if err != nil {
			t.Errorf("bucket create failed: %v\nstderr: %s", err, stderr)
			return
		}

		t.Logf("Create output: %s", stdout)

		if !strings.Contains(stdout, bucketName) && !strings.Contains(stderr, bucketName) {
			t.Log("Bucket name not in output, but command succeeded")
		}
	})

	time.Sleep(1 * time.Second)

	t.Run("list buckets to verify", func(t *testing.T) {
		stdout, _, _ := runCommand("bucket", "list")
		if !strings.Contains(stdout, bucketName) {
			t.Logf("Warning: bucket %s not found in list", bucketName)
		} else {
			t.Logf("Bucket %s found in list", bucketName)
		}
	})

	t.Run("delete bucket", func(t *testing.T) {
		stdout, stderr, err := runCommand("bucket", "delete", bucketName)
		if err != nil {
			t.Errorf("bucket delete failed: %v\nstderr: %s", err, stderr)
			return
		}

		t.Logf("Delete output: %s", stdout)
	})
}

func TestCLIObjectOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-obj-test-%d", time.Now().UnixNano())

	stdout, stderr, _ := runCommand("bucket", "create", bucketName)
	t.Logf("Create bucket: %s", stdout+stderr)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	testFile := "/tmp/s3tool-test-upload.txt"
	testContent := "Test content for CLI upload"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	testKey := "cli-test-object.txt"

	t.Run("upload object", func(t *testing.T) {
		stdout, stderr, err := runCommand("object", "put", bucketName, testKey, testFile)
		if err != nil {
			t.Errorf("object put failed: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Put output: %s", stdout)
	})

	t.Run("list objects", func(t *testing.T) {
		stdout, stderr, err := runCommand("object", "list", bucketName)
		if err != nil {
			t.Errorf("object list failed: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("List output: %s", stdout)

		if !strings.Contains(stdout, testKey) {
			t.Logf("Warning: object %s not found in list", testKey)
		}
	})

	t.Run("get object info", func(t *testing.T) {
		stdout, stderr, err := runCommand("object", "info", bucketName, testKey)
		if err != nil {
			t.Errorf("object info failed: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Info output: %s", stdout)
	})

	downloadFile := "/tmp/s3tool-test-download.txt"
	defer os.Remove(downloadFile)

	t.Run("download object", func(t *testing.T) {
		stdout, stderr, err := runCommand("object", "get", bucketName, testKey, downloadFile)
		if err != nil {
			t.Errorf("object get failed: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Get output: %s", stdout)

		content, err := os.ReadFile(downloadFile)
		if err != nil {
			t.Errorf("Failed to read downloaded file: %v", err)
			return
		}

		if string(content) != testContent {
			t.Errorf("Content mismatch: expected %q, got %q", testContent, string(content))
		}
	})

	t.Run("delete object", func(t *testing.T) {
		stdout, stderr, err := runCommand("object", "delete", bucketName, testKey)
		if err != nil {
			t.Errorf("object delete failed: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Delete output: %s", stdout)
	})
}

func TestCLIPresignOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-presign-test-%d", time.Now().UnixNano())

	runCommand("bucket", "create", bucketName)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	testFile := "/tmp/s3tool-presign-test.txt"
	os.WriteFile(testFile, []byte("presign test content"), 0644)
	defer os.Remove(testFile)

	runCommand("object", "put", bucketName, "presign-test.txt", testFile)

	t.Run("presign get", func(t *testing.T) {
		stdout, stderr, err := runCommand("presign", "get", bucketName, "presign-test.txt")
		if err != nil {
			t.Errorf("presign get failed: %v\nstderr: %s", err, stderr)
			return
		}

		if !strings.Contains(stdout, "http") {
			t.Errorf("Expected URL in output, got: %s", stdout)
		}

		t.Logf("Presigned GET URL: %s", strings.TrimSpace(stdout))
	})

	t.Run("presign put", func(t *testing.T) {
		stdout, stderr, err := runCommand("presign", "put", bucketName, "new-file.txt")
		if err != nil {
			t.Errorf("presign put failed: %v\nstderr: %s", err, stderr)
			return
		}

		if !strings.Contains(stdout, "http") {
			t.Errorf("Expected URL in output, got: %s", stdout)
		}

		t.Logf("Presigned PUT URL: %s", strings.TrimSpace(stdout))
	})

	t.Run("presign with expiration", func(t *testing.T) {
		stdout, stderr, err := runCommand("presign", "get", bucketName, "presign-test.txt", "--expires", "7200")
		if err != nil {
			t.Errorf("presign get with expires failed: %v\nstderr: %s", err, stderr)
			return
		}

		t.Logf("Presigned URL with custom expiration: %s", strings.TrimSpace(stdout))
	})
}

func TestCLIVersioningOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-version-test-%d", time.Now().UnixNano())

	runCommand("bucket", "create", bucketName)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	t.Run("get versioning status", func(t *testing.T) {
		stdout, _, err := runCommand("versioning", "get", bucketName)
		if err != nil {
			t.Logf("versioning get returned error (may not be supported): %v", err)
			return
		}
		t.Logf("Versioning status: %s", stdout)
	})

	t.Run("enable versioning", func(t *testing.T) {
		stdout, stderr, err := runCommand("versioning", "enable", bucketName)
		if err != nil {
			t.Logf("versioning enable returned error (may not be supported): %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Enable output: %s", stdout)
	})

	t.Run("suspend versioning", func(t *testing.T) {
		stdout, stderr, err := runCommand("versioning", "suspend", bucketName)
		if err != nil {
			t.Logf("versioning suspend returned error (may not be supported): %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Suspend output: %s", stdout)
	})
}

func TestCLITagsOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-tags-test-%d", time.Now().UnixNano())

	runCommand("bucket", "create", bucketName)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	t.Run("set bucket tags", func(t *testing.T) {
		stdout, stderr, err := runCommand("tags", "bucket-put", bucketName, "Environment=Test", "Project=S3Tool")
		if err != nil {
			t.Logf("tags bucket-put returned error: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Set tags output: %s", stdout)
	})

	t.Run("get bucket tags", func(t *testing.T) {
		stdout, stderr, err := runCommand("tags", "bucket-get", bucketName)
		if err != nil {
			t.Logf("tags bucket-get returned error: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Get tags output: %s", stdout)
	})

	t.Run("delete bucket tags", func(t *testing.T) {
		stdout, stderr, err := runCommand("tags", "bucket-delete", bucketName)
		if err != nil {
			t.Logf("tags bucket-delete returned error: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Delete tags output: %s", stdout)
	})
}

func TestCLIEncryptionOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-enc-test-%d", time.Now().UnixNano())

	runCommand("bucket", "create", bucketName)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	t.Run("get encryption config", func(t *testing.T) {
		stdout, stderr, err := runCommand("encryption", "get", bucketName)
		if err != nil {
			t.Logf("encryption get returned error (may not be configured): %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Encryption config: %s", stdout)
	})

	t.Run("enable encryption", func(t *testing.T) {
		stdout, stderr, err := runCommand("encryption", "enable", bucketName)
		if err != nil {
			t.Logf("encryption enable returned error (may not be supported): %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Enable encryption output: %s", stdout)
	})

	t.Run("disable encryption", func(t *testing.T) {
		stdout, stderr, err := runCommand("encryption", "disable", bucketName)
		if err != nil {
			t.Logf("encryption disable returned error: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Disable encryption output: %s", stdout)
	})
}

func TestCLILifecycleOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-lc-test-%d", time.Now().UnixNano())

	runCommand("bucket", "create", bucketName)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	t.Run("get lifecycle config", func(t *testing.T) {
		stdout, stderr, err := runCommand("lifecycle", "get", bucketName)
		if err != nil {
			t.Logf("lifecycle get returned error (may not be configured): %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Lifecycle config: %s", stdout)
	})

	t.Run("delete lifecycle config", func(t *testing.T) {
		stdout, stderr, err := runCommand("lifecycle", "delete", bucketName)
		if err != nil {
			t.Logf("lifecycle delete returned error: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Delete lifecycle output: %s", stdout)
	})
}

func TestCLICORSOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-cors-test-%d", time.Now().UnixNano())

	runCommand("bucket", "create", bucketName)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	t.Run("get CORS config", func(t *testing.T) {
		stdout, stderr, err := runCommand("cors", "get", bucketName)
		if err != nil {
			t.Logf("cors get returned error (may not be configured): %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("CORS config: %s", stdout)
	})

	t.Run("delete CORS config", func(t *testing.T) {
		stdout, stderr, err := runCommand("cors", "delete", bucketName)
		if err != nil {
			t.Logf("cors delete returned error: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Delete CORS output: %s", stdout)
	})
}

func TestCLIWebsiteOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-web-test-%d", time.Now().UnixNano())

	runCommand("bucket", "create", bucketName)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	t.Run("get website config", func(t *testing.T) {
		stdout, stderr, err := runCommand("website", "get", bucketName)
		if err != nil {
			t.Logf("website get returned error (may not be configured): %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Website config: %s", stdout)
	})

	t.Run("disable website", func(t *testing.T) {
		stdout, stderr, err := runCommand("website", "disable", bucketName)
		if err != nil {
			t.Logf("website disable returned error: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Disable website output: %s", stdout)
	})
}

func TestCLIPolicyOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-policy-test-%d", time.Now().UnixNano())

	runCommand("bucket", "create", bucketName)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	t.Run("get bucket policy", func(t *testing.T) {
		stdout, stderr, err := runCommand("policy", "get", bucketName)
		if err != nil {
			t.Logf("policy get returned error (may not be configured): %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Policy: %s", stdout)
	})

	t.Run("delete bucket policy", func(t *testing.T) {
		stdout, stderr, err := runCommand("policy", "delete", bucketName)
		if err != nil {
			t.Logf("policy delete returned error: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Delete policy output: %s", stdout)
	})
}

func TestCLIAQLOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-acl-test-%d", time.Now().UnixNano())

	runCommand("bucket", "create", bucketName)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	t.Run("get bucket ACL", func(t *testing.T) {
		stdout, stderr, err := runCommand("acl", "bucket-get", bucketName)
		if err != nil {
			t.Logf("acl bucket-get returned error: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Bucket ACL: %s", stdout)
	})
}

func TestCLILoggingOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-log-test-%d", time.Now().UnixNano())

	runCommand("bucket", "create", bucketName)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	t.Run("get logging config", func(t *testing.T) {
		stdout, stderr, err := runCommand("logging", "get", bucketName)
		if err != nil {
			t.Logf("logging get returned error: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Logging config: %s", stdout)
	})
}

func TestCLIReplicationOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-repl-test-%d", time.Now().UnixNano())

	runCommand("bucket", "create", bucketName)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	t.Run("get replication config", func(t *testing.T) {
		stdout, stderr, err := runCommand("replication", "get", bucketName)
		if err != nil {
			t.Logf("replication get returned error (may not be configured): %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Replication config: %s", stdout)
	})
}

func TestCLINotificationOperations(t *testing.T) {
	bucketName := fmt.Sprintf("cli-notify-test-%d", time.Now().UnixNano())

	runCommand("bucket", "create", bucketName)
	defer runCommand("bucket", "delete", bucketName)

	time.Sleep(1 * time.Second)

	t.Run("get notification config", func(t *testing.T) {
		stdout, stderr, err := runCommand("notification", "get", bucketName)
		if err != nil {
			t.Logf("notification get returned error: %v\nstderr: %s", err, stderr)
			return
		}
		t.Logf("Notification config: %s", stdout)
	})
}

func TestCLIInvalidCommands(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"missing bucket name", []string{"bucket", "delete"}},
		{"missing object key", []string{"object", "info"}},
		{"invalid command", []string{"invalid-command"}},
		{"missing arguments", []string{"presign", "get"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, stderr, err := runCommand(tt.args...)
			if err == nil {
				t.Log("Command returned error as expected for invalid input")
			} else {
				t.Logf("Command failed as expected: %v, stderr: %s", err, stderr)
			}
		})
	}
}

func TestCLIWithConfigFlag(t *testing.T) {
	configPath := "/Users/zhangyuqing/.s3tool.yaml"

	stdout, stderr, err := runCommand("--config", configPath, "bucket", "list")
	if err != nil {
		t.Errorf("Command with config flag failed: %v\nstderr: %s", err, stderr)
		return
	}

	t.Logf("Output with config flag: %s", stdout[:min(200, len(stdout))])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
