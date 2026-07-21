<?php

namespace Bootstrap\Console;

/**
 * Artisan 命令基类（对应 Laravel Illuminate\Console\Command）
 */
abstract class Command implements CommandInterface
{
    protected Output $output;

    protected ArgvInput $input;

    protected Prompt $prompt;

    /**
     * 由 CLI 运行时调用的入口（对应 Symfony Command::execute）
     */
    public function execute(): void
    {
        $this->output = new Output();
        $this->input = ArgvInput::fromArgv($this->defineInput());
        $this->prompt = new Prompt($this->output);
        $this->handle();
    }

    protected function defineInput(): InputDefinition
    {
        return new InputDefinition();
    }

    abstract public function handle(): void;

    protected function table(): Table
    {
        return new Table($this->output);
    }

    protected function createProgressBar(int $max): ProgressBar
    {
        return new ProgressBar($this->output, $max);
    }
}
