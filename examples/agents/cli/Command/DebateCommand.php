<?php

namespace Console\Command;

use App\Service\AgentRunnerService;
use Cli\Annotation\Command;

#[Command(name: "debate", description: "辩论工作流 — 正方 vs 反方 → 主持人总结")]
class DebateCommand
{
    public function execute(): void
    {
        $runner = new AgentRunnerService();
        $topic = "AI 辅助编程是否应该成为开发者的默认工作方式？";

        \Log::info("=== 辩论工作流 (MODEL: {$runner->getModel()}) ===");
        \Log::info("辩题: {$topic}\n");

        $result = $runner->run("debate", $topic);

        \Log::info("=== 主持人总结 ===");
        \Log::info($result["final"]);
    }
}
