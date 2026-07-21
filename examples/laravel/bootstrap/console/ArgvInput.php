<?php

namespace Bootstrap\Console;

/**
 * 从 $argv / $_SERVER['argv'] 解析命令行输入
 */
class ArgvInput
{
    /** @var array<string, mixed> */
    private array $arguments = [];

    /** @var array<string, mixed> */
    private array $options = [];

    public static function fromArgv(InputDefinition $definition, ?array $rawArgv = null): self
    {
        $raw = $rawArgv ?? ($_SERVER['argv'] ?? []);
        $tokens = array_slice($raw, 1);

        return new self($definition, $tokens);
    }

    private function __construct(InputDefinition $definition, array $tokens)
    {
        $this->bootstrapDefaults($definition);
        $this->parse($definition, $tokens);
    }

    private function bootstrapDefaults(InputDefinition $definition): void
    {
        foreach ($definition->getArguments() as $name => $meta) {
            $this->arguments[$name] = $meta['default'];
        }

        foreach ($definition->getOptions() as $name => $meta) {
            $this->options[$name] = $meta['default'];
        }
    }

    private function parse(InputDefinition $definition, array $tokens): void
    {
        $positional = [];
        $options = $definition->getOptions();
        $optionByShortcut = [];

        foreach ($options as $name => $meta) {
            if ($meta['shortcut'] !== null && $meta['shortcut'] !== '') {
                $optionByShortcut[$meta['shortcut']] = $name;
            }
        }

        $index = 0;
        while ($index < count($tokens)) {
            $token = $tokens[$index];

            if (str_starts_with($token, '--')) {
                $this->parseLongOption($token, $tokens, $index, $options);
                continue;
            }

            if (str_starts_with($token, '-') && strlen($token) > 1) {
                $this->parseShortOption($token, $tokens, $index, $options, $optionByShortcut);
                continue;
            }

            $positional[] = $token;
            $index++;
        }

        $argumentNames = array_keys($definition->getArguments());
        foreach ($argumentNames as $position => $name) {
            if (array_key_exists($position, $positional)) {
                $this->arguments[$name] = $positional[$position];
            }
        }
    }

    private function parseLongOption(string $token, array $tokens, int &$index, array $options): void
    {
        $body = substr($token, 2);
        $value = null;

        if (str_contains($body, '=')) {
            [$name, $value] = explode('=', $body, 2);
        } else {
            $name = $body;
        }

        if (!array_key_exists($name, $options)) {
            $index++;
            return;
        }

        $meta = $options[$name];
        if ($meta['accept_value']) {
            if ($value === null) {
                $index++;
                if ($index >= count($tokens)) {
                    return;
                }
                $value = $tokens[$index];
            }
            $this->options[$name] = $value;
        } else {
            $this->options[$name] = true;
        }

        $index++;
    }

    private function parseShortOption(
        string $token,
        array $tokens,
        int &$index,
        array $options,
        array $optionByShortcut,
    ): void {
        $flags = substr($token, 1);

        if (strlen($flags) === 1) {
            $name = $optionByShortcut[$flags] ?? null;
            if ($name === null) {
                $index++;
                return;
            }

            $meta = $options[$name];
            if ($meta['accept_value']) {
                $index++;
                if ($index >= count($tokens)) {
                    return;
                }
                $this->options[$name] = $tokens[$index];
            } else {
                $this->options[$name] = true;
            }

            $index++;
            return;
        }

        for ($i = 0; $i < strlen($flags); $i++) {
            $shortcut = $flags[$i];
            $name = $optionByShortcut[$shortcut] ?? null;
            if ($name === null) {
                continue;
            }

            $meta = $options[$name];
            if ($meta['accept_value']) {
                $this->options[$name] = substr($flags, $i + 1);
                break;
            }

            $this->options[$name] = true;
        }

        $index++;
    }

    public function getArgument(string $name): mixed
    {
        return $this->arguments[$name] ?? null;
    }

    public function getOption(string $name): mixed
    {
        return $this->options[$name] ?? null;
    }

    public function hasOption(string $name): bool
    {
        $value = $this->getOption($name);

        return $value !== null && $value !== false && $value !== '';
    }

    /** @return array<string, mixed> */
    public function getArguments(): array
    {
        return $this->arguments;
    }

    /** @return array<string, mixed> */
    public function getOptions(): array
    {
        return $this->options;
    }
}
