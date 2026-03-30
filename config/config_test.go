package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func resetViper() {
	viper.Reset()
	viper.SetConfigType("yaml")
	viper.SetDefault("region", "us-east-1")
	viper.SetDefault("force_path_style", false)
	viper.SetDefault("use_accelerate", false)
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		setup      func(t *testing.T) string
		wantErr    bool
	}{
		{
			name:       "加载有效配置文件",
			configFile: "",
			setup: func(t *testing.T) string {
				resetViper()
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, "test_config.yaml")
				content := `region: us-west-2
access_key: test_access_key
secret_key: test_secret_key
endpoint: http://localhost:9000
force_path_style: true
use_accelerate: false`
				if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
				return configPath
			},
			wantErr: false,
		},
		{
			name:       "配置文件不存在",
			configFile: "/nonexistent/path/config.yaml",
			setup: func(t *testing.T) string {
				resetViper()
				return "/nonexistent/path/config.yaml"
			},
			wantErr: true,
		},
		{
			name:       "无效的YAML格式",
			configFile: "",
			setup: func(t *testing.T) string {
				resetViper()
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, "invalid_config.yaml")
				content := `region: us-west-2
access_key: [invalid: yaml: format`
				if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
				return configPath
			},
			wantErr: true,
		},
		{
			name:       "空配置文件",
			configFile: "",
			setup: func(t *testing.T) string {
				resetViper()
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, "empty_config.yaml")
				if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
					t.Fatal(err)
				}
				return configPath
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := tt.setup(t)
			err := LoadConfig(configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadConfig_Values(t *testing.T) {
	resetViper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")
	content := `region: eu-west-1
access_key: AKIAIOSFODNN7EXAMPLE
secret_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
endpoint: https://s3.custom.com
force_path_style: true
use_accelerate: true
profile: custom-profile`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if GlobalConfig.Region != "eu-west-1" {
		t.Errorf("Region = %v, want eu-west-1", GlobalConfig.Region)
	}
	if GlobalConfig.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("AccessKey = %v, want AKIAIOSFODNN7EXAMPLE", GlobalConfig.AccessKey)
	}
	if GlobalConfig.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("SecretKey = %v, want wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", GlobalConfig.SecretKey)
	}
	if GlobalConfig.Endpoint != "https://s3.custom.com" {
		t.Errorf("Endpoint = %v, want https://s3.custom.com", GlobalConfig.Endpoint)
	}
	if !GlobalConfig.ForcePathStyle {
		t.Error("ForcePathStyle should be true")
	}
	if !GlobalConfig.UseAccelerate {
		t.Error("UseAccelerate should be true")
	}
	if GlobalConfig.Profile != "custom-profile" {
		t.Errorf("Profile = %v, want custom-profile", GlobalConfig.Profile)
	}
}

func TestLoadConfig_EnvOverride(t *testing.T) {
	resetViper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")
	content := `region: us-east-1
access_key: config_access_key
secret_key: config_secret_key`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	os.Setenv("AWS_REGION", "ap-northeast-1")
	os.Setenv("AWS_ACCESS_KEY_ID", "env_access_key")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "env_secret_key")
	os.Setenv("AWS_ENDPOINT", "http://env-endpoint:9000")
	defer func() {
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
		os.Unsetenv("AWS_ENDPOINT")
	}()

	err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if GlobalConfig.Region != "ap-northeast-1" {
		t.Errorf("Region = %v, want ap-northeast-1", GlobalConfig.Region)
	}
	if GlobalConfig.AccessKey != "env_access_key" {
		t.Errorf("AccessKey = %v, want env_access_key", GlobalConfig.AccessKey)
	}
	if GlobalConfig.SecretKey != "env_secret_key" {
		t.Errorf("SecretKey = %v, want env_secret_key", GlobalConfig.SecretKey)
	}
	if GlobalConfig.Endpoint != "http://env-endpoint:9000" {
		t.Errorf("Endpoint = %v, want http://env-endpoint:9000", GlobalConfig.Endpoint)
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	resetViper()
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_ENDPOINT")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "minimal_config.yaml")
	content := `access_key: test_key
secret_key: test_secret`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if GlobalConfig.Region != "us-east-1" {
		t.Errorf("Region = %v, want us-east-1 (default)", GlobalConfig.Region)
	}
	if GlobalConfig.ForcePathStyle {
		t.Error("ForcePathStyle should be false (default)")
	}
	if GlobalConfig.UseAccelerate {
		t.Error("UseAccelerate should be false (default)")
	}
}

func TestGetAWSConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "有效配置",
			config: Config{
				Region:    "us-east-1",
				AccessKey: "test_access_key",
				SecretKey: "test_secret_key",
			},
			wantErr: false,
		},
		{
			name: "只有Region",
			config: Config{
				Region: "us-west-2",
			},
			wantErr: false,
		},
		{
			name: "无效Region",
			config: Config{
				Region: "invalid-region-123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GlobalConfig = tt.config
			ctx := context.Background()
			_, err := GetAWSConfig(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAWSConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetS3Client(t *testing.T) {
	GlobalConfig = Config{
		Region:    "us-east-1",
		AccessKey: "test_access_key",
		SecretKey: "test_secret_key",
	}

	ctx := context.Background()
	client, err := GetS3Client(ctx)
	if err != nil {
		t.Errorf("GetS3Client() error = %v", err)
	}
	if client == nil {
		t.Error("GetS3Client() returned nil client")
	}
}

func TestGetS3Client_WithEndpoint(t *testing.T) {
	GlobalConfig = Config{
		Region:         "us-east-1",
		AccessKey:      "test_access_key",
		SecretKey:      "test_secret_key",
		Endpoint:       "http://localhost:9000",
		ForcePathStyle: true,
	}

	ctx := context.Background()
	client, err := GetS3Client(ctx)
	if err != nil {
		t.Errorf("GetS3Client() error = %v", err)
	}
	if client == nil {
		t.Error("GetS3Client() returned nil client")
	}
}

func TestGetS3ClientWithBucket(t *testing.T) {
	GlobalConfig = Config{
		Region:    "us-east-1",
		AccessKey: "test_access_key",
		SecretKey: "test_secret_key",
	}

	ctx := context.Background()
	client, err := GetS3ClientWithBucket(ctx, "test-bucket")
	if err != nil {
		t.Errorf("GetS3ClientWithBucket() error = %v", err)
	}
	if client == nil {
		t.Error("GetS3ClientWithBucket() returned nil client")
	}
}

func TestConfigStruct(t *testing.T) {
	config := Config{
		Profile:        "test-profile",
		Region:         "us-west-2",
		Endpoint:       "http://minio:9000",
		AccessKey:      "test_access",
		SecretKey:      "test_secret",
		ForcePathStyle: true,
		UseAccelerate:  false,
	}

	if config.Profile != "test-profile" {
		t.Errorf("Profile = %v, want test-profile", config.Profile)
	}
	if config.Region != "us-west-2" {
		t.Errorf("Region = %v, want us-west-2", config.Region)
	}
	if config.Endpoint != "http://minio:9000" {
		t.Errorf("Endpoint = %v, want http://minio:9000", config.Endpoint)
	}
	if config.AccessKey != "test_access" {
		t.Errorf("AccessKey = %v, want test_access", config.AccessKey)
	}
	if config.SecretKey != "test_secret" {
		t.Errorf("SecretKey = %v, want test_secret", config.SecretKey)
	}
	if !config.ForcePathStyle {
		t.Error("ForcePathStyle should be true")
	}
	if config.UseAccelerate {
		t.Error("UseAccelerate should be false")
	}
}
