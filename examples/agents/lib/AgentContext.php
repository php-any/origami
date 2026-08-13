<?php

namespace Agent;

class AgentContext
{
    private array $data = [];

    public function set(string $key, mixed $value): void
    {
        $this->data[$key] = $value;
    }

    public function get(string $key, mixed $default = null): mixed
    {
        if (isset($this->data[$key])) {
            return $this->data[$key];
        }
        return $default;
    }

    public function has(string $key): bool
    {
        return isset($this->data[$key]);
    }

    public function getAll(): array
    {
        return $this->data;
    }

    public function merge(array $items): void
    {
        foreach ($items as $k => $v) {
            $this->data[$k] = $v;
        }
    }

    public function toPromptSection(): string
    {
        if (count($this->data) === 0) {
            return "";
        }
        $lines = ["## 共享上下文"];
        foreach ($this->data as $key => $value) {
            if (is_string($value) || is_numeric($value)) {
                $lines[] = "- {$key}: {$value}";
            } else {
                $lines[] = "- {$key}: " . json_encode($value);
            }
        }
        return implode("\n", $lines);
    }
}
