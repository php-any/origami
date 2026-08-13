<?php

namespace Agent;

class AgentMessage
{
    public string $from;
    public string $to;
    public string $content;
    public string $type;

    public function __construct(string $from, string $to, string $content, string $type = "message")
    {
        $this->from = $from;
        $this->to = $to;
        $this->content = $content;
        $this->type = $type;
    }

    public function format(): string
    {
        return "[{$this->from} → {$this->to}] {$this->content}";
    }
}
