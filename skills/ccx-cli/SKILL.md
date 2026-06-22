---
name: ccx-cli
description: "CCX API 代理网关命令行管理工具 — 渠道 CRUD、API 密钥、模型映射、熔断器、健康检查。Go + Cobra 实现的运维 CLI。"
version: 1.1.0
tags: [ccx, cli, api-proxy, channel-management, devops, go]
related_skills:
  - ccx-api-proxy
trigger: 当用户说 ccx-cli、ccx 命令行、管理 CCX、渠道管理、ccx channel、ccx 运维时加载
---

# ccx-cli — CCX API 代理网关命令行管理工具

> 源码：`github.com/sixiang-world/ccx-cli` | 技术栈：Go 1.23 + Cobra + Viper

`ccx-cli` 是 CCX API 代理网关的命令行管理工具，让运维人员通过终端管理上游渠道、API 密钥、模型映射和运行时设置，无需 Web UI。

---

## 快速开始

```bash
export CCX_SERVER=http://localhost:3000
export CCX_API_KEY=*** health
ccx channel list
```

> 配置优先级：命令行参数 > 环境变量 `CCX_*` > `~/.config/ccx/config.json` > 默认值。

---

## 安装

### 方式 A：安装脚本（推荐）

```bash
git clone https://github.com/sixiang-world/ccx-cli.git
cd ccx-cli
sudo ./scripts/install.sh
# 自定义路径：./scripts/install.sh /opt/bin
```

🔴 **CHECKPOINT**：安装需要 sudo 权限，确认不会覆盖现有二进制。默认安装到 /usr/local/bin/ccx。

### 方式 B：手动编译

```bash
go build -ldflags="-X ccx-cli/internal/version.Version=$(git describe --tags --always) -X ccx-cli/internal/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X ccx-cli/internal/version.GitCommit=$(git rev-parse HEAD)" -o ccx-cli .
sudo mv ccx-cli /usr/local/bin/
```

### 方式 C：直接下载

从 [Releases](https://github.com/sixiang-world/ccx-cli/releases) 下载预编译二进制并放入 `$PATH`。

### 验证安装

```bash
ccx version
```

---

## 常用工作流

### 1. 渠道全生命周期

> **输入**：渠道名 + type/url/key → **输出**：渠道创建/更新/删除成功

```bash
ccx channel create my-channel --type chat --base-url https://api.openai.com/v1 --api-key sk-proj-xxx --model gpt-4o
```

🔴 **CHECKPOINT**：确认 base-url 和 api-key 正确，创建后立即可用。

```bash
ccx channel update my-channel --type chat --priority 0
```

🔴 **CHECKPOINT**：更新 base-url 或 service-type 会影响调度中的请求，建议低峰期操作。

```bash
ccx channel delete my-channel --type chat
```

🛑 **STOP**：删除渠道不可恢复。确认渠道名称后再执行 delete。如有备份先 `ccx config backup`。

### 2. API 密钥管理

> **输入**：渠道名 + 密钥值 → **输出**：密钥已添加/移除/恢复

```bash
ccx channel key list my-channel --type chat
ccx channel key add my-channel sk-xxxxx --type chat
ccx channel key list my-channel --type chat --show-keys
ccx channel key move my-channel <key-hash> --position 0
ccx channel key restore my-channel <key-hash> --type chat
```

### 3. 模型映射管理

> **输入**：渠道名 + source 模型名 + target 模型名 → **输出**：映射已设置（⚠️ 全量替换）

```bash
ccx channel mapping list my-channel --type chat
ccx channel mapping set my-channel source-model target-model --type chat  # 🔴 全量替换
```

### 4. 渠道状态与诊断

> **输入**：渠道名 + 状态值 → **输出**：渠道状态已变更 / 连通性结果

```bash
ccx channel status set my-channel active     # 激活
ccx channel status set my-channel suspended  # 暂停
ccx channel status set my-channel disabled   # 禁用
ccx channel resume my-channel
ccx health                                    # 服务健康
ccx ping                                      # 全局连通性
```

### 5. 配置管理

> **输入**：配置文件路径 → **输出**：配置已查看/备份/应用/恢复

```bash
ccx config show    # 查看
ccx config backup  # 备份
ccx config apply config.json   # 应用（diff + 确认）
ccx config save                # 持久化
ccx config restore backup.json # 恢复
```

🔴 **CHECKPOINT**：`config apply` / `config restore` 会覆盖当前配置，操作前先用 `config backup` 备份。

### 6. 运行时设置

> **输入**：设置参数 → **输出**：运行时行为已调整

```bash
ccx settings fuzzy set true                 # 模糊匹配
ccx settings circuit-breaker set            # 熔断器（影响生产流量）
ccx settings image-turn-limit set 5
ccx settings conversations get
```

🔴 **CHECKPOINT**：更改熔断器参数影响生产故障转移行为。先用 `circuit-breaker get` 查看当前配置。

### 7. Shell 自动补全

> **输入**：shell 类型（bash/zsh/fish）→ **输出**：补全脚本已安装

```bash
ccx completion bash > /etc/bash_completion.d/ccx
source /etc/bash_completion.d/ccx
```

---

## 全局参数速查

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-s, --server` | `http://localhost:3000` | CCX 管理 API 地址（非代理端口） |
| `-k, --key` | — | CCX 管理密钥（PROXY_ACCESS_KEY） |
| `-t, --type` | `messages` | 渠道类型：messages/chat/responses/gemini/images |
| `-o, --output` | `table` | 输出格式：table/json/yaml |
| `--timeout` | `30s` | 请求超时 |
| `--retry` | `3` | 最大重试次数 |
| `--show-keys` | `false` | 显示完整 API Key |
| `--verbose` | `false` | 显示详细请求信息 |
| `--insecure-skip-verify` | `false` | 跳过 TLS 验证 |
| `--ca-cert` | — | 自定义 CA 证书 |

---

## 渠道类型说明

| `--type` 值 | 对应协议 | 适用场景 |
|-------------|---------|----------|
| `messages`（默认） | Claude Messages API | Claude 系列模型 |
| `chat` | OpenAI Chat Completions | OpenAI / DeepSeek / 开源模型 |
| `responses` | OpenAI Responses API | Codex CLI / OpenAI 新协议 |
| `gemini` | Google Gemini API | Gemini 系列模型 |
| `images` | OpenAI Images API | DALL·E 图片生成 |

---

## 故障模式与修复

| 错误现象 | 一线修复 | 兜底 |
|----------|---------|------|
| `Connection refused` | `docker ps \| grep ccx` | `docker compose up -d` |
| `401 Unauthorized` | 确认用的是 PROXY_ACCESS_KEY | 检查 `~/.config/ccx/config.json` |
| 渠道列表为空 | 确认 `--type` 匹配渠道类型 | 不加 `--type` 或用 `-o json` 看原始数据 |
| `404 渠道未找到` | 检查渠道名拼写和 `--type` | `ccx channel list --type <type>` |
| key 显示为 `***` | 正常脱敏，加 `--show-keys` | 确认 `CCX_API_KEY` 完整 |
| 创建后 key 为空 | key 可能未正确保存 | 创建后手动 `key add` |

---

## 反例与黑名单

| 反模式 | 为什么有害 | 正确做法 |
|--------|-----------|---------|
| 不传 `--type` 操作 chat 渠道 | 默认去 messages 路由，找不到 | 每次显式传 `--type chat` |
| 把上游 key 当 `--key` 传 | `--key` 需要 PROXY_ACCESS_KEY | 检查 `.env` 的 PROXY_ACCESS_KEY |
| `--server` 指向代理端口 `:2342` | 代理端口无 Admin API | 指向管理端口 `:3000` |
| 以为 `mapping set` 是追加 | 全量替换，丢失之前映射 | 先 `mapping list` 再合并后 set |
| 跳过 `ccx health` 直接操作 | 连不上时所有操作报错 | 操作前先 `ccx health` |
| 用过期 ccx-cli 管新版 CCX | API 格式不兼容 | 保持版本对齐 |

---

## 输出格式技巧

```bash
ccx channel list -o json | jq '.[] | {name, status}'
ccx channel list -o json | jq '.[] | select(.circuitBreaker.tripped == true) | .name'
ccx channel list -o yaml > all-channels.yaml
```

---

## 相关文档

- `docs/user-manual.md` — 完整用户手册
- `docs/ccx-cli-design.md` — 架构设计文档
- 关联 skill `ccx-api-proxy` — 配置 AI 工具使用 CCX 代理
