package cmd

import (
	"context"
	"fmt"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

var corsCmd = &cobra.Command{
	Use:   "cors",
	Short: "CORS配置管理",
	Long:  "管理存储桶的跨域资源共享(CORS)配置",
}

var corsGetCmd = &cobra.Command{
	Use:   "get [bucket]",
	Short: "获取CORS配置",
	Long:  "获取指定存储桶的CORS配置",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetBucketCors(ctx, &s3.GetBucketCorsInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("获取CORS配置失败: %w", err)
		}

		if output.CORSRules == nil || len(output.CORSRules) == 0 {
			cmd.Println("该存储桶没有CORS配置")
			return nil
		}

		for i, rule := range output.CORSRules {
			cmd.Printf("规则 %d:\n", i+1)
			cmd.Printf("  允许的来源: %v\n", rule.AllowedOrigins)
			cmd.Printf("  允许的方法: %v\n", rule.AllowedMethods)
			cmd.Printf("  允许的头部: %v\n", rule.AllowedHeaders)
			cmd.Printf("  暴露的头部: %v\n", rule.ExposeHeaders)
			if rule.MaxAgeSeconds != nil {
				cmd.Printf("  缓存时间: %d秒\n", *rule.MaxAgeSeconds)
			}
			cmd.Println()
		}
		return nil
	},
}

var corsDeleteCmd = &cobra.Command{
	Use:   "delete [bucket]",
	Short: "删除CORS配置",
	Long:  "删除指定存储桶的CORS配置",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.DeleteBucketCors(ctx, &s3.DeleteBucketCorsInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("删除CORS配置失败: %w", err)
		}

		cmd.Printf("存储桶 %s CORS配置已删除\n", bucket)
		return nil
	},
}

func init() {
	corsCmd.AddCommand(corsGetCmd, corsDeleteCmd)
}
