package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

var reorderOrder string

// channelReorderCmd 表示 channel reorder 子命令
var channelReorderCmd = &cobra.Command{
	Use:   "reorder",
	Short: "重新排序渠道",
	Long: `按指定顺序重新排列渠道。渠道顺序影响故障转移时的优先级。

示例：
  ccx channel reorder --order ch1,ch2,ch3
  ccx channel reorder --order 2,0,1
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chType := getChannelType()

		if reorderOrder == "" {
			return fmt.Errorf("--order 是必填参数")
		}

		c := NewClient()

		// 获取渠道列表
		resp, err := c.Get(client.ChannelPath(chType), nil)
		if err != nil {
			return err
		}

		var list models.ChannelListResponse
		if err := client.DecodeResponse(resp, &list); err != nil {
			return err
		}

		// 构建名称到索引的映射
		nameToIndex := make(map[string]int)
		for _, ch := range list.Channels {
			nameToIndex[ch.Name] = ch.Index
		}

		// 解析 order 参数
		parts := strings.Split(reorderOrder, ",")
		order := make([]int, 0, len(parts))

		for _, part := range parts {
			part = strings.TrimSpace(part)
			// 尝试按名称查找
			if idx, ok := nameToIndex[part]; ok {
				order = append(order, idx)
				continue
			}
			// 尝试按索引解析
			idx, err := strconv.Atoi(part)
			if err != nil {
				return fmt.Errorf("无效的排序值 %q（应为渠道名称或索引数字）", part)
			}
			order = append(order, idx)
		}

		body := models.ReorderRequest{Order: order}
		resp2, err := c.Post(client.ChannelReorderPath(chType), body)
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp2, &result); err != nil {
			return err
		}

		PrintSuccess("渠道顺序已更新")
		return nil
	},
}

func init() {
	channelReorderCmd.Flags().StringVar(&reorderOrder, "order", "", "渠道顺序（逗号分隔的名称或索引）")
	channelCmd.AddCommand(channelReorderCmd)
}
