<?php

namespace Bootstrap\Console;

/**
 * Artisan 命令契约（handle 约定，与 Illuminate\Console\Command 一致）
 */
interface CommandInterface
{
    public function handle(): void;
}
