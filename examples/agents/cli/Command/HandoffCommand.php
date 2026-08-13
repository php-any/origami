<?php

namespace Console\Command;

use App\Service\AgentRunnerService;
use Cli\Annotation\Command;

#[Command(name: "handoff", description: "接力工作流 — 规划师制定计划 → 执行者完成任务")]
class HandoffCommand
{
    public function execute(): void
    {
        $runner = new AgentRunnerService();
        $goal = "为一个小型博客系统设计 REST API 端点列表";

        \Log::info("=== 接力工作流 (MODEL: {$runner->getModel()}) ===");
        \Log::info("目标: {$goal}\n");

        $result = $runner->run("handoff", $goal);

        \Log::info("=== 最终结果 ===");
        \Log::info($result["final"]);
    }
}
