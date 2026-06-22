// Package config 管理 CLI 自身配置（~/.config/ccx/config.json）
package config

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"ccx-cli/internal/errors"
)

// CLIConfig CLI 自身配置文件结构
type CLIConfig struct {
	Server      string `json:"server"`
	APIKey      string `json:"apiKey"`
	DefaultType string `json:"defaultType"`
	Output      string `json:"output"`
	Timeout     string `json:"timeout"`
}

// DefaultCLIConfig 返回默认 CLI 配置
func DefaultCLIConfig() CLIConfig {
	return CLIConfig{
		Server:      "http://localhost:3000",
		DefaultType: "messages",
		Output:      "table",
		Timeout:     "30s",
	}
}

// ConfigPath 返回默认配置文件路径
func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "ccx", "config.json")
}

// LoadCLIConfig 加载 CLI 配置文件
func LoadCLIConfig(path string) (*CLIConfig, error) {
	if path == "" {
		path = ConfigPath()
	}
	if path == "" {
		return nil, errors.New(errors.ExitCodeUserError, "无法确定配置文件路径（$HOME 未设置）")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if stderrors.Is(err, fs.ErrNotExist) {
			// 配置文件不存在，返回默认配置
			cfg := DefaultCLIConfig()
			return &cfg, nil
		}
		return nil, errors.NewWithDetail(errors.ExitCodeUserError,
			fmt.Sprintf("读取配置文件 %s 失败", path), err.Error())
	}

	var cfg CLIConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, errors.NewWithDetail(errors.ExitCodeUserError,
			fmt.Sprintf("解析配置文件 %s 失败", path), err.Error())
	}

	// 使用默认值填充缺失字段
	if cfg.Server == "" {
		cfg.Server = "http://localhost:3000"
	}
	if cfg.DefaultType == "" {
		cfg.DefaultType = "messages"
	}
	if cfg.Output == "" {
		cfg.Output = "table"
	}
	if cfg.Timeout == "" {
		cfg.Timeout = "30s"
	}

	return &cfg, nil
}

// SaveCLIConfig 保存 CLI 配置文件
func SaveCLIConfig(cfg *CLIConfig, path string) error {
	if path == "" {
		path = ConfigPath()
	}
	if path == "" {
		return errors.New(errors.ExitCodeUserError, "无法确定配置文件路径")
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return errors.NewWithDetail(errors.ExitCodeUserError,
			fmt.Sprintf("创建配置目录 %s 失败", dir), err.Error())
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return errors.NewWithDetail(errors.ExitCodeUserError, "序列化配置失败", err.Error())
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return errors.NewWithDetail(errors.ExitCodeUserError,
			fmt.Sprintf("写入配置文件 %s 失败", path), err.Error())
	}

	fmt.Printf("✓ 配置已保存到 %s\n", path)
	fmt.Println("⚠ 注意：文件中包含明文 API Key，建议设置权限 chmod 600")
	return nil
}
