<?php

namespace Console;

use Cli\Annotation\CliApplication;

/**
 * CLI 应用入口
 *
 * #[CliApplication] 扫描 cli/Command 目录下的所有 #[Command]。
 */
#[CliApplication(name: "agents", version: "1.0.0", scan: __DIR__ . "/Command")]
class AgentApp
{
    public static function boot(): void
    {
        \Log::info("Origami Multi-Agent CLI v1.0.0");
    }

    public static function exit(): void {}
}
