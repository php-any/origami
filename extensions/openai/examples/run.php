<?php
/**
 * OpenAI 扩展测试 CLI
 *
 * 项目结构:
 *   src/App/TestRunnerApp.php        应用入口 (#[CliApplication])
 *   src/App/Command/ChatCommand.php   命令类 (#[Command])
 *   tests/chat_test.php              测试脚本
 *
 * 用法:
 *   ./testrunner examples/run.php chat      # Chat Completions
 *   ./testrunner examples/run.php json      # JSON 输出
 *   ./testrunner examples/run.php schema    # JSON Schema
 *   ./testrunner examples/run.php error     # 错误处理
 *   ./testrunner examples/run.php           # 全部测试（默认）
 */

// 加载应用入口，触发 #[CliApplication] 扫描 Command/ 目录，注册所有 #[Command]
include __DIR__ . "/src/App/TestRunnerApp.php";

// 命令路由由 Go 层 annotation.ExecuteCommand() 处理
// 无需在此做 switch/if-else
