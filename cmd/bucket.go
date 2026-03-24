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

		// 优先使用命令行指定的区域，其次使用配置文件中的区域
		region, _ := cmd.Flags().GetString("region")
		if region == "" {
			region = config.GlobalConfig.Region
		}

		// 只有在区域不是 us-east-1 时才需要设置 LocationConstraint
		// us-east-1 是默认区域，不需要指定 LocationConstraint
		if region != "" && region != "us-east-1" {
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

		// 获取存储桶区域
		locationOutput, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			location := string(locationOutput.LocationConstraint)
			if location == "" {
				location = "us-east-1"
			}
			t.AppendRow([]interface{}{"区域", location})
		}

		// 获取存储桶创建时间
		listOutput, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err == nil {
			for _, bucket := range listOutput.Buckets {
				if *bucket.Name == bucketName {
					t.AppendRow([]interface{}{"创建时间", bucket.CreationDate.Format("2006-01-02 15:04:05")})
					break
				}
			}
		}

		// 获取版本控制状态
		versionOutput, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			status := "未启用"
			if versionOutput.Status == types.BucketVersioningStatusEnabled {
				status = "已启用"
			} else if versionOutput.Status == types.BucketVersioningStatusSuspended {
				status = "已暂停"
			}
			t.AppendRow([]interface{}{"版本控制", status})
		}

		// 获取标签
		tagsOutput, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil && len(tagsOutput.TagSet) > 0 {
			tags := ""
			for i, tag := range tagsOutput.TagSet {
				if i > 0 {
					tags += ", "
				}
				tags += fmt.Sprintf("%s=%s", *tag.Key, *tag.Value)
			}
			t.AppendRow([]interface{}{"标签", tags})
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

	// 为 bucket create 命令添加 region 参数
	bucketCreateCmd.Flags().StringP("region", "r", "", "指定存储桶区域(覆盖配置文件)")
}
