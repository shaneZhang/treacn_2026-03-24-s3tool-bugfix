package config

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

type Config struct {
	Profile      string `mapstructure:"profile"`
	Region       string `mapstructure:"region"`
	Endpoint     string `mapstructure:"endpoint"`
	AccessKey    string `mapstructure:"access_key"`
	SecretKey    string `mapstructure:"secret_key"`
	ForcePathStyle bool `mapstructure:"force_path_style"`
	UseAccelerate bool `mapstructure:"use_accelerate"`
}

var GlobalConfig Config

func LoadConfig(configFile string) error {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	viper.SetDefault("region", "us-east-1")
	viper.SetDefault("force_path_style", false)
	viper.SetDefault("use_accelerate", false)

	envProfile := os.Getenv("AWS_PROFILE")
	envRegion := os.Getenv("AWS_REGION")
	envAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	envSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	envEndpoint := os.Getenv("AWS_ENDPOINT")

	if envProfile != "" {
		viper.Set("profile", envProfile)
	}
	if envRegion != "" {
		viper.Set("region", envRegion)
	}
	if envAccessKey != "" {
		viper.Set("access_key", envAccessKey)
	}
	if envSecretKey != "" {
		viper.Set("secret_key", envSecretKey)
	}
	if envEndpoint != "" {
		viper.Set("endpoint", envEndpoint)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func GetAWSConfig(ctx context.Context) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error

	opts = append(opts, config.WithRegion(GlobalConfig.Region))

	if GlobalConfig.Profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(GlobalConfig.Profile))
	}

	if GlobalConfig.AccessKey != "" && GlobalConfig.SecretKey != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				GlobalConfig.AccessKey,
				GlobalConfig.SecretKey,
				"",
			),
		))
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return awsConfig, nil
}

func GetS3Client(ctx context.Context) (*s3.Client, error) {
	awsConfig, err := GetAWSConfig(ctx)
	if err != nil {
		return nil, err
	}

	opts := []func(*s3.Options){}

	if GlobalConfig.Endpoint != "" {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(GlobalConfig.Endpoint)
		})
	}

	if GlobalConfig.ForcePathStyle {
		opts = append(opts, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}

	if GlobalConfig.UseAccelerate {
		opts = append(opts, func(o *s3.Options) {
			o.UseAccelerate = true
		})
	}

	client := s3.NewFromConfig(awsConfig, opts...)
	return client, nil
}

func GetS3ClientWithBucket(ctx context.Context, bucket string) (*s3.Client, error) {
	client, err := GetS3Client(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
