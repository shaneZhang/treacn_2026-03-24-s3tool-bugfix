package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configData  string
		expectError bool
		checkFunc   func(t *testing.T)
	}{
		{
			name: "valid config file",
			configData: `region: us-west-2
access_key: test_access_key
secret_key: test_secret_key
endpoint: https://test.example.com
force_path_style: true
`,
			expectError: false,
			checkFunc: func(t *testing.T) {
				if GlobalConfig.Region != "us-west-2" {
					t.Errorf("expected region us-west-2, got %s", GlobalConfig.Region)
				}
				if GlobalConfig.AccessKey != "test_access_key" {
					t.Errorf("expected access_key test_access_key, got %s", GlobalConfig.AccessKey)
				}
				if GlobalConfig.SecretKey != "test_secret_key" {
					t.Errorf("expected secret_key test_secret_key, got %s", GlobalConfig.SecretKey)
				}
				if GlobalConfig.Endpoint != "https://test.example.com" {
					t.Errorf("expected endpoint https://test.example.com, got %s", GlobalConfig.Endpoint)
				}
				if !GlobalConfig.ForcePathStyle {
					t.Errorf("expected force_path_style true, got %v", GlobalConfig.ForcePathStyle)
				}
			},
		},
		{
			name:        "empty config file",
			configData:  ``,
			expectError: false,
			checkFunc: func(t *testing.T) {
				if GlobalConfig.Region != "us-east-1" {
					t.Errorf("expected default region us-east-1, got %s", GlobalConfig.Region)
				}
			},
		},
		{
			name: "partial config",
			configData: `region: eu-west-1
access_key: my_key
`,
			expectError: false,
			checkFunc: func(t *testing.T) {
				if GlobalConfig.Region != "eu-west-1" {
					t.Errorf("expected region eu-west-1, got %s", GlobalConfig.Region)
				}
				if GlobalConfig.AccessKey != "my_key" {
					t.Errorf("expected access_key my_key, got %s", GlobalConfig.AccessKey)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "test_config.yaml")

			err := os.WriteFile(configPath, []byte(tt.configData), 0644)
			if err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			err = LoadConfig(configPath)
			if (err != nil) != tt.expectError {
				t.Errorf("LoadConfig() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t)
			}
		})
	}
}

func TestLoadConfigWithEnvVars(t *testing.T) {
	viper.Reset()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")

	configData := `region: original-region
access_key: original_key
secret_key: original_secret
`
	err := os.WriteFile(configPath, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	os.Setenv("AWS_REGION", "env-region")
	os.Setenv("AWS_ACCESS_KEY_ID", "env_access_key")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "env_secret_key")
	os.Setenv("AWS_ENDPOINT", "https://env.example.com")
	defer func() {
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
		os.Unsetenv("AWS_ENDPOINT")
	}()

	err = LoadConfig(configPath)
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
		return
	}

	if GlobalConfig.Region != "env-region" {
		t.Errorf("expected region env-region, got %s", GlobalConfig.Region)
	}
	if GlobalConfig.AccessKey != "env_access_key" {
		t.Errorf("expected access_key env_access_key, got %s", GlobalConfig.AccessKey)
	}
	if GlobalConfig.SecretKey != "env_secret_key" {
		t.Errorf("expected secret_key env_secret_key, got %s", GlobalConfig.SecretKey)
	}
	if GlobalConfig.Endpoint != "https://env.example.com" {
		t.Errorf("expected endpoint https://env.example.com, got %s", GlobalConfig.Endpoint)
	}
}

func TestLoadConfigNonExistentFile(t *testing.T) {
	err := LoadConfig("/non/existent/path/config.yaml")
	if err == nil {
		t.Error("expected error for non-existent config file, got nil")
	}
}

func TestGetAWSConfig(t *testing.T) {
	GlobalConfig = Config{
		Region:    "us-west-2",
		AccessKey: "test_key",
		SecretKey: "test_secret",
	}

	ctx := context.Background()
	awsConfig, err := GetAWSConfig(ctx)
	if err != nil {
		t.Errorf("GetAWSConfig() error = %v", err)
		return
	}

	if awsConfig.Region != "us-west-2" {
		t.Errorf("expected region us-west-2, got %s", awsConfig.Region)
	}
}

func TestGetAWSConfigWithProfile(t *testing.T) {
	originalProfile := GlobalConfig.Profile
	defer func() { GlobalConfig.Profile = originalProfile }()

	GlobalConfig = Config{
		Region:    "us-east-1",
		Profile:   "",
		AccessKey: "test_key",
		SecretKey: "test_secret",
	}

	ctx := context.Background()
	awsConfig, err := GetAWSConfig(ctx)
	if err != nil {
		t.Errorf("GetAWSConfig() error = %v", err)
		return
	}

	if awsConfig.Region != "us-east-1" {
		t.Errorf("expected region us-east-1, got %s", awsConfig.Region)
	}
}

func TestConfigDefaults(t *testing.T) {
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_ENDPOINT")
	os.Unsetenv("AWS_PROFILE")
	viper.Reset()

	GlobalConfig = Config{}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "empty_config.yaml")
	err := os.WriteFile(configPath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	err = LoadConfig(configPath)
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
		return
	}

	if GlobalConfig.Region != "us-east-1" {
		t.Errorf("expected default region us-east-1, got %s", GlobalConfig.Region)
	}

	if GlobalConfig.ForcePathStyle {
		t.Errorf("expected default force_path_style false, got %v", GlobalConfig.ForcePathStyle)
	}

	if GlobalConfig.UseAccelerate {
		t.Errorf("expected default use_accelerate false, got %v", GlobalConfig.UseAccelerate)
	}
}

func TestInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid_config.yaml")

	invalidYAML := `
region: [invalid
  unclosed bracket
`
	err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	err = LoadConfig(configPath)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestConfigWithExtraFields(t *testing.T) {
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_ENDPOINT")
	os.Unsetenv("AWS_PROFILE")
	viper.Reset()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "extra_fields_config.yaml")

	configData := `region: ap-northeast-1
access_key: my_key
secret_key: my_secret
extra_field: this_should_be_ignored
another_field: 12345
`
	err := os.WriteFile(configPath, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	err = LoadConfig(configPath)
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
		return
	}

	if GlobalConfig.Region != "ap-northeast-1" {
		t.Errorf("expected region ap-northeast-1, got %s", GlobalConfig.Region)
	}
}

func TestConfigBooleanValues(t *testing.T) {
	tests := []struct {
		name          string
		configData    string
		expectedPath  bool
		expectedAccel bool
	}{
		{
			name: "force_path_style true",
			configData: `force_path_style: true
`,
			expectedPath:  true,
			expectedAccel: false,
		},
		{
			name: "use_accelerate true",
			configData: `use_accelerate: true
`,
			expectedPath:  false,
			expectedAccel: true,
		},
		{
			name: "both true",
			configData: `force_path_style: true
use_accelerate: true
`,
			expectedPath:  true,
			expectedAccel: true,
		},
		{
			name: "both false",
			configData: `force_path_style: false
use_accelerate: false
`,
			expectedPath:  false,
			expectedAccel: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "test_config.yaml")

			err := os.WriteFile(configPath, []byte(tt.configData), 0644)
			if err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			err = LoadConfig(configPath)
			if err != nil {
				t.Errorf("LoadConfig() error = %v", err)
				return
			}

			if GlobalConfig.ForcePathStyle != tt.expectedPath {
				t.Errorf("expected ForcePathStyle %v, got %v", tt.expectedPath, GlobalConfig.ForcePathStyle)
			}

			if GlobalConfig.UseAccelerate != tt.expectedAccel {
				t.Errorf("expected UseAccelerate %v, got %v", tt.expectedAccel, GlobalConfig.UseAccelerate)
			}
		})
	}
}
