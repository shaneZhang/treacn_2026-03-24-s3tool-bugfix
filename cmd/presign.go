package cmd

import (
	"context"
	"fmt"
	"time"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

var presignCmd = &cobra.Command{
	Use:   "presign",
	Short: "预签名URL",
	Long:  "生成预签名 URL 用于临时访问",
}

var presignGetCmd = &cobra.Command{
	Use:   "get [bucket] [key]",
	Short: "生成GET预签名URL",
	Long:  "生成用于下载对象的预签名URL",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key := args[0], args[1]
		ctx := context.Background()

		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		expires, _ := cmd.Flags().GetInt("expires")
		if expires == 0 {
			expires = 3600
		}

		presignClient := s3.NewPresignClient(client)
		url, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(expires) * time.Second
		})
		if err != nil {
			return fmt.Errorf("生成预签名URL失败: %w", err)
		}

		cmd.Println(url.URL)
		return nil
	},
}

var presignPutCmd = &cobra.Command{
	Use:   "put [bucket] [key]",
	Short: "生成PUT预签名URL",
	Long:  "生成用于上传对象的预签名URL",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key := args[0], args[1]
		ctx := context.Background()

		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		expires, _ := cmd.Flags().GetInt("expires")
		if expires == 0 {
			expires = 3600
		}

		presignClient := s3.NewPresignClient(client)
		url, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(expires) * time.Second
		})
		if err != nil {
			return fmt.Errorf("生成预签名URL失败: %w", err)
		}

		cmd.Println(url.URL)
		return nil
	},
}

var presignDeleteCmd = &cobra.Command{
	Use:   "delete [bucket] [key]",
	Short: "生成DELETE预签名URL",
	Long:  "生成用于删除对象的预签名URL",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key := args[0], args[1]
		ctx := context.Background()

		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		expires, _ := cmd.Flags().GetInt("expires")
		if expires == 0 {
			expires = 3600
		}

		presignClient := s3.NewPresignClient(client)
		url, err := presignClient.PresignDeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(expires) * time.Second
		})
		if err != nil {
			return fmt.Errorf("生成预签名URL失败: %w", err)
		}

		cmd.Println(url.URL)
		return nil
	},
}

func init() {
	presignCmd.AddCommand(presignGetCmd, presignPutCmd, presignDeleteCmd)
	presignGetCmd.Flags().Int("expires", 3600, "URL过期时间(秒)")
	presignPutCmd.Flags().Int("expires", 3600, "URL过期时间(秒)")
	presignDeleteCmd.Flags().Int("expires", 3600, "URL过期时间(秒)")
}
