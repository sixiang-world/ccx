package cmd

import (
	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// channelSchedulerStatsCmd 表示 channel scheduler-stats 子命令
var channelSchedulerStatsCmd = &cobra.Command{
	Use:     "scheduler-stats",
	Aliases: []string{"sched-stats"},
	Short:   "查看调度器统计信息",
	Long: `查看渠道调度器的统计信息（如调度次数、失败率等）。
调度器统计仅 messages 和 chat 渠道类型支持。

示例：
  ccx channel scheduler-stats
  ccx channel scheduler-stats --type chat
  ccx channel scheduler-stats -o json
`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		chType := getChannelType()

		c := NewClient()

		queryParams := map[string]string{}
		if chType != "" {
			queryParams["type"] = chType
		}

		resp, err := c.Get(client.ChannelSchedulerStatsPath(), queryParams)
		if err != nil {
			return err
		}

		var result models.SchedulerStatsResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintOutput(result)
		return nil
	},
}

func init() {
	channelCmd.AddCommand(channelSchedulerStatsCmd)
}
