// Package cmd 渠道管理命令
package cmd

import (
	"github.com/spf13/cobra"
)

// channelCmd 表示 channel 子命令
var channelCmd = &cobra.Command{
	Use:     "channel",
	Aliases: []string{"ch", "channels"},
	Short:   "管理上游渠道",
	Long: `管理 CCX 上游渠道（上游 API 代理目标）。

支持 5 种渠道类型（通过 --type 指定）：
  - messages（默认）：Claude Messages API
  - responses：OpenAI Responses API
  - chat：OpenAI Chat Completions API
  - gemini：Google Gemini API
  - images：OpenAI Images API

子命令：
  list, get, create, update, delete, key, status, mapping,
  reorder, promotion, ping, resume, metrics, logs, dashboard,
  scheduler-stats, capability

示例：
  ccx channel list --type messages
  ccx channel create my-channel --type chat --base-url https://api.openai.com
  ccx channel get my-channel --type responses
`,
}

func init() {
	RootCmd.AddCommand(channelCmd)
}
