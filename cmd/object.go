package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var objectCmd = &cobra.Command{
	Use:   "object",
	Short: "对象管理",
	Long:  "对象相关操作: 上传、下载、列表、删除等",
}

var objectListCmd = &cobra.Command{
	Use:   "list [bucket]",
	Short: "列出对象",
	Long:  "列出存储桶中的对象",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		prefix, _ := cmd.Flags().GetString("prefix")
		recursive, _ := cmd.Flags().GetBool("recursive")
		maxKeys, _ := cmd.Flags().GetInt("max-keys")

		t := table.NewWriter()
		t.AppendHeader(table.Row{"键名", "大小", "最后修改", "存储类型"})

		count := 0
		var continuationToken *string
		for {
			page, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
				Bucket:            aws.String(bucket),
				Prefix:            aws.String(prefix),
				MaxKeys:           aws.Int32(int32(maxKeys)),
				ContinuationToken: continuationToken,
			})
			if err != nil {
				return fmt.Errorf("列举对象失败: %w", err)
			}

			for _, obj := range page.Contents {
				key := *obj.Key
				if !recursive && strings.Contains(key, "/") && key != prefix {
					continue
				}

				t.AppendRow([]interface{}{
					key,
					formatBytes(*obj.Size),
					obj.LastModified.Format("2006-01-02 15:04:05"),
					obj.StorageClass,
				})
				count++
			}

			if page.IsTruncated == nil || !*page.IsTruncated {
				break
			}
			continuationToken = page.NextContinuationToken
		}

		cmd.Printf("总对象数量: %d\n", count)
		cmd.Println(t.Render())
		return nil
	},
}

var objectPutCmd = &cobra.Command{
	Use:   "put [bucket] [key] [file]",
	Short: "上传对象",
	Long:  "上传本地文件到 S3 存储桶",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key, filePath := args[0], args[1], args[2]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("打开文件失败: %w", err)
		}
		defer file.Close()

		contentType, _ := cmd.Flags().GetString("content-type")
		storageClass, _ := cmd.Flags().GetString("storage-class")

		input := &s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(key),
			Body:        file,
			ContentType: aws.String(contentType),
		}

		if storageClass != "" {
			input.StorageClass = types.StorageClass(storageClass)
		}

		_, err = client.PutObject(ctx, input)
		if err != nil {
			return fmt.Errorf("上传对象失败: %w", err)
		}

		cmd.Printf("文件 %s 上传成功到 s3://%s/%s\n", filePath, bucket, key)
		return nil
	},
}

var objectGetCmd = &cobra.Command{
	Use:   "get [bucket] [key] [local-file]",
	Short: "下载对象",
	Long:  "从 S3 存储桶下载对象到本地",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key, filePath := args[0], args[1], args[2]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("下载对象失败: %w", err)
		}
		defer output.Body.Close()

		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("创建文件失败: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(file, output.Body)
		if err != nil {
			return fmt.Errorf("写入文件失败: %w", err)
		}

		cmd.Printf("对象 s3://%s/%s 下载成功到 %s\n", bucket, key, filePath)
		return nil
	},
}

var objectDeleteCmd = &cobra.Command{
	Use:   "delete [bucket] [key]",
	Short: "删除对象",
	Long:  "从存储桶删除指定对象",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key := args[0], args[1]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("删除对象失败: %w", err)
		}

		cmd.Printf("对象 s3://%s/%s 已删除\n", bucket, key)
		return nil
	},
}

var objectCopyCmd = &cobra.Command{
	Use:   "copy [source-bucket] [source-key] [dest-bucket] [dest-key]",
	Short: "复制对象",
	Long:  "在 S3 存储桶之间复制对象",
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcBucket, srcKey, destBucket, destKey := args[0], args[1], args[2], args[3]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		copySource := srcBucket + "/" + srcKey
		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(destBucket),
			Key:        aws.String(destKey),
			CopySource: aws.String(copySource),
		})
		if err != nil {
			return fmt.Errorf("复制对象失败: %w", err)
		}

		cmd.Printf("对象从 s3://%s/%s 复制到 s3://%s/%s\n", srcBucket, srcKey, destBucket, destKey)
		return nil
	},
}

var objectInfoCmd = &cobra.Command{
	Use:   "info [bucket] [key]",
	Short: "获取对象信息",
	Long:  "获取对象的详细信息",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key := args[0], args[1]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("获取对象信息失败: %w", err)
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"属性", "值"})
		t.AppendRow([]interface{}{"键名", key})
		t.AppendRow([]interface{}{"大小", formatBytes(*output.ContentLength)})
		t.AppendRow([]interface{}{"内容类型", *output.ContentType})
		t.AppendRow([]interface{}{"最后修改", output.LastModified.Format("2006-01-02 15:04:05")})
		t.AppendRow([]interface{}{"ETag", *output.ETag})
		if output.StorageClass != "" {
			t.AppendRow([]interface{}{"存储类型", string(output.StorageClass)})
		}

		cmd.Println(t.Render())
		return nil
	},
}

var objectUrlCmd = &cobra.Command{
	Use:   "url [bucket] [key]",
	Short: "生成对象 URL",
	Long:  "生成对象的直接访问 URL",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key := args[0], args[1]
		ctx := context.Background()

		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		presignClient := s3.NewPresignClient(client)
		url, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("生成 URL 失败: %w", err)
		}

		cmd.Println(url.URL)
		return nil
	},
}

var objectMvCmd = &cobra.Command{
	Use:   "mv [bucket] [source-key] [dest-key]",
	Short: "移动对象",
	Long:  "移动/重命名对象",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, srcKey, destKey := args[0], args[1], args[2]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		copySource := bucket + "/" + srcKey
		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bucket),
			Key:        aws.String(destKey),
			CopySource: aws.String(copySource),
		})
		if err != nil {
			return fmt.Errorf("复制对象失败: %w", err)
		}

		_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(srcKey),
		})
		if err != nil {
			return fmt.Errorf("删除源对象失败: %w", err)
		}

		cmd.Printf("对象从 %s 移动到 %s\n", srcKey, destKey)
		return nil
	},
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func init() {
	objectCmd.AddCommand(objectListCmd, objectPutCmd, objectGetCmd, objectDeleteCmd, objectCopyCmd, objectInfoCmd, objectUrlCmd, objectMvCmd)

	objectListCmd.Flags().StringP("prefix", "p", "", "对象前缀过滤")
	objectListCmd.Flags().BoolP("recursive", "r", false, "递归列出所有对象")
	objectListCmd.Flags().Int("max-keys", 1000, "最大返回对象数")

	objectPutCmd.Flags().StringP("content-type", "t", "application/octet-stream", "内容类型")
	objectPutCmd.Flags().StringP("storage-class", "s", "", "存储类型 (STANDARD, REDUCED_REDUNDANCY, etc.)")
}
