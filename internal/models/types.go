// Package models 定义与 CCX 后端对齐的数据结构
package models

import "time"

// ============== 5 种渠道类型 ==============

// ChannelType 渠道类型
type ChannelType string

const (
	ChannelTypeMessages  ChannelType = "messages"
	ChannelTypeResponses ChannelType = "responses"
	ChannelTypeChat      ChannelType = "chat"
	ChannelTypeGemini    ChannelType = "gemini"
	ChannelTypeImages    ChannelType = "images"
)

// AllChannelTypes 所有渠道类型列表
var AllChannelTypes = []ChannelType{
	ChannelTypeMessages,
	ChannelTypeResponses,
	ChannelTypeChat,
	ChannelTypeGemini,
	ChannelTypeImages,
}

// ChannelTypeToAPI 渠道类型到 API 路径段的映射
func ChannelTypeToPathSegment(t ChannelType) string {
	return string(t)
}

// ChannelTypeToUpstreamField 渠道类型到 Config 中的字段名
func ChannelTypeToUpstreamField(t ChannelType) string {
	switch t {
	case ChannelTypeMessages:
		return "upstream"
	case ChannelTypeResponses:
		return "responsesUpstream"
	case ChannelTypeChat:
		return "chatUpstream"
	case ChannelTypeGemini:
		return "geminiUpstream"
	case ChannelTypeImages:
		return "imagesUpstream"
	default:
		return "upstream"
	}
}

// ============== 上游配置（与 CCX backend-go/internal/config/config.go 对齐） ==============

// UpstreamConfig 上游渠道配置
type UpstreamConfig struct {
	BaseURL            string              `json:"baseUrl"`
	BaseURLs           []string            `json:"baseUrls,omitempty"`
	APIKeys            []string            `json:"apiKeys"`
	APIKeyConfigs      []APIKeyConfig      `json:"apiKeyConfigs,omitempty"`
	HistoricalAPIKeys  []string            `json:"historicalApiKeys,omitempty"`
	DisabledAPIKeys    []DisabledKeyInfo   `json:"disabledApiKeys,omitempty"`
	ServiceType        string              `json:"serviceType"` // gemini, openai, claude
	AuthHeader         string              `json:"authHeader,omitempty"`
	Name               string              `json:"name,omitempty"`
	Description        string              `json:"description,omitempty"`
	Website            string              `json:"website,omitempty"`
	InsecureSkipVerify bool                `json:"insecureSkipVerify,omitempty"`
	ModelMapping       map[string]string   `json:"modelMapping,omitempty"`
	ModelCapabilities  map[string]any      `json:"modelCapabilities,omitempty"`
	PromotionUntil     *time.Time          `json:"promotionUntil,omitempty"`
	LowQuality         bool                `json:"lowQuality,omitempty"`
	Priority           int                 `json:"priority"`
	Status             string              `json:"status"` // active, suspended, disabled
	CustomHeaders      map[string]string   `json:"customHeaders,omitempty"`
	ProxyURL           string              `json:"proxyUrl,omitempty"`
	SupportedModels    []string            `json:"supportedModels,omitempty"`
	RoutePrefix        string              `json:"routePrefix,omitempty"`
	NoVision           bool                `json:"noVision,omitempty"`
	NoVisionModels     []string            `json:"noVisionModels,omitempty"`
	VisionFallbackModel string             `json:"visionFallbackModel,omitempty"`
	HistoricalImageTurnLimit int           `json:"historicalImageTurnLimit,omitempty"`
	RequestTimeoutMs        int             `json:"requestTimeoutMs,omitempty"`
	ResponseHeaderTimeoutMs int             `json:"responseHeaderTimeoutMs,omitempty"`
	RateLimitRPM            int             `json:"rateLimitRpm,omitempty"`
	RateLimitWindowMinutes  int             `json:"rateLimitWindowMinutes,omitempty"`
	RateLimitMaxConcurrent  int             `json:"rateLimitMaxConcurrent,omitempty"`
	// Claude 协议兼容
	PassbackReasoningContent      bool `json:"passbackReasoningContent,omitempty"`
	PassbackThinkingBlocks        bool `json:"passbackThinkingBlocks,omitempty"`
	StripEmptyTextBlocks          bool `json:"stripEmptyTextBlocks,omitempty"`
	NormalizeSystemRoleToTopLevel bool `json:"normalizeSystemRoleToTopLevel,omitempty"`
	InjectDummyThoughtSignature   bool `json:"injectDummyThoughtSignature,omitempty"`
	StripThoughtSignature         bool `json:"stripThoughtSignature,omitempty"`
	// 其他布尔开关
	NormalizeNonstandardChatRoles bool `json:"normalizeNonstandardChatRoles,omitempty"`
	CodexNativeToolPassthrough    bool `json:"codexNativeToolPassthrough,omitempty"`
	StripImageGenerationTool      bool `json:"stripImageGenerationTool,omitempty"`
	ConvertImageURLToB64JSON      bool `json:"convertImageUrlToB64Json,omitempty"`
}

// APIKeyConfig API Key 附加配置
type APIKeyConfig struct {
	Key    string `json:"key"`
	Name   string `json:"name,omitempty"`
	Weight int    `json:"weight,omitempty"`
}

// DisabledKeyInfo 被拉黑的 API Key 信息
type DisabledKeyInfo struct {
	Key        string `json:"key"`
	Reason     string `json:"reason"`
	Message    string `json:"message"`
	DisabledAt string `json:"disabledAt"`
	RecoverAt  string `json:"recoverAt,omitempty"`
}

// ============== API 响应结构 ==============

// ChannelListResponse 渠道列表响应
type ChannelListResponse struct {
	Channels []ChannelView `json:"channels"`
}

// ChannelView 渠道视图（来自 backend-go BuildChannelView）
type ChannelView struct {
	Index            int               `json:"index"`
	Name             string            `json:"name"`
	ServiceType      string            `json:"serviceType"`
	BaseURL          string            `json:"baseUrl"`
	BaseURLs         []string          `json:"baseUrls"`
	APIKeys          []string          `json:"apiKeys"`
	APIKeyConfigs    []APIKeyConfig    `json:"apiKeyConfigs"`
	Description      string            `json:"description"`
	Website          string            `json:"website"`
	ModelMapping     map[string]string `json:"modelMapping"`
	Status           string            `json:"status"`
	AdminState       string            `json:"adminState"`
	EffectiveState   string            `json:"effectiveState"`
	RuntimeState     string            `json:"runtimeState"`
	Priority         int               `json:"priority"`
	PromotionUntil   *time.Time        `json:"promotionUntil"`
	LowQuality       bool              `json:"lowQuality"`
	ProxyURL         string            `json:"proxyUrl"`
	RoutePrefix      string            `json:"routePrefix"`
	SupportedModels  []string          `json:"supportedModels"`
	DisabledAPIKeys  []DisabledKeyInfo `json:"disabledApiKeys"`
	CustomHeaders    map[string]string `json:"customHeaders"`
	// 熔断/指标相关
	Latency  *int64 `json:"latency"`
	Requests *int   `json:"requests,omitempty"`
	Errors   *int   `json:"errors,omitempty"`
}

// ErrorResponse API 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse 通用成功响应
type SuccessResponse struct {
	Message string `json:"message"`
	Success *bool  `json:"success,omitempty"`
}

// ============== 健康检查 ==============

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Uptime    float64 `json:"uptime"`
	Mode      string `json:"mode"`
	Version   struct {
		Version   string `json:"version"`
		BuildTime string `json:"buildTime"`
		GitCommit string `json:"gitCommit"`
	} `json:"version"`
	Config struct {
		UpstreamCount int `json:"upstreamCount"`
	} `json:"config"`
}

// ============== Ping 响应 ==============

// PingResult 单个渠道 Ping 结果
type PingResult struct {
	Success    bool   `json:"success"`
	Latency    int64  `json:"latency"`
	Status     string `json:"status"`
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
	Name       string `json:"name,omitempty"`
	Index      int    `json:"index,omitempty"`
}

// ============== 设置相关 ==============

// FuzzyModeResponse Fuzzy 模式响应
type FuzzyModeResponse struct {
	FuzzyModeEnabled bool `json:"fuzzyModeEnabled"`
}

// FuzzyModeRequest Fuzzy 模式请求
type FuzzyModeRequest struct {
	Enabled bool `json:"enabled"`
}

// CircuitBreakerResponse 熔断器配置响应
type CircuitBreakerResponse struct {
	WindowSize                   *int     `json:"windowSize"`
	FailureThreshold             *float64 `json:"failureThreshold"`
	ConsecutiveFailuresThreshold *int     `json:"consecutiveFailuresThreshold"`
	RequestTimeoutMs             *int     `json:"requestTimeoutMs"`
	ResponseHeaderTimeoutMs      *int     `json:"responseHeaderTimeoutMs"`
	StreamFirstContentTimeoutMs  *int     `json:"streamFirstContentTimeoutMs"`
	StreamInactivityTimeoutMs    *int     `json:"streamInactivityTimeoutMs"`
	StreamToolCallIdleTimeoutMs  *int     `json:"streamToolCallIdleTimeoutMs"`
}

// CircuitBreakerRequest 熔断器更新请求
type CircuitBreakerRequest struct {
	WindowSize                   *int     `json:"windowSize,omitempty"`
	FailureThreshold             *float64 `json:"failureThreshold,omitempty"`
	ConsecutiveFailuresThreshold *int     `json:"consecutiveFailuresThreshold,omitempty"`
	RequestTimeoutMs             *int     `json:"requestTimeoutMs,omitempty"`
	ResponseHeaderTimeoutMs      *int     `json:"responseHeaderTimeoutMs,omitempty"`
	StreamFirstContentTimeoutMs  *int     `json:"streamFirstContentTimeoutMs,omitempty"`
	StreamInactivityTimeoutMs    *int     `json:"streamInactivityTimeoutMs,omitempty"`
	StreamToolCallIdleTimeoutMs  *int     `json:"streamToolCallIdleTimeoutMs,omitempty"`
}

// HistoricalImageTurnLimitResponse 图片轮次限制响应
type HistoricalImageTurnLimitResponse struct {
	HistoricalImageTurnLimit int `json:"historicalImageTurnLimit"`
}

// HistoricalImageTurnLimitRequest 图片轮次限制请求
type HistoricalImageTurnLimitRequest struct {
	Limit int `json:"limit"`
}

// ============== 对话相关 ==============

// ConversationSettings 对话设置
type ConversationSettings struct {
	Enabled *bool `json:"enabled,omitempty"`
	MaxAge  *int  `json:"maxAge,omitempty"`
}

// ============== 密钥管理请求 ==============

// AddKeyRequest 添加密钥请求
type AddKeyRequest struct {
	APIKey string `json:"apiKey"`
}

// RestoreKeyRequest 恢复密钥请求
type RestoreKeyRequest struct {
	APIKey string `json:"apiKey"`
}

// ============== 渠道操作请求 ==============

// ReorderRequest 渠道重排序请求
type ReorderRequest struct {
	Order []int `json:"order"`
}

// PromotionRequest 促销期请求
type PromotionRequest struct {
	Duration int `json:"duration"` // 秒数，<=0 时清除促销期
}

// StatusUpdateRequest 状态更新请求
type StatusUpdateRequest struct {
	Status string `json:"status"`
}

// ============== 渠道指标 & 监控 ==============

// ChannelMetricsResponse 渠道性能指标（动态结构，使用 map 灵活适配后端返回）
type ChannelMetricsResponse map[string]any

// ChannelLogEntry 单条日志条目
type ChannelLogEntry struct {
	Time        string `json:"time"`
	Level       string `json:"level"`
	Message     string `json:"message"`
	Source      string `json:"source,omitempty"`
	ChannelName string `json:"channelName,omitempty"`
	Latency     int64  `json:"latency,omitempty"`
	StatusCode  int    `json:"statusCode,omitempty"`
	Model       string `json:"model,omitempty"`
}

// ChannelLogsResponse 渠道日志响应
type ChannelLogsResponse []ChannelLogEntry

// ChannelDashboardResponse 仪表盘响应（动态结构）
type ChannelDashboardResponse map[string]any

// SchedulerStatsResponse 调度器统计响应（动态结构）
type SchedulerStatsResponse map[string]any

// CapabilityTestResponse 能力测试响应（动态结构）
type CapabilityTestResponse map[string]any

// ============== Config Apply ==============

// FullConfig 完整配置（用于 config show/apply/backup/restore）
type FullConfig struct {
	Upstream                   []UpstreamConfig              `json:"upstream"`
	ResponsesUpstream          []UpstreamConfig              `json:"responsesUpstream"`
	GeminiUpstream             []UpstreamConfig              `json:"geminiUpstream"`
	ChatUpstream               []UpstreamConfig              `json:"chatUpstream"`
	ImagesUpstream             []UpstreamConfig              `json:"imagesUpstream"`
	FuzzyModeEnabled           bool                          `json:"fuzzyModeEnabled"`
	HistoricalImageTurnLimit   int                           `json:"historicalImageTurnLimit"`
	CircuitBreaker             *CircuitBreakerConfig         `json:"circuitBreaker,omitempty"`
	UpstreamModelCapabilities  map[string]UpstreamModelCapability `json:"upstreamModelCapabilities,omitempty"`
}

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	WindowSize                   *int     `json:"windowSize,omitempty"`
	FailureThreshold             *float64 `json:"failureThreshold,omitempty"`
	ConsecutiveFailuresThreshold *int     `json:"consecutiveFailuresThreshold,omitempty"`
	RequestTimeoutMs             *int     `json:"requestTimeoutMs,omitempty"`
	ResponseHeaderTimeoutMs      *int     `json:"responseHeaderTimeoutMs,omitempty"`
	StreamFirstContentTimeoutMs  *int     `json:"streamFirstContentTimeoutMs,omitempty"`
	StreamInactivityTimeoutMs    *int     `json:"streamInactivityTimeoutMs,omitempty"`
	StreamToolCallIdleTimeoutMs  *int     `json:"streamToolCallIdleTimeoutMs,omitempty"`
}

// UpstreamModelCapability 模型能力描述
type UpstreamModelCapability struct {
	ContextWindowTokens int    `json:"contextWindowTokens,omitempty"`
	MaxOutputTokens     int    `json:"maxOutputTokens,omitempty"`
	Provider            string `json:"provider,omitempty"`
	DisplayName         string `json:"displayName,omitempty"`
}
