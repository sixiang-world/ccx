package cmd

import (
	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

var channelMetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "查看渠道性能指标",
	Long: `查看渠道的性能指标和统计数据。

示例：
  ccx channel metrics                 # 查看 messages 渠道指标
  ccx channel metrics --type chat     # 查看 chat 渠道指标
  ccx channel metrics -o json         # JSON 格式输出
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := NewClient()
		chType := getChannelType()

		resp, err := c.Get(client.ChannelMetricsPath(chType), nil)
		if err != nil {
			return err
		}
		var result models.ChannelMetricsResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}
		PrintOutput(result)
		return nil
	},
}

func init() {
	channelCmd.AddCommand(channelMetricsCmd)
}
