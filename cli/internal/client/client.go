// Package client 提供 CCX API 客户端（含认证、重试、超时、限流处理）
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ccx-cli/internal/errors"
	"ccx-cli/internal/version"
)

// Client CCX API 客户端
type Client struct {
	serverURL       string
	apiKey          string
	authType        string // "x-api-key" | "bearer" | "x-goog-api-key"
	httpClient      *http.Client
	requestTimeout  time.Duration
	tlsConfig       *tls.Config
	retryMax        int
	retryWait       time.Duration
	retryMaxWait    time.Duration
	verbose         bool
	userAgent       string
}

// ClientOption 客户端选项
type ClientOption func(*Client)

// WithRetry 设置重试参数
func WithRetry(max int) ClientOption {
	return func(c *Client) {
		c.retryMax = max
	}
}

// WithTimeout 设置超时
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.requestTimeout = timeout
		c.httpClient.Timeout = timeout
	}
}

// WithVerbose 设置 verbose 模式
func WithVerbose(verbose bool) ClientOption {
	return func(c *Client) {
		c.verbose = verbose
	}
}

// WithTLS 设置 TLS 选项
func WithTLS(caCert string, insecureSkipVerify bool) ClientOption {
	return func(c *Client) {
		if caCert != "" || insecureSkipVerify {
			c.tlsConfig = &tls.Config{
				InsecureSkipVerify: insecureSkipVerify,
			}
		}
	}
}

// NewClient 创建 CCX API 客户端
func NewClient(serverURL, apiKey string, opts ...ClientOption) *Client {
	// 去除尾部 /
	serverURL = strings.TrimRight(serverURL, "/")

	// 校验 URL scheme
	if !strings.HasPrefix(serverURL, "http://") && !strings.HasPrefix(serverURL, "https://") {
		serverURL = "http://" + serverURL
	}

	// 解析认证类型
	authType := "x-api-key" // 默认
	if strings.HasPrefix(apiKey, "AIzaSy") {
		authType = "x-goog-api-key"
	}

	c := &Client{
		serverURL:      serverURL,
		apiKey:         apiKey,
		authType:       authType,
		requestTimeout: 30 * time.Second,
		retryMax:       3,
		retryWait:      1 * time.Second,
		retryMaxWait:   4 * time.Second,
		userAgent:      "ccx-cli/" + version.Version,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	// 根据最终配置构建 Transport（使 dial timeout 跟随 --timeout）
	dialTimeout := c.requestTimeout / 2
	if dialTimeout < 5*time.Second {
		dialTimeout = 5 * time.Second
	}
	if dialTimeout > 30*time.Second {
		dialTimeout = 30 * time.Second
	}
	c.httpClient.Transport = &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		TLSClientConfig:     c.tlsConfig,
		DialContext: (&net.Dialer{
			Timeout:   dialTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	}

	return c
}

// ============== 认证设置 ==============

func (c *Client) setAuthHeader(req *http.Request) {
	switch c.authType {
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	case "x-goog-api-key":
		req.Header.Set("X-Goog-Api-Key", c.apiKey)
	default: // x-api-key
		req.Header.Set("X-Api-Key", c.apiKey)
	}
}

// ============== HTTP 方法 ==============

// Get 发送 GET 请求
func (c *Client) Get(path string, queryParams map[string]string) (*http.Response, error) {
	return c.doRequest("GET", path, queryParams, nil)
}

// Post 发送 POST 请求
func (c *Client) Post(path string, body any) (*http.Response, error) {
	return c.doRequest("POST", path, nil, body)
}

// Put 发送 PUT 请求
func (c *Client) Put(path string, body any) (*http.Response, error) {
	return c.doRequest("PUT", path, nil, body)
}

// Patch 发送 PATCH 请求
func (c *Client) Patch(path string, body any) (*http.Response, error) {
	return c.doRequest("PATCH", path, nil, body)
}

// Delete 发送 DELETE 请求
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.doRequest("DELETE", path, nil, nil)
}

// ============== 核心请求方法 ==============

func (c *Client) doRequest(method, path string, queryParams map[string]string, body any) (*http.Response, error) {
	// 构建 URL
	reqURL := c.serverURL + path
	if queryParams != nil {
		params := url.Values{}
		for k, v := range queryParams {
			if v != "" {
				params.Set(k, v)
			}
		}
		if len(params) > 0 {
			reqURL += "?" + params.Encode()
		}
	}

	// 序列化请求体
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, errors.NewWithDetail(errors.ExitCodeUserError,
				"序列化请求体失败", err.Error())
		}
		bodyReader = bytes.NewReader(jsonData)
		if c.verbose {
			fmt.Printf("[DEBUG] %s %s\n请求体: %s\n", method, reqURL, string(jsonData))
		}
	} else if c.verbose {
		fmt.Printf("[DEBUG] %s %s\n", method, reqURL)
	}

	// 创建请求
	req, err := http.NewRequest(method, reqURL, bodyReader)
	if err != nil {
		return nil, errors.NewWithDetail(errors.ExitCodeNetworkError,
			fmt.Sprintf("创建请求失败：%s %s", method, reqURL), err.Error())
	}

	// 设置请求头
	c.setAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	// 执行请求（含重试）
	var resp *http.Response
	var lastErr error

	for attempt := 0; attempt <= c.retryMax; attempt++ {
		if attempt > 0 {
			// 指数退避 + 抖动
			wait := c.retryWait * (1 << (attempt - 1))
			if wait > c.retryMaxWait {
				wait = c.retryMaxWait
			}
			jitter := time.Duration(rand.Int64N(int64(wait) / 4))
			wait += jitter

			if c.verbose {
				fmt.Printf("[DEBUG] 重试 %d/%d（等待 %v）...\n", attempt, c.retryMax, wait)
			}
			time.Sleep(wait)
		}

		// 每个重试都需重新创建 body
		if body != nil {
			jsonData, _ := json.Marshal(body)
			req.Body = io.NopCloser(bytes.NewReader(jsonData))
			req.ContentLength = int64(len(jsonData))
		}

		resp, lastErr = c.httpClient.Do(req)
		if lastErr == nil {
			// 检查是否需要重试
			if c.shouldRetry(resp.StatusCode, method) {
				resp.Body.Close()
				continue
			}
			break
		}

		// 网络错误重试
		if c.isNetworkError(lastErr) && attempt < c.retryMax {
			continue
		}
		break
	}

	if lastErr != nil {
		return nil, c.classifyError(lastErr)
	}

	if c.verbose {
		fmt.Printf("[DEBUG] ← %d %s\n", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

// shouldRetry 判断是否应重试
func (c *Client) shouldRetry(statusCode int, method string) bool {
	// 429 限流可重试
	if statusCode == http.StatusTooManyRequests {
		return true
	}
	// 非幂等操作（POST）不重试 502/503/504
	if method == "POST" {
		return false
	}
	// 502/503/504 网关错误可重试
	if statusCode == http.StatusBadGateway ||
		statusCode == http.StatusServiceUnavailable ||
		statusCode == http.StatusGatewayTimeout {
		return true
	}
	return false
}

// isNetworkError 判断是否为网络错误
func (c *Client) isNetworkError(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(net.Error); ok {
		return true
	}
	if strings.Contains(err.Error(), "connection refused") {
		return true
	}
	if strings.Contains(err.Error(), "no such host") {
		return true
	}
	if strings.Contains(err.Error(), "reset by peer") {
		return true
	}
	return false
}

// classifyError 将错误分类为 CLIError
func (c *Client) classifyError(err error) error {
	if err == nil {
		return nil
	}
	// 超时
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return errors.ErrTimeout()
	}
	// 连接失败
	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "no such host") ||
		strings.Contains(err.Error(), "dial tcp") {
		return errors.ErrConnectionFailed(c.serverURL, err)
	}
	return errors.NewWithDetail(errors.ExitCodeNetworkError, "网络请求失败", err.Error())
}

// maxResponseSize 最大响应体大小 (10MB)
const maxResponseSize = 10 << 20

// ============== 响应解析 ==============

// DecodeResponse 解码响应体到目标结构
func DecodeResponse(resp *http.Response, target any) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return errors.NewWithDetail(errors.ExitCodeServerError, "读取响应体失败", err.Error())
	}

	if resp.StatusCode >= 400 {
		return handleErrorResponse(resp.StatusCode, body)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return errors.NewWithDetail(errors.ExitCodeServerError,
			"解析响应体失败", fmt.Sprintf("body: %s, error: %v", string(body), err))
	}
	return nil
}

// ReadBody 读取响应体（直接返回字节）
func ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
}

// handleErrorResponse 处理错误响应
func handleErrorResponse(statusCode int, body []byte) error {
	// 尝试解析错误体
	var errResp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
		Detail  string `json:"detail"`
	}
	if err := json.Unmarshal(body, &errResp); err != nil {
		errResp.Error = string(body)
	}

	detail := errResp.Error
	if errResp.Message != "" {
		detail = errResp.Message
	}
	if errResp.Detail != "" {
		detail += ": " + errResp.Detail
	}

	switch statusCode {
	case http.StatusUnauthorized:
		return errors.ErrAuthFailed(detail)
	case http.StatusNotFound:
		return errors.NewWithHTTP(errors.ExitCodeUserError,
			fmt.Sprintf("✗ 资源未找到 (404)：%s", detail), statusCode)
	case http.StatusConflict:
		return errors.ErrConflict(detail)
	case http.StatusBadRequest:
		return errors.NewWithHTTP(errors.ExitCodeUserError,
			fmt.Sprintf("✗ 请求参数错误：%s", detail), statusCode)
	case http.StatusTooManyRequests:
		return errors.ErrRateLimited("")
	case http.StatusInternalServerError:
		return errors.ErrServerError(detail)
	default:
		if statusCode >= 500 {
			return errors.NewWithHTTP(errors.ExitCodeServerError,
				fmt.Sprintf("✗ 服务器错误 (%d)：%s", statusCode, detail), statusCode)
		}
		return errors.NewWithHTTP(errors.ExitCodeUserError,
			fmt.Sprintf("✗ 请求失败 (%d)：%s", statusCode, detail), statusCode)
	}
}

// ============== API 路径构建 ==============

// ChannelPath 构建渠道 API 路径
// 例如: /api/messages/channels, /api/chat/channels
func ChannelPath(channelType string) string {
	return fmt.Sprintf("/api/%s/channels", channelType)
}

// ChannelIDPath 构建指定渠道的 API 路径
func ChannelIDPath(channelType string, id int) string {
	return fmt.Sprintf("/api/%s/channels/%d", channelType, id)
}

// ChannelKeyPath 构建密钥 API 路径
func ChannelKeyPath(channelType string, id int) string {
	return fmt.Sprintf("/api/%s/channels/%d/keys", channelType, id)
}

// ChannelKeyIDPath 构建指定密钥的 API 路径
func ChannelKeyIDPath(channelType string, id int, apiKey string) string {
	return fmt.Sprintf("/api/%s/channels/%d/keys/%s", channelType, id, url.PathEscape(apiKey))
}

// ChannelKeyMovePath 构建密钥移动路径
func ChannelKeyMovePath(channelType string, id int, apiKey string, position string) string {
	return fmt.Sprintf("/api/%s/channels/%d/keys/%s/%s", channelType, id, url.PathEscape(apiKey), url.PathEscape(position))
}

// ChannelKeyRestorePath 构建密钥恢复路径
func ChannelKeyRestorePath(channelType string, id int) string {
	return fmt.Sprintf("/api/%s/channels/%d/keys/restore", channelType, id)
}

// ChannelMappingPath 构建模型映射路径
func ChannelMappingPath(channelType string, id int) string {
	return fmt.Sprintf("/api/%s/channels/%d/mappings", channelType, id)
}

// ChannelReorderPath 构建重排序路径
func ChannelReorderPath(channelType string) string {
	return fmt.Sprintf("/api/%s/channels/reorder", channelType)
}

// ChannelStatusPath 构建状态设置路径
func ChannelStatusPath(channelType string, id int) string {
	return fmt.Sprintf("/api/%s/channels/%d/status", channelType, id)
}

// ChannelResumePath 构建恢复路径
func ChannelResumePath(channelType string, id int) string {
	return fmt.Sprintf("/api/%s/channels/%d/resume", channelType, id)
}

// ChannelPromotionPath 构建促销期路径
func ChannelPromotionPath(channelType string, id int) string {
	return fmt.Sprintf("/api/%s/channels/%d/promotion", channelType, id)
}

// ChannelPingPath 构建 Ping 路径
func ChannelPingPath(channelType string, id int) string {
	return fmt.Sprintf("/api/%s/ping/%d", channelType, id)
}

// ChannelPingAllPath 构建全局 Ping 路径
func ChannelPingAllPath(channelType string) string {
	return fmt.Sprintf("/api/%s/ping", channelType)
}

// HealthPath 健康检查路径
func HealthPath() string {
	return "/health"
}

// ChannelModelsPath 模型列表路径
func ChannelModelsPath(channelType string, id int) string {
	return fmt.Sprintf("/api/%s/channels/%d/models", channelType, id)
}

// ChannelMetricsPath 指标路径
func ChannelMetricsPath(channelType string) string {
	return fmt.Sprintf("/api/%s/channels/metrics", channelType)
}

// ChannelDashboardPath 仪表盘路径
func ChannelDashboardPath() string {
	return "/api/messages/channels/dashboard"
}

// ChannelLogsPath 渠道日志路径
func ChannelLogsPath(channelType string, id int) string {
	return fmt.Sprintf("/api/%s/channels/%d/logs", channelType, id)
}

// ChannelSchedulerStatsPath 调度器统计路径
func ChannelSchedulerStatsPath() string {
	return "/api/messages/channels/scheduler/stats"
}

// ChannelCapabilityTestPath 能力测试路径
func ChannelCapabilityTestPath(channelType string, id int) string {
	return fmt.Sprintf("/api/%s/channels/%d/capability-test", channelType, id)
}

// ChannelCapabilitySnapshotPath 能力快照路径
func ChannelCapabilitySnapshotPath(channelType string, id int) string {
	return fmt.Sprintf("/api/%s/channels/%d/capability-snapshot", channelType, id)
}

// SettingsFuzzyPath Fuzzy 模式设置路径
func SettingsFuzzyPath() string {
	return "/api/settings/fuzzy-mode"
}

// SettingsCircuitBreakerPath 熔断器设置路径
func SettingsCircuitBreakerPath() string {
	return "/api/settings/circuit-breaker"
}

// SettingsImageTurnLimitPath 图片轮次限制设置路径
func SettingsImageTurnLimitPath() string {
	return "/api/settings/historical-image-turn-limit"
}

// SettingsConversationsPath 对话设置路径
func SettingsConversationsPath() string {
	return "/api/conversations/settings"
}

// ConversationsListPath 对话列表路径
func ConversationsListPath() string {
	return "/api/conversations"
}

// ConversationOverridePath 对话覆盖路径
func ConversationOverridePath(id string) string {
	return fmt.Sprintf("/api/conversations/%s/override", url.PathEscape(id))
}

// ConfigSavePath 配置保存路径
func ConfigSavePath() string {
	return "/admin/config/save"
}
