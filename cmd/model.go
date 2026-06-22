package cmd

import (
	"fmt"

	"ccx-cli/internal/client"

	"github.com/spf13/cobra"
)

var modelChannel string

// modelCmd 表示 model 子命令
var modelCmd = &cobra.Command{
	Use:     "model",
	Aliases: []string{"models"},
	Short:   "管理模型",
	Long: `列出或管理上游渠道支持的模型。

示例：
  ccx model list --channel my-channel --type messages
  ccx model list --channel my-channel --type chat
`,
}

// modelListCmd 列出模型
var modelListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出渠道支持的模型",
	Long: `列出指定渠道支持的模型列表。需要 --channel 指定渠道名称。

示例：
  ccx model list --channel my-channel
  ccx model list --channel my-channel --type chat
  ccx model list --channel my-channel --output json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if modelChannel == "" {
			return fmt.Errorf("--channel 是必填参数，请指定渠道名称")
		}
		chType := getChannelType()

		c := NewClient()

		// 获取渠道 ID
		id, err := getChannelIDByName(c, chType, modelChannel)
		if err != nil {
			return notFoundError(modelChannel)
		}

		// 发送模型列表请求（空 body 即可触发后端查询）
		resp, err := c.Post(client.ChannelModelsPath(chType, id), map[string]string{})
		if err != nil {
			return err
		}

		body, err := client.ReadBody(resp)
		if err != nil {
			return err
		}

		// 尝试解析为模型列表
		var modelList []string
		if err := parseJSON(body, &modelList); err == nil {
			PrintOutput(modelList)
			return nil
		}

		// 否则直接输出原始响应
		var raw any
		if err := parseJSON(body, &raw); err == nil {
			PrintOutput(raw)
			return nil
		}

		fmt.Println(string(body))
		return nil
	},
}

func init() {
	modelListCmd.Flags().StringVar(&modelChannel, "channel", "", "渠道名称（必填）")
	modelCmd.AddCommand(modelListCmd)
	RootCmd.AddCommand(modelCmd)
}
