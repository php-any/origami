<?php

namespace Agent;

class Orchestrator
{
    /** @var array<string, Agent> */
    private array $agents = [];

    /** @var AgentMessage[] */
    private array $messageLog = [];

    private AgentContext $context;

    public function __construct()
    {
        $this->context = new AgentContext();
    }

    public function addAgent(Agent $agent): self
    {
        $this->agents[$agent->name] = $agent;
        return $this;
    }

    public function getAgent(string $name): Agent
    {
        if (!isset($this->agents[$name])) {
            throw new \Exception("Agent 不存在: {$name}");
        }
        return $this->agents[$name];
    }

    public function hasAgent(string $name): bool
    {
        return isset($this->agents[$name]);
    }

    public function getAgentNames(): array
    {
        return array_keys($this->agents);
    }

    public function getContext(): AgentContext
    {
        return $this->context;
    }

    public function send(string $from, string $to, string $content): string
    {
        $this->messageLog[] = new AgentMessage($from, $to, $content);

        $agent = $this->getAgent($to);
        $prompt = "来自 Agent「{$from}」的消息:\n\n{$content}";
        $reply = $agent->chat($prompt, $this->context);

        $this->messageLog[] = new AgentMessage($to, $from, $reply, "reply");
        $this->context->set("last_reply_from_{$to}", $reply);

        return $reply;
    }

    public function broadcast(string $from, string $content): array
    {
        $replies = [];
        foreach ($this->agents as $name => $agent) {
            if ($name === $from) {
                continue;
            }
            $replies[$name] = $this->send($from, $name, $content);
        }
        return $replies;
    }

    public function getMessageLog(): array
    {
        return $this->messageLog;
    }

    public function dumpLog(): void
    {
        foreach ($this->messageLog as $msg) {
            \Log::info($msg->format());
        }
    }
}
