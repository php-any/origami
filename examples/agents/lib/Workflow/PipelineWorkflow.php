<?php

namespace Agent\Workflow;

use Agent\Orchestrator;

class PipelineWorkflow
{
    private Orchestrator $orchestrator;

    /** @var string[] */
    private array $steps;

    public function __construct(Orchestrator $orchestrator, array $steps)
    {
        $this->orchestrator = $orchestrator;
        $this->steps = $steps;
    }

    public function run(string $task): array
    {
        $this->orchestrator->getContext()->set("task", $task);
        $output = $task;
        $steps = [];

        foreach ($this->steps as $agentName) {
            $agent = $this->orchestrator->getAgent($agentName);
            \Log::info("--- [{$agentName}] 开始处理 ---");

            $output = $agent->chat($output, $this->orchestrator->getContext());
            $this->orchestrator->getContext()->set("{$agentName}_output", $output);
            $steps[] = ["agent" => $agentName, "output" => $output];

            \Log::info("[{$agentName}] 输出:\n{$output}\n");
        }

        $this->orchestrator->getContext()->set("final_output", $output);
        return ["steps" => $steps, "final" => $output];
    }
}
