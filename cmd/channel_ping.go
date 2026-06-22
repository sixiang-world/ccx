package cmd

import (
	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// channelPingCmd 表示 channel ping 子命令
var channelPingCmd = &cobra.Command{
	Use:   "ping <name>",
	Short: "测试渠道连通性",
	Long: `测试指定渠道的连通性（延迟检测）。

示例：
  ccx channel ping my-channel
  ccx channel ping my-channel --type responses
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		chType := getChannelType()

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		resp, err := c.Get(client.ChannelPingPath(chType, id), nil)
		if err != nil {
			return err
		}

		var result models.PingResult
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintOutput(result)
		return nil
	},
}

func init() {
	channelCmd.AddCommand(channelPingCmd)
}
