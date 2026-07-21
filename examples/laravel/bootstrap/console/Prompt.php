<?php

namespace Bootstrap\Console;

/**
 * 交互式用户输入（ask / confirm / choice / secret）
 */
class Prompt
{
    public function __construct(
        private readonly Output $output,
    ) {}

    public function ask(string $question, ?string $default = null): string
    {
        $suffix = $default !== null ? ' [' . $default . ']' : '';
        $this->output->write("  {$question}{$suffix}: ");

        $answer = $this->readLine();
        if ($answer === '' && $default !== null) {
            return $default;
        }

        return $answer;
    }

    public function confirm(string $question, bool $default = true): bool
    {
        $hint = $default ? 'Y/n' : 'y/N';
        $this->output->write("  {$question} ({$hint}): ");

        $answer = strtolower(trim($this->readLine()));

        if ($answer === '') {
            return $default;
        }

        return in_array($answer, ['y', 'yes', '1', 'true'], true);
    }

    public function choice(string $question, array $choices, ?string $default = null): string
    {
        $this->output->writeln('  ' . $question);
        foreach ($choices as $choice) {
            $marker = ($default !== null && $choice === $default) ? '*' : ' ';
            $this->output->writeln("    [{$marker}] {$choice}");
        }

        $suffix = $default !== null ? " [{$default}]" : '';
        $this->output->write("  请选择{$suffix}: ");

        $answer = trim($this->readLine());
        if ($answer === '' && $default !== null) {
            return $default;
        }

        if (in_array($answer, $choices, true)) {
            return $answer;
        }

        $this->output->warning('无效选项，使用默认值。');
        return $default ?? $choices[0];
    }

    public function secret(string $question): string
    {
        $this->output->write("  {$question}: ");

        if (function_exists('readline')) {
            $value = readline('');
            return is_string($value) ? trim($value) : '';
        }

        return trim($this->readLine());
    }

    private function readLine(): string
    {
        $stdin = new \SplFileObject('php://stdin');
        $line = $stdin->fgets();

        return is_string($line) ? rtrim($line, "\r\n") : '';
    }
}
