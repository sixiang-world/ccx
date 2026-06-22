# ccx-cli

> CCX API 代理网关命令行管理工具

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> **位置**：本工具现位于 `sixiang-world/ccx` fork 仓库的 [`cli/`](https://github.com/sixiang-world/ccx/tree/main/cli) 子目录下。
> 通过 `git subtree` 保持与独立仓库 `sixiang-world/ccx-cli` 同步。
>
> **跟踪上游**：运行 `./scripts/check-upstream.sh` 检查 `BenedictKing/ccx` 上游是否有 API 变更需要 ccx-cli 适配。

`ccx-cli` 是 [CCX](https://github.com/BenedictKing/ccx) API 代理网关的命令行管理工具，为运维人员提供**无需 Web UI** 的完整配置管理能力。

支持渠道（Channel）的 CRUD、API 密钥管理、模型映射、熔断器配置、健康检查、连通性检测等功能，输出支持 `table` / `json` / `yaml` 三种格式。

---

## 功能特性

- **渠道全生命周期管理** — 创建、查看、更新、删除上游渠道
- **API 密钥管理** — 添加、移除、排序、恢复被拉黑的密钥
- **调度与状态控制** — 设置渠道状态（active/suspended/disabled）、重排序、促销期
- **监控与诊断** — 健康检查、连通性检测（Ping）、性能指标、请求日志、仪表盘
- **配置管理** — 查看/应用/备份/恢复全局配置，支持热重载
- **运行时设置** — Fuzzy 模式、熔断器、图片轮次限制、对话设置
- **能力测试** — 对渠道运行模型能力测试，查看能力快照
- **三种输出格式** — Table（默认，适合人眼阅读）、JSON（适合管道处理）、YAML（适合配置管理）
- **多种认证方式** — `X-Api-Key`、`Authorization: Bearer`、`X-Goog-Api-Key`
- **自动重试** — 支持指数退避 + 抖动，可配置重试次数
- **TLS 支持** — 自定义 CA 证书、跳过验证模式

---

## 安装方式

### 前置要求

- Go 1.23 或更高版本

### 从源码编译

```bash
git clone https://github.com/BenedictKing/ccx-cli.git
cd ccx-cli
go build -o ccx-cli .
```

### 直接下载

从 [Releases](https://github.com/BenedictKing/ccx-cli/releases) 页面下载对应平台的预编译二进制文件。

---

## 快速开始

### 1. 配置服务地址和 API 密钥

```bash
# 方式一：命令行参数（优先级最高）
ccx channel list --server http://your-server:3000 --key your-admin-key

# 方式二：环境变量
export CCX_SERVER=http://your-server:3000
export CCX_API_KEY=your-admin-key
ccx channel list

# 方式三：配置文件（默认路径 ~/.config/ccx/config.json）
ccx channel list
```

### 2. 常用操作一览

```bash
# 查看所有渠道
ccx channel list

# 创建渠道
ccx channel create my-channel \
  --base-url https://api.anthropic.com \
  --api-key sk-ant-xxx \
  --service-type claude

# 查看渠道详情
ccx channel get my-channel

# 测试渠道连通性
ccx channel ping my-channel

# 查看渠道性能指标
ccx channel metrics -o json

# 全局健康检查
ccx health

# 查看帮助
ccx channel --help
```

---

## 命令参考

### 全局参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--server`, `-s` | `http://localhost:3000` | CCX 服务地址 |
| `--key`, `-k` | — | 管理 API 密钥 |
| `--type`, `-t` | `messages` | 渠道类型（messages/responses/chat/gemini/images） |
| `--output`, `-o` | `table` | 输出格式（table/json/yaml） |
| `--timeout` | `30s` | 请求超时时间 |
| `--retry` | `3` | 最大重试次数 |
| `--no-retry` | `false` | 关闭自动重试 |
| `--config` | `~/.config/ccx/config.json` | CLI 配置文件路径 |
| `--verbose`, `-v` | `false` | 显示详细请求信息 |
| `--show-keys` | `false` | 显示完整 API Key |
| `--ca-cert` | — | 自定义 CA 证书路径 |
| `--insecure-skip-verify` | `false` | 跳过 TLS 证书验证 |
| `--prefix` | — | 路由前缀 |

### 顶层命令

| 命令 | 说明 |
|------|------|
| `ccx channel` | 渠道管理（核心命令） |
| `ccx model` | 模型管理 |
| `ccx config` | 全局配置管理 |
| `ccx settings` | 运行时设置管理 |
| `ccx health` | 服务健康检查 |
| `ccx ping` | 全局连通性检测 |
| `ccx completion` | Shell 自动补全脚本生成 |
| `ccx help` | 查看帮助 |

### `ccx channel` 子命令

| 子命令 | 说明 |
|--------|------|
| `list` | 列出所有渠道 |
| `get <name>` | 获取渠道详情 |
| `create <name>` | 创建新渠道 |
| `update <name>` | 更新渠道配置（merge 语义） |
| `delete <name>` | 删除渠道 |
| `key add <name> <key>` | 添加 API 密钥 |
| `key remove <name> <key>` | 移除 API 密钥 |
| `key list <name>` | 列出渠道的 API 密钥 |
| `key move <name> <key>` | 调整密钥优先级 |
| `key restore <name> <key>` | 恢复被拉黑的密钥 |
| `status set <name> <state>` | 设置渠道状态（active/suspended/disabled） |
| `mapping list <name>` | 列出模型映射 |
| `mapping set <name> <source> <target>` | 设置模型映射（全量替换） |
| `reorder` | 重新排序渠道 |
| `promotion set <name> <duration>` | 设置促销期 |
| `promotion clear <name>` | 清除促销期 |
| `ping <name>` | 测试渠道连通性 |
| `resume <name>` | 恢复被熔断/拉黑的渠道 |
| `metrics` | 查看渠道性能指标 |
| `logs <name>` | 查看渠道请求日志 |
| `dashboard` | 查看统一仪表盘 |
| `scheduler-stats` | 查看调度器统计信息 |
| `capability snapshot <name>` | 查看渠道能力快照 |
| `capability test <name>` | 运行渠道能力测试 |
| `capability test-status <name> <jobId>` | 查看测试任务状态 |
| `capability test-cancel <name> <jobId>` | 取消测试任务 |
| `capability test-retry <name> <jobId>` | 重试测试失败的模型 |

### `ccx config` 子命令

| 子命令 | 说明 |
|--------|------|
| `show` | 显示全局配置聚合 |
| `apply <file>` | 应用配置文件（diff + 确认 + 逐项提交） |
| `save` | 强制持久化运行时配置到磁盘 |
| `backup` | 下载完整配置备份到本地 |
| `restore <file>` | 从备份文件恢复配置 |

### `ccx settings` 子命令

| 子命令 | 说明 |
|--------|------|
| `fuzzy get` | 查看 Fuzzy 模式状态 |
| `fuzzy set <true\|false>` | 设置 Fuzzy 模式 |
| `circuit-breaker get` | 查看熔断器配置 |
| `circuit-breaker set [flags]` | 设置熔断器配置 |
| `image-turn-limit get` | 查看图片轮次限制 |
| `image-turn-limit set <limit>` | 设置图片轮次限制 |
| `conversations get` | 查看对话设置 |
| `conversations set [flags]` | 更新对话设置 |

### `ccx model` 子命令

| 子命令 | 说明 |
|--------|------|
| `list --channel <name>` | 列出指定渠道支持的模型 |

---

## 配置文件说明

### CLI 配置文件

默认路径：`~/.config/ccx/config.json`

```json
{
  "server": "http://localhost:3000",
  "apiKey": "sk-xxx",
  "type": "messages",
  "output": "table",
  "timeout": "30s"
}
```

配置优先级（从高到低）：
1. 命令行参数（`--server`、`--key` 等）
2. 环境变量（`CCX_SERVER`、`CCX_API_KEY` 等）
3. 配置文件（`~/.config/ccx/config.json`）
4. 默认值

### 渠道类型

| 类型 | `--type` 取值 | 对应协议 |
|------|-------------|---------|
| Messages | `messages`（默认） | Claude Messages API |
| Responses | `responses` | OpenAI Responses API |
| Chat | `chat` | OpenAI Chat Completions API |
| Gemini | `gemini` | Google Gemini API |
| Images | `images` | OpenAI Images API |

---

## 开发指南

### 项目结构

```
ccx-cli/
├── main.go                  # 入口文件
├── cmd/                     # CLI 命令（Cobra）
│   ├── root.go              # 根命令 + 全局参数
│   ├── channel.go           # channel 父命令
│   ├── channel_list.go      # channel list
│   ├── channel_*.go         # 其他 channel 子命令
│   ├── config_cmd.go        # config 命令
│   ├── settings.go          # settings 命令
│   ├── model.go             # model 命令
│   ├── health.go            # health 命令
│   ├── ping.go              # ping 命令
│   └── completion.go        # shell 补全
├── internal/
│   ├── client/              # CCX API 客户端
│   │   └── client.go        # HTTP 客户端 + API 路径构建
│   ├── config/              # CLI 自身配置管理
│   │   └── config.go
│   ├── models/              # 数据结构
│   │   └── types.go         # 与 CCX 后端对齐的类型
│   ├── formatter/           # 输出格式化
│   │   └── formatter.go     # table/json/yaml 输出
│   ├── validator/           # 参数校验
│   │   └── validator.go
│   ├── errors/              # 自定义错误
│   │   └── errors.go
│   └── diff/                # JSON diff 工具
│       └── diff.go
```

### 技术栈

- **语言**: Go 1.23+
- **CLI 框架**: [Cobra](https://github.com/spf13/cobra) + [Viper](https://github.com/spf13/viper)
- **HTTP 客户端**: Go `net/http`（自定义重试、限流、认证）
- **输出格式化**: [tablewriter](https://github.com/olekukonko/tablewriter) + `encoding/json` + `gopkg.in/yaml.v3`

### 添加新命令

1. 在 `internal/models/types.go` 中添加对应的响应类型
2. 在 `internal/client/client.go` 中添加 API 路径构建函数（如需要）
3. 在 `cmd/` 下创建 `channel_xxx.go`，定义 Cobra 命令并在 `init()` 中注册到父命令
4. 运行 `go build ./...` 确认编译通过

### 构建

```bash
# 编译
go build -o ccx-cli .

# 注入版本信息
go build -ldflags="-X ccx-cli/internal/version.Version=1.0.0 -X ccx-cli/internal/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X ccx-cli/internal/version.GitCommit=$(git rev-parse HEAD)" -o ccx-cli .

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o ccx-cli-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o ccx-cli-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o ccx-cli-darwin-arm64 .
```

---

## 许可

MIT License
