# Origami 多 Agent 框架

基于 `extensions/openai` 的多 Agent 协作框架，支持 Web 可视化开发与 PHP 热更新。

参照 `examples/spring` 的组织方式：**框架代码**、**业务代码**、**命令行**三者分离，
命名空间通过标准 PHP `spl_autoload_register`（`autoload.php`）解析，与 `composer.json` 中的 `psr-4` 保持一致。

## 目录结构

```
examples/agents/
├── autoload.php            # PSR-4 自动加载（spl_autoload_register）
├── composer.json           # PSR-4 映射（Agent\/App\/Console\）
├── index.php               # Web 入口
├── run.php                 # CLI 入口
├── main.go / dev.go        # Go 运行时（buildVM + 热更新开发服务器）
├── config/
│   └── app.php             # 返回配置数组（读取 .env）
├── lib/                    # 框架代码（Agent\ 命名空间）
│   ├── Agent.php
│   ├── AgentContext.php
│   ├── AgentMessage.php
│   ├── Orchestrator.php
│   └── Workflow/
│       ├── PipelineWorkflow.php
│       ├── DebateWorkflow.php
│       └── HandoffWorkflow.php
├── src/                    # 业务代码（App\ 命名空间）
│   ├── AgentsApplication.php
│   ├── Controller/
│   └── Service/
├── cli/                    # 命令行（Console\ 命名空间）
│   ├── AgentApp.php
│   └── Command/
└── public/                 # 前端静态资源
```

## 快速开始

```bash
cd examples/agents
cp .env.example .env        # 填入 OpenAI API Key
go mod tidy
go build -o agents .

./agents dev                # Web 开发服务（热更新，默认 8080）
./agents dev 8090           # 指定端口
./agents index.php          # Web 生产服务
./agents run.php pipeline   # CLI：流水线
./agents run.php debate     # CLI：辩论
./agents run.php handoff    # CLI：接力
```

## 分层职责

| 目录 | 命名空间 | 职责 |
|------|----------|------|
| `lib/` | `Agent\` | 可复用框架：Agent、Orchestrator、Workflow |
| `src/` | `App\` | 业务应用：控制器、服务、应用引导 |
| `cli/` | `Console\` | 命令行入口与示例命令 |
| `config/` | — | `app.php` 返回配置数组 |

框架能力放 `lib/`，业务逻辑放 `src/`，命令放 `cli/`。三者互不混杂。

## 加载机制

- `autoload.php` 用 `spl_autoload_register` 注册 `Agent\ => lib/`、`App\ => src/`、`Console\ => cli/`，按需加载类文件。
- Go 侧 `buildVM()` 只负责装配运行时（与标准 `zy` 一致 + OpenAI 扩展），**不涉及任何命名空间硬编码**。
- 入口文件 `index.php` / `run.php` 只 `require autoload.php`，其余类均自动加载。

## 热更新

`./agents dev` 监听 `src/`、`lib/`、`cli/`、`config/`、`public/` 及 `index.php`、`autoload.php` 的变更，
自动以全新的 VM 重新加载 PHP 代码，无需重启 Go 进程。每次重载会清理进程级全局缓存
（include 缓存、autoload 注册表、HTTP 路由与注解扫描状态），避免旧状态污染。

前端通过轮询 `/api/dev/status` 的版本号感知重载。

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/workflows` | 工作流列表 |
| POST | `/api/agents/run` | 运行工作流（body: `{"workflow": "...", "input": "..."}`） |
| GET | `/api/dev/status` | 热更新版本号 |
