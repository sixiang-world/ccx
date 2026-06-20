// 模型名优先级排序：按预定义优先级模式降序排列，同优先级组内按自然降序
// 规则顺序：先新后旧、先精确后宽松；同家族新版本在前，带 codex/pro/max 等精确后缀优先于通用名
// 数据基线：2026-05 各家官方在售模型

const modelPriorityPatterns: RegExp[] = [
  // Anthropic Claude（Fable 5 / 4.8 旗舰 / 4.7 / 4.6 Sonnet / 4.5 Haiku）
  /fable-5/i,
  /opus-4-8/i,
  /opus-4-7/i,
  /sonnet-4-7/i,
  /haiku-4-7/i,
  /opus-4-6/i,
  /sonnet-4-6/i,
  /haiku-4-6/i,
  /opus-4-5/i,
  /sonnet-4-5/i,
  /haiku-4-5/i,

  // OpenAI GPT-5 系列（pro / codex 变体优先匹配，再降级到主版本）
  /gpt-5\.5-pro/i,
  /gpt-5\.5/i,
  /gpt-5\.4-pro/i,
  /gpt-5\.4-mini/i,
  /gpt-5\.4-nano/i,
  /gpt-5\.4/i,
  /gpt-5\.3-codex/i,
  /gpt-5\.3/i,
  /gpt-5\.2-codex/i,
  /gpt-5\.2-pro/i,
  /gpt-5\.2/i,
  /gpt-5\.1-codex/i,
  /gpt-5\.1/i,
  /gpt-5-codex/i,
  /gpt-5-pro/i,
  /gpt-5/i,

  // Google Gemini（3.5 Flash → 3.1 Pro Preview → 3 Pro / Flash Preview → 3.1 Flash Lite → 2.5 系列）
  /gemini-3\.5-flash/i,
  /gemini-3\.1-pro/i,
  /gemini-3\.1-flash-lite/i,
  /gemini-3-pro/i,
  /gemini-3-flash/i,
  /gemini-3/i,
  /gemini-2\.5-pro/i,
  /gemini-2\.5-flash-lite/i,
  /gemini-2\.5-flash/i,

  // xAI Grok（4.3 当前旗舰；保留 4.2/4.1 以兼容旧 channel 命名）
  /grok-4\.3/i,
  /grok-4-3/i,
  /grok-4\.2/i,
  /grok-4\.1/i,
  /grok-4/i,

  // 智谱 GLM
  /glm-?5\.2/i,
  /glm-?5\.1/i,
  /glm-?5/i,
  /glm-?4\.7-flash/i,
  /glm-?4\.7/i,
  /glm-?4\.6/i,

  // 阿里 Qwen（3.6 / 3.5 / 3-Max）
  /qwen-?3\.6-plus/i,
  /qwen-?3\.6/i,
  /qwen-?3\.5/i,
  /qwen-?3-max/i,
  /qwen-?3-coder/i,
  /qwen-?3/i,

  // DeepSeek（V4 已发布；deepseek-chat / deepseek-reasoner 对应 V3.2）
  /deepseek-v4-pro/i,
  /deepseek-v4-flash/i,
  /deepseek-v4/i,
  /deepseek-v3\.2/i,
  /deepseek-reasoner/i,
  /deepseek-chat/i,
  /deepseek-v3/i,

  // Moonshot Kimi / MiniMax（带版本号 → 通用简写）
  /kimi-?k2\.7/i,
  /kimi-?k2\.6/i,
  /kimi-?k2\.5/i,
  /kimi-?k2-thinking/i,
  /minimax-?m3/i,
  /minimax-?m2\.7/i,
  /minimax-?m2\.5/i,
  /mimo-v2\.5/i,
  /doubao-seed-2-0/i,
  /ernie-4\.5/i,
  /baichuan-m2/i,
  /yi-34b-200k/i,
  /k2\.7/i,
  /k2\.6/i,
  /k2\.5/i,
  /m3/i,
  /m2\.7/i,
  /m2\.5/i,

  // DeepSeek 兜底（匹配各种 deepseek- 前缀变体）
  /deepseek-/i,
]

const modelNameCollator = new Intl.Collator('en', { numeric: true, sensitivity: 'base' })

export const getModelPriority = (name: string): number => {
  for (let i = 0; i < modelPriorityPatterns.length; i++) {
    if (modelPriorityPatterns[i].test(name)) return i
  }
  return modelPriorityPatterns.length
}

export const sortModelNamesDesc = (models: string[]): string[] => {
  return [...models].sort((a, b) => {
    const pa = getModelPriority(a)
    const pb = getModelPriority(b)
    if (pa !== pb) return pa - pb
    // 同优先级组内按自然降序
    return modelNameCollator.compare(b, a)
  })
}
