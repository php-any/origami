<?php

namespace Agent\Workflow;

use Agent\Orchestrator;

class DebateWorkflow
{
    private Orchestrator $orchestrator;

    /** @var string[] */
    private array $debaters;

    private string $moderator;

    public function __construct(Orchestrator $orchestrator, array $debaters, string $moderator)
    {
        $this->orchestrator = $orchestrator;
        $this->debaters = $debaters;
        $this->moderator = $moderator;
    }

    public function run(string $topic, int $maxRounds = 2): array
    {
        $this->orchestrator->getContext()->set("topic", $topic);
        $transcript = "辩题: {$topic}\n";
        $steps = [];

        for ($round = 1; $round <= $maxRounds; $round++) {
            \Log::info("=== 第 {$round} 轮辩论 ===");
            foreach ($this->debaters as $name) {
                $agent = $this->orchestrator->getAgent($name);
                $prompt = "辩题: {$topic}\n\n";
                if ($transcript !== "") {
                    $prompt .= "之前的讨论记录:\n{$transcript}\n\n";
                }
                $prompt .= "请发表你的观点（简洁，200字以内）。";

                $reply = $agent->chat($prompt, $this->orchestrator->getContext());
                $transcript .= "\n[{$name}]: {$reply}";
                $steps[] = ["agent" => $name, "round" => $round, "output" => $reply];
                \Log::info("[{$name}]: {$reply}\n");
            }
        }

        \Log::info("=== 主持人总结 ===");
        $moderator = $this->orchestrator->getAgent($this->moderator);
        $summaryPrompt = "辩题: {$topic}\n\n完整讨论记录:\n{$transcript}\n\n请作为主持人，客观总结双方观点并给出你的判断。";
        $summary = $moderator->chat($summaryPrompt, $this->orchestrator->getContext());
        $steps[] = ["agent" => $this->moderator, "round" => "summary", "output" => $summary];

        $this->orchestrator->getContext()->set("transcript", $transcript);
        $this->orchestrator->getContext()->set("summary", $summary);

        \Log::info("[{$this->moderator} 总结]:\n{$summary}");
        return ["steps" => $steps, "final" => $summary, "transcript" => $transcript];
    }
}
