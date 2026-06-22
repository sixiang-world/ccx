// Package version 存放 ccx-cli 编译时注入的版本信息
package version

// 版本信息，编译时通过 -ldflags 注入
var (
	Version   = "1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)
