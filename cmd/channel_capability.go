package cmd

import (
	"fmt"
	"net/url"

	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// channelCapabilityCmd 表示 channel capability 子命令
var channelCapabilityCmd = &cobra.Command{
	Use:     "capability",
	Aliases: []string{"cap"},
	Short:   "管理渠道能力测试",
	Long: `查看和测试渠道的模型能力。

子命令：
  snapshot <name>       查看渠道能力快照
  test <name>           运行渠道能力测试
  test-status <name> <jobId>  查看测试任务状态
  test-cancel <name> <jobId>  取消测试任务
  test-retry <name> <jobId>   重试测试失败的模型

示例：
  ccx channel capability snapshot my-channel
  ccx channel capability test my-channel
  ccx channel capability test-status my-channel <jobId>
`,
}

// channelCapabilitySnapshotCmd 查看能力快照
var channelCapabilitySnapshotCmd = &cobra.Command{
	Use:   "snapshot <name>",
	Short: "查看渠道能力快照",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		chType := getChannelType()

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		resp, err := c.Get(client.ChannelCapabilitySnapshotPath(chType, id), nil)
		if err != nil {
			return err
		}

		var result models.CapabilityTestResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintOutput(result)
		return nil
	},
}

// channelCapabilityTestCmd 运行能力测试
var channelCapabilityTestCmd = &cobra.Command{
	Use:   "test <name>",
	Short: "运行渠道能力测试",
	Long: `对指定渠道运行模型能力测试（发送测试请求验证各模型可用性）。

示例：
  ccx channel capability test my-channel
  ccx channel capability test my-channel --type responses
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

		resp, err := c.Post(client.ChannelCapabilityTestPath(chType, id), nil)
		if err != nil {
			return err
		}

		var result models.CapabilityTestResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintOutput(result)
		return nil
	},
}

// channelCapabilityTestStatusCmd 查看测试任务状态
var channelCapabilityTestStatusCmd = &cobra.Command{
	Use:   "test-status <name> <jobId>",
	Short: "查看能力测试任务状态",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, jobID := args[0], args[1]
		chType := getChannelType()

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		path := fmt.Sprintf("%s/%s", client.ChannelCapabilityTestPath(chType, id), url.PathEscape(jobID))
		resp, err := c.Get(path, nil)
		if err != nil {
			return err
		}

		var result models.CapabilityTestResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintOutput(result)
		return nil
	},
}

// channelCapabilityTestCancelCmd 取消测试任务
var channelCapabilityTestCancelCmd = &cobra.Command{
	Use:   "test-cancel <name> <jobId>",
	Short: "取消能力测试任务",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, jobID := args[0], args[1]
		chType := getChannelType()

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		path := fmt.Sprintf("%s/%s", client.ChannelCapabilityTestPath(chType, id), url.PathEscape(jobID))
		resp, err := c.Delete(path)
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("能力测试任务 %s 已取消", jobID)
		return nil
	},
}

// channelCapabilityTestRetryCmd 重试测试失败的模型
var channelCapabilityTestRetryCmd = &cobra.Command{
	Use:   "test-retry <name> <jobId>",
	Short: "重试能力测试中失败的模型",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, jobID := args[0], args[1]
		chType := getChannelType()

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		path := fmt.Sprintf("%s/%s/retry", client.ChannelCapabilityTestPath(chType, id), url.PathEscape(jobID))
		resp, err := c.Post(path, nil)
		if err != nil {
			return err
		}

		var result models.CapabilityTestResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintOutput(result)
		return nil
	},
}

func init() {
	channelCapabilityCmd.AddCommand(channelCapabilitySnapshotCmd)
	channelCapabilityCmd.AddCommand(channelCapabilityTestCmd)
	channelCapabilityCmd.AddCommand(channelCapabilityTestStatusCmd)
	channelCapabilityCmd.AddCommand(channelCapabilityTestCancelCmd)
	channelCapabilityCmd.AddCommand(channelCapabilityTestRetryCmd)
	channelCmd.AddCommand(channelCapabilityCmd)
}
