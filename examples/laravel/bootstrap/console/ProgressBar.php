<?php

namespace Bootstrap\Console;

/**
 * 终端进度条
 */
class ProgressBar
{
    private int $current = 0;

    private bool $started = false;

    public function __construct(
        private readonly Output $output,
        private readonly int $max = 100,
        private readonly int $width = 36,
    ) {}

    public function start(): void
    {
        $this->started = true;
        $this->current = 0;
        $this->draw();
    }

    public function advance(int $step = 1): void
    {
        $this->current = min($this->max, $this->current + $step);
        if ($this->started) {
            $this->draw();
        }
    }

    public function finish(): void
    {
        $this->current = $this->max;
        if ($this->started) {
            $this->draw();
            $this->output->newLine();
        }
    }

    private function draw(): void
    {
        $percent = $this->max > 0 ? (int) floor(($this->current / $this->max) * 100) : 100;
        $filled = $this->max > 0 ? (int) floor(($this->current / $this->max) * $this->width) : $this->width;
        $bar = str_repeat('=', max(0, $filled - 1)) . '>' . str_repeat('-', max(0, $this->width - $filled));

        $this->output->write("\r  {$this->current}/{$this->max} [{$bar}] {$percent}%");
    }
}
