// Package validator 提供参数校验功能
package validator

import (
	"fmt"
	"net/url"
	"strings"
)

// ValidateAPIKey 校验 API Key 格式
// 参考设计文档 §3.6 密钥格式前置校验
func ValidateAPIKey(key string) error {
	if key == "" {
		return fmt.Errorf("API Key 不能为空")
	}

	// 针对已知类型的 Key 做前缀校验
	if strings.HasPrefix(key, "sk-ant-") {
		if len(key) < 20 {
			return fmt.Errorf("Anthropic API Key 长度不足（预期 >= 20，当前 %d）", len(key))
		}
	} else if strings.HasPrefix(key, "sk-proj-") {
		if len(key) < 20 {
			return fmt.Errorf("OpenAI 项目 API Key 长度不足（预期 >= 20，当前 %d）", len(key))
		}
	} else if strings.HasPrefix(key, "sk-") {
		if len(key) < 20 {
			return fmt.Errorf("OpenAI API Key 长度不足（预期 >= 20，当前 %d）", len(key))
		}
	} else if strings.HasPrefix(key, "AIzaSy") {
		if len(key) < 20 {
			return fmt.Errorf("Gemini API Key 长度不足（预期 >= 20，当前 %d）", len(key))
		}
	} else {
		// 自定义 Key，仅检查非空
		if len(key) < 8 {
			return fmt.Errorf("API Key 长度不足（建议 >= 8，当前 %d）", len(key))
		}
	}
	return nil
}

// ValidateURL 校验 URL 格式
func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL 不能为空")
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("无效的 URL 格式：%w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("不支持的协议 %q（仅支持 http/https）", u.Scheme)
	}
	if u.Host == "" {
		return fmt.Errorf("URL 缺少主机名")
	}
	return nil
}

// ValidateChannelType 校验渠道类型
func ValidateChannelType(t string) error {
	switch t {
	case "messages", "responses", "chat", "gemini", "images":
		return nil
	default:
		return fmt.Errorf("不支持的渠道类型 %q（可选值：messages|responses|chat|gemini|images）", t)
	}
}

// ValidateChannelStatus 校验渠道状态
func ValidateChannelStatus(s string) error {
	switch s {
	case "active", "suspended", "disabled":
		return nil
	default:
		return fmt.Errorf("不支持的状态 %q（可选值：active|suspended|disabled）", s)
	}
}

// ValidateServiceType 校验服务类型
func ValidateServiceType(s string) error {
	switch s {
	case "", "claude", "openai", "gemini", "custom":
		return nil
	default:
		return fmt.Errorf("不支持的服务类型 %q（可选值：claude|openai|gemini|custom）", s)
	}
}

// ValidateOutputFormat 校验输出格式
func ValidateOutputFormat(f string) error {
	switch f {
	case "table", "json", "yaml":
		return nil
	default:
		return fmt.Errorf("不支持的输出格式 %q（可选值：table|json|yaml）", f)
	}
}

// ValidateDuration 校验持续时间字符串（秒数）
func ValidateDuration(d string) (int, error) {
	if d == "" {
		return 0, fmt.Errorf("持续时间不能为空")
	}
	// 纯数字视为秒
	var seconds int
	n, err := fmt.Sscanf(d, "%d", &seconds)
	if err != nil || n != 1 {
		return 0, fmt.Errorf("无效的持续时间 %q（请使用纯数字秒数）", d)
	}
	// 检查是否有尾部非数字字符
	trimmed := strings.TrimSpace(d)
	for _, c := range trimmed {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("无效的持续时间 %q（请使用纯数字秒数）", d)
		}
	}
	if seconds < 0 {
		return 0, fmt.Errorf("持续时间不能为负数")
	}
	return seconds, nil
}

// MaskAPIKey 脱敏 API Key，只显示前4位和后4位
func MaskAPIKey(key string) string {
	if len(key) <= 8 {
		if len(key) <= 4 {
			return key
		}
		return key[:4] + strings.Repeat("*", len(key)-4)
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

// MaskAPIKeys 脱敏 API Key 列表
func MaskAPIKeys(keys []string) []string {
	masked := make([]string, len(keys))
	for i, k := range keys {
		masked[i] = MaskAPIKey(k)
	}
	return masked
}

// ParseModelMapping 解析 key=value 格式的模型映射
func ParseModelMapping(pairs []string) (map[string]string, error) {
	if len(pairs) == 0 {
		return nil, nil
	}
	mapping := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("无效的模型映射 %q（格式：source=target）", pair)
		}
		mapping[parts[0]] = parts[1]
	}
	return mapping, nil
}
