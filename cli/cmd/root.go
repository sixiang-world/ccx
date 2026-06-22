// Package cmd 定义所有 CLI 命令
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"ccx-cli/internal/client"
	"ccx-cli/internal/config"
	"ccx-cli/internal/errors"
	"ccx-cli/internal/formatter"
	"ccx-cli/internal/validator"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 全局变量
var (
	cfgFile       string
	serverURL     string
	apiKey        string
	channelType   string
	outputFormat  string
	timeoutStr    string
	retryCount    int
	noRetry       bool
	verbose       bool
	showKeys      bool
	yesFlag       bool
	caCertPath    string
	insecureSkip  bool
	prefix        string
)

// RootCmd 根命令
var RootCmd = &cobra.Command{
	Use:   "ccx",
	Short: "CCX API 代理网关命令行管理工具",
	Long: `ccx 是 CCX API 代理网关的命令行管理工具。

提供渠道（channel）的 CRUD、API 密钥管理、配置管理、
健康检查、连通性检测等功能。支持 table/json/yaml 三种输出格式。

文档: https://github.com/BenedictKing/ccx
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 初始化 Viper 配置
		initViper()
		return nil
	},
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		// 不带参数时显示帮助
		cmd.Help()
	},
}

// Execute 执行根命令
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if cliErr, ok := err.(*errors.CLIError); ok {
			fmt.Fprintln(os.Stderr, errors.FormatCLIError(cliErr))
			os.Exit(cliErr.Code)
		}
		fmt.Fprintf(os.Stderr, "✗ 错误：%v\n", err)
		os.Exit(1)
	}
}

// init 初始化全局 flag
func init() {
	// 全局 flag
	RootCmd.PersistentFlags().StringVarP(&serverURL, "server", "s", "", "CCX 服务地址（默认 http://localhost:3000）")
	RootCmd.PersistentFlags().StringVarP(&apiKey, "key", "k", "", "管理 API 密钥")
	RootCmd.PersistentFlags().StringVarP(&channelType, "type", "t", "messages", "渠道类型（messages|responses|chat|gemini|images）")
	RootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "输出格式（table|json|yaml）")
	RootCmd.PersistentFlags().StringVar(&timeoutStr, "timeout", "30s", "请求超时时间")
	RootCmd.PersistentFlags().IntVar(&retryCount, "retry", 3, "最大重试次数")
	RootCmd.PersistentFlags().BoolVar(&noRetry, "no-retry", false, "关闭自动重试")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "CLI 配置文件路径（默认 ~/.config/ccx/config.json）")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "显示详细请求信息")
	RootCmd.PersistentFlags().BoolVar(&yesFlag, "yes", false, "跳过确认提示（用于自动化脚本）")
	RootCmd.PersistentFlags().BoolVar(&showKeys, "show-keys", false, "显示完整的 API Key")
	RootCmd.PersistentFlags().StringVar(&caCertPath, "ca-cert", "", "自定义 CA 证书路径")
	RootCmd.PersistentFlags().BoolVar(&insecureSkip, "insecure-skip-verify", false, "跳过 TLS 证书验证")
	RootCmd.PersistentFlags().StringVar(&prefix, "prefix", "", "路由前缀（对应 /:routePrefix/health）")

	// 绑定 Viper
	viper.BindPFlag("server", RootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("apiKey", RootCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("type", RootCmd.PersistentFlags().Lookup("type"))
	viper.BindPFlag("output", RootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("timeout", RootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("retry", RootCmd.PersistentFlags().Lookup("retry"))
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))

	// 环境变量绑定
	viper.SetEnvPrefix("CCX")
	viper.BindEnv("server", "CCX_SERVER")
	viper.BindEnv("apiKey", "CCX_API_KEY")
	viper.BindEnv("type", "CCX_TYPE")
	viper.BindEnv("output", "CCX_OUTPUT")
	viper.BindEnv("timeout", "CCX_TIMEOUT")
}

// initViper 初始化 Viper 配置
func initViper() {
	// 配置文件
	cfgPath := cfgFile
	if cfgPath == "" {
		cfgPath = config.ConfigPath()
	}
	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
		viper.ReadInConfig() // 忽略错误，配置文件不存在也可以
	}

	// 从 Viper 读取值（如果未通过 flag 设置）
	if serverURL == "" {
		serverURL = viper.GetString("server")
	}
	if apiKey == "" {
		apiKey = viper.GetString("apiKey")
	}
	if outputFormat == "" {
		outputFormat = viper.GetString("output")
	}
	if channelType == "" {
		channelType = viper.GetString("type")
	}
	if timeoutStr == "" {
		timeoutStr = viper.GetString("timeout")
	}

	// 默认值
	if serverURL == "" {
		serverURL = "http://localhost:3000"
	}
	if outputFormat == "" {
		outputFormat = "table"
	}
	if channelType == "" {
		channelType = "messages"
	}
	if timeoutStr == "" {
		timeoutStr = "30s"
	}

	// 校验输出格式
	if err := validator.ValidateOutputFormat(outputFormat); err != nil {
		fmt.Fprintf(os.Stderr, "⚠ 警告：%v，使用默认格式 table\n", err)
		outputFormat = "table"
	}
}

// NewClient 创建 API 客户端（供子命令使用）
func NewClient() *client.Client {
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠ 警告：无效的超时设置 %q，使用 30s\n", timeoutStr)
		timeout = 30 * time.Second
	}

	maxRetries := retryCount
	if noRetry {
		maxRetries = 0
	}

	opts := []client.ClientOption{
		client.WithTimeout(timeout),
		client.WithRetry(maxRetries),
		client.WithVerbose(verbose),
	}

	if caCertPath != "" || insecureSkip {
		opts = append(opts, client.WithTLS(caCertPath, insecureSkip))
	}

	return client.NewClient(serverURL, apiKey, opts...)
}

// GetOutputFormat 获取输出格式
func GetOutputFormat() formatter.Format {
	return formatter.Format(outputFormat)
}

// getChannelType 获取并校验渠道类型
func getChannelType() string {
	if err := validator.ValidateChannelType(channelType); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		os.Exit(1)
	}
	return channelType
}

// GetShowKeys 获取是否显示完整密钥
func GetShowKeys() bool {
	return showKeys
}

// GetYesFlag 获取是否跳过确认
func GetYesFlag() bool {
	return yesFlag
}

// PrintOutput 统一输出处理
func PrintOutput(data any) {
	if err := formatter.PrintToStdout(data, GetOutputFormat(), GetShowKeys()); err != nil {
		fmt.Fprintf(os.Stderr, "⚠ 输出格式化失败：%v\n", err)
	}
}

// PrintSuccess 输出成功消息
func PrintSuccess(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if !strings.HasPrefix(msg, "✓") {
		msg = "✓ " + msg
	}
	fmt.Println(msg)
}

// PrintWarning 输出警告消息
func PrintWarning(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if !strings.HasPrefix(msg, "⚠") {
		msg = "⚠ " + msg
	}
	fmt.Fprintln(os.Stderr, msg)
}

// maskKey 脱敏 API Key（仅显示前4后4）
func maskKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) <= 8 {
		if len(key) <= 4 {
			return key
		}
		return key[:4] + "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

// getChannelIDByName 根据渠道名称获取渠道索引
func getChannelIDByName(c *client.Client, chType, name string) (int, error) {
	resp, err := c.Get(client.ChannelPath(chType), nil)
	if err != nil {
		return -1, err
	}

	var result struct {
		Channels []struct {
			Index int    `json:"index"`
			Name  string `json:"name"`
		} `json:"channels"`
	}
	if err := client.DecodeResponse(resp, &result); err != nil {
		return -1, err
	}

	for _, ch := range result.Channels {
		if ch.Name == name {
			return ch.Index, nil
		}
	}
	return -1, fmt.Errorf("未找到名为 %q 的渠道", name)
}

// notFoundError 返回渠道未找到错误
func notFoundError(name string) error {
	return fmt.Errorf("✗ 未找到名为 %q 的渠道", name)
}

// parseJSON 解析 JSON 字节到目标结构
func parseJSON(data []byte, target any) error {
	return json.Unmarshal(data, target)
}
