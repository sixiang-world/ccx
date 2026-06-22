package cmd

import (
	"fmt"

	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// channelLogsCmd 表示 channel logs 子命令
var channelLogsCmd = &cobra.Command{
	Use:   "logs <name>",
	Short: "查看渠道请求日志",
	Long: `查看指定渠道的请求日志。

示例：
  ccx channel logs my-channel
  ccx channel logs my-channel --type chat
  ccx channel logs my-channel -o json
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

		resp, err := c.Get(client.ChannelLogsPath(chType, id), nil)
		if err != nil {
			return err
		}

		var result models.ChannelLogsResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		if len(result) == 0 {
			fmt.Println("（无日志记录）")
			return nil
		}
		PrintOutput(result)
		return nil
	},
}

func init() {
	channelCmd.AddCommand(channelLogsCmd)
}
