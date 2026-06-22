1|# ccx-cli — CCX API 代理网关命令行管理工具
2|
3|`ccx-cli` 是 [CCX](https://github.com/BenedictKing/ccx) 的命令行管理工具，用于通过终端管理上游渠道、API 密钥、模型映射和全局设置，无需打开 Web UI。
4|
5|> **位置变动**：ccx-cli 现已合并到 `sixiang-world/ccx` 仓库的 [`cli/`](https://github.com/sixiang-world/ccx/tree/main/cli) 子目录下。
6|
7|---
8|
9|## 安装
10|
11|### 方式一：下载预编译二进制
12|
13|```bash
14|# Linux (amd64)
15|# TODO: 替换为实际仓库地址
16|wget https://github.com/sixiang-world/ccx/releases/latest/download/ccx-cli-linux-amd64.tar.gz
17|wget https://github.com/sixiang-world/ccx/releases/latest/download/ccx-cli-linux-amd64.tar.gz.sha256
18|sha256sum -c ccx-cli-linux-amd64.tar.gz.sha256
19|tar xzf ccx-cli-linux-amd64.tar.gz
20|sudo mv ccx-cli /usr/local/bin/
21|
22|# macOS (arm64)
23|# TODO: 替换为实际仓库地址
24|wget https://github.com/sixiang-world/ccx/releases/latest/download/ccx-cli-darwin-arm64.tar.gz
25|wget https://github.com/sixiang-world/ccx/releases/latest/download/ccx-cli-darwin-arm64.tar.gz.sha256
26|shasum -a 256 -c ccx-cli-darwin-arm64.tar.gz.sha256
27|tar xzf ccx-cli-darwin-arm64.tar.gz
28|sudo mv ccx-cli /usr/local/bin/
29|```
30|
31|### 方式二：使用 Go 安装
32|
33|```bash
34|# TODO: 替换为实际仓库地址
35|go install github.com/sixiang-world/ccx@latest
36|```
37|
38|### 方式三：从源码构建
39|
40|```bash
41|# TODO: 替换为实际仓库地址
42|git clone https://github.com/sixiang-world/ccx.git
43|cd ccx-cli
44|make build
45|sudo mv ccx-cli /usr/local/bin/
46|```
47|
48|### 方式四：使用 Docker
49|
50|```bash
51|# TODO: 替换为实际仓库地址
52|docker pull ghcr.io/sixiang-world/ccx-cli:latest
53|docker run --rm ghcr.io/sixiang-world/ccx-cli ccx health --server http://host.docker.internal:3000
54|```
55|
56|---
57|
58|## 快速开始
59|
60|### 1. 配置连接信息
61|
62|```bash
63|# 方式 A：命令行参数
64|ccx --server http://localhost:3000 --key your-admin-key health
65|
66|# 方式 B：环境变量
67|export CCX_SERVER=http://localhost:3000
68|export CCX_API_KEY=your-admin-key
69|ccx health
70|
71|# 方式 C：CLI 配置文件
72|mkdir -p ~/.config/ccx
73|chmod 700 ~/.config/ccx
74|cat > ~/.config/ccx/config.json << 'EOF'
75|{
76|  "server": "http://localhost:3000",
77|  "apiKey": "your-admin-key",
78|  "defaultType": "messages",
79|  "output": "table"
80|}
81|EOF
82|chmod 600 ~/.config/ccx/config.json
83|ccx health
84|```
85|
86|### 2. 验证连接
87|
88|```bash
89|$ ccx health
90|
91|  状态       运行模式   版本     运行时间   通道数
92| ─────────── ────────── ──────── ────────── ────────
93|  healthy    production v2.9.14  1234.56s   5
94|
95|# JSON 格式查看详细信息
96|$ ccx health --output json
97|{
98|  "status": "healthy",
99|  "timestamp": "2026-06-22T10:00:00+08:00",
100|  "uptime": 1234.56,
101|  "mode": "production",
102|  "version": {
103|    "version": "v2.9.14",
104|    "buildTime": "2026-06-20T12:00:00Z",
105|    "gitCommit": "abc1234"
106|  },
107|  "config": {
108|    "upstreamCount": 5
109|  }
110|}
111|```
112|
113|### 3. 查看现有通道
114|
115|```bash
116|$ ccx channel list
117|
118| 索引  名称             类型      状态      BaseURL                              优先级
119|───── ──────────────── ──────── ──────── ──────────────────────────────────── ────────
120| 0    claude-direct    claude   active   https://api.anthropic.com             0
121| 1    openai-gpt       openai   active   https://api.openai.com/v1             1
122| 2    gemini-v2        gemini   suspended https://generativelanguage.googleapis.com 2
123| 3    deepseek-chat    openai   active   https://api.deepseek.com/v1           3
124|```
125|
126|---
127|
128|## 命令参考
129|
130|### `ccx channel` — 通道管理
131|
132|通道管理是核心功能。所有通道操作都需要通过 `--type` / `-t` 指定渠道类型。
133|
134|#### 列出通道
135|
136|```bash
137|# 列出所有 Messages 通道（默认）
138|ccx channel list
139|
140|# 列出 Responses 通道
141|ccx channel list --type responses
142|
143|# 列出 Chat 通道（JSON 格式）
144|ccx channel list --type chat --output json
145|
146|# 列出 Gemini 通道（YAML 格式）
147|ccx channel list --type gemini -o yaml
148|```
149|
150|#### 查看通道详情
151|
152|```bash
153|ccx channel get claude-direct
154|ccx channel get openai-gpt --type chat
155|ccx channel get gemini-v2 --type gemini -o json
156|```
157|
158|#### 创建通道
159|
160|```bash
161|# 创建 Claude Messages 通道
162|ccx channel create my-claude \
163|  --type messages \
164|  --base-url https://api.anthropic.com \
165|  --api-key *** \
166|  --service-type claude \
167|  --description "主 Claude 通道"
168|
169|# 创建 OpenAI Chat 通道（多 Key）
170|ccx channel create my-openai \
171|  --type chat \
172|  --base-url https://api.openai.com/v1 \
173|  --api-key sk-proj-xxx \
174|  --api-key sk-proj-yyy \
175|  --service-type openai \
176|  --priority 0 \
177|  --status active
178|
179|# 创建 Gemini 通道
180|ccx channel create my-gemini \
181|  --type gemini \
182|  --base-url https://generativelanguage.googleapis.com \
183|  --api-key AIzaSyXXX \
184|  --service-type gemini
185|
186|# 创建带模型映射的通道
187|ccx channel create my-custom \
188|  --type messages \
189|  --base-url https://custom-api.example.com \
190|  --api-key cus-xxx \
191|  --service-type openai \
192|  --model-mapping "claude-sonnet-4=gpt-4o" \
193|  --model-mapping "claude-opus-4=gpt-4-turbo"
194|```
195|
196|#### 更新通道
197|
198|```bash
199|# 更新基础 URL
200|ccx channel update claude-direct --base-url https://api.anthropic.com/v1
201|
202|# 更新描述和优先级
203|ccx channel update my-channel --description "主通道" --priority 0
204|
205|# 启用代理
206|ccx channel update my-channel --proxy-url http://proxy.example.com:8080
207|
208|# 设置模型映射（⚠️ 全量替换，会覆盖所有现有映射）
209|ccx channel update my-channel --model-mapping "claude-sonnet-4=gpt-4o"
210|```
211|
212|#### 删除通道
213|
214|```bash
215|# 提示确认（默认）
216|ccx channel delete my-channel
217|
218|# 强制删除（跳过确认）
219|ccx channel delete my-channel --force
220|ccx channel delete my-channel -f
221|```
222|
223|#### 通道状态管理
224|
225|```bash
226|# 暂停通道（不再调度新请求，已有连接不受影响）
227|ccx channel status set my-channel suspended
228|
229|# 激活通道
230|ccx channel status set my-channel active
231|
232|# 移入备用池
233|ccx channel status set my-channel disabled
234|```
235|
236|#### 通道重排序
237|
238|```bash
239|# 按名称顺序
240|ccx channel reorder --order claude-direct,openai-gpt,gemini-v2
241|```
242|
243|#### 恢复熔断/拉黑渠道
244|
245|```bash
246|# 恢复被熔断或被拉黑的渠道（重置熔断状态、恢复拉黑 Key，保留历史统计）
247|ccx channel resume my-channel
248|```
249|
250|#### 通道促销
251|
252|```bash
253|# 设置促销期（指定持续时长，以秒为单位，如 7200 秒=2 小时）
254|ccx channel promotion set my-channel 7200
255|
256|# 清除促销（发送 {"duration": 0}，与 set 共用同一 POST 端点）
257|ccx channel promotion clear my-channel
258|```
259|
260|#### 通道指标
261|
262|```bash
263|# 查看通道指标
264|ccx channel metrics my-channel
265|
266|# 查看通道历史指标
267|ccx channel metrics history my-channel
268|
269|# 查看指定类型的通道指标
270|ccx channel metrics my-channel --type chat
271|```
272|
273|#### 通道日志
274|
275|```bash
276|# 查看通道请求日志
277|ccx channel logs my-channel
278|```
279|
280|#### 通道仪表盘
281|
282|```bash
283|# 查看统一仪表盘（支持 --type 过滤）
284|ccx channel dashboard
285|
286|# 仅查看 Chat 类型的仪表盘
287|ccx channel dashboard --type chat
288|```
289|
290|#### 调度器统计
291|
292|```bash
293|# 查看调度器统计
294|ccx channel scheduler-stats
295|```
296|
297|#### 通道能力测试
298|
299|```bash
300|# 查看渠道能力快照
301|ccx channel capability snapshot my-channel
302|
303|# 运行渠道能力测试
304|ccx channel capability test my-channel
305|
306|# 查看能力测试任务状态
307|ccx channel capability test-status my-channel <jobId>
308|
309|# 取消能力测试任务
310|ccx channel capability test-cancel my-channel <jobId>
311|
312|# 重试能力测试中失败的模型
313|ccx channel capability test-retry my-channel <jobId>
314|```
315|
316|### `ccx channel key` — API 密钥管理
317|
318|```bash
319|# 列出通道的 API 密钥
320|ccx channel key list my-channel
321|
322|# 添加密钥
323|ccx channel key add my-channel sk-ant...wkey
324|
325|# 删除密钥
326|ccx channel key remove my-channel sk-ant...dkey
327|
328|# 调整密钥优先级（移到顶部）— 使用独立端点 /keys/:apiKey/top
329|ccx channel key move my-channel *** --position top
330|
331|# 调整密钥优先级（移到底部）— 使用独立端点 /keys/:apiKey/bottom
332|ccx channel key move my-channel *** --position bottom
333|
334|# 恢复被自动拉黑的密钥 — API Key 在请求体中传递，不在 URL 路径中
335|ccx channel key restore my-channel sk-ant...sted
336|```
337|
338|### `ccx channel mapping` — 模型映射管理
339|
340|```bash
341|# 查看当前模型映射
342|ccx channel mapping list my-channel
343|
344|# 添加/更新映射
345|ccx channel mapping set my-channel claude-sonnet-4 claude-sonnet-4-20250514
346|
347|# 添加带 reasoning 级别的映射
348|ccx channel mapping set my-channel claude-opus-4 claude-opus-4-20250514 --reasoning high
349|```
350|
351|### `ccx model` — 模型查询
352|
353|```bash
354|# 查看特定渠道的模型（需指定 --channel）
355|ccx model list --channel my-channel
356|
357|# 查看 Chat 渠道的模型（JSON 格式）
358|ccx model list --type chat --channel my-channel --output json
359|```
360|
361|### `ccx config` — 配置管理
362|
363|```bash
364|# 查看当前配置摘要
365|ccx config show
366|
367|# 查看完整配置（JSON）
368|ccx config show --output json
369|
370|# 从文件应用配置（预览模式）
371|ccx config apply path/to/config.json --dry-run
372|
373|# 从文件应用配置
374|ccx config apply path/to/config.json
375|
376|# 强制保存当前配置到磁盘
377|ccx config save
378|
379|# 备份当前配置（下载到本地）
380|ccx config backup
381|
382|# 从备份恢复
383|ccx config restore path/to/backup.json
384|```
385|
386|### `ccx settings` — 运行时设置
387|
388|```bash
389|# Fuzzy 模式
390|ccx settings fuzzy get
391|ccx settings fuzzy set true
392|ccx settings fuzzy set false
393|
394|# 熔断器配置
395|ccx settings circuit-breaker get
396|ccx settings circuit-breaker set \
397|  --window-size 10 \
398|  --failure-threshold 0.5 \
399|  --consecutive-failures 5
400|
401|# 历史图片轮次限制
402|ccx settings image-turn-limit get
403|ccx settings image-turn-limit set 5
404|
405|# 对话设置
406|ccx settings conversations get
407|ccx settings conversations set [flags]
408|```
409|
410|### `ccx conversation` — 对话管理
411|
412|```bash
413|# 列出活跃对话
414|ccx conversation list
415|
416|# 设置对话覆盖
417|ccx conversation override set <conversation-id>
418|
419|# 移除对话覆盖
420|ccx conversation override remove <conversation-id>
421|```
422|
423|### `ccx health` — 健康检查
424|
425|```bash
426|# 基础健康检查
427|ccx health
428|
429|# 详细健康信息
430|ccx health --output json
431|
432|# 特定路由前缀的健康检查
433|ccx health --prefix my-prefix
434|```
435|
436|### `ccx ping` — 连通性检测
437|
438|```bash
439|# Ping 所有类型的所有通道
440|ccx ping
441|
442|# Ping 特定类型的所有通道
443|ccx ping --type chat
444|
445|# Ping 单个指定名称的通道（自动匹配所有类型）
446|ccx ping my-channel
447|
448|# JSON 格式输出
449|ccx ping --output json
450|```
451|
452|> `ccx ping <name>` 与 `ccx channel ping <name>` 功能相同，前者会自动在所有类型中查找匹配的通道，后者需要指定 `--type`。
453|
454|---
455|
456|## 输出格式详解
457|
458|### table（默认）
459|
460|适合人眼阅读的表格，列宽自动调整。
461|
462|```bash
463|$ ccx channel list --type responses -o table
464|
465| 索引  名称             状态      BaseURL                 延迟    优先级
466|───── ──────────────── ──────── ─────────────────────── ──────── ────────
467| 0    openai-responses active   https://api.openai.com   45ms    0
468| 1    responses-fallback active https://api.backup.com   120ms   1
469|```
470|
471|### json
472|
473|适合程序化消费，可通过 `jq` 进一步处理。
474|
475|```bash
476|$ ccx channel list --type messages -o json | jq '.channels[0].apiKeys | length'
477|3
478|```
479|
480|### yaml
481|
482|适合用于版本控制或配置生成。
483|
484|```bash
485|$ ccx channel get my-channel -o yaml > channel-backup.yaml
486|```
487|
488|---
489|
490|## 环境变量
491|
492|| 变量 | 默认值 | 说明 |
493||------|--------|------|
494|| `CCX_SERVER` | `http://localhost:3000` | CCX 服务地址 |
495|| `CCX_API_KEY` | — | 管理 API 密钥 |
496|| `CCX_DEFAULT_TYPE` | `messages` | 默认渠道类型 |
497|| `CCX_OUTPUT` | `table` | 默认输出格式 |
498|| `CCX_TIMEOUT` | `30s` | 请求超时时间 |
499|| `CCX_RETRY` | `3` | 瞬时故障最大重试次数 |
500|| `CCX_VERBOSE` | `false` | 详细模式 |
501|