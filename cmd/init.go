package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化配置文件",
	Long:  "创建默认的配置文件模板",
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("无法获取用户主目录: %w", err)
		}

		configPath := homeDir + "/.s3tool.yaml"
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("配置文件已存在: %s\n", configPath)
			return nil
		}

		configContent := `# S3 工具配置文件
# 配置说明:
# - region: AWS 区域
# - profile: AWS 凭证 profile
# - access_key: AWS Access Key ID
# - secret_key: AWS Secret Access Key
# - endpoint: 自定义 S3 兼容端点 (可选, 用于 MinIO 等)
# - force_path_style: 强制使用路径样式 (MinIO 需要开启)
# - use_accelerate: 启用 S3 Transfer Acceleration

region: us-east-1
profile: ""
access_key: ""
secret_key: ""
endpoint: ""
force_path_style: false
use_accelerate: false
`
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			return fmt.Errorf("创建配置文件失败: %w", err)
		}

		fmt.Printf("配置文件已创建: %s\n", configPath)
		fmt.Println("请编辑配置文件并填入您的 AWS 凭证")
		return nil
	},
}
