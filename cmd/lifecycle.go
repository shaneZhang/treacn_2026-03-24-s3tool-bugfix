package cmd

import (
	"context"
	"fmt"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var lifecycleCmd = &cobra.Command{
	Use:   "lifecycle",
	Short: "生命周期规则管理",
	Long:  "管理存储桶的生命周期规则",
}

var lifecycleGetCmd = &cobra.Command{
	Use:   "get [bucket]",
	Short: "获取生命周期规则",
	Long:  "获取指定存储桶的生命周期规则",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetBucketLifecycleConfiguration(ctx, &s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("获取生命周期规则失败: %w", err)
		}

		if output.Rules == nil || len(output.Rules) == 0 {
			cmd.Println("该存储桶没有设置生命周期规则")
			return nil
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"ID", "状态", "前缀", "过期天数", "转换天数", "存储类型"})

		for _, rule := range output.Rules {
			expirationDays := ""
			transitionDays := ""
			storageClass := ""

			if rule.Expiration != nil && rule.Expiration.Days != nil {
				expirationDays = fmt.Sprintf("%d", *rule.Expiration.Days)
			}

			if len(rule.Transitions) > 0 {
				transitionDays = fmt.Sprintf("%d", *rule.Transitions[0].Days)
				storageClass = string(rule.Transitions[0].StorageClass)
			}

			t.AppendRow([]interface{}{
				*rule.ID,
				rule.Status,
				rule.Prefix,
				expirationDays,
				transitionDays,
				storageClass,
			})
		}

		cmd.Println(t.Render())
		return nil
	},
}

var lifecycleDeleteCmd = &cobra.Command{
	Use:   "delete [bucket]",
	Short: "删除生命周期规则",
	Long:  "删除指定存储桶的所有生命周期规则",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.DeleteBucketLifecycle(ctx, &s3.DeleteBucketLifecycleInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("删除生命周期规则失败: %w", err)
		}

		cmd.Printf("存储桶 %s 生命周期规则已删除\n", bucket)
		return nil
	},
}

func init() {
	lifecycleCmd.AddCommand(lifecycleGetCmd, lifecycleDeleteCmd)
}
