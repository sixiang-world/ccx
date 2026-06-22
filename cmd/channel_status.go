package cmd

import (
	"ccx-cli/internal/client"
	"ccx-cli/internal/models"
	"ccx-cli/internal/validator"

	"github.com/spf13/cobra"
)

// channelStatusCmd 表示 channel status 子命令
var channelStatusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st"},
	Short:   "管理渠道状态",
	Long: `设置渠道的运行状态。

状态值：
  - active：正常（默认）
  - suspended：暂停（不参与调度）
  - disabled：禁用（放入备用池）

示例：
  ccx channel status set my-channel active
  ccx channel status set my-channel suspended
  ccx channel status set my-channel disabled
`,
}

// channelStatusSetCmd 设置渠道状态
var channelStatusSetCmd = &cobra.Command{
	Use:   "set <name> <active|suspended|disabled>",
	Short: "设置渠道状态",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, status := args[0], args[1]
		chType := getChannelType()

		if err := validator.ValidateChannelStatus(status); err != nil {
			return err
		}

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		body := models.StatusUpdateRequest{Status: status}
		resp, err := c.Patch(client.ChannelStatusPath(chType, id), body)
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("渠道 %q 状态已设为 %s", name, status)
		return nil
	},
}

func init() {
	channelStatusCmd.AddCommand(channelStatusSetCmd)
	channelCmd.AddCommand(channelStatusCmd)
}
