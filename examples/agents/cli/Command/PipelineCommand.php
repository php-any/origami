<?php

namespace Console\Command;

use App\Service\AgentRunnerService;
use Cli\Annotation\Command;

#[Command(name: "pipeline", description: "流水线工作流 — 研究员 → 写作者 → 编辑")]
class PipelineCommand
{
    public function execute(): void
    {
        $runner = new AgentRunnerService();
        $topic = "折言（Origami）是一门用 Go 实现的 PHP 风格脚本语言";

        \Log::info("=== 流水线工作流 (MODEL: {$runner->getModel()}) ===");
        \Log::info("主题: {$topic}\n");

        $result = $runner->run("pipeline", $topic);

        \Log::info("=== 最终结果 ===");
        \Log::info($result["final"]);
    }
}
