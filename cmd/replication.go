package cmd

import (
	"context"
	"fmt"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

var replicationCmd = &cobra.Command{
	Use:   "replication",
	Short: "复制配置管理",
	Long:  "管理存储桶的跨区域复制配置",
}

var replicationGetCmd = &cobra.Command{
	Use:   "get [bucket]",
	Short: "获取复制配置",
	Long:  "获取指定存储桶的复制配置",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetBucketReplication(ctx, &s3.GetBucketReplicationInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("获取复制配置失败: %w", err)
		}

		if output.ReplicationConfiguration == nil || len(output.ReplicationConfiguration.Rules) == 0 {
			cmd.Println("该存储桶没有配置复制规则")
			return nil
		}

		cmd.Printf("角色ARN: %s\n", *output.ReplicationConfiguration.Role)
		for i, rule := range output.ReplicationConfiguration.Rules {
			cmd.Printf("规则 %d:\n", i+1)
			if rule.ID != nil {
				cmd.Printf("  ID: %s\n", *rule.ID)
			}
			cmd.Printf("  状态: %s\n", rule.Status)
			if rule.Prefix != nil {
				cmd.Printf("  前缀: %s\n", *rule.Prefix)
			}
			if dest := rule.Destination; dest != nil && dest.Bucket != nil {
				cmd.Printf("  目标存储桶: %s\n", *dest.Bucket)
			}
		}
		return nil
	},
}

var replicationDeleteCmd = &cobra.Command{
	Use:   "delete [bucket]",
	Short: "删除复制配置",
	Long:  "删除指定存储桶的复制配置",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.DeleteBucketReplication(ctx, &s3.DeleteBucketReplicationInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("删除复制配置失败: %w", err)
		}

		cmd.Printf("存储桶 %s 复制配置已删除\n", bucket)
		return nil
	},
}

func init() {
	replicationCmd.AddCommand(replicationGetCmd, replicationDeleteCmd)
}
