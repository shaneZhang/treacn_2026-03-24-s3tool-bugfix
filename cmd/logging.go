package cmd

import (
	"context"
	"fmt"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/cobra"
)

var loggingCmd = &cobra.Command{
	Use:   "logging",
	Short: "日志记录配置",
	Long:  "管理存储桶的服务器访问日志记录",
}

var loggingGetCmd = &cobra.Command{
	Use:   "get [bucket]",
	Short: "获取日志记录配置",
	Long:  "获取指定存储桶的日志记录配置",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetBucketLogging(ctx, &s3.GetBucketLoggingInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("获取日志记录配置失败: %w", err)
		}

		if output.LoggingEnabled == nil {
			cmd.Println("该存储桶没有启用日志记录")
			return nil
		}

		cmd.Printf("日志目标存储桶: %s\n", *output.LoggingEnabled.TargetBucket)
		cmd.Printf("日志前缀: %s\n", *output.LoggingEnabled.TargetPrefix)
		return nil
	},
}

var loggingDisableCmd = &cobra.Command{
	Use:   "disable [bucket]",
	Short: "禁用日志记录",
	Long:  "禁用存储桶的服务器访问日志记录",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.PutBucketLogging(ctx, &s3.PutBucketLoggingInput{
			Bucket: aws.String(bucket),
			BucketLoggingStatus: &types.BucketLoggingStatus{
				LoggingEnabled: nil,
			},
		})
		if err != nil {
			return fmt.Errorf("禁用日志记录失败: %w", err)
		}

		cmd.Printf("存储桶 %s 日志记录已禁用\n", bucket)
		return nil
	},
}

func init() {
	loggingCmd.AddCommand(loggingGetCmd, loggingDisableCmd)
}
