package cmd

import (
	"fmt"
	"net/url"

	"ccx-cli/internal/client"
	"ccx-cli/internal/models"

	"github.com/spf13/cobra"
)

// healthCmd 表示 health 子命令
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "服务健康检查",
	Long: `检查 CCX 服务的健康状态。

使用 --prefix 可指定路由前缀（如 Docker 部署中路由前缀非空时使用）。

示例：
  ccx health
  ccx health --prefix my-prefix
  ccx health --output json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := NewClient()

		path := client.HealthPath()
		if prefix != "" {
			path = fmt.Sprintf("/%s/health", url.PathEscape(prefix))
		}

		resp, err := c.Get(path, nil)
		if err != nil {
			return err
		}

		var result models.HealthResponse
		if err := client.DecodeResponse(resp, &result); err != nil {
			return err
		}

		PrintOutput(result)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(healthCmd)
}
