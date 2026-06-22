package cmd

import (
	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// pingCmd 表示 ping 子命令
var pingCmd = &cobra.Command{
	Use:   "ping [name]",
	Short: "全局连通性检测",
	Long: `检测所有渠道或指定渠道的连通性。

不传参数时检测所有渠道，传渠道名称时检测单个渠道。

示例：
  ccx ping                          # 检测所有 messages 类型渠道
  ccx ping --type chat              # 检测所有 chat 类型渠道
  ccx ping my-channel               # 检测单个渠道
  ccx ping --output json            # JSON 格式输出
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chType := getChannelType()
		c := NewClient()

		if len(args) > 0 {
			// 单个渠道 ping
			name := args[0]
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
		}

		// 所有渠道 ping
		resp, err := c.Get(client.ChannelPingAllPath(chType), nil)
		if err != nil {
			return err
		}

		var rawResult any
		if err := client.DecodeResponse(resp, &rawResult); err != nil {
			return err
		}

		PrintOutput(rawResult)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(pingCmd)
}
