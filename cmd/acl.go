package cmd

import (
	"context"
	"fmt"

	"s3tool/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/cobra"
)

var aclCmd = &cobra.Command{
	Use:   "acl",
	Short: "访问控制列表管理",
	Long:  "管理存储桶和对象的访问控制列表(ACL)",
}

var aclBucketGetCmd = &cobra.Command{
	Use:   "bucket-get [bucket]",
	Short: "获取存储桶ACL",
	Long:  "获取指定存储桶的ACL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("获取存储桶ACL失败: %w", err)
		}

		// 安全地获取所有者信息
		ownerDisplayName := "未知"
		ownerID := "未知"
		if output.Owner != nil {
			if output.Owner.DisplayName != nil {
				ownerDisplayName = *output.Owner.DisplayName
			}
			if output.Owner.ID != nil {
				ownerID = *output.Owner.ID
			}
		}
		cmd.Printf("所有者: %s (%s)\n", ownerDisplayName, ownerID)
		cmd.Println("授权:")
		for _, grant := range output.Grants {
			grantee := grant.Grantee
			granteeDisplayName := "未知"
			granteeID := "未知"
			if grantee != nil {
				if grantee.DisplayName != nil {
					granteeDisplayName = *grantee.DisplayName
				}
				if grantee.ID != nil {
					granteeID = *grantee.ID
				}
				if grantee.URI != nil {
					granteeDisplayName = *grantee.URI
				}
			}
			cmd.Printf("  - %s: %s (%s)\n", grant.Permission, granteeDisplayName, granteeID)
		}
		return nil
	},
}

var aclBucketSetCmd = &cobra.Command{
	Use:   "bucket-set [bucket] [acl]",
	Short: "设置存储桶ACL",
	Long:  "设置存储桶的ACL (private | public-read | public-read-write | authenticated-read)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, acl := args[0], args[1]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.PutBucketAcl(ctx, &s3.PutBucketAclInput{
			Bucket: aws.String(bucket),
			ACL:    types.BucketCannedACL(acl),
		})
		if err != nil {
			return fmt.Errorf("设置存储桶ACL失败: %w", err)
		}

		cmd.Printf("存储桶 %s ACL设置为 %s\n", bucket, acl)
		return nil
	},
}

var aclObjectGetCmd = &cobra.Command{
	Use:   "object-get [bucket] [key]",
	Short: "获取对象ACL",
	Long:  "获取指定对象的ACL",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key := args[0], args[1]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		output, err := client.GetObjectAcl(ctx, &s3.GetObjectAclInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("获取对象ACL失败: %w", err)
		}

		// 安全地获取所有者信息
		ownerDisplayName := "未知"
		ownerID := "未知"
		if output.Owner != nil {
			if output.Owner.DisplayName != nil {
				ownerDisplayName = *output.Owner.DisplayName
			}
			if output.Owner.ID != nil {
				ownerID = *output.Owner.ID
			}
		}
		cmd.Printf("所有者: %s (%s)\n", ownerDisplayName, ownerID)
		cmd.Println("授权:")
		for _, grant := range output.Grants {
			grantee := grant.Grantee
			granteeDisplayName := "未知"
			granteeID := "未知"
			if grantee != nil {
				if grantee.DisplayName != nil {
					granteeDisplayName = *grantee.DisplayName
				}
				if grantee.ID != nil {
					granteeID = *grantee.ID
				}
				if grantee.URI != nil {
					granteeDisplayName = *grantee.URI
				}
			}
			cmd.Printf("  - %s: %s (%s)\n", grant.Permission, granteeDisplayName, granteeID)
		}
		return nil
	},
}

var aclObjectSetCmd = &cobra.Command{
	Use:   "object-set [bucket] [key] [acl]",
	Short: "设置对象ACL",
	Long:  "设置对象的ACL (private | public-read | public-read-write | authenticated-read | bucket-owner-read | bucket-owner-full-control)",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, key, acl := args[0], args[1], args[2]
		ctx := context.Background()
		client, err := config.GetS3Client(ctx)
		if err != nil {
			return err
		}

		_, err = client.PutObjectAcl(ctx, &s3.PutObjectAclInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			ACL:    types.ObjectCannedACL(acl),
		})
		if err != nil {
			return fmt.Errorf("设置对象ACL失败: %w", err)
		}

		cmd.Printf("对象 s3://%s/%s ACL设置为 %s\n", bucket, key, acl)
		return nil
	},
}

func init() {
	aclCmd.AddCommand(aclBucketGetCmd, aclBucketSetCmd, aclObjectGetCmd, aclObjectSetCmd)
}
