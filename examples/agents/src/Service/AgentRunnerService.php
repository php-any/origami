<?php

namespace App\Service;

use Agent\Agent;
use Agent\Orchestrator;
use Agent\Workflow\DebateWorkflow;
use Agent\Workflow\HandoffWorkflow;
use Agent\Workflow\PipelineWorkflow;
use Container\Annotation\Singleton;

/**
 * Agent 运行服务
 *
 * 从 config/app.php 读取配置，构建 OpenAI 客户端，
 * 提供三种工作流的统一执行入口，供 Web 控制器与 CLI 命令复用。
 */
#[Singleton]
class AgentRunnerService
{
    private $client;
    private string $model;

    public function __construct()
    {
        $config = require dirname(dirname(__DIR__)) . "/config/app.php";
        $openai = $config["openai"];

        $this->client = new \OpenAI\Client($openai["key"], $openai["base_url"]);
        $this->model = $openai["model"];
    }

    public function getModel(): string
    {
        return $this->model;
    }

    public function listWorkflows(): array
    {
        return [
            ["id" => "pipeline", "name" => "流水线", "description" => "研究员 → 写作者 → 编辑"],
            ["id" => "debate", "name" => "辩论", "description" => "正方 vs 反方 → 主持人总结"],
            ["id" => "handoff", "name" => "接力", "description" => "规划师 → 执行者"],
        ];
    }

    public function run(string $workflow, string $input): array
    {
        return match ($workflow) {
            "pipeline" => $this->runPipeline($input),
            "debate" => $this->runDebate($input),
            "handoff" => $this->runHandoff($input),
            default => throw new \Exception("未知工作流: {$workflow}"),
        };
    }

    private function runPipeline(string $task): array
    {
        $orchestrator = new Orchestrator();
        $orchestrator
            ->addAgent(new Agent("研究员", "你是专业研究员。根据用户主题，列出 3-5 个关键要点，每个要点一句话。只输出要点列表。", $this->client, $this->model))
            ->addAgent(new Agent("写作者", "你是技术博客写作者。根据研究员提供的要点，写一篇 150 字以内的短文。语言流畅、结构清晰。", $this->client, $this->model))
            ->addAgent(new Agent("编辑", "你是资深编辑。审阅文章，修正语病并润色，保持原意。输出最终版本。", $this->client, $this->model));

        $result = (new PipelineWorkflow($orchestrator, ["研究员", "写作者", "编辑"]))->run($task);
        $result["workflow"] = "pipeline";
        $result["input"] = $task;
        return $result;
    }

    private function runDebate(string $topic): array
    {
        $orchestrator = new Orchestrator();
        $orchestrator
            ->addAgent(new Agent("正方", "你是辩论正方。支持 AI 辅助编程能显著提升开发效率。观点要有理有据，简洁有力。", $this->client, $this->model))
            ->addAgent(new Agent("反方", "你是辩论反方。认为过度依赖 AI 会降低开发者基本功、引入安全隐患。观点要有理有据，简洁有力。", $this->client, $this->model))
            ->addAgent(new Agent("主持人", "你是中立主持人。客观总结双方观点，指出各自优劣，给出平衡的判断。", $this->client, $this->model));

        $result = (new DebateWorkflow($orchestrator, ["正方", "反方"], "主持人"))->run($topic, 2);
        $result["workflow"] = "debate";
        $result["input"] = $topic;
        return $result;
    }

    private function runHandoff(string $goal): array
    {
        $orchestrator = new Orchestrator();
        $orchestrator
            ->addAgent(new Agent("规划师", "你是项目规划师。将目标拆解为 3-5 个可执行步骤，每步一句话。只输出步骤列表。", $this->client, $this->model))
            ->addAgent(new Agent("执行者", "你是执行者。严格按照规划师的步骤完成任务，输出简洁的最终结果。", $this->client, $this->model));

        $result = (new HandoffWorkflow($orchestrator, "规划师", "执行者"))->run($goal);
        $result["workflow"] = "handoff";
        $result["input"] = $goal;
        return $result;
    }
}
