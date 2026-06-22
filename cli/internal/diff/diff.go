// Package diff 提供 JSON diff 工具（给 config apply 使用）
package diff

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// ChangeType 变更类型
type ChangeType int

const (
	ChangeAdd    ChangeType = iota // 新增
	ChangeRemove                   // 删除
	ChangeModify                   // 修改
)

// Change 描述一项变更
type Change struct {
	Type     ChangeType `json:"type"`
	Path     string     `json:"path"`     // JSON 路径，如 "upstream[0].name"
	OldValue any        `json:"oldValue"` // 旧值（nil 表示新增）
	NewValue any        `json:"newValue"` // 新值（nil 表示删除）
}

// String 返回变更的可读描述
func (c Change) String() string {
	switch c.Type {
	case ChangeAdd:
		return fmt.Sprintf("+ %s = %v", c.Path, c.NewValue)
	case ChangeRemove:
		return fmt.Sprintf("- %s (原值: %v)", c.Path, c.OldValue)
	case ChangeModify:
		return fmt.Sprintf("~ %s: %v → %v", c.Path, c.OldValue, c.NewValue)
	default:
		return ""
	}
}

// DiffResult diff 结果
type DiffResult struct {
	HasChanges bool
	Changes    []Change
}

// HasBreakingChanges 是否有破坏性变更（删除或修改）
func (d *DiffResult) HasBreakingChanges() bool {
	for _, c := range d.Changes {
		if c.Type == ChangeRemove {
			return true
		}
	}
	return false
}

// Summary 返回变更摘要
func (d *DiffResult) Summary() string {
	if !d.HasChanges {
		return "无变更"
	}
	adds, removes, modifies := 0, 0, 0
	for _, c := range d.Changes {
		switch c.Type {
		case ChangeAdd:
			adds++
		case ChangeRemove:
			removes++
		case ChangeModify:
			modifies++
		}
	}
	parts := make([]string, 0, 3)
	if adds > 0 {
		parts = append(parts, fmt.Sprintf("+%d 新增", adds))
	}
	if removes > 0 {
		parts = append(parts, fmt.Sprintf("-%d 删除", removes))
	}
	if modifies > 0 {
		parts = append(parts, fmt.Sprintf("~%d 修改", modifies))
	}
	return strings.Join(parts, ", ")
}

// CompareJSON 比较两个 JSON 对象的差异
func CompareJSON(oldObj, newObj any) *DiffResult {
	result := &DiffResult{}
	compareValues("", oldObj, newObj, result)
	result.HasChanges = len(result.Changes) > 0
	return result
}

// CompareJSONBytes 比较两个 JSON 字节数组的差异
func CompareJSONBytes(oldData, newData []byte) (*DiffResult, error) {
	var oldObj, newObj any

	if len(oldData) > 0 {
		if err := json.Unmarshal(oldData, &oldObj); err != nil {
			return nil, fmt.Errorf("解析旧 JSON 失败：%w", err)
		}
	}
	if len(newData) > 0 {
		if err := json.Unmarshal(newData, &newObj); err != nil {
			return nil, fmt.Errorf("解析新 JSON 失败：%w", err)
		}
	}

	return CompareJSON(oldObj, newObj), nil
}

// compareValues 递归比较两个值
func compareValues(path string, oldVal, newVal any, result *DiffResult) {
	// 两者都为 nil，无变化
	if oldVal == nil && newVal == nil {
		return
	}
	// 新增
	if oldVal == nil && newVal != nil {
		addChange(ChangeAdd, path, nil, newVal, result)
		return
	}
	// 删除
	if oldVal != nil && newVal == nil {
		addChange(ChangeRemove, path, oldVal, nil, result)
		return
	}

	// 类型不同，视为修改
	if reflect.TypeOf(oldVal) != reflect.TypeOf(newVal) {
		addChange(ChangeModify, path, oldVal, newVal, result)
		return
	}

	switch v := oldVal.(type) {
	case map[string]any:
		compareMaps(path, v, newVal.(map[string]any), result)
	case []any:
		compareSlices(path, v, newVal.([]any), result)
	default:
		// 基本类型比较
		if oldVal != newVal {
			addChange(ChangeModify, path, oldVal, newVal, result)
		}
	}
}

// compareMaps 比较两个 map
func compareMaps(path string, oldMap, newMap map[string]any, result *DiffResult) {
	// 收集所有 key
	keys := make(map[string]bool)
	for k := range oldMap {
		keys[k] = true
	}
	for k := range newMap {
		keys[k] = true
	}

	// 排序以保持输出稳定
	sortedKeys := make([]string, 0, len(keys))
	for k := range keys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, k := range sortedKeys {
		childPath := path
		if childPath != "" {
			childPath += "." + k
		} else {
			childPath = k
		}

		oldChild, oldExists := oldMap[k]
		newChild, newExists := newMap[k]

		if oldExists && !newExists {
			addChange(ChangeRemove, childPath, oldChild, nil, result)
		} else if !oldExists && newExists {
			addChange(ChangeAdd, childPath, nil, newChild, result)
		} else {
			compareValues(childPath, oldChild, newChild, result)
		}
	}
}

// compareSlices 比较两个切片（仅当长度不同时视为变更，或逐元素比较）
func compareSlices(path string, oldSlice, newSlice []any, result *DiffResult) {
	if len(oldSlice) != len(newSlice) {
		addChange(ChangeModify, path, oldSlice, newSlice, result)
		return
	}
	for i := 0; i < len(oldSlice); i++ {
		childPath := fmt.Sprintf("%s[%d]", path, i)
		compareValues(childPath, oldSlice[i], newSlice[i], result)
	}
}

func addChange(changeType ChangeType, path string, oldVal, newVal any, result *DiffResult) {
	// 简化值的显示
	if oldVal != nil {
		oldVal = simplifyValue(oldVal)
	}
	if newVal != nil {
		newVal = simplifyValue(newVal)
	}
	result.Changes = append(result.Changes, Change{
		Type:     changeType,
		Path:     path,
		OldValue: oldVal,
		NewValue: newVal,
	})
}

// simplifyValue 简化值的显示（长字符串截断等）
func simplifyValue(val any) any {
	switch v := val.(type) {
	case string:
		if len(v) > 60 {
			return v[:57] + "..."
		}
		return v
	case map[string]any:
		return fmt.Sprintf("{...} (%d fields)", len(v))
	case []any:
		return fmt.Sprintf("[...] (%d items)", len(v))
	default:
		return v
	}
}
