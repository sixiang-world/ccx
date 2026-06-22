package cmd

import (
	"fmt"

	"ccx-cli/internal/client"
	"ccx-cli/internal/models"
	"ccx-cli/internal/validator"

	"github.com/spf13/cobra"
)

var (
	createBaseURL      string
	createAPIKeys      []string
	createServiceType  string
	createStatus       string
	createPriority     int
	createDescription  string
	createWebsite      string
	createModelMapping []string
	createProxyURL     string
	createAuthHeader   string
)

// channelCreateCmd 表示 channel create 子命令
var channelCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "创建新渠道",
	Long: `创建指定类型（--type）的新上游渠道。

示例：
  ccx channel create my-channel --base-url https://api.anthropic.com --service-type claude
  ccx channel create my-channel --type chat --base-url https://api.openai.com --api-key sk-xxx
  ccx channel create my-channel --type responses --base-url https://api.openai.com --api-key sk-xxx --priority 1
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		chType := getChannelType()

		// 校验必要参数
		if createBaseURL == "" {
			return fmt.Errorf("--base-url/-u 是必填参数")
		}
		if err := validator.ValidateURL(createBaseURL); err != nil {
			return fmt.Errorf("无效的 --base-url：%v", err)
		}
		if createServiceType != "" {
			if err := validator.ValidateServiceType(createServiceType); err != nil {
				return err
			}
		}
		if createStatus != "" {
			if err := validator.ValidateChannelStatus(createStatus); err != nil {
				return err
			}
		}

		// 校验 API Key 格式
		for _, key := range createAPIKeys {
			if err := validator.ValidateAPIKey(key); err != nil {
				return fmt.Errorf("无效的 API Key %q：%v", maskKey(key), err)
			}
		}

		// 构建请求体
		upstream := models.UpstreamConfig{
			Name:        name,
			BaseURL:     createBaseURL,
			APIKeys:     createAPIKeys,
			ServiceType: createServiceType,
			Status:      createStatus,
			Priority:    createPriority,
			Description: createDescription,
			Website:     createWebsite,
			ProxyURL:    createProxyURL,
			AuthHeader:  createAuthHeader,
		}

		// 解析 model-mapping（key=value 格式）
		if len(createModelMapping) > 0 {
			mapping, err := validator.ParseModelMapping(createModelMapping)
			if err != nil {
				return err
			}
			upstream.ModelMapping = mapping
		}

		c := NewClient()
		resp, err := c.Post(client.ChannelPath(chType), upstream)
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("渠道 %q 已创建", name)
		return nil
	},
}

func init() {
	channelCreateCmd.Flags().StringVarP(&createBaseURL, "base-url", "u", "", "上游 BaseURL（必填）")
	channelCreateCmd.Flags().StringArrayVarP(&createAPIKeys, "api-key", "a", nil, "API Key（可重复指定）")
	channelCreateCmd.Flags().StringVar(&createServiceType, "service-type", "", "服务类型（claude|openai|gemini|custom）")
	channelCreateCmd.Flags().StringVar(&createStatus, "status", "", "初始状态（active|suspended|disabled）")
	channelCreateCmd.Flags().IntVar(&createPriority, "priority", 0, "优先级（数字越小优先级越高）")
	channelCreateCmd.Flags().StringVar(&createDescription, "description", "", "描述信息")
	channelCreateCmd.Flags().StringVar(&createWebsite, "website", "", "网站")
	channelCreateCmd.Flags().StringArrayVar(&createModelMapping, "model-mapping", nil, "模型映射（key=value 格式，可重复）")
	channelCreateCmd.Flags().StringVar(&createProxyURL, "proxy-url", "", "代理地址")
	channelCreateCmd.Flags().StringVar(&createAuthHeader, "auth-header", "", "认证头覆盖（auto|bearer|x-api-key）")
	channelCmd.AddCommand(channelCreateCmd)
}
