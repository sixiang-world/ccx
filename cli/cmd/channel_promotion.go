package cmd

import (
	"fmt"
	"strconv"

	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// channelPromotionCmd 表示 channel promotion 子命令
var channelPromotionCmd = &cobra.Command{
	Use:     "promotion",
	Aliases: []string{"promo"},
	Short:   "管理渠道促销期",
	Long: `设置或清除渠道的促销期。促销期内的渠道会被调度器优先选择。

示例：
  ccx channel promotion set my-channel 3600    # 促销1小时
  ccx channel promotion clear my-channel         # 清除促销期
`,
}

// channelPromotionSetCmd 设置促销期
var channelPromotionSetCmd = &cobra.Command{
	Use:   "set <name> <duration>",
	Short: "设置促销期（秒）",
	Long: `设置渠道的促销期。促销期内的渠道会被优先选择（忽略 trace 亲和性）。
duration 为秒数，<=0 时清除促销期。

示例：
  ccx channel promotion set my-channel 3600    # 促销1小时
  ccx channel promotion set my-channel 0        # 清除促销期（同 clear）
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, durationStr := args[0], args[1]
		chType := getChannelType()

		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			return fmt.Errorf("无效的持续时间 %q（应为秒数）", durationStr)
		}

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		body := models.PromotionRequest{Duration: duration}
		resp, err := c.Post(client.ChannelPromotionPath(chType, id), body)
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		if duration <= 0 {
			PrintSuccess("渠道 %q 的促销期已清除", name)
		} else {
			PrintSuccess("渠道 %q 促销期已设为 %d 秒", name, duration)
		}
		return nil
	},
}

// channelPromotionClearCmd 清除促销期
var channelPromotionClearCmd = &cobra.Command{
	Use:   "clear <name>",
	Short: "清除促销期",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		chType := getChannelType()

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		body := models.PromotionRequest{Duration: 0}
		resp, err := c.Post(client.ChannelPromotionPath(chType, id), body)
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("渠道 %q 的促销期已清除", name)
		return nil
	},
}

func init() {
	channelPromotionCmd.AddCommand(channelPromotionSetCmd)
	channelPromotionCmd.AddCommand(channelPromotionClearCmd)
	channelCmd.AddCommand(channelPromotionCmd)
}
