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

var encryptionCmd = &cobra.Command{
	Use:   "encryption",
	Short: "加密配置管理",
	Long:  "管理存储桶的服务器端加密配置",
}

var encryptionGetCmd = &cobra.Command{
	Use:   "get [bucket]",
	Short: "获取加密配置",
	Long:  "获取指定存储桶的加密配置",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("获取加密配置失败: %w", err)
		}

		if output.ServerSideEncryptionConfiguration == nil || len(output.ServerSideEncryptionConfiguration.Rules) == 0 {
			cmd.Println("该存储桶没有启用加密")
			return nil
		}

		for i, rule := range output.ServerSideEncryptionConfiguration.Rules {
			cmd.Printf("规则 %d:\n", i+1)
			if apply := rule.ApplyServerSideEncryptionByDefault; apply != nil {
				cmd.Printf("  加密算法: %s\n", apply.SSEAlgorithm)
				if apply.KMSMasterKeyID != nil {
					cmd.Printf("  KMS密钥ID: %s\n", *apply.KMSMasterKeyID)
				}
			}
		}
		return nil
	},
}

var encryptionEnableCmd = &cobra.Command{
	Use:   "enable [bucket]",
	Short: "启用加密",
	Long:  "为存储桶启用服务器端加密(SSE-S3)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.PutBucketEncryption(ctx, &s3.PutBucketEncryptionInput{
			Bucket: aws.String(bucket),
			ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
				Rules: []types.ServerSideEncryptionRule{
					{
						ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
							SSEAlgorithm: types.ServerSideEncryptionAes256,
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("启用加密失败: %w", err)
		}

		cmd.Printf("存储桶 %s 加密已启用(SSE-S3)\n", bucket)
		return nil
	},
}

var encryptionDisableCmd = &cobra.Command{
	Use:   "disable [bucket]",
	Short: "禁用加密",
	Long:  "禁用存储桶的服务器端加密",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.DeleteBucketEncryption(ctx, &s3.DeleteBucketEncryptionInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("禁用加密失败: %w", err)
		}

		cmd.Printf("存储桶 %s 加密已禁用\n", bucket)
		return nil
	},
}

func init() {
	encryptionCmd.AddCommand(encryptionGetCmd, encryptionEnableCmd, encryptionDisableCmd)
}
