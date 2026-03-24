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

var versioningCmd = &cobra.Command{
	Use:   "versioning",
	Short: "版本控制管理",
	Long:  "管理存储桶的版本控制配置",
}

var versioningGetCmd = &cobra.Command{
	Use:   "get [bucket]",
	Short: "获取版本控制状态",
	Long:  "获取指定存储桶的版本控制状态",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("获取版本控制状态失败: %w", err)
		}

		status := string(output.Status)
		if status == "" {
			status = "未启用"
		}
		mfaDelete := string(output.MFADelete)
		if mfaDelete == "" {
			mfaDelete = "未启用"
		}

		cmd.Printf("存储桶: %s\n", bucket)
		cmd.Printf("版本控制状态: %s\n", status)
		cmd.Printf("MFA Delete: %s\n", mfaDelete)
		return nil
	},
}

var versioningEnableCmd = &cobra.Command{
	Use:   "enable [bucket]",
	Short: "启用版本控制",
	Long:  "启用指定存储桶的版本控制",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucket),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: "Enabled",
			},
		})
		if err != nil {
			return fmt.Errorf("启用版本控制失败: %w", err)
		}

		cmd.Printf("存储桶 %s 版本控制已启用\n", bucket)
		return nil
	},
}

var versioningSuspendCmd = &cobra.Command{
	Use:   "suspend [bucket]",
	Short: "暂停版本控制",
	Long:  "暂停指定存储桶的版本控制",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucket),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: "Suspended",
			},
		})
		if err != nil {
			return fmt.Errorf("暂停版本控制失败: %w", err)
		}

		cmd.Printf("存储桶 %s 版本控制已暂停\n", bucket)
		return nil
	},
}

func init() {
	versioningCmd.AddCommand(versioningGetCmd, versioningEnableCmd, versioningSuspendCmd)
}
