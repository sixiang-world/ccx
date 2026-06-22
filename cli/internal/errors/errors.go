// Package errors 定义自定义错误类型与退出码
package errors

import (
	"fmt"
)

// 退出码定义（参考设计文档 §9.3）
const (
	ExitCodeSuccess       = 0  // 成功
	ExitCodeUserError     = 1  // 用户错误（参数错误、认证失败、资源不存在）
	ExitCodeServerError   = 2  // 服务器错误（后端异常、配置保存失败）
	ExitCodeNetworkError  = 3  // 网络错误（连接失败、DNS 解析失败）
	ExitCodeTimeout       = 4  // 超时
	ExitCodeRateLimit     = 5  // 限流（HTTP 429）
	ExitCodeConfigInvalid = 6  // 配置校验错误（config apply 校验失败）
)

// CLIError CLI 自定义错误
type CLIError struct {
	Code     int    // 退出码
	Message  string // 用户可见的消息
	Detail   string // 详细错误信息（仅 --verbose 显示）
	HTTPCode int    // 关联的 HTTP 状态码（0 表示无）
	Action   string // 修复建议
}

func (e *CLIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	}
	return e.Message
}

// New 创建 CLIError
func New(code int, message string) *CLIError {
	return &CLIError{Code: code, Message: message}
}

// NewWithDetail 创建带详细信息的 CLIError
func NewWithDetail(code int, message, detail string) *CLIError {
	return &CLIError{Code: code, Message: message, Detail: detail}
}

// NewWithHTTP 创建关联 HTTP 状态的 CLIError
func NewWithHTTP(code int, message string, httpCode int) *CLIError {
	return &CLIError{Code: code, Message: message, HTTPCode: httpCode}
}

// NewWithAction 创建带修复建议的 CLIError
func NewWithAction(code int, message, action string) *CLIError {
	return &CLIError{Code: code, Message: message, Action: action}
}

// Wrap 包装已有 error 为 CLIError
func Wrap(code int, err error) *CLIError {
	if err == nil {
		return nil
	}
	if cliErr, ok := err.(*CLIError); ok {
		return cliErr
	}
	return &CLIError{Code: code, Message: err.Error()}
}

// 常用错误构造器

// ErrAuthFailed 认证失败
func ErrAuthFailed(detail string) *CLIError {
	return &CLIError{
		Code:    ExitCodeUserError,
		Message: "✗ 认证失败：请检查 --key 或 CCX_API_KEY",
		Detail:  detail,
		Action:  "确认 ADMIN_ACCESS_KEY 或 PROXY_ACCESS_KEY 已在服务端正确配置",
	}
}

// ErrNotFound 资源未找到
func ErrNotFound(name string) *CLIError {
	return &CLIError{
		Code:    ExitCodeUserError,
		Message: fmt.Sprintf("✗ 未找到名为 %q 的渠道", name),
	}
}

// ErrConflict 资源冲突
func ErrConflict(msg string) *CLIError {
	return &CLIError{
		Code:    ExitCodeUserError,
		Message: fmt.Sprintf("✗ 资源冲突：%s", msg),
	}
}

// ErrRateLimited 限流
func ErrRateLimited(retryAfter string) *CLIError {
	msg := "✗ 请求过于频繁"
	if retryAfter != "" {
		msg += fmt.Sprintf("，请在 %s 后重试", retryAfter)
	}
	return &CLIError{
		Code:    ExitCodeRateLimit,
		Message: msg,
	}
}

// ErrServerError 服务器错误
func ErrServerError(detail string) *CLIError {
	return &CLIError{
		Code:    ExitCodeServerError,
		Message: fmt.Sprintf("✗ 服务器内部错误：%s", detail),
	}
}

// ErrConnectionFailed 连接失败
func ErrConnectionFailed(url string, err error) *CLIError {
	return &CLIError{
		Code:    ExitCodeNetworkError,
		Message: fmt.Sprintf("✗ 无法连接到服务器 %s", url),
		Detail:  err.Error(),
		Action:  "检查服务器地址是否正确，以及服务是否正常运行",
	}
}

// ErrTimeout 超时
func ErrTimeout() *CLIError {
	return &CLIError{
		Code:    ExitCodeTimeout,
		Message: "✗ 请求超时",
	}
}

// ErrConfigInvalid 配置校验失败
func ErrConfigInvalid(msg string) *CLIError {
	return &CLIError{
		Code:    ExitCodeConfigInvalid,
		Message: fmt.Sprintf("✗ 配置校验失败：%s", msg),
	}
}

// FormatCLIError 格式化 CLI 错误输出
func FormatCLIError(err error) string {
	if err == nil {
		return ""
	}
	if cliErr, ok := err.(*CLIError); ok {
		s := cliErr.Message
		if cliErr.Action != "" {
			s += "\n  → 建议：" + cliErr.Action
		}
		return s
	}
	return fmt.Sprintf("✗ 错误：%v", err)
}
