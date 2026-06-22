# ccx-cli — CCX API 代理网关命令行管理工具

`ccx-cli` 是 [CCX](https://github.com/BenedictKing/ccx) 的命令行管理工具，用于通过终端管理上游渠道、API 密钥、模型映射和全局设置，无需打开 Web UI。

> **🚧 注意**：`ccx-cli` 仓库尚在开发中，以下安装方式中的 GitHub 链接（`github.com/BenedictKing/ccx-cli`）为预留地址，正式发布前将替换为实际仓库地址。

---

## 安装

### 方式一：下载预编译二进制

```bash
# Linux (amd64)
# TODO: 替换为实际仓库地址
wget https://github.com/BenedictKing/ccx-cli/releases/latest/download/ccx-cli-linux-amd64.tar.gz
wget https://github.com/BenedictKing/ccx-cli/releases/latest/download/ccx-cli-linux-amd64.tar.gz.sha256
sha256sum -c ccx-cli-linux-amd64.tar.gz.sha256
tar xzf ccx-cli-linux-amd64.tar.gz
sudo mv ccx-cli /usr/local/bin/

# macOS (arm64)
# TODO: 替换为实际仓库地址
wget https://github.com/BenedictKing/ccx-cli/releases/latest/download/ccx-cli-darwin-arm64.tar.gz
wget https://github.com/BenedictKing/ccx-cli/releases/latest/download/ccx-cli-darwin-arm64.tar.gz.sha256
shasum -a 256 -c ccx-cli-darwin-arm64.tar.gz.sha256
tar xzf ccx-cli-darwin-arm64.tar.gz
sudo mv ccx-cli /usr/local/bin/
```

### 方式二：使用 Go 安装

```bash
# TODO: 替换为实际仓库地址
go install github.com/BenedictKing/ccx-cli@latest
```

### 方式三：从源码构建

```bash
# TODO: 替换为实际仓库地址
git clone https://github.com/BenedictKing/ccx-cli.git
cd ccx-cli
make build
sudo mv ccx-cli /usr/local/bin/
```

### 方式四：使用 Docker

```bash
# TODO: 替换为实际仓库地址
docker pull ghcr.io/benedictking/ccx-cli:latest
docker run --rm ghcr.io/benedictking/ccx-cli ccx health --server http://host.docker.internal:3000
```

---

## 快速开始

### 1. 配置连接信息

```bash
# 方式 A：命令行参数
ccx --server http://localhost:3000 --key your-admin-key health

# 方式 B：环境变量
export CCX_SERVER=http://localhost:3000
export CCX_API_KEY=your-admin-key
ccx health

# 方式 C：CLI 配置文件
mkdir -p ~/.config/ccx
chmod 700 ~/.config/ccx
cat > ~/.config/ccx/config.json << 'EOF'
{
  "server": "http://localhost:3000",
  "apiKey": "your-admin-key",
  "defaultType": "messages",
  "output": "table"
}
EOF
chmod 600 ~/.config/ccx/config.json
ccx health
```

### 2. 验证连接

```bash
$ ccx health

  状态       运行模式   版本     运行时间   通道数
 ─────────── ────────── ──────── ────────── ────────
  healthy    production v2.9.14  1234.56s   5

# JSON 格式查看详细信息
$ ccx health --output json
{
  "status": "healthy",
  "timestamp": "2026-06-22T10:00:00+08:00",
  "uptime": 1234.56,
  "mode": "production",
  "version": {
    "version": "v2.9.14",
    "buildTime": "2026-06-20T12:00:00Z",
    "gitCommit": "abc1234"
  },
  "config": {
    "upstreamCount": 5
  }
}
```

### 3. 查看现有通道

```bash
$ ccx channel list

 索引  名称             类型      状态      BaseURL                              优先级
───── ──────────────── ──────── ──────── ──────────────────────────────────── ────────
 0    claude-direct    claude   active   https://api.anthropic.com             0
 1    openai-gpt       openai   active   https://api.openai.com/v1             1
 2    gemini-v2        gemini   suspended https://generativelanguage.googleapis.com 2
 3    deepseek-chat    openai   active   https://api.deepseek.com/v1           3
```

---

## 命令参考

### `ccx channel` — 通道管理

通道管理是核心功能。所有通道操作都需要通过 `--type` / `-t` 指定渠道类型。

#### 列出通道

```bash
# 列出所有 Messages 通道（默认）
ccx channel list

# 列出 Responses 通道
ccx channel list --type responses

# 列出 Chat 通道（JSON 格式）
ccx channel list --type chat --output json

# 列出 Gemini 通道（YAML 格式）
ccx channel list --type gemini -o yaml
```

#### 查看通道详情

```bash
ccx channel get claude-direct
ccx channel get openai-gpt --type chat
ccx channel get gemini-v2 --type gemini -o json
```

#### 创建通道

```bash
# 创建 Claude Messages 通道
ccx channel create my-claude \
  --type messages \
  --base-url https://api.anthropic.com \
  --api-key sk-ant-api03-xxx \
  --service-type claude \
  --description "主 Claude 通道"

# 创建 OpenAI Chat 通道（多 Key）
ccx channel create my-openai \
  --type chat \
  --base-url https://api.openai.com/v1 \
  --api-key sk-proj-xxx \
  --api-key sk-proj-yyy \
  --service-type openai \
  --priority 0 \
  --status active

# 创建 Gemini 通道
ccx channel create my-gemini \
  --type gemini \
  --base-url https://generativelanguage.googleapis.com \
  --api-key AIzaSyXXX \
  --service-type gemini

# 创建带模型映射的通道
ccx channel create my-custom \
  --type messages \
  --base-url https://custom-api.example.com \
  --api-key cus-xxx \
  --service-type openai \
  --model-mapping "claude-sonnet-4=gpt-4o" \
  --model-mapping "claude-opus-4=gpt-4-turbo"
```

#### 更新通道

```bash
# 更新基础 URL
ccx channel update claude-direct --base-url https://api.anthropic.com/v1

# 更新描述和优先级
ccx channel update my-channel --description "主通道" --priority 0

# 启用代理
ccx channel update my-channel --proxy-url http://proxy.example.com:8080

# 设置模型映射（⚠️ 全量替换，会覆盖所有现有映射）
ccx channel update my-channel --model-mapping "claude-sonnet-4=gpt-4o"
```

#### 删除通道

```bash
# 提示确认（默认）
ccx channel delete my-channel

# 强制删除（跳过确认）
ccx channel delete my-channel --force
ccx channel delete my-channel -f
```

#### 通道状态管理

```bash
# 暂停通道（不再调度新请求，已有连接不受影响）
ccx channel status set my-channel suspended

# 激活通道
ccx channel status set my-channel active

# 移入备用池
ccx channel status set my-channel disabled
```

#### 通道重排序

```bash
# 按名称顺序
ccx channel reorder --order claude-direct,openai-gpt,gemini-v2
```

#### 恢复熔断/拉黑渠道

```bash
# 恢复被熔断或被拉黑的渠道（重置熔断状态、恢复拉黑 Key，保留历史统计）
ccx channel resume my-channel
```

#### 通道促销

```bash
# 设置促销期（指定持续时长，以秒为单位，如 7200 秒=2 小时）
ccx channel promotion set my-channel 7200

# 清除促销（发送 {"duration": 0}，与 set 共用同一 POST 端点）
ccx channel promotion clear my-channel
```

#### 通道指标

```bash
# 查看通道指标
ccx channel metrics my-channel

# 查看通道历史指标
ccx channel metrics history my-channel

# 查看指定类型的通道指标
ccx channel metrics my-channel --type chat
```

#### 通道日志

```bash
# 查看通道请求日志
ccx channel logs my-channel
```

#### 通道仪表盘

```bash
# 查看统一仪表盘（支持 --type 过滤）
ccx channel dashboard

# 仅查看 Chat 类型的仪表盘
ccx channel dashboard --type chat
```

#### 调度器统计

```bash
# 查看调度器统计
ccx channel scheduler-stats
```

#### 通道能力测试

```bash
# 查看渠道能力快照
ccx channel capability snapshot my-channel

# 运行渠道能力测试
ccx channel capability test my-channel

# 查看能力测试任务状态
ccx channel capability test-status my-channel <jobId>

# 取消能力测试任务
ccx channel capability test-cancel my-channel <jobId>

# 重试能力测试中失败的模型
ccx channel capability test-retry my-channel <jobId>
```

### `ccx channel key` — API 密钥管理

```bash
# 列出通道的 API 密钥
ccx channel key list my-channel

# 添加密钥
ccx channel key add my-channel sk-ant-api03-newkey

# 删除密钥
ccx channel key remove my-channel sk-ant-api03-oldkey

# 调整密钥优先级（移到顶部）— 使用独立端点 /keys/:apiKey/top
ccx channel key move my-channel sk-ant-api03-key --position top

# 调整密钥优先级（移到底部）— 使用独立端点 /keys/:apiKey/bottom
ccx channel key move my-channel sk-ant-api03-key --position bottom

# 恢复被自动拉黑的密钥 — API Key 在请求体中传递，不在 URL 路径中
ccx channel key restore my-channel sk-ant-api03-blacklisted
```

### `ccx channel mapping` — 模型映射管理

```bash
# 查看当前模型映射
ccx channel mapping list my-channel

# 添加/更新映射
ccx channel mapping set my-channel claude-sonnet-4 claude-sonnet-4-20250514

# 添加带 reasoning 级别的映射
ccx channel mapping set my-channel claude-opus-4 claude-opus-4-20250514 --reasoning high
```

### `ccx model` — 模型查询

```bash
# 查看特定渠道的模型（需指定 --channel）
ccx model list --channel my-channel

# 查看 Chat 渠道的模型（JSON 格式）
ccx model list --type chat --channel my-channel --output json
```

### `ccx config` — 配置管理

```bash
# 查看当前配置摘要
ccx config show

# 查看完整配置（JSON）
ccx config show --output json

# 从文件应用配置（预览模式）
ccx config apply path/to/config.json --dry-run

# 从文件应用配置
ccx config apply path/to/config.json

# 强制保存当前配置到磁盘
ccx config save

# 备份当前配置（下载到本地）
ccx config backup

# 从备份恢复
ccx config restore path/to/backup.json
```

### `ccx settings` — 运行时设置

```bash
# Fuzzy 模式
ccx settings fuzzy get
ccx settings fuzzy set true
ccx settings fuzzy set false

# 熔断器配置
ccx settings circuit-breaker get
ccx settings circuit-breaker set \
  --window-size 10 \
  --failure-threshold 0.5 \
  --consecutive-failures 5

# 历史图片轮次限制
ccx settings image-turn-limit get
ccx settings image-turn-limit set 5

# 对话设置
ccx settings conversations get
ccx settings conversations set [flags]
```

### `ccx conversation` — 对话管理

```bash
# 列出活跃对话
ccx conversation list

# 设置对话覆盖
ccx conversation override set <conversation-id>

# 移除对话覆盖
ccx conversation override remove <conversation-id>
```

### `ccx health` — 健康检查

```bash
# 基础健康检查
ccx health

# 详细健康信息
ccx health --output json

# 特定路由前缀的健康检查
ccx health --prefix my-prefix
```

### `ccx ping` — 连通性检测

```bash
# Ping 所有类型的所有通道
ccx ping

# Ping 特定类型的所有通道
ccx ping --type chat

# Ping 单个指定名称的通道（自动匹配所有类型）
ccx ping my-channel

# JSON 格式输出
ccx ping --output json
```

> `ccx ping <name>` 与 `ccx channel ping <name>` 功能相同，前者会自动在所有类型中查找匹配的通道，后者需要指定 `--type`。

---

## 输出格式详解

### table（默认）

适合人眼阅读的表格，列宽自动调整。

```bash
$ ccx channel list --type responses -o table

 索引  名称             状态      BaseURL                 延迟    优先级
───── ──────────────── ──────── ─────────────────────── ──────── ────────
 0    openai-responses active   https://api.openai.com   45ms    0
 1    responses-fallback active https://api.backup.com   120ms   1
```

### json

适合程序化消费，可通过 `jq` 进一步处理。

```bash
$ ccx channel list --type messages -o json | jq '.channels[0].apiKeys | length'
3
```

### yaml

适合用于版本控制或配置生成。

```bash
$ ccx channel get my-channel -o yaml > channel-backup.yaml
```

---

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `CCX_SERVER` | `http://localhost:3000` | CCX 服务地址 |
| `CCX_API_KEY` | — | 管理 API 密钥 |
| `CCX_DEFAULT_TYPE` | `messages` | 默认渠道类型 |
| `CCX_OUTPUT` | `table` | 默认输出格式 |
| `CCX_TIMEOUT` | `30s` | 请求超时时间 |
| `CCX_RETRY` | `3` | 瞬时故障最大重试次数 |
| `CCX_VERBOSE` | `false` | 详细模式 |
| `CCX_NO_MASK` | `false` | 是否显示完整的 API Key（谨慎使用） |
| `CCX_CONFIG` | `~/.config/ccx/config.json` | CLI 配置文件路径 |
| `CCX_CA_CERT` | — | 自定义 CA 证书路径 |
| `CCX_INSECURE_SKIP_VERIFY` | `false` | 跳过 TLS 证书验证 |
| `CCX_PREFIX` | — | 路由前缀（对应 `/:routePrefix/health` 健康检查） |

---

## Shell 自动补全

```bash
# Bash
source <(ccx completion bash)

# Zsh
source <(ccx completion zsh)

# Fish
ccx completion fish | source

# 永久安装（以 bash 为例）
ccx completion bash > /etc/bash_completion.d/ccx
```

---

## 常见问题

### Q: 连接被拒绝

```bash
$ ccx health
✗ 连接失败：无法连接到 http://localhost:3000
```

**解决方案：**
1. 确认 CCX 服务正在运行：`curl http://localhost:3000/health`
2. 检查端口号：`ccx --server http://localhost:8080 health`
3. 如果使用 Docker，确认端口映射正确

### Q: 认证失败

```bash
$ ccx health
✗ 认证失败：401 Unauthorized
```

**解决方案：**
1. 确认已设置正确的 API Key：`export CCX_API_KEY=your-key`
2. 服务器可能配置了 `ADMIN_ACCESS_KEY`，需要使用管理密钥而非代理密钥
3. 检查密钥是否包含特殊字符（建议用单引号包裹）

### Q: 数据不一致

**场景：** CLI 显示与 Web UI 不一致

**解决方案：**
1. 确认 CLI 和 Web UI 连接的是同一台服务器（比较 `--server` 参数和环境变量）
2. 强制保存配置：`ccx config save`（调用 `POST /admin/config/save`）
3. 检查 `.config/config.json` 文件是否有未保存的改动
4. 如果使用 `--output json`，缓存可能导致差异，可尝试 `--timeout 5s` 缩短超时

### Q: 如何批量导入通道？

使用 `ccx config apply`：

```bash
# 1. 先导出当前配置作为模板
ccx config show -o json > template.json

# 2. 编辑 template.json，添加需要的通道

# 3. 预览变更
ccx config apply template.json --dry-run

# 4. 确认后应用
ccx config apply template.json
```

### Q: 如何查看 API Key 的完整内容？

默认情况下，所有输出格式（table/json/yaml）都会对 API Key 进行脱敏，只显示前 4 位和后 4 位。

```bash
# 通过 --show-keys 标志显示完整密钥（需要确认）
ccx channel key list my-channel --show-keys

# 通过环境变量全局关闭脱敏（谨慎使用）
export CCX_NO_MASK=true
ccx channel key list my-channel
```

### Q: 支持哪些渠道类型？有什么区别？

| 类型 | 对应端点 | 典型上游 |
|------|---------|---------|
| `messages` | `/v1/messages` | Anthropic Claude API |
| `responses` | `/v1/responses` | OpenAI Responses API |
| `chat` | `/v1/chat/completions` | OpenAI Chat 兼容 API |
| `gemini` | `/v1beta/models/*` | Google Gemini API |
| `images` | `/v1/images/*` | OpenAI Images API |

### Q: 为什么有些通道创建后状态是 "unknown"？

如果通道缺少必要的配置（如 `apiKeys` 为空），CCX 会自动将其状态标记为异常。建议创建时至少提供 `--base-url`、`--api-key` 和 `--service-type`。

---

## 最佳实践

1. **使用配置文件而非命令行参数**：将 `CCX_SERVER` 和 `CCX_API_KEY` 写入 `~/.config/ccx/config.json`，并设置 `chmod 600` 保护密钥
2. **区分管理密钥**：生产环境务必设置 `ADMIN_ACCESS_KEY`，与 `PROXY_ACCESS_KEY` 隔离
3. **批量操作前先 dry-run**：`ccx config apply file.json --dry-run`
4. **定期备份配置**：`ccx config backup`
5. **谨慎删除**：删除通道前建议先 `status set suspended` 观察影响
6. **管道友好**：善用 `-o json | jq` 组合进行复杂查询
7. **版本匹配**：保持 CLI 版本与 CCX 后端版本一致（可用 `ccx health` 查看后端版本）
8. **CI/CD 集成**：在自动化流水线中使用 `--no-retry` 和 `--dry-run` 避免意外变更

---

## 构建与开发

```bash
# TODO: 替换为实际仓库地址
git clone https://github.com/BenedictKing/ccx-cli.git
cd ccx-cli

# 构建
make build

# 开发（热重载）
make dev

# 测试
make test

# 安装到 $GOPATH/bin
make install
```
