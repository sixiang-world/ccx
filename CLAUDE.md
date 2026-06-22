# ccx-cli — CCX 命令行管理工具

## 项目定位
ccx-cli 是 CCX API 代理网关的命令行管理工具，Go + Cobra 实现，用于通过终端管理上游渠道、API 密钥、模型映射和全局设置。

## 技术栈
- Go (1.23+)
- Cobra CLI 框架
- Viper 配置管理
- 输出格式：tablewriter + json + yaml

## 关键约定
- CLI 子命令与 CCX Admin REST API 一一对应
- 支持 --output json|yaml|table 三种格式
- 敏感信息（API Key）默认脱敏
- source ~/ccx-cli/.claude-env 设置环境变量
- 使用 简体中文

## 设计文档
- /root/ccx-source/docs/ccx-cli-design.md — 完整设计文档（666行）
- /root/ccx-source/cli/README.md — 用户手册（640行）
- /root/task-archive/ccx-cli-design/ — 历史版本存档
