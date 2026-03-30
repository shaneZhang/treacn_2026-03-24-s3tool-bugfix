package tests

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestLifecycleGet(t *testing.T) {
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

	t.Run("get lifecycle configuration (empty)", func(t *testing.T) {
		_, err := tc.Client.GetBucketLifecycleConfiguration(tc.Context, &s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NoSuchLifecycleConfiguration") {
				t.Log("No lifecycle configuration (expected for new bucket)")
			} else {
				t.Logf("GetBucketLifecycleConfiguration() returned: %v", err)
			}
		}
	})
}

func TestLifecycleSetAndDelete(t *testing.T) {
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

	t.Run("set lifecycle configuration", func(t *testing.T) {
		rules := []types.LifecycleRule{
			{
				ID:         aws.String("test-rule-1"),
				Status:     types.ExpirationStatusEnabled,
				Prefix:     aws.String("logs/"),
				Expiration: &types.LifecycleExpiration{Days: aws.Int32(30)},
			},
			{
				ID:         aws.String("test-rule-2"),
				Status:     types.ExpirationStatusEnabled,
				Prefix:     aws.String("temp/"),
				Expiration: &types.LifecycleExpiration{Days: aws.Int32(7)},
			},
		}

		_, err := tc.Client.PutBucketLifecycleConfiguration(tc.Context, &s3.PutBucketLifecycleConfigurationInput{
			Bucket: aws.String(tc.BucketName),
			LifecycleConfiguration: &types.BucketLifecycleConfiguration{
				Rules: rules,
			},
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Lifecycle configuration not supported by this S3 implementation")
			}
			t.Errorf("PutBucketLifecycleConfiguration() error = %v", err)
			return
		}
		t.Log("Lifecycle configuration set successfully")
	})

	t.Run("get lifecycle configuration", func(t *testing.T) {
		output, err := tc.Client.GetBucketLifecycleConfiguration(tc.Context, &s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Lifecycle configuration not supported by this S3 implementation")
			}
			t.Errorf("GetBucketLifecycleConfiguration() error = %v", err)
			return
		}

		for _, rule := range output.Rules {
			t.Logf("Rule: ID=%s, Status=%s, Prefix=%s", *rule.ID, rule.Status, *rule.Prefix)
			if rule.Expiration != nil && rule.Expiration.Days != nil {
				t.Logf("  Expiration: %d days", *rule.Expiration.Days)
			}
		}
	})

	t.Run("delete lifecycle configuration", func(t *testing.T) {
		_, err := tc.Client.DeleteBucketLifecycle(tc.Context, &s3.DeleteBucketLifecycleInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Lifecycle configuration not supported by this S3 implementation")
			}
			t.Errorf("DeleteBucketLifecycle() error = %v", err)
			return
		}
		t.Log("Lifecycle configuration deleted successfully")
	})
}

func TestLifecycleWithTransitions(t *testing.T) {
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

	t.Run("set lifecycle with transitions", func(t *testing.T) {
		rules := []types.LifecycleRule{
			{
				ID:     aws.String("transition-rule"),
				Status: types.ExpirationStatusEnabled,
				Prefix: aws.String("archive/"),
				Transitions: []types.Transition{
					{
						Days:         aws.Int32(30),
						StorageClass: types.TransitionStorageClassStandardIa,
					},
					{
						Days:         aws.Int32(90),
						StorageClass: types.TransitionStorageClassGlacier,
					},
				},
			},
		}

		_, err := tc.Client.PutBucketLifecycleConfiguration(tc.Context, &s3.PutBucketLifecycleConfigurationInput{
			Bucket: aws.String(tc.BucketName),
			LifecycleConfiguration: &types.BucketLifecycleConfiguration{
				Rules: rules,
			},
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Lifecycle transitions not supported by this S3 implementation")
			}
			t.Errorf("PutBucketLifecycleConfiguration() error = %v", err)
			return
		}
		t.Log("Lifecycle with transitions set successfully")
	})

	t.Run("verify transitions", func(t *testing.T) {
		output, err := tc.Client.GetBucketLifecycleConfiguration(tc.Context, &s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Lifecycle configuration not supported by this S3 implementation")
			}
			t.Errorf("GetBucketLifecycleConfiguration() error = %v", err)
			return
		}

		for _, rule := range output.Rules {
			t.Logf("Rule: %s", *rule.ID)
			for _, transition := range rule.Transitions {
				t.Logf("  Transition: %d days -> %s", *transition.Days, transition.StorageClass)
			}
		}
	})

	tc.Client.DeleteBucketLifecycle(tc.Context, &s3.DeleteBucketLifecycleInput{
		Bucket: aws.String(tc.BucketName),
	})
}

func TestLifecycleWithNoncurrentVersionExpiration(t *testing.T) {
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

	t.Run("set lifecycle with noncurrent version expiration", func(t *testing.T) {
		rules := []types.LifecycleRule{
			{
				ID:     aws.String("noncurrent-expiration-rule"),
				Status: types.ExpirationStatusEnabled,
				NoncurrentVersionExpiration: &types.NoncurrentVersionExpiration{
					NoncurrentDays: aws.Int32(30),
				},
			},
		}

		_, err := tc.Client.PutBucketLifecycleConfiguration(tc.Context, &s3.PutBucketLifecycleConfigurationInput{
			Bucket: aws.String(tc.BucketName),
			LifecycleConfiguration: &types.BucketLifecycleConfiguration{
				Rules: rules,
			},
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Noncurrent version expiration not supported by this S3 implementation")
			}
			t.Errorf("PutBucketLifecycleConfiguration() error = %v", err)
			return
		}
		t.Log("Lifecycle with noncurrent version expiration set successfully")
	})

	tc.Client.DeleteBucketLifecycle(tc.Context, &s3.DeleteBucketLifecycleInput{
		Bucket: aws.String(tc.BucketName),
	})
}

func TestCORSGet(t *testing.T) {
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

	t.Run("get CORS configuration (empty)", func(t *testing.T) {
		_, err := tc.Client.GetBucketCors(tc.Context, &s3.GetBucketCorsInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NoSuchCORSConfiguration") {
				t.Log("No CORS configuration (expected for new bucket)")
			} else {
				t.Logf("GetBucketCors() returned: %v", err)
			}
		}
	})
}

func TestCORSSetAndDelete(t *testing.T) {
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

	t.Run("set CORS configuration", func(t *testing.T) {
		rules := []types.CORSRule{
			{
				AllowedHeaders: []string{"*"},
				AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
				AllowedOrigins: []string{"https://example.com", "https://app.example.com"},
				ExposeHeaders:  []string{"ETag", "x-amz-request-id"},
				MaxAgeSeconds:  aws.Int32(3600),
			},
		}

		_, err := tc.Client.PutBucketCors(tc.Context, &s3.PutBucketCorsInput{
			Bucket: aws.String(tc.BucketName),
			CORSConfiguration: &types.CORSConfiguration{
				CORSRules: rules,
			},
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("CORS configuration not supported by this S3 implementation")
			}
			t.Errorf("PutBucketCors() error = %v", err)
			return
		}
		t.Log("CORS configuration set successfully")
	})

	t.Run("get CORS configuration", func(t *testing.T) {
		output, err := tc.Client.GetBucketCors(tc.Context, &s3.GetBucketCorsInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("CORS configuration not supported by this S3 implementation")
			}
			t.Errorf("GetBucketCors() error = %v", err)
			return
		}

		for i, rule := range output.CORSRules {
			t.Logf("Rule %d:", i+1)
			t.Logf("  AllowedOrigins: %v", rule.AllowedOrigins)
			t.Logf("  AllowedMethods: %v", rule.AllowedMethods)
			t.Logf("  AllowedHeaders: %v", rule.AllowedHeaders)
			t.Logf("  ExposeHeaders: %v", rule.ExposeHeaders)
			if rule.MaxAgeSeconds != nil {
				t.Logf("  MaxAgeSeconds: %d", *rule.MaxAgeSeconds)
			}
		}
	})

	t.Run("delete CORS configuration", func(t *testing.T) {
		_, err := tc.Client.DeleteBucketCors(tc.Context, &s3.DeleteBucketCorsInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("CORS configuration not supported by this S3 implementation")
			}
			t.Errorf("DeleteBucketCors() error = %v", err)
			return
		}
		t.Log("CORS configuration deleted successfully")
	})
}

func TestWebsiteGet(t *testing.T) {
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

	t.Run("get website configuration (empty)", func(t *testing.T) {
		_, err := tc.Client.GetBucketWebsite(tc.Context, &s3.GetBucketWebsiteInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NoSuchWebsiteConfiguration") {
				t.Log("No website configuration (expected for new bucket)")
			} else {
				t.Logf("GetBucketWebsite() returned: %v", err)
			}
		}
	})
}

func TestWebsiteSetAndDelete(t *testing.T) {
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

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String("index.html"),
		Body:   bytes.NewReader([]byte("<html><body>Index</body></html>")),
	})
	if err != nil {
		t.Fatalf("Failed to upload index.html: %v", err)
	}

	_, err = tc.Client.PutObject(tc.Context, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String("error.html"),
		Body:   bytes.NewReader([]byte("<html><body>Error</body></html>")),
	})
	if err != nil {
		t.Fatalf("Failed to upload error.html: %v", err)
	}

	t.Run("set website configuration", func(t *testing.T) {
		_, err := tc.Client.PutBucketWebsite(tc.Context, &s3.PutBucketWebsiteInput{
			Bucket: aws.String(tc.BucketName),
			WebsiteConfiguration: &types.WebsiteConfiguration{
				IndexDocument: &types.IndexDocument{
					Suffix: aws.String("index.html"),
				},
				ErrorDocument: &types.ErrorDocument{
					Key: aws.String("error.html"),
				},
			},
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Website configuration not supported by this S3 implementation")
			}
			t.Errorf("PutBucketWebsite() error = %v", err)
			return
		}
		t.Log("Website configuration set successfully")
	})

	t.Run("get website configuration", func(t *testing.T) {
		output, err := tc.Client.GetBucketWebsite(tc.Context, &s3.GetBucketWebsiteInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Website configuration not supported by this S3 implementation")
			}
			t.Errorf("GetBucketWebsite() error = %v", err)
			return
		}

		if output.IndexDocument != nil {
			t.Logf("Index document: %s", *output.IndexDocument.Suffix)
		}
		if output.ErrorDocument != nil {
			t.Logf("Error document: %s", *output.ErrorDocument.Key)
		}
	})

	t.Run("delete website configuration", func(t *testing.T) {
		_, err := tc.Client.DeleteBucketWebsite(tc.Context, &s3.DeleteBucketWebsiteInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Website configuration not supported by this S3 implementation")
			}
			t.Errorf("DeleteBucketWebsite() error = %v", err)
			return
		}
		t.Log("Website configuration deleted successfully")
	})
}

func TestPolicyGet(t *testing.T) {
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

	t.Run("get bucket policy (empty)", func(t *testing.T) {
		_, err := tc.Client.GetBucketPolicy(tc.Context, &s3.GetBucketPolicyInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NoSuchBucketPolicy") {
				t.Log("No bucket policy (expected for new bucket)")
			} else {
				t.Logf("GetBucketPolicy() returned: %v", err)
			}
		}
	})
}

func TestPolicySetAndDelete(t *testing.T) {
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

	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": "*"},
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::` + tc.BucketName + `/*"
			}
		]
	}`

	t.Run("set bucket policy", func(t *testing.T) {
		_, err := tc.Client.PutBucketPolicy(tc.Context, &s3.PutBucketPolicyInput{
			Bucket: aws.String(tc.BucketName),
			Policy: aws.String(policy),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Bucket policy not supported by this S3 implementation")
			}
			t.Errorf("PutBucketPolicy() error = %v", err)
			return
		}
		t.Log("Bucket policy set successfully")
	})

	t.Run("get bucket policy", func(t *testing.T) {
		output, err := tc.Client.GetBucketPolicy(tc.Context, &s3.GetBucketPolicyInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Bucket policy not supported by this S3 implementation")
			}
			t.Errorf("GetBucketPolicy() error = %v", err)
			return
		}

		t.Logf("Policy: %s", *output.Policy)
	})

	t.Run("delete bucket policy", func(t *testing.T) {
		_, err := tc.Client.DeleteBucketPolicy(tc.Context, &s3.DeleteBucketPolicyInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("Bucket policy not supported by this S3 implementation")
			}
			t.Errorf("DeleteBucketPolicy() error = %v", err)
			return
		}
		t.Log("Bucket policy deleted successfully")
	})
}

func TestLoggingGet(t *testing.T) {
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

	t.Run("get logging configuration", func(t *testing.T) {
		output, err := tc.Client.GetBucketLogging(tc.Context, &s3.GetBucketLoggingInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			t.Errorf("GetBucketLogging() error = %v", err)
			return
		}

		if output.LoggingEnabled == nil {
			t.Log("No logging enabled (expected for new bucket)")
		} else {
			t.Logf("Target bucket: %s", *output.LoggingEnabled.TargetBucket)
			t.Logf("Target prefix: %s", *output.LoggingEnabled.TargetPrefix)
		}
	})
}

func TestReplicationGet(t *testing.T) {
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

	t.Run("get replication configuration", func(t *testing.T) {
		_, err := tc.Client.GetBucketReplication(tc.Context, &s3.GetBucketReplicationInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "ReplicationConfigurationNotFoundError") {
				t.Log("No replication configuration (expected for new bucket)")
			} else {
				t.Logf("GetBucketReplication() returned: %v", err)
			}
		}
	})
}

func TestNotificationGet(t *testing.T) {
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

	t.Run("get notification configuration", func(t *testing.T) {
		output, err := tc.Client.GetBucketNotificationConfiguration(tc.Context, &s3.GetBucketNotificationConfigurationInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			t.Errorf("GetBucketNotificationConfiguration() error = %v", err)
			return
		}

		hasConfig := len(output.TopicConfigurations) > 0 ||
			len(output.QueueConfigurations) > 0 ||
			len(output.LambdaFunctionConfigurations) > 0

		if !hasConfig {
			t.Log("No notification configuration (expected for new bucket)")
		}
	})
}

func TestACLGet(t *testing.T) {
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

	t.Run("get bucket ACL", func(t *testing.T) {
		output, err := tc.Client.GetBucketAcl(tc.Context, &s3.GetBucketAclInput{
			Bucket: aws.String(tc.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("ACL not supported by this S3 implementation")
			}
			t.Errorf("GetBucketAcl() error = %v", err)
			return
		}

		t.Logf("Owner: %s", *output.Owner.DisplayName)
		for _, grant := range output.Grants {
			t.Logf("Grant: %s", grant.Permission)
		}
	})
}

func TestACLSet(t *testing.T) {
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

	t.Run("set bucket ACL to private", func(t *testing.T) {
		_, err := tc.Client.PutBucketAcl(tc.Context, &s3.PutBucketAclInput{
			Bucket: aws.String(tc.BucketName),
			ACL:    types.BucketCannedACLPrivate,
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") || strings.Contains(err.Error(), "not supported") {
				t.Skip("ACL not supported by this S3 implementation")
			}
			t.Errorf("PutBucketAcl() error = %v", err)
			return
		}
		t.Log("Bucket ACL set to private")
	})
}
