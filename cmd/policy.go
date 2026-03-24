package cmd

import (
	"context"
	"fmt"
	"os"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "存储桶策略管理",
	Long:  "管理存储桶的访问策略",
}

var policyGetCmd = &cobra.Command{
	Use:   "get [bucket]",
	Short: "获取存储桶策略",
	Long:  "获取指定存储桶的访问策略",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("获取存储桶策略失败: %w", err)
		}

		if output.Policy != nil {
			cmd.Println(*output.Policy)
		} else {
			cmd.Println("该存储桶没有设置策略")
		}
		return nil
	},
}

var policySetCmd = &cobra.Command{
	Use:   "set [bucket] [policy-file]",
	Short: "设置存储桶策略",
	Long:  "从JSON文件设置存储桶的访问策略",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, policyFile := args[0], args[1]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		policyContent, err := os.ReadFile(policyFile)
		if err != nil {
			return fmt.Errorf("读取策略文件失败: %w", err)
		}

		_, err = client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
			Bucket: aws.String(bucket),
			Policy: aws.String(string(policyContent)),
		})
		if err != nil {
			return fmt.Errorf("设置存储桶策略失败: %w", err)
		}

		cmd.Printf("存储桶 %s 策略设置成功\n", bucket)
		return nil
	},
}

var policyDeleteCmd = &cobra.Command{
	Use:   "delete [bucket]",
	Short: "删除存储桶策略",
	Long:  "删除指定存储桶的访问策略",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.DeleteBucketPolicy(ctx, &s3.DeleteBucketPolicyInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("删除存储桶策略失败: %w", err)
		}

		cmd.Printf("存储桶 %s 策略已删除\n", bucket)
		return nil
	},
}

func init() {
	policyCmd.AddCommand(policyGetCmd, policySetCmd, policyDeleteCmd)
}
