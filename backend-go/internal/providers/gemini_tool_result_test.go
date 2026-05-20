package providers

import (
	"testing"

	"github.com/BenedictKing/ccx/internal/config"
	"github.com/BenedictKing/ccx/internal/types"
	"github.com/stretchr/testify/assert"
)

// makeToolUseAssistant 构造一条只带 tool_use 的 assistant 消息，
// 配合后续 tool_result 模拟真实的多轮历史。
func makeToolUseAssistant(toolUseID, name string, input map[string]interface{}) types.ClaudeMessage {
	return types.ClaudeMessage{
		Role: "assistant",
		Content: []interface{}{
			map[string]interface{}{
				"type":  "tool_use",
				"id":    toolUseID,
				"name":  name,
				"input": input,
			},
		},
	}
}

// extractToolResultMsg 取出转换后第二条 user 消息（functionResponse 所在）。
func extractToolResultMsg(t *testing.T, geminiMsgs []map[string]interface{}) map[string]interface{} {
	t.Helper()
	assert.Len(t, geminiMsgs, 2, "应同时产出 assistant(model) 和 user(functionResponse) 两条消息")
	return geminiMsgs[1]
}

func TestGeminiProvider_ConvertMessages_ToolResultArray(t *testing.T) {
	provider := &GeminiProvider{}

	// 测试场景：tool_result 的 content 是一个 Content Blocks 数组。
	// 关键断言：functionResponse.name 必须等于前面 functionCall.name（函数名 get_weather），
	// 否则 Gemini 无法匹配，会沉默返回空内容。
	msgs := []types.ClaudeMessage{
		makeToolUseAssistant("toolu_0", "get_weather", map[string]interface{}{"location": "Tokyo"}),
		{
			Role: "user",
			Content: []interface{}{
				map[string]interface{}{
					"type":        "tool_result",
					"tool_use_id": "toolu_0",
					"content": []interface{}{
						map[string]interface{}{"type": "text", "text": "Tokyo is sunny."},
						map[string]interface{}{"type": "text", "text": "Temperature is 22C."},
					},
				},
			},
		},
	}

	geminiMsgs := provider.convertMessages(msgs)
	toolResultMsg := extractToolResultMsg(t, geminiMsgs)
	assert.Equal(t, "user", toolResultMsg["role"])

	parts, ok := toolResultMsg["parts"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, parts, 1)

	part, ok := parts[0].(map[string]interface{})
	assert.True(t, ok)

	funcResp, ok := part["functionResponse"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "get_weather", funcResp["name"])

	response, ok := funcResp["response"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Tokyo is sunny.\nTemperature is 22C.", response["result"])
}

func TestGeminiProvider_ConvertMessages_ToolResultString(t *testing.T) {
	provider := &GeminiProvider{}

	// 测试场景：tool_result 的 content 是一个简单字符串
	msgs := []types.ClaudeMessage{
		makeToolUseAssistant("toolu_1", "get_weather", map[string]interface{}{"location": "Tokyo"}),
		{
			Role: "user",
			Content: []interface{}{
				map[string]interface{}{
					"type":        "tool_result",
					"tool_use_id": "toolu_1",
					"content":     "Tokyo is sunny.",
				},
			},
		},
	}

	geminiMsgs := provider.convertMessages(msgs)
	toolResultMsg := extractToolResultMsg(t, geminiMsgs)

	parts, ok := toolResultMsg["parts"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, parts, 1)

	part, ok := parts[0].(map[string]interface{})
	assert.True(t, ok)

	funcResp, ok := part["functionResponse"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "get_weather", funcResp["name"])

	response, ok := funcResp["response"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Tokyo is sunny.", response["result"])
}

func TestGeminiProvider_ConvertMessages_ToolResultObject(t *testing.T) {
	provider := &GeminiProvider{}

	// 测试场景：tool_result 的 content 是一个 JSON 对象
	msgs := []types.ClaudeMessage{
		makeToolUseAssistant("toolu_2", "get_weather", map[string]interface{}{"location": "Tokyo"}),
		{
			Role: "user",
			Content: []interface{}{
				map[string]interface{}{
					"type":        "tool_result",
					"tool_use_id": "toolu_2",
					"content": map[string]interface{}{
						"temperature": 22,
						"condition":   "sunny",
					},
				},
			},
		},
	}

	geminiMsgs := provider.convertMessages(msgs)
	toolResultMsg := extractToolResultMsg(t, geminiMsgs)

	parts, ok := toolResultMsg["parts"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, parts, 1)

	part, ok := parts[0].(map[string]interface{})
	assert.True(t, ok)

	funcResp, ok := part["functionResponse"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "get_weather", funcResp["name"])

	response, ok := funcResp["response"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 22, response["temperature"])
	assert.Equal(t, "sunny", response["condition"])
}

func TestGeminiProvider_ConvertMessages_ToolResultFallbackToID(t *testing.T) {
	provider := &GeminiProvider{}

	// 边界场景：tool_result 出现时历史中没有对应 tool_use（例如客户端裁剪了历史）。
	// 此时回退使用 tool_use_id 作为 name，避免丢失字段。
	msgs := []types.ClaudeMessage{
		{
			Role: "user",
			Content: []interface{}{
				map[string]interface{}{
					"type":        "tool_result",
					"tool_use_id": "toolu_orphan",
					"content":     "orphaned result",
				},
			},
		},
	}

	geminiMsgs := provider.convertMessages(msgs)
	assert.Len(t, geminiMsgs, 1)

	parts, ok := geminiMsgs[0]["parts"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, parts, 1)

	part, ok := parts[0].(map[string]interface{})
	assert.True(t, ok)

	funcResp, ok := part["functionResponse"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "toolu_orphan", funcResp["name"], "无映射时回退用 tool_use_id")
}

func TestGeminiProvider_ConvertMessage_SkipsEmptyTextBlock(t *testing.T) {
	provider := &GeminiProvider{}

	// 测试场景：Claude assistant 消息常在 tool_use 前附带空 text 块。
	// 转换为 Gemini parts 时必须跳过该空文本，否则上游会返回 400:
	// "contents[X].parts[Y].data: required oneof field 'data' must have one initialized field"
	msg := types.ClaudeMessage{
		Role: "assistant",
		Content: []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": "",
			},
			map[string]interface{}{
				"type":  "tool_use",
				"id":    "toolu_1",
				"name":  "get_weather",
				"input": map[string]interface{}{"location": "Tokyo"},
			},
		},
	}

	geminiMsg := provider.convertMessage(msg, nil)
	assert.NotNil(t, geminiMsg)
	assert.Equal(t, "model", geminiMsg["role"])

	parts, ok := geminiMsg["parts"].([]interface{})
	assert.True(t, ok, "parts should be []interface{}")
	assert.Len(t, parts, 1, "空 text 块应被跳过，仅保留 functionCall part")

	part, ok := parts[0].(map[string]interface{})
	assert.True(t, ok)
	_, hasFuncCall := part["functionCall"]
	assert.True(t, hasFuncCall, "保留的 part 必须是 functionCall")
}

func TestGeminiProvider_ConvertMessage_KeepsNonEmptyTextBlock(t *testing.T) {
	provider := &GeminiProvider{}

	// 反向验证：非空 text 块仍应被保留。
	msg := types.ClaudeMessage{
		Role: "assistant",
		Content: []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": "我将调用工具",
			},
			map[string]interface{}{
				"type":  "tool_use",
				"id":    "toolu_2",
				"name":  "get_weather",
				"input": map[string]interface{}{"location": "Tokyo"},
			},
		},
	}

	geminiMsg := provider.convertMessage(msg, nil)
	assert.NotNil(t, geminiMsg)

	parts, ok := geminiMsg["parts"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, parts, 2, "非空 text 块和 functionCall 都应保留")

	textPart, ok := parts[0].(map[string]string)
	assert.True(t, ok)
	assert.Equal(t, "我将调用工具", textPart["text"])
}

func TestGeminiProvider_ConvertToGeminiRequest_InjectDummyThoughtSignature(t *testing.T) {
	provider := &GeminiProvider{}

	// 场景：messages 接口 → service_type=gemini 上游，渠道开启 injectDummyThoughtSignature。
	// 期望：转换出的 functionCall part 在 part 层级补一个 thoughtSignature=DummyThoughtSignature，
	// 避免严格校验的上游（如 vip.undyingapi.com）返回
	// "Function call is missing a thought_signature in functionCall parts"。
	claudeReq := &types.ClaudeRequest{
		Model: "gemini-3.5-flash",
		Messages: []types.ClaudeMessage{
			{
				Role: "assistant",
				Content: []interface{}{
					map[string]interface{}{
						"type":  "tool_use",
						"id":    "toolu_1",
						"name":  "Bash",
						"input": map[string]interface{}{"command": "ls"},
					},
				},
			},
		},
	}
	upstream := &config.UpstreamConfig{InjectDummyThoughtSignature: true}

	geminiReq := provider.convertToGeminiRequest(claudeReq, upstream)
	contents, ok := geminiReq["contents"].([]map[string]interface{})
	assert.True(t, ok, "contents 必须是 []map[string]interface{}")
	assert.Len(t, contents, 1)

	parts, ok := contents[0]["parts"].([]interface{})
	assert.True(t, ok, "parts 必须是 []interface{}")
	assert.Len(t, parts, 1)

	part, ok := parts[0].(map[string]interface{})
	assert.True(t, ok)
	_, hasFuncCall := part["functionCall"]
	assert.True(t, hasFuncCall, "保留的 part 必须是 functionCall")

	sig, ok := part["thoughtSignature"].(string)
	assert.True(t, ok, "应在 part 层级注入 thoughtSignature")
	assert.Equal(t, types.DummyThoughtSignature, sig)
}

func TestGeminiProvider_ConvertToGeminiRequest_DefaultNoSignature(t *testing.T) {
	provider := &GeminiProvider{}

	// 默认开关关闭：不注入 thoughtSignature，保持原有行为（与原生 Gemini 入口一致）。
	claudeReq := &types.ClaudeRequest{
		Model: "gemini-3.5-flash",
		Messages: []types.ClaudeMessage{
			{
				Role: "assistant",
				Content: []interface{}{
					map[string]interface{}{
						"type":  "tool_use",
						"id":    "toolu_2",
						"name":  "Bash",
						"input": map[string]interface{}{"command": "ls"},
					},
				},
			},
		},
	}
	upstream := &config.UpstreamConfig{}

	geminiReq := provider.convertToGeminiRequest(claudeReq, upstream)
	contents, _ := geminiReq["contents"].([]map[string]interface{})
	parts, _ := contents[0]["parts"].([]interface{})
	part, _ := parts[0].(map[string]interface{})
	_, hasSig := part["thoughtSignature"]
	assert.False(t, hasSig, "默认不应注入 thoughtSignature")
}

func TestGeminiProvider_ConvertToGeminiRequest_StripThoughtSignatureNoOp(t *testing.T) {
	provider := &GeminiProvider{}

	// StripThoughtSignature 在 Claude→Gemini 场景下是 no-op：Claude 协议本来就不带签名，
	// 没东西可剥；同时该开关必须能压制 InjectDummyThoughtSignature 注入（开关互斥）。
	claudeReq := &types.ClaudeRequest{
		Model: "gemini-3.5-flash",
		Messages: []types.ClaudeMessage{
			{
				Role: "assistant",
				Content: []interface{}{
					map[string]interface{}{
						"type":  "tool_use",
						"id":    "toolu_3",
						"name":  "Bash",
						"input": map[string]interface{}{"command": "ls"},
					},
				},
			},
		},
	}
	upstream := &config.UpstreamConfig{
		StripThoughtSignature:       true,
		InjectDummyThoughtSignature: true,
	}

	geminiReq := provider.convertToGeminiRequest(claudeReq, upstream)
	contents, _ := geminiReq["contents"].([]map[string]interface{})
	parts, _ := contents[0]["parts"].([]interface{})
	part, _ := parts[0].(map[string]interface{})
	_, hasSig := part["thoughtSignature"]
	assert.False(t, hasSig, "StripThoughtSignature 优先生效，不应注入签名")
}
