# ccx-cli — CCX 命令行管理工具

## 项目定位
ccx-cli 是 CCX API 代理网关的命令行管理工具，Go + Cobra 实现，用于通过终端管理上游渠道、API 密钥、模型映射和全局设置。

## 仓库位置
本工具位于 `sixiang-world/ccx` 仓库的 [`cli/`](https://github.com/sixiang-world/ccx/tree/main/cli) 子目录下。

## 技术栈
- Go (1.23+)
- Cobra CLI 框架
- Viper 配置管理
- 输出格式：tablewriter + json + yaml

## 关键约定
- CLI 子命令与 CCX Admin REST API 一一对应
- 支持 --output json|yaml|table 三种格式
- 敏感信息（API Key）默认脱敏
- 使用 简体中文

## 开发环境
```bash
cd /root/ccx/cli          # 进入 cli 目录
go build -o ccx-cli .     # 编译
```

## 跟踪上游变更
```bash
./scripts/check-upstream.sh  # 检查 BenedictKing/ccx 上游是否有 API 变更
```

## 设计文档
- `docs/ccx-cli-design.md` — 架构设计文档（本目录下）
- `docs/user-manual.md` — 用户手册（本目录下）
