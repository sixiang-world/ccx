// Package formatter 提供 table/json/yaml 三种输出格式
package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

// Format 输出格式
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

// Print 根据格式输出数据
func Print(w io.Writer, data any, format Format, showKeys bool) error {
	switch format {
	case FormatJSON:
		return printJSON(w, data, showKeys)
	case FormatYAML:
		return printYAML(w, data, showKeys)
	case FormatTable:
		return printTable(w, data, showKeys)
	default:
		return printTable(w, data, showKeys)
	}
}

// PrintToStdout 输出到标准输出
func PrintToStdout(data any, format Format, showKeys bool) error {
	return Print(os.Stdout, data, format, showKeys)
}

// ============== JSON 输出 ==============

func printJSON(w io.Writer, data any, showKeys bool) error {
	// 先对数据做脱敏处理
	if !showKeys {
		data = maskSensitiveData(data)
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// ============== YAML 输出 ==============

func printYAML(w io.Writer, data any, showKeys bool) error {
	if !showKeys {
		data = maskSensitiveData(data)
	}
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	defer encoder.Close()
	return encoder.Encode(data)
}

// ============== Table 输出 ==============

func printTable(w io.Writer, data any, showKeys bool) error {
	// 根据数据类型选择不同的表格渲染方式
	switch d := data.(type) {
	case []map[string]any:
		return printMapSliceTable(w, d, showKeys)
	case map[string]any:
		return printSingleMapTable(w, d, showKeys)
	case []any:
		// json.Unmarshal 产生的 []any 可能包含 map[string]any 元素
		if len(d) > 0 {
			if _, ok := d[0].(map[string]any); ok {
				rows := make([]map[string]any, len(d))
				for i, item := range d {
					if m, ok := item.(map[string]any); ok {
						rows[i] = m
					}
				}
				return printMapSliceTable(w, rows, showKeys)
			}
		}
		// 非对象数组：渲染为单列表格
		return printAnySliceTable(w, d, showKeys)
	default:
		// 尝试 JSON 序列化后渲染
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("table 输出失败：%w", err)
		}
		var parsed any
		if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
			return errJSON(w, data)
		}
		return printTable(w, parsed, showKeys)
	}
}

// printSingleMapTable 输出单条记录的表格
func printSingleMapTable(w io.Writer, data map[string]any, showKeys bool) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"字段", "值"})
	table.SetAutoWrapText(false)
	table.SetRowLine(false)
	table.SetColumnSeparator("  ")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")

	for k, v := range data {
		vStr := formatValue(v, showKeys)
		if !showKeys && isSensitiveKey(k) {
			vStr = maskString(vStr)
		}
		table.Append([]string{k, vStr})
	}
	table.Render()
	return nil
}

// printMapSliceTable 输出多条记录的表格
func printMapSliceTable(w io.Writer, rows []map[string]any, showKeys bool) error {
	if len(rows) == 0 {
		fmt.Fprintln(w, "（空）")
		return nil
	}

	// 收集所有列名
	headers := make([]string, 0)
	seen := make(map[string]bool)
	for _, row := range rows {
		for k := range row {
			if !seen[k] {
				headers = append(headers, k)
				seen[k] = true
			}
		}
	}

	table := tablewriter.NewWriter(w)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetRowLine(false)
	table.SetColumnSeparator("  ")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, row := range rows {
		cells := make([]string, len(headers))
		for i, h := range headers {
			val := ""
			if v, ok := row[h]; ok {
				val = formatValue(v, showKeys)
				if !showKeys && isSensitiveKey(h) {
					val = maskString(val)
				}
			}
			cells[i] = val
		}
		table.Append(cells)
	}
	table.Render()
	return nil
}

// printAnySliceTable 输出非对象数组为单列值表格
func printAnySliceTable(w io.Writer, items []any, showKeys bool) error {
	if len(items) == 0 {
		fmt.Fprintln(w, "（空）")
		return nil
	}
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"#", "Value"})
	table.SetAutoWrapText(false)
	table.SetRowLine(false)
	table.SetColumnSeparator("  ")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for idx, item := range items {
		table.Append([]string{fmt.Sprintf("%d", idx+1), formatValue(item, showKeys)})
	}
	table.Render()
	return nil
}

// ============== 辅助函数 ==============

func formatValue(v any, showKeys bool) string {
	if v == nil {
		return "—"
	}
	switch val := v.(type) {
	case string:
		if !showKeys && isAPIKeyString(val) {
			return maskAPIKey(val)
		}
		if val == "" {
			return "—"
		}
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%v", val)
	case int, int64, int32:
		return fmt.Sprintf("%d", val)
	case bool:
		if val {
			return "✓"
		}
		return "✗"
	case []any:
		if len(val) == 0 {
			return "—"
		}
		parts := make([]string, 0, len(val))
		for _, item := range val {
			parts = append(parts, fmt.Sprintf("%v", item))
		}
		joined := strings.Join(parts, ", ")
		if len(joined) > 60 {
			return joined[:57] + "..."
		}
		return joined
	case []string:
		if len(val) == 0 {
			return "—"
		}
		if !showKeys {
			masked := make([]string, len(val))
			for i, k := range val {
				if isAPIKeyString(k) {
					masked[i] = maskAPIKey(k)
				} else {
					masked[i] = k
				}
			}
			val = masked
		}
		joined := strings.Join(val, ", ")
		if len(joined) > 60 {
			return joined[:57] + "..."
		}
		return joined
	case map[string]any:
		if len(val) == 0 {
			return "—"
		}
		jsonBytes, _ := json.Marshal(val)
		return string(jsonBytes)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// isSensitiveKey 判断是否为敏感字段名
func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	return strings.Contains(lower, "apikey") ||
		strings.Contains(lower, "api_key") ||
		strings.Contains(lower, "api-key") ||
		strings.Contains(lower, "apikeys") ||
		strings.Contains(lower, "secret") ||
		strings.Contains(lower, "token") ||
		strings.Contains(lower, "credential") ||
		strings.Contains(lower, "authorization")
}

// isAPIKeyString 判断是否为 API Key 字符串
func isAPIKeyString(s string) bool {
	// 已知前缀直接匹配
	if strings.HasPrefix(s, "sk-") || strings.HasPrefix(s, "AIzaSy") || strings.HasPrefix(s, "key-") {
		return true
	}
	// 长字符串启发式：纯字母数字+连字符/下划线且长度 >=20 可能是 API Key
	if len(s) >= 20 {
		alphanum := 0
		for _, c := range s {
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
				alphanum++
			}
		}
		if float64(alphanum)/float64(len(s)) > 0.9 {
			return true
		}
	}
	return false
}

// maskAPIKey 脱敏 API Key
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		if len(key) <= 4 {
			return key
		}
		return key[:4] + "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

// maskString 通用字符串脱敏
func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}

// maskSensitiveData 递归脱敏敏感数据
func maskSensitiveData(data any) any {
	switch v := data.(type) {
	case map[string]any:
		masked := make(map[string]any, len(v))
		for key, val := range v {
			if isSensitiveKey(key) {
				if str, ok := val.(string); ok {
					masked[key] = maskString(str)
					continue
				}
				if arr, ok := val.([]any); ok {
					maskedArr := make([]any, len(arr))
					for i, item := range arr {
						if str, ok := item.(string); ok {
							maskedArr[i] = maskAPIKey(str)
						} else {
							maskedArr[i] = item
						}
					}
					masked[key] = maskedArr
					continue
				}
			}
			masked[key] = maskSensitiveData(val)
		}
		return masked
	case []any:
		masked := make([]any, len(v))
		for i, item := range v {
			masked[i] = maskSensitiveData(item)
		}
		return masked
	default:
		return v
	}
}

// errJSON 降级输出 JSON
func errJSON(w io.Writer, data any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
