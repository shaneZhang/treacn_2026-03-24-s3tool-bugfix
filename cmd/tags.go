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

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "标签管理",
	Long:  "管理存储桶和对象的标签",
}

var tagsBucketGetCmd = &cobra.Command{
	Use:   "bucket-get [bucket]",
	Short: "获取存储桶标签",
	Long:  "获取指定存储桶的标签",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("获取存储桶标签失败: %w", err)
		}

		if len(output.TagSet) == 0 {
			cmd.Println("该存储桶没有标签")
			return nil
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"键", "值"})

		for _, tag := range output.TagSet {
			t.AppendRow([]interface{}{*tag.Key, *tag.Value})
		}

		cmd.Println(t.Render())
		return nil
	},
}

var tagsBucketPutCmd = &cobra.Command{
	Use:   "bucket-put [bucket] [key1=value1] [key2=value2]...",
	Short: "设置存储桶标签",
	Long:  "为存储桶设置标签",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		tagArgs := args[1:]

		var tags []types.Tag
		for _, tagArg := range tagArgs {
			parts := splitTag(tagArg)
			if len(parts) == 2 {
				tags = append(tags, types.Tag{
					Key:   aws.String(parts[0]),
					Value: aws.String(parts[1]),
				})
			}
		}

		if len(tags) == 0 {
			return fmt.Errorf("没有有效的标签")
		}

		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.PutBucketTagging(ctx, &s3.PutBucketTaggingInput{
			Bucket: aws.String(bucket),
			Tagging: &types.Tagging{
				TagSet: tags,
			},
		})
		if err != nil {
			return fmt.Errorf("设置存储桶标签失败: %w", err)
		}

		cmd.Printf("存储桶 %s 标签设置成功\n", bucket)
		return nil
	},
}

var tagsBucketDeleteCmd = &cobra.Command{
	Use:   "bucket-delete [bucket]",
	Short: "删除存储桶标签",
	Long:  "删除指定存储桶的所有标签",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.DeleteBucketTagging(ctx, &s3.DeleteBucketTaggingInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("删除存储桶标签失败: %w", err)
		}

		cmd.Printf("存储桶 %s 标签已删除\n", bucket)
		return nil
	},
}

var tagsObjectGetCmd = &cobra.Command{
	Use:   "object-get [bucket] [key]",
	Short: "获取对象标签",
	Long:  "获取指定对象的标签",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key := args[0], args[1]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("获取对象标签失败: %w", err)
		}

		if len(output.TagSet) == 0 {
			cmd.Println("该对象没有标签")
			return nil
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"键", "值"})

		for _, tag := range output.TagSet {
			t.AppendRow([]interface{}{*tag.Key, *tag.Value})
		}

		cmd.Println(t.Render())
		return nil
	},
}

func splitTag(tag string) []string {
	for i := 0; i < len(tag); i++ {
		if tag[i] == '=' {
			return []string{tag[:i], tag[i+1:]}
		}
	}
	return nil
}

func init() {
	tagsCmd.AddCommand(tagsBucketGetCmd, tagsBucketPutCmd, tagsBucketDeleteCmd, tagsObjectGetCmd)
}
