package cmd

import (
	"context"
	"fmt"
	"os"

	"s3tool/config"

	"github.com/spf13/cobra"
)

var (
	configFile string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "s3tool",
	Short: "S3 命令行管理工具",
	Long: `一个功能强大的 S3 命令行管理工具，支持存储桶管理、对象操作、
预签名URL、多部分上传、版本控制、生命周期规则等功能。`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if configFile == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("无法获取用户主目录: %w", err)
			}
			configFile = homeDir + "/.s3tool.yaml"
		}

		if err := config.LoadConfig(configFile); err != nil {
			return fmt.Errorf("加载配置文件失败: %w", err)
		}

		ctx := context.Background()
		_, err := config.GetS3Client(ctx)
		if err != nil {
			return fmt.Errorf("初始化 S3 客户端失败: %w", err)
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件路径 (默认: ~/.s3tool.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "详细输出模式")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(bucketCmd)
	rootCmd.AddCommand(objectCmd)
	rootCmd.AddCommand(presignCmd)
	rootCmd.AddCommand(multipartCmd)
	rootCmd.AddCommand(policyCmd)
	rootCmd.AddCommand(lifecycleCmd)
	rootCmd.AddCommand(versioningCmd)
	rootCmd.AddCommand(tagsCmd)
	rootCmd.AddCommand(websiteCmd)
	rootCmd.AddCommand(corsCmd)
	rootCmd.AddCommand(aclCmd)
	rootCmd.AddCommand(loggingCmd)
	rootCmd.AddCommand(encryptionCmd)
	rootCmd.AddCommand(replicationCmd)
	rootCmd.AddCommand(notificationCmd)
}
