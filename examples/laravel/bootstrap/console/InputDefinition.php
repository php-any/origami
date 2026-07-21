<?php

namespace Bootstrap\Console;

/**
 * 命令行输入定义（位置参数 + 选项）
 */
class InputDefinition
{
    /** @var array<string, array{required: bool, description: string, default: mixed}> */
    private array $arguments = [];

    /** @var array<string, array{shortcut: ?string, accept_value: bool, description: string, default: mixed}> */
    private array $options = [];

    public function addArgument(string $name, bool $required = false, string $description = '', mixed $default = null): self
    {
        $this->arguments[$name] = [
            'required' => $required,
            'description' => $description,
            'default' => $default,
        ];

        return $this;
    }

    public function addOption(
        string $name,
        ?string $shortcut = null,
        bool $acceptValue = false,
        mixed $default = null,
        string $description = '',
    ): self {
        $this->options[$name] = [
            'shortcut' => $shortcut,
            'accept_value' => $acceptValue,
            'description' => $description,
            'default' => $default,
        ];

        return $this;
    }

    /** @return array<string, array{required: bool, description: string, default: mixed}> */
    public function getArguments(): array
    {
        return $this->arguments;
    }

    /** @return array<string, array{shortcut: ?string, accept_value: bool, description: string, default: mixed}> */
    public function getOptions(): array
    {
        return $this->options;
    }
}
