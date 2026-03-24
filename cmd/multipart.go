package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var multipartCmd = &cobra.Command{
	Use:   "multipart",
	Short: "多部分上传",
	Long:  "多部分上传相关操作",
}

var multipartInitCmd = &cobra.Command{
	Use:   "init [bucket] [key]",
	Short: "初始化多部分上传",
	Long:  "初始化一个大文件的多部分上传",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key := args[0], args[1]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		resp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("初始化多部分上传失败: %w", err)
		}

		cmd.Printf("UploadId: %s\n", *resp.UploadId)
		cmd.Println("请保存此 UploadId 用于后续操作")
		return nil
	},
}

var multipartUploadCmd = &cobra.Command{
	Use:   "upload [bucket] [key] [upload-id] [part-number] [file]",
	Short: "上传分片",
	Long:  "上传多部分上传的一个分片",
	Args:  cobra.ExactArgs(5),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key, uploadID, filePath := args[0], args[1], args[2], args[4]
		partNum, err := strconv.Atoi(args[3])
		if err != nil {
			return fmt.Errorf("分片编号必须是数字: %w", err)
		}

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

		partSize, _ := cmd.Flags().GetInt("part-size")

		// 创建指定大小的缓冲区并读取文件内容
		buf := make([]byte, partSize)
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("读取文件失败: %w", err)
		}

		// 只上传实际读取到的字节数
		resp, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(bucket),
			Key:        aws.String(key),
			UploadId:   aws.String(uploadID),
			PartNumber: aws.Int32(int32(partNum)),
			Body:       bytes.NewReader(buf[:n]),
		})
		if err != nil {
			return fmt.Errorf("上传分片失败: %w", err)
		}

		cmd.Printf("分片 %d 上传成功, ETag: %s\n", partNum, *resp.ETag)
		return nil
	},
}

var multipartListCmd = &cobra.Command{
	Use:   "list [bucket] [key] [upload-id]",
	Short: "列出已上传的分片",
	Long:  "列出多部分上传的所有已上传分片",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key, uploadID := args[0], args[1], args[2]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		resp, err := client.ListParts(ctx, &s3.ListPartsInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String(key),
			UploadId: aws.String(uploadID),
		})
		if err != nil {
			return fmt.Errorf("列出分片失败: %w", err)
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"分片编号", "ETag", "大小"})

		for _, part := range resp.Parts {
			t.AppendRow([]interface{}{
				*part.PartNumber,
				*part.ETag,
				formatBytes(*part.Size),
			})
		}

		cmd.Println(t.Render())
		return nil
	},
}

var multipartCompleteCmd = &cobra.Command{
	Use:   "complete [bucket] [key] [upload-id]",
	Short: "完成多部分上传",
	Long:  "完成多部分上传并合并所有分片",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key, uploadID := args[0], args[1], args[2]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		listResp, err := client.ListParts(ctx, &s3.ListPartsInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String(key),
			UploadId: aws.String(uploadID),
		})
		if err != nil {
			return fmt.Errorf("获取分片列表失败: %w", err)
		}

		var parts []types.CompletedPart
		for _, part := range listResp.Parts {
			parts = append(parts, types.CompletedPart{
				ETag:       part.ETag,
				PartNumber: part.PartNumber,
			})
		}

		_, err = client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String(key),
			UploadId: aws.String(uploadID),
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: parts,
			},
		})
		if err != nil {
			return fmt.Errorf("完成多部分上传失败: %w", err)
		}

		cmd.Printf("多部分上传完成: s3://%s/%s\n", bucket, key)
		return nil
	},
}

var multipartAbortCmd = &cobra.Command{
	Use:   "abort [bucket] [key] [upload-id]",
	Short: "中止多部分上传",
	Long:  "中止一个进行中的多部分上传",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key, uploadID := args[0], args[1], args[2]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String(key),
			UploadId: aws.String(uploadID),
		})
		if err != nil {
			return fmt.Errorf("中止多部分上传失败: %w", err)
		}

		cmd.Printf("多部分上传已中止: %s\n", uploadID)
		return nil
	},
}

func init() {
	multipartCmd.AddCommand(multipartInitCmd, multipartUploadCmd, multipartListCmd, multipartCompleteCmd, multipartAbortCmd)
	multipartUploadCmd.Flags().Int("part-size", 5*1024*1024, "分片大小(字节)")
}
