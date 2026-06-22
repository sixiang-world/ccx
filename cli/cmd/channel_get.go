package cmd

import (
	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// channelGetCmd 表示 channel get 子命令
var channelGetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "获取渠道详情",
	Long: `获取指定名称的渠道详情。

示例：
  ccx channel get my-channel
  ccx channel get my-channel --type chat
  ccx channel get my-channel --output json
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		c := NewClient()
		chType := getChannelType()

		// 获取渠道列表并在本地过滤
		resp, err := c.Get(client.ChannelPath(chType), nil)
		if err != nil {
			return err
		}

		var list models.ChannelListResponse
		if err := client.DecodeResponse(resp, &list); err != nil {
			return err
		}

		// 本地查找
		for _, ch := range list.Channels {
			if ch.Name == name {
				// 脱敏密钥
				if !GetShowKeys() {
					masked := make([]string, len(ch.APIKeys))
					for i, k := range ch.APIKeys {
						masked[i] = maskKey(k)
					}
					ch.APIKeys = masked
				}
				PrintOutput(ch)
				return nil
			}
		}

		return notFoundError(name)
	},
}

func init() {
	channelCmd.AddCommand(channelGetCmd)
}
