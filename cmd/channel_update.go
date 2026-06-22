package cmd

import (
	"fmt"

	"ccx-cli/internal/client"
	"ccx-cli/internal/models"
	"ccx-cli/internal/validator"

	"github.com/spf13/cobra"
)

var (
	updateBaseURL      string
	updateAPIKeys      []string
	updateServiceType  string
	updateStatus       string
	updatePriority     int
	updateDescription  string
	updateWebsite      string
	updateModelMapping []string
	updateProxyURL     string
	updateName         string
	updateClearFields  []string
)

// channelUpdateCmd 表示 channel update 子命令
var channelUpdateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "更新渠道配置（merge 语义）",
	Long: `更新指定渠道的配置。使用 merge 语义，仅更新命令行指定的字段，
未指定的字段保持不变。

注意：--model-mapping 是全量替换操作，会覆盖渠道上所有已有映射。
如需增量添加/修改单个映射，使用 ccx channel mapping set。

示例：
  ccx channel update my-channel --base-url https://new-url.com
  ccx channel update my-channel --status suspended
  ccx channel update my-channel --priority 1 --description "主渠道"
  ccx channel update my-channel --api-key sk-new-key
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		chType := getChannelType()

		c := NewClient()

		// 获取渠道 ID
		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		// 校验参数
		if updateBaseURL != "" {
			if err := validator.ValidateURL(updateBaseURL); err != nil {
				return fmt.Errorf("无效的 --base-url：%v", err)
			}
		}
		if updateServiceType != "" {
			if err := validator.ValidateServiceType(updateServiceType); err != nil {
				return err
			}
		}
		if updateStatus != "" {
			if err := validator.ValidateChannelStatus(updateStatus); err != nil {
				return err
			}
		}
		for _, key := range updateAPIKeys {
			if err := validator.ValidateAPIKey(key); err != nil {
				return fmt.Errorf("无效的 API Key %q：%v", maskKey(key), err)
			}
		}

		// 构建更新请求体（仅设置需要变更的字段）
		update := map[string]any{}

		if updateBaseURL != "" {
			update["baseUrl"] = updateBaseURL
		}
		if len(updateAPIKeys) > 0 {
			update["apiKeys"] = updateAPIKeys
		}
		if updateServiceType != "" {
			update["serviceType"] = updateServiceType
		}
		if updateStatus != "" {
			update["status"] = updateStatus
		}
		if cmd.Flags().Changed("priority") {
			update["priority"] = updatePriority
		}
		if updateDescription != "" {
			update["description"] = updateDescription
		}
		if updateWebsite != "" {
			update["website"] = updateWebsite
		}
		if updateProxyURL != "" {
			update["proxyUrl"] = updateProxyURL
		}
		if len(updateModelMapping) > 0 {
			mapping, err := validator.ParseModelMapping(updateModelMapping)
			if err != nil {
				return err
			}
			update["modelMapping"] = mapping
		}

		if len(update) == 0 {
			return fmt.Errorf("请指定至少一个要更新的字段（如 --base-url）")
		}

		resp, err := c.Put(client.ChannelIDPath(chType, id), update)
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("渠道 %q 已更新", name)
		return nil
	},
}

func init() {
	channelUpdateCmd.Flags().StringVarP(&updateBaseURL, "base-url", "u", "", "上游 BaseURL")
	channelUpdateCmd.Flags().StringArrayVarP(&updateAPIKeys, "api-key", "a", nil, "API Key（可重复指定，全量替换）")
	channelUpdateCmd.Flags().StringVar(&updateServiceType, "service-type", "", "服务类型（claude|openai|gemini|custom）")
	channelUpdateCmd.Flags().StringVar(&updateStatus, "status", "", "状态（active|suspended|disabled）")
	channelUpdateCmd.Flags().IntVar(&updatePriority, "priority", 0, "优先级")
	channelUpdateCmd.Flags().StringVar(&updateDescription, "description", "", "描述信息")
	channelUpdateCmd.Flags().StringVar(&updateWebsite, "website", "", "网站")
	channelUpdateCmd.Flags().StringArrayVar(&updateModelMapping, "model-mapping", nil, "模型映射（key=value 格式，全量替换）")
	channelUpdateCmd.Flags().StringVar(&updateProxyURL, "proxy-url", "", "代理地址")
	channelCmd.AddCommand(channelUpdateCmd)
}
