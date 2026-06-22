package cmd

import (
	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// resumeCmd 恢复渠道
var resumeCmd = &cobra.Command{
	Use:   "resume <name>",
	Short: "恢复被熔断/拉黑的渠道",
	Long: `恢复被熔断或拉黑的渠道（重置熔断状态、恢复拉黑 Key，保留历史统计）。

示例：
  ccx channel resume my-channel
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

		resp, err := c.Post(client.ChannelResumePath(chType, id), nil)
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("渠道 %q 已恢复", name)
		return nil
	},
}

func init() {
	channelCmd.AddCommand(resumeCmd)
}
