package cmd

import (
	"sort"

	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// channelListCmd 表示 channel list 子命令
var channelListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有渠道",
	Long: `列出指定类型（--type）的所有上游渠道。

示例：
  ccx channel list                    # 列出 messages 类型渠道
  ccx channel list --type responses   # 列出 responses 类型渠道
  ccx channel list --type chat        # 列出 chat 类型渠道
  ccx channel list --output json      # JSON 格式输出
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := NewClient()
		chType := getChannelType()

		// 发送请求
		resp, err := c.Get(client.ChannelPath(chType), nil)
		if err != nil {
			return err
		}

		// 解析响应
		var result models.ChannelListResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		// 按索引排序
		sort.Slice(result.Channels, func(i, j int) bool {
			return result.Channels[i].Index < result.Channels[j].Index
		})

		// 转换为可输出的格式
		type channelRow struct {
			Index       int    `json:"index"`
			Name        string `json:"name"`
			ServiceType string `json:"serviceType"`
			Status      string `json:"status"`
			BaseURL     string `json:"baseUrl"`
			Priority    int    `json:"priority"`
			APIKeys     any    `json:"apiKeys,omitempty"`
			ModelCount  int    `json:"modelMappingCount"`
		}

		rows := make([]channelRow, 0, len(result.Channels))
		for _, ch := range result.Channels {
			apiKeys := ch.APIKeys
			if !GetShowKeys() {
				masked := make([]string, len(ch.APIKeys))
				for i, k := range ch.APIKeys {
					masked[i] = maskKey(k)
				}
				apiKeys = masked
			}
			rows = append(rows, channelRow{
				Index:       ch.Index,
				Name:        ch.Name,
				ServiceType: ch.ServiceType,
				Status:      ch.EffectiveState,
				BaseURL:     ch.BaseURL,
				Priority:    ch.Priority,
				APIKeys:     apiKeys,
				ModelCount:  len(ch.ModelMapping),
			})
		}

		PrintOutput(rows)
		return nil
	},
}

func init() {
	channelCmd.AddCommand(channelListCmd)
}
