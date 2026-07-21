<?php

namespace Bootstrap\Console;

/**
 * Artisan 命令契约（对应 Laravel Illuminate\Contracts\Console\Isolatable 等核心 handle 约定）
 */
interface CommandInterface
{
    public function handle(): void;
}
