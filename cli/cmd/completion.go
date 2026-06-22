package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd 表示 completion 子命令
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "生成 shell 自动补全脚本",
	Long: `生成 shell 自动补全脚本，支持 bash、zsh、fish 和 powershell。

示例：
  ccx completion bash > /etc/bash_completion.d/ccx
  ccx completion zsh > /usr/local/share/zsh/site-functions/_ccx
  ccx completion fish > ~/.config/fish/completions/ccx.fish
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := args[0]

		switch shell {
		case "bash":
			return RootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return RootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return RootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return RootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return fmt.Errorf("不支持的 shell %q（支持：bash|zsh|fish|powershell）", shell)
		}
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
}
