package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"ccx-cli/internal/client"
	"ccx-cli/internal/diff"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// configCmd 表示 config 子命令
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "管理全局配置",
	Long: `查看、应用、保存、备份和恢复 CCX 全局配置。

子命令：
  show, apply, save, backup, restore

示例：
  ccx config show
  ccx config apply config.json
  ccx config save
  ccx config backup
  ccx config restore backup.json
`,
}

// configShowCmd 显示全局配置
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "显示全局配置聚合",
	Long: `聚合显示所有类型的渠道配置和全局设置。

示例：
  ccx config show
  ccx config show --output json
  ccx config show --output yaml
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := NewClient()
		cfg, err := fetchFullConfig(c)
		if err != nil {
			return err
		}
		PrintOutput(cfg)
		return nil
	},
}

// configApplyCmd 应用配置
var configApplyCmd = &cobra.Command{
	Use:   "apply <file>",
	Short: "应用配置文件（diff + 确认 + 逐项提交）",
	Long: `读取本地配置文件，与服务器当前配置比较差异，
展示变更预览后逐项提交 API 变更。

示例：
  ccx config apply config.json
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		c := NewClient()

		// 读取本地配置文件
		localData, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("读取配置文件 %s 失败：%w", filePath, err)
		}

		// 获取服务器当前配置
		currentCfg, err := fetchFullConfig(c)
		if err != nil {
			return err
		}
		currentData, err := json.MarshalIndent(currentCfg, "", "  ")
		if err != nil {
			return fmt.Errorf("序列化当前配置失败：%w", err)
		}

		// 计算差异
		d, err := diff.CompareJSONBytes(currentData, localData)
		if err != nil {
			return fmt.Errorf("差异比较失败：%w", err)
		}

		if !d.HasChanges {
			fmt.Println("✓ 配置无变更")
			return nil
		}

		fmt.Printf("配置变更预览（%s）：\n", d.Summary())
		// 展示变更列表
		for _, change := range d.Changes {
			fmt.Println("  " + change.String())
		}
		fmt.Println()

		// 确认
		if !GetYesFlag() {
			fmt.Print("确定应用以上变更？(y/N): ")
			var confirm string
			if _, err := fmt.Scanln(&confirm); err != nil || (confirm != "y" && confirm != "Y") {
				fmt.Println("已取消")
				return nil
			}
		}

		// 发送本地配置到服务器
		resp, err := c.Post(client.ConfigSavePath(), json.RawMessage(localData))
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("配置已保存")
		return nil
	},
}

// configSaveCmd 保存配置
var configSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "强制持久化运行时配置到磁盘",
	Long: `强制将当前运行时配置持久化到磁盘上的 config.json 文件。

示例：
  ccx config save
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := NewClient()

		resp, err := c.Post(client.ConfigSavePath(), nil)
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("配置已持久化到磁盘")
		return nil
	},
}

// configBackupCmd 备份配置
var configBackupCmd = &cobra.Command{
	Use:   "backup [output-file]",
	Short: "下载完整配置备份到本地",
	Long: `下载服务器的完整配置到本地 JSON 文件。

示例：
  ccx config backup                    # 输出到 ccx-backup-{timestamp}.json
  ccx config backup my-backup.json
`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := NewClient()

		cfg, err := fetchFullConfig(c)
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return fmt.Errorf("序列化配置失败：%w", err)
		}

		outputFile := "ccx-backup.json"
		if len(args) > 0 {
			outputFile = args[0]
		}

		if err := os.WriteFile(outputFile, data, 0600); err != nil {
			return fmt.Errorf("写入备份文件 %s 失败：%w", outputFile, err)
		}

		PrintSuccess("配置已备份到 %s", outputFile)
		return nil
	},
}

// configRestoreCmd 恢复配置
var configRestoreCmd = &cobra.Command{
	Use:   "restore <backup-file>",
	Short: "从备份文件恢复配置",
	Long: `从本地备份文件恢复配置（复用 apply 流程）。

示例：
  ccx config restore ccx-backup.json
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 复用 config apply 的逻辑
		return configApplyCmd.RunE(cmd, args)
	},
}

// fetchFullConfig 获取完整配置
func fetchFullConfig(c *client.Client) (*models.FullConfig, error) {
	cfg := &models.FullConfig{}

	// 遍历所有类型获取渠道配置
	for _, chType := range models.AllChannelTypes {
		path := client.ChannelPath(string(chType))
		resp, err := c.Get(path, nil)
		if err != nil {
			return nil, fmt.Errorf("获取 %s 渠道列表失败：%w", chType, err)
		}

		var list models.ChannelListResponse
		if err := client.DecodeResponse(resp, &list); err != nil {
			return nil, err
		}

		// 转换为 UpstreamConfig
		upstreams := make([]models.UpstreamConfig, 0, len(list.Channels))
		for _, cv := range list.Channels {
			upstreams = append(upstreams, models.UpstreamConfig{
				Name:       cv.Name,
				BaseURL:    cv.BaseURL,
				BaseURLs:   cv.BaseURLs,
				APIKeys:    cv.APIKeys,
				ServiceType: cv.ServiceType,
				Status:     cv.Status,
				Priority:   cv.Priority,
				ModelMapping: cv.ModelMapping,
				ProxyURL:   cv.ProxyURL,
			})
		}

		switch chType {
		case models.ChannelTypeMessages:
			cfg.Upstream = upstreams
		case models.ChannelTypeResponses:
			cfg.ResponsesUpstream = upstreams
		case models.ChannelTypeGemini:
			cfg.GeminiUpstream = upstreams
		case models.ChannelTypeChat:
			cfg.ChatUpstream = upstreams
		case models.ChannelTypeImages:
			cfg.ImagesUpstream = upstreams
		}
	}

	return cfg, nil
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configApplyCmd)
	configCmd.AddCommand(configSaveCmd)
	configCmd.AddCommand(configBackupCmd)
	configCmd.AddCommand(configRestoreCmd)
	RootCmd.AddCommand(configCmd)
}
