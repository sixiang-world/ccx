// Package cmd 定义所有 CLI 命令
package cmd

import (
	"fmt"

	"ccx-cli/internal/version"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示 ccx-cli 的版本号、构建时间和 Git 提交哈希。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ccx-cli %s\n", version.Version)
		fmt.Printf("  BuildTime: %s\n", version.BuildTime)
		fmt.Printf("  GitCommit: %s\n", version.GitCommit)
	},
}
