package cmd

import (
	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// channelDashboardCmd 表示 channel dashboard 子命令
var channelDashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "查看统一仪表盘",
	Long: `查看渠道统一仪表盘（聚合展示所有渠道类型的概览信息）。

示例：
  ccx channel dashboard
  ccx channel dashboard --type messages
  ccx channel dashboard -o json
`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		chType := getChannelType()

		c := NewClient()

		// dashboard 端点总是指向 /api/messages/channels/dashboard
		// 并通过 ?type= 查询参数区分渠道类型
		queryParams := map[string]string{}
		if chType != "" {
			queryParams["type"] = chType
		}

		resp, err := c.Get(client.ChannelDashboardPath(), queryParams)
		if err != nil {
			return err
		}

		var result models.ChannelDashboardResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintOutput(result)
		return nil
	},
}

func init() {
	channelCmd.AddCommand(channelDashboardCmd)
}
