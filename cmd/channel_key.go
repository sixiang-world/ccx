package cmd

import (
	"fmt"

	"ccx-cli/internal/client"
	"ccx-cli/internal/models"
	"ccx-cli/internal/validator"

	"github.com/spf13/cobra"
)

// channelKeyCmd 表示 channel key 子命令
var channelKeyCmd = &cobra.Command{
	Use:     "key",
	Aliases: []string{"keys"},
	Short:   "管理渠道 API 密钥",
	Long: `管理渠道的 API 密钥（添加、删除、列出、移动优先级、恢复拉黑密钥）。

子命令：
  add, remove, list, move, restore

示例：
  ccx channel key add my-channel sk-xxx
  ccx channel key remove my-channel sk-xxx
  ccx channel key list my-channel
  ccx channel key move my-channel sk-xxx --position top
  ccx channel key restore my-channel sk-xxx
`,
}

// channelKeyAddCmd 添加密钥
var channelKeyAddCmd = &cobra.Command{
	Use:   "add <name> <key>",
	Short: "添加 API 密钥到渠道",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, key := args[0], args[1]
		chType := getChannelType()

		if err := validator.ValidateAPIKey(key); err != nil {
			return fmt.Errorf("无效的 API Key：%v", err)
		}

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		body := models.AddKeyRequest{APIKey: key}
		resp, err := c.Post(client.ChannelKeyPath(chType, id), body)
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("密钥 %s 已添加到渠道 %q", maskKey(key), name)
		return nil
	},
}

// channelKeyRemoveCmd 删除密钥
var channelKeyRemoveCmd = &cobra.Command{
	Use:   "remove <name> <key>",
	Short: "从渠道移除 API 密钥",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, key := args[0], args[1]
		chType := getChannelType()

		if err := validator.ValidateAPIKey(key); err != nil {
			return fmt.Errorf("无效的 API Key：%v", err)
		}

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		resp, err := c.Delete(client.ChannelKeyIDPath(chType, id, key))
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("密钥 %s 已从渠道 %q 移除", maskKey(key), name)
		return nil
	},
}

// channelKeyListCmd 列出密钥
var channelKeyListCmd = &cobra.Command{
	Use:   "list <name>",
	Short: "列出渠道的所有 API 密钥",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		chType := getChannelType()

		c := NewClient()

		// 获取渠道列表并本地过滤
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
				type keyEntry struct {
					Index int    `json:"index"`
					Key   string `json:"key"`
					State string `json:"state"`
				}
				entries := make([]keyEntry, 0, len(ch.APIKeys)+len(ch.DisabledAPIKeys))
				for i, k := range ch.APIKeys {
					displayKey := k
					if !GetShowKeys() {
						displayKey = maskKey(k)
					}
					entries = append(entries, keyEntry{Index: i, Key: displayKey, State: "active"})
				}
				for _, dk := range ch.DisabledAPIKeys {
					displayKey := dk.Key
					if !GetShowKeys() {
						displayKey = maskKey(dk.Key)
					}
					entries = append(entries, keyEntry{Key: displayKey, State: "disabled (" + dk.Reason + ")"})
				}
				PrintOutput(entries)
				return nil
			}
		}

		return notFoundError(name)
	},
}

var keyPosition string

// channelKeyMoveCmd 移动密钥优先级
var channelKeyMoveCmd = &cobra.Command{
	Use:   "move <name> <key>",
	Short: "调整 API 密钥优先级",
	Long: `将 API 密钥移动到顶部或底部以调整优先级。

示例：
  ccx channel key move my-channel sk-xxx --position top
  ccx channel key move my-channel sk-xxx --position bottom
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, key := args[0], args[1]
		chType := getChannelType()

		if keyPosition != "top" && keyPosition != "bottom" {
			return fmt.Errorf("--position 必须为 top 或 bottom")
		}

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		resp, err := c.Post(client.ChannelKeyMovePath(chType, id, key, keyPosition), nil)
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("密钥 %s 已移到%s", maskKey(key), map[string]string{"top": "顶部", "bottom": "底部"}[keyPosition])
		return nil
	},
}

// channelKeyRestoreCmd 恢复被拉黑的密钥
var channelKeyRestoreCmd = &cobra.Command{
	Use:   "restore <name> <key>",
	Short: "恢复被拉黑的 API 密钥",
	Long: `恢复被自动拉黑的 API 密钥到活跃列表。

示例：
  ccx channel key restore my-channel sk-xxx
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, key := args[0], args[1]
		chType := getChannelType()

		c := NewClient()

		id, err := getChannelIDByName(c, chType, name)
		if err != nil {
			return notFoundError(name)
		}

		body := models.RestoreKeyRequest{APIKey: key}
		resp, err := c.Post(client.ChannelKeyRestorePath(chType, id), body)
		if err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("密钥 %s 已恢复", maskKey(key))
		return nil
	},
}

func init() {
	channelKeyMoveCmd.Flags().StringVar(&keyPosition, "position", "top", "目标位置（top|bottom）")

	channelKeyCmd.AddCommand(channelKeyAddCmd)
	channelKeyCmd.AddCommand(channelKeyRemoveCmd)
	channelKeyCmd.AddCommand(channelKeyListCmd)
	channelKeyCmd.AddCommand(channelKeyMoveCmd)
	channelKeyCmd.AddCommand(channelKeyRestoreCmd)
	channelCmd.AddCommand(channelKeyCmd)
}
