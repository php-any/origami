<?php

namespace Agent\Workflow;

use Agent\Orchestrator;

class HandoffWorkflow
{
    private Orchestrator $orchestrator;
    private string $planner;
    private string $executor;

    public function __construct(Orchestrator $orchestrator, string $planner, string $executor)
    {
        $this->orchestrator = $orchestrator;
        $this->planner = $planner;
        $this->executor = $executor;
    }

    public function run(string $goal): array
    {
        $this->orchestrator->getContext()->set("goal", $goal);
        $steps = [];

        \Log::info("--- [{$this->planner}] 制定计划 ---");
        $planner = $this->orchestrator->getAgent($this->planner);
        $plan = $planner->chat(
            "目标: {$goal}\n\n请制定一个简洁的执行计划（分步骤列出）。",
            $this->orchestrator->getContext()
        );
        $this->orchestrator->getContext()->set("plan", $plan);
        $steps[] = ["agent" => $this->planner, "phase" => "plan", "output" => $plan];
        \Log::info("计划:\n{$plan}\n");

        \Log::info("--- [{$this->executor}] 执行计划 ---");
        $executor = $this->orchestrator->getAgent($this->executor);
        $result = $executor->chat(
            "目标: {$goal}\n\n执行计划:\n{$plan}\n\n请按照计划完成任务，直接给出最终结果。",
            $this->orchestrator->getContext()
        );
        $this->orchestrator->getContext()->set("result", $result);
        $steps[] = ["agent" => $this->executor, "phase" => "execute", "output" => $result];
        \Log::info("结果:\n{$result}");

        return ["steps" => $steps, "final" => $result, "plan" => $plan];
    }
}
