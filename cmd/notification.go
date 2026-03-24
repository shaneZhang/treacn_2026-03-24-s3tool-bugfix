package cmd

import (
	"context"
	"fmt"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "事件通知配置",
	Long:  "管理存储桶的事件通知配置",
}

var notificationGetCmd = &cobra.Command{
	Use:   "get [bucket]",
	Short: "获取通知配置",
	Long:  "获取指定存储桶的事件通知配置",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetBucketNotificationConfiguration(ctx, &s3.GetBucketNotificationConfigurationInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("获取通知配置失败: %w", err)
		}

		hasConfig := false

		if len(output.TopicConfigurations) > 0 {
			hasConfig = true
			cmd.Println("SNS 主题配置:")
			for _, config := range output.TopicConfigurations {
				cmd.Printf("  - ID: %s\n", *config.Id)
				cmd.Printf("    主题ARN: %s\n", *config.TopicArn)
				cmd.Printf("    事件: %v\n", config.Events)
			}
		}

		if len(output.QueueConfigurations) > 0 {
			hasConfig = true
			cmd.Println("SQS 队列配置:")
			for _, config := range output.QueueConfigurations {
				cmd.Printf("  - ID: %s\n", *config.Id)
				cmd.Printf("    队列ARN: %s\n", *config.QueueArn)
				cmd.Printf("    事件: %v\n", config.Events)
			}
		}

		if len(output.LambdaFunctionConfigurations) > 0 {
			hasConfig = true
			cmd.Println("Lambda 函数配置:")
			for _, config := range output.LambdaFunctionConfigurations {
				cmd.Printf("  - ID: %s\n", *config.Id)
				cmd.Printf("    函数ARN: %s\n", *config.LambdaFunctionArn)
				cmd.Printf("    事件: %v\n", config.Events)
			}
		}

		if !hasConfig {
			cmd.Println("该存储桶没有配置事件通知")
		}
		return nil
	},
}

func init() {
	notificationCmd.AddCommand(notificationGetCmd)
}
