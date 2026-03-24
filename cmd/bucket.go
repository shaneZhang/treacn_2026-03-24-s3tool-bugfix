package cmd

import (
	"context"
	"fmt"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "存储桶管理",
	Long:  "存储桶相关操作: 列出、创建、删除等",
}

var bucketListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有存储桶",
	Long:  "列出当前账户下所有 S3 存储桶",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			return fmt.Errorf("列出存储桶失败: %w", err)
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"桶名称", "创建时间"})

		for _, bucket := range output.Buckets {
			t.AppendRow([]interface{}{
				*bucket.Name,
				bucket.CreationDate.Format("2006-01-02 15:04:05"),
			})
		}

		cmd.Println("总存储桶数量:", len(output.Buckets))
		cmd.Println(t.Render())
		return nil
	},
}

var bucketCreateCmd = &cobra.Command{
	Use:   "create [bucket-name]",
	Short: "创建存储桶",
	Long:  "创建一个新的 S3 存储桶",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucketName := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		input := &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		}

		if region := config.GlobalConfig.Region; region != "" && region != "us-east-1" {
			input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(region),
			}
		}

		_, err = client.CreateBucket(ctx, input)
		if err != nil {
			return fmt.Errorf("创建存储桶失败: %w", err)
		}

		cmd.Printf("存储桶 %s 创建成功\n", bucketName)
		return nil
	},
}

var bucketDeleteCmd = &cobra.Command{
	Use:   "delete [bucket-name]",
	Short: "删除存储桶",
	Long:  "删除指定的 S3 存储桶 (桶必须为空)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucketName := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("删除存储桶失败: %w", err)
		}

		cmd.Printf("存储桶 %s 删除成功\n", bucketName)
		return nil
	},
}

var bucketInfoCmd = &cobra.Command{
	Use:   "info [bucket-name]",
	Short: "获取存储桶信息",
	Long:  "获取指定存储桶的详细信息",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucketName := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"属性", "值"})

		locOutput, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			location := string(locOutput.LocationConstraint)
			if location == "" {
				location = "us-east-1"
			}
			t.AppendRow([]interface{}{"区域", location})
		}

		versioningOutput, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			status := string(versioningOutput.Status)
			if status == "" {
				status = "未启用"
			}
			t.AppendRow([]interface{}{"版本控制", status})
			if versioningOutput.MFADelete != "" {
				t.AppendRow([]interface{}{"MFA删除", string(versioningOutput.MFADelete)})
			}
		}

		tagsOutput, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil && len(tagsOutput.TagSet) > 0 {
			for _, tag := range tagsOutput.TagSet {
				t.AppendRow([]interface{}{"标签: " + *tag.Key, *tag.Value})
			}
		}

		cmd.Printf("存储桶: %s\n", bucketName)
		cmd.Println(t.Render())
		return nil
	},
}

var bucketLocationCmd = &cobra.Command{
	Use:   "location [bucket-name]",
	Short: "获取存储桶区域",
	Long:  "获取指定存储桶所在的区域",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucketName := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("获取存储桶区域失败: %w", err)
		}

		location := string(output.LocationConstraint)
		if location == "" {
			location = "us-east-1"
		}
		cmd.Printf("存储桶 %s 区域: %s\n", bucketName, location)
		return nil
	},
}

var bucketEmptyCmd = &cobra.Command{
	Use:   "empty [bucket-name]",
	Short: "清空存储桶",
	Long:  "删除存储桶中的所有对象",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucketName := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		deletedCount := 0
		for {
			listOutput, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
				Bucket: aws.String(bucketName),
			})
			if err != nil {
				return fmt.Errorf("列举对象失败: %w", err)
			}

			if len(listOutput.Contents) == 0 {
				break
			}

			objectsToDelete := make([]types.ObjectIdentifier, len(listOutput.Contents))
			for i, obj := range listOutput.Contents {
				objectsToDelete[i] = types.ObjectIdentifier{Key: obj.Key}
			}

			_, err = client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket: aws.String(bucketName),
				Delete: &types.Delete{
					Objects: objectsToDelete,
				},
			})
			if err != nil {
				return fmt.Errorf("删除对象失败: %w", err)
			}

			deletedCount += len(listOutput.Contents)

			if listOutput.IsTruncated == nil || !*listOutput.IsTruncated {
				break
			}
		}

		cmd.Printf("存储桶 %s 已清空，删除了 %d 个对象\n", bucketName, deletedCount)
		return nil
	},
}

func init() {
	bucketCmd.AddCommand(bucketListCmd)
	bucketCmd.AddCommand(bucketCreateCmd)
	bucketCmd.AddCommand(bucketDeleteCmd)
	bucketCmd.AddCommand(bucketInfoCmd)
	bucketCmd.AddCommand(bucketLocationCmd)
	bucketCmd.AddCommand(bucketEmptyCmd)
}
