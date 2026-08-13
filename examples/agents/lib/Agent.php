<?php

namespace Agent;

/**
 * 单个 LLM Agent
 */
class Agent
{
    public string $name;
    public string $role;

    private $client;
    private string $model;
    private array $history = [];

    public function __construct(string $name, string $role, $client, string $model)
    {
        $this->name = $name;
        $this->role = $role;
        $this->client = $client;
        $this->model = $model;
    }

    public function chat(string $input, AgentContext $context = null): string
    {
        $messages = [
            new \OpenAI\SystemMessage($this->buildSystemPrompt($context)),
        ];

        foreach ($this->history as $msg) {
            $messages[] = $msg;
        }
        $messages[] = new \OpenAI\UserMessage($input);

        $result = $this->client->chat($this->model, $messages, [
            \OpenAI\ChatOption::TEMPERATURE => 0.7,
        ]);

        $this->history[] = new \OpenAI\UserMessage($input);
        $this->history[] = new \OpenAI\AssistantMessage($result->content);

        return $result->content;
    }

    public function reset(): void
    {
        $this->history = [];
    }

    private function buildSystemPrompt(AgentContext $context = null): string
    {
        $parts = [$this->role];
        if ($context !== null) {
            $section = $context->toPromptSection();
            if ($section !== "") {
                $parts[] = $section;
            }
        }
        $parts[] = "你是 Agent「{$this->name}」，请专注于你的职责。";
        return implode("\n\n", $parts);
    }
}
