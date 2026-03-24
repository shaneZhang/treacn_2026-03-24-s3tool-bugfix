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

var websiteCmd = &cobra.Command{
	Use:   "website",
	Short: "静态网站配置",
	Long:  "管理存储桶的静态网站配置",
}

var websiteGetCmd = &cobra.Command{
	Use:   "get [bucket]",
	Short: "获取静态网站配置",
	Long:  "获取指定存储桶的静态网站配置",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetBucketWebsite(ctx, &s3.GetBucketWebsiteInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("获取静态网站配置失败: %w", err)
		}

		if output.IndexDocument != nil {
			cmd.Printf("索引文档: %s\n", *output.IndexDocument.Suffix)
		}
		if output.ErrorDocument != nil {
			cmd.Printf("错误文档: %s\n", *output.ErrorDocument.Key)
		}
		return nil
	},
}

var websiteEnableCmd = &cobra.Command{
	Use:   "enable [bucket] [index-document] [error-document]",
	Short: "启用静态网站",
	Long:  "为存储桶启用静态网站托管",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		indexDoc := args[1]
		errorDoc := ""
		if len(args) > 2 {
			errorDoc = args[2]
		}

		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		input := &s3.PutBucketWebsiteInput{
			Bucket: aws.String(bucket),
			WebsiteConfiguration: &types.WebsiteConfiguration{
				IndexDocument: &types.IndexDocument{
					Suffix: aws.String(indexDoc),
				},
			},
		}

		if errorDoc != "" {
			input.WebsiteConfiguration.ErrorDocument = &types.ErrorDocument{
				Key: aws.String(errorDoc),
			}
		}

		_, err = client.PutBucketWebsite(ctx, input)
		if err != nil {
			return fmt.Errorf("启用静态网站失败: %w", err)
		}

		cmd.Printf("存储桶 %s 静态网站已启用\n", bucket)
		return nil
	},
}

var websiteDisableCmd = &cobra.Command{
	Use:   "disable [bucket]",
	Short: "禁用静态网站",
	Long:  "禁用存储桶的静态网站托管",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.DeleteBucketWebsite(ctx, &s3.DeleteBucketWebsiteInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("禁用静态网站失败: %w", err)
		}

		cmd.Printf("存储桶 %s 静态网站已禁用\n", bucket)
		return nil
	},
}

func init() {
	websiteCmd.AddCommand(websiteGetCmd, websiteEnableCmd, websiteDisableCmd)
}
