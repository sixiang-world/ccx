package cmd

import (
	"fmt"
	"strconv"

	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// settingsCmd 表示 settings 子命令
var settingsCmd = &cobra.Command{
	Use:     "settings",
	Aliases: []string{"setting"},
	Short:   "管理运行时设置",
	Long: `管理 CCX 的运行时设置，包括 Fuzzy 模式、熔断器、图片轮次限制等。

子命令：
  fuzzy, circuit-breaker, image-turn-limit, conversations

示例：
  ccx settings fuzzy get
  ccx settings fuzzy set true
  ccx settings circuit-breaker get
  ccx settings circuit-breaker set --window-size 10
  ccx settings image-turn-limit get
  ccx settings image-turn-limit set 5
`,
}

// ============== Fuzzy 模式 ==============

// settingsFuzzyCmd Fuzzy 模式子命令
var settingsFuzzyCmd = &cobra.Command{
	Use:   "fuzzy",
	Short: "管理 Fuzzy 模式",
	Long: `管理 Fuzzy 模式。启用后所有非 2xx 错误都尝试 failover。

示例：
  ccx settings fuzzy get
  ccx settings fuzzy set true
  ccx settings fuzzy set false
`,
}

var settingsFuzzyGetCmd = &cobra.Command{
	Use:   "get",
	Short: "查看 Fuzzy 模式状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := NewClient()
		resp, err := c.Get(client.SettingsFuzzyPath(), nil)
		if err != nil {
			return err
		}

		var result models.FuzzyModeResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintOutput(result)
		return nil
	},
}

var settingsFuzzySetCmd = &cobra.Command{
	Use:   "set <true|false>",
	Short: "设置 Fuzzy 模式",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		enabled := args[0] == "true" || args[0] == "1" || args[0] == "yes"

		c := NewClient()
		body := models.FuzzyModeRequest{Enabled: enabled}
		resp, err := c.Put(client.SettingsFuzzyPath(), body)
		if err != nil {
			return err
		}

		var result models.FuzzyModeResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("Fuzzy 模式已%s", map[bool]string{true: "启用", false: "关闭"}[enabled])
		return nil
	},
}

// ============== 熔断器 ==============

// settingsCircuitBreakerCmd 熔断器子命令
var settingsCircuitBreakerCmd = &cobra.Command{
	Use:   "circuit-breaker",
	Short: "管理熔断器配置",
	Long: `管理熔断器运行时配置。

示例：
  ccx settings circuit-breaker get
  ccx settings circuit-breaker set --window-size 10 --failure-threshold 0.5
`,
}

var settingsCircuitBreakerGetCmd = &cobra.Command{
	Use:   "get",
	Short: "查看熔断器配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := NewClient()
		resp, err := c.Get(client.SettingsCircuitBreakerPath(), nil)
		if err != nil {
			return err
		}

		var result models.CircuitBreakerResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintOutput(result)
		return nil
	},
}

var (
	cbWindowSize                   int
	cbFailureThreshold             float64
	cbConsecutiveFailuresThreshold int
	cbRequestTimeoutMs             int
	cbResponseHeaderTimeoutMs      int
	cbStreamFirstContentTimeoutMs  int
	cbStreamInactivityTimeoutMs    int
	cbStreamToolCallIdleTimeoutMs  int
)

var settingsCircuitBreakerSetCmd = &cobra.Command{
	Use:   "set",
	Short: "设置熔断器配置",
	Long: `设置熔断器运行时参数。仅设置指定的字段，未设置字段保持不变。

示例：
  ccx settings circuit-breaker set --window-size 10 --failure-threshold 0.5
  ccx settings circuit-breaker set --request-timeout-ms 60000
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := NewClient()

		body := models.CircuitBreakerRequest{}
		if cmd.Flags().Changed("window-size") {
			v := cbWindowSize
			body.WindowSize = &v
		}
		if cmd.Flags().Changed("failure-threshold") {
			v := cbFailureThreshold
			body.FailureThreshold = &v
		}
		if cmd.Flags().Changed("consecutive-failures") {
			v := cbConsecutiveFailuresThreshold
			body.ConsecutiveFailuresThreshold = &v
		}
		if cmd.Flags().Changed("request-timeout-ms") {
			v := cbRequestTimeoutMs
			body.RequestTimeoutMs = &v
		}
		if cmd.Flags().Changed("response-header-timeout-ms") {
			v := cbResponseHeaderTimeoutMs
			body.ResponseHeaderTimeoutMs = &v
		}
		if cmd.Flags().Changed("stream-first-content-timeout-ms") {
			v := cbStreamFirstContentTimeoutMs
			body.StreamFirstContentTimeoutMs = &v
		}
		if cmd.Flags().Changed("stream-inactivity-timeout-ms") {
			v := cbStreamInactivityTimeoutMs
			body.StreamInactivityTimeoutMs = &v
		}
		if cmd.Flags().Changed("stream-tool-call-idle-timeout-ms") {
			v := cbStreamToolCallIdleTimeoutMs
			body.StreamToolCallIdleTimeoutMs = &v
		}

		resp, err := c.Put(client.SettingsCircuitBreakerPath(), body)
		if err != nil {
			return err
		}

		var result struct {
			Success       bool                           `json:"success"`
			CircuitBreaker models.CircuitBreakerResponse `json:"circuitBreaker"`
		}
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("熔断器配置已更新")
		return nil
	},
}

// ============== 图片轮次限制 ==============

// settingsImageLimitCmd 图片轮次限制子命令
var settingsImageLimitCmd = &cobra.Command{
	Use:   "image-turn-limit",
	Short: "管理历史图片轮次限制",
	Long: `管理全局历史图片轮次限制。
超过此轮次的历史图片将替换为占位符。

示例：
  ccx settings image-turn-limit get
  ccx settings image-turn-limit set 5
`,
}

var settingsImageLimitGetCmd = &cobra.Command{
	Use:   "get",
	Short: "查看图片轮次限制",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := NewClient()
		resp, err := c.Get(client.SettingsImageTurnLimitPath(), nil)
		if err != nil {
			return err
		}

		var result models.HistoricalImageTurnLimitResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintOutput(result)
		return nil
	},
}

var settingsImageLimitSetCmd = &cobra.Command{
	Use:   "set <limit>",
	Short: "设置图片轮次限制",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("无效的限制值 %q（应为数字）", args[0])
		}

		c := NewClient()
		body := models.HistoricalImageTurnLimitRequest{Limit: limit}
		resp, err := c.Put(client.SettingsImageTurnLimitPath(), body)
		if err != nil {
			return err
		}

		var result models.HistoricalImageTurnLimitResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("历史图片轮次限制已设为 %d", result.HistoricalImageTurnLimit)
		return nil
	},
}

// ============== 对话设置 ==============

// settingsConversationsCmd 对话设置子命令
var settingsConversationsCmd = &cobra.Command{
	Use:   "conversations",
	Short: "管理对话设置",
	Long: `管理对话相关的设置。

示例：
  ccx settings conversations get
  ccx settings conversations set --enabled true --max-age 3600
`,
}

var settingsConversationsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "查看对话设置",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := NewClient()
		resp, err := c.Get(client.SettingsConversationsPath(), nil)
		if err != nil {
			return err
		}

		var result any
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintOutput(result)
		return nil
	},
}

var (
	convEnabled bool
	convMaxAge  int
)

var settingsConversationsSetCmd = &cobra.Command{
	Use:   "set",
	Short: "更新对话设置",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := NewClient()

		body := map[string]any{}
		if cmd.Flags().Changed("enabled") {
			body["enabled"] = convEnabled
		}
		if cmd.Flags().Changed("max-age") {
			body["maxAge"] = convMaxAge
		}

		if len(body) == 0 {
			return fmt.Errorf("请指定至少一个要设置的参数（--enabled 或 --max-age）")
		}

		resp, err := c.Put(client.SettingsConversationsPath(), body)
		if err != nil {
			return err
		}

		var result any
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintSuccess("对话设置已更新")
		return nil
	},
}

func init() {
	// Fuzzy
	settingsFuzzyCmd.AddCommand(settingsFuzzyGetCmd)
	settingsFuzzyCmd.AddCommand(settingsFuzzySetCmd)
	settingsCmd.AddCommand(settingsFuzzyCmd)

	// Circuit Breaker
	settingsCircuitBreakerSetCmd.Flags().IntVar(&cbWindowSize, "window-size", 0, "滑动窗口大小 (3-100)")
	settingsCircuitBreakerSetCmd.Flags().Float64Var(&cbFailureThreshold, "failure-threshold", 0, "失败率阈值 (0.01-1.0)")
	settingsCircuitBreakerSetCmd.Flags().IntVar(&cbConsecutiveFailuresThreshold, "consecutive-failures", 0, "连续失败阈值 (1-100)")
	settingsCircuitBreakerSetCmd.Flags().IntVar(&cbRequestTimeoutMs, "request-timeout-ms", 0, "请求超时 (ms)")
	settingsCircuitBreakerSetCmd.Flags().IntVar(&cbResponseHeaderTimeoutMs, "response-header-timeout-ms", 0, "响应头超时 (ms)")
	settingsCircuitBreakerSetCmd.Flags().IntVar(&cbStreamFirstContentTimeoutMs, "stream-first-content-timeout-ms", 0, "流式首字超时 (ms)")
	settingsCircuitBreakerSetCmd.Flags().IntVar(&cbStreamInactivityTimeoutMs, "stream-inactivity-timeout-ms", 0, "流式空闲超时 (ms)")
	settingsCircuitBreakerSetCmd.Flags().IntVar(&cbStreamToolCallIdleTimeoutMs, "stream-tool-call-idle-timeout-ms", 0, "工具调用空闲超时 (ms)")
	settingsCircuitBreakerCmd.AddCommand(settingsCircuitBreakerGetCmd)
	settingsCircuitBreakerCmd.AddCommand(settingsCircuitBreakerSetCmd)
	settingsCmd.AddCommand(settingsCircuitBreakerCmd)

	// Image Turn Limit
	settingsImageLimitCmd.AddCommand(settingsImageLimitGetCmd)
	settingsImageLimitCmd.AddCommand(settingsImageLimitSetCmd)
	settingsCmd.AddCommand(settingsImageLimitCmd)

	// Conversations
	settingsConversationsGetCmd.Flags().BoolVar(&convEnabled, "enabled", false, "启用对话追踪")
	settingsConversationsSetCmd.Flags().BoolVar(&convEnabled, "enabled", false, "启用对话追踪")
	settingsConversationsSetCmd.Flags().IntVar(&convMaxAge, "max-age", 0, "对话最大存活时间（秒）")
	settingsConversationsCmd.AddCommand(settingsConversationsGetCmd)
	settingsConversationsCmd.AddCommand(settingsConversationsSetCmd)
	settingsCmd.AddCommand(settingsConversationsCmd)

	RootCmd.AddCommand(settingsCmd)
}
