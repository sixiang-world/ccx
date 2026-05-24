package configservice

import (
	"encoding/json"
	"fmt"
	"strings"
)

// computeTextDiff 逐行对比 before/after，生成 git-style diff 行。
func computeTextDiff(path, before, after string) FileDiff {
	oldLines := splitLines(before)
	newLines := splitLines(after)

	action := "modify"
	if before == "" && after != "" {
		action = "create"
	} else if before != "" && after == "" {
		action = "delete"
	} else if before == after {
		action = "modify"
	}

	lines := lcsDiff(oldLines, newLines)
	return FileDiff{Path: path, Action: action, Lines: lines}
}

// computeJSONDiff 将两个 map 格式化为 JSON 后逐行对比。
func computeJSONDiff(path string, before, after map[string]any) FileDiff {
	oldContent := ""
	if before != nil {
		oldContent = formatJSON(before)
	}
	newContent := ""
	if after != nil {
		newContent = formatJSON(after)
	}
	return computeTextDiff(path, oldContent, newContent)
}

// formatJSON 将 map 格式化为缩进 JSON 字符串。
func formatJSON(data map[string]any) string {
	if data == nil {
		return ""
	}
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", data)
	}
	return string(content) + "\n"
}

// maskSensitiveValue 对单个值进行脱敏。
// 短于 12 字符显示为 "***"；否则保留前 3 后 4，中间 "***"。
func maskSensitiveValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	runes := []rune(value)
	if len(runes) < 12 {
		return "***"
	}
	prefix := string(runes[:3])
	suffix := string(runes[len(runes)-4:])
	return prefix + "***" + suffix
}

// sensitiveFieldKeys 需要脱敏的配置字段名。
var sensitiveFieldKeys = []string{
	"ANTHROPIC_API_KEY",
	"ANTHROPIC_AUTH_TOKEN",
	"OPENAI_API_KEY",
}

// maskMapSensitiveKeys 对 map 中指定 key 的值进行脱敏（返回新 map，不修改原 map）。
func maskMapSensitiveKeys(data map[string]any, keys ...string) map[string]any {
	if data == nil {
		return nil
	}
	result := make(map[string]any, len(data))
	for k, v := range data {
		result[k] = v
	}
	for _, key := range keys {
		if val, ok := result[key]; ok {
			if s, ok := val.(string); ok && s != "" {
				result[key] = maskSensitiveValue(s)
			}
		}
	}
	return result
}

// maskJSONSensitiveKeys 对 JSON map 中嵌套的 env map 内的敏感字段进行脱敏。
func maskJSONSensitiveKeys(data map[string]any) map[string]any {
	if data == nil {
		return nil
	}
	result := make(map[string]any, len(data))
	for k, v := range data {
		result[k] = v
	}
	if env, ok := result["env"].(map[string]any); ok {
		result["env"] = maskMapSensitiveKeys(env, sensitiveFieldKeys...)
	}
	return result
}

// maskTextSensitiveValues 对文本内容中出现的敏感值进行行内脱敏。
// 用于 TOML / JSON 文本 diff 的 before/after 内容。
func maskTextSensitiveValues(content string, keyValues map[string]string) string {
	if len(keyValues) == 0 {
		return content
	}
	result := content
	for _, value := range keyValues {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		masked := maskSensitiveValue(value)
		result = strings.ReplaceAll(result, value, masked)
	}
	return result
}

// splitLines 将文本按换行符分割为行切片。空文本返回空切片。
func splitLines(text string) []string {
	if text == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	// 去除末尾空行（由末尾换行符产生）
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// lcsDiff 使用 LCS 算法生成 diff 行序列。
func lcsDiff(oldLines, newLines []string) []DiffLine {
	m, n := len(oldLines), len(newLines)
	if m == 0 && n == 0 {
		return nil
	}

	// 特殊情况优化
	if m == 0 {
		lines := make([]DiffLine, n)
		for i, l := range newLines {
			lines[i] = DiffLine{Type: "added", Content: l}
		}
		return lines
	}
	if n == 0 {
		lines := make([]DiffLine, m)
		for i, l := range oldLines {
			lines[i] = DiffLine{Type: "removed", Content: l}
		}
		return lines
	}

	// LCS DP 表
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if oldLines[i-1] == newLines[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] > dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}

	// 回溯生成 diff（逆序收集，最后反转）
	var reversed []DiffLine
	i, j := m, n
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && oldLines[i-1] == newLines[j-1] {
			reversed = append(reversed, DiffLine{Type: "context", Content: oldLines[i-1]})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			reversed = append(reversed, DiffLine{Type: "added", Content: newLines[j-1]})
			j--
		} else {
			reversed = append(reversed, DiffLine{Type: "removed", Content: oldLines[i-1]})
			i--
		}
	}

	// 反转得到正确顺序
	for left, right := 0, len(reversed)-1; left < right; left, right = left+1, right-1 {
		reversed[left], reversed[right] = reversed[right], reversed[left]
	}
	return reversed
}
