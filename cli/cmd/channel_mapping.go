package cmd

import (
	"fmt"

	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// channelMappingCmd 表示 channel mapping 子命令
var channelMappingCmd = &cobra.Command{
	Use:     "mapping",
	Aliases: []string{"map", "mappings"},
	Short:   "管理模型映射",
	Long: `管理渠道的模型映射（Model Mapping）规则。

模型映射用于将客户端请求中的模型名映射为上游实际模型名。
支持前缀匹配、后缀匹配和包含匹配。

示例：
  ccx channel mapping list my-channel
  ccx channel mapping set my-channel claude-sonnet-4 claude-sonnet-4-20250514
`,
}

// channelMappingListCmd 列出模型映射
var channelMappingListCmd = &cobra.Command{
	Use:   "list <name>",
	Short: "列出渠道的模型映射",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		chType := getChannelType()

		c := NewClient()

		resp, err := c.Get(client.ChannelPath(chType), nil)
		if err != nil {
			return err
		}

		var list models.ChannelListResponse
		if err := client.DecodeResponse(resp, &list); err != nil {
			return err
		}

		for _, ch := range list.Channels {
			if ch.Name == name {
				if len(ch.ModelMapping) == 0 {
					fmt.Println("（无模型映射）")
					return nil
				}
				type mappingEntry struct {
					Source string `json:"source"`
					Target string `json:"target"`
				}
				entries := make([]mappingEntry, 0, len(ch.ModelMapping))
				for source, target := range ch.ModelMapping {
					entries = append(entries, mappingEntry{Source: source, Target: target})
				}
				PrintOutput(entries)
				return nil
			}
		}

		return notFoundError(name)
	},
}

var mappingReasoning string

// channelMappingSetCmd 设置模型映射（全量替换）
var channelMappingSetCmd = &cobra.Command{
	Use:   "set <name> <source> <target>",
	Short: "设置模型映射（全量替换）",
	Long: `设置渠道的模型映射。此为全量替换操作，会覆盖渠道上所有已有映射。
会提示当前已有的映射数量。

示例：
  ccx channel mapping set my-channel claude-sonnet-4 claude-sonnet-4-20250514
`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, source, target := args[0], args[1], args[2]
		chType := getChannelType()

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		// 构建请求体（全量替换当前 mapping）
		body := map[string]any{
			"source_pattern": source,
			"target_model":   target,
		}
		if mappingReasoning != "" {
			body["reasoning"] = mappingReasoning
		}

		resp, err := c.Put(client.ChannelMappingPath(chType, id), body)
		if err != nil {
			return err
		}

		var result struct {
			Message string `json:"message"`
		}
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("模型映射已更新：%s → %s", source, target)
		return nil
	},
}

func init() {
	channelMappingSetCmd.Flags().StringVar(&mappingReasoning, "reasoning", "", "推理级别（low|medium|high）")

	channelMappingCmd.AddCommand(channelMappingListCmd)
	channelMappingCmd.AddCommand(channelMappingSetCmd)
	channelCmd.AddCommand(channelMappingCmd)
}
