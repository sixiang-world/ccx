package cmd

import (
	"fmt"

	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// channelDeleteCmd 表示 channel delete 子命令
var channelDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "删除渠道",
	Long: `删除指定名称的渠道。

示例：
  ccx channel delete my-channel
  ccx channel delete my-channel --type chat
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		chType := getChannelType()

		c := NewClient()

		// 获取渠道 ID
		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		// 确认
		if !GetYesFlag() {
			fmt.Printf("⚠ 确定要删除渠道 %q（索引 %d）吗？(y/N): ", name, id)
			var confirm string
			if _, err := fmt.Scanln(&confirm); err != nil || (confirm != "y" && confirm != "Y") {
				fmt.Println("已取消")
				return nil
			}
		}

		resp, err := c.Delete(client.ChannelIDPath(chType, id))
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("渠道 %q 已删除", name)
		return nil
	},
}

func init() {
	channelCmd.AddCommand(channelDeleteCmd)
}
