<?php

namespace App;

use Cli\Annotation\CliApplication;

/**
 * CLI 应用入口
 *
 * scan 指向 Command/ 目录，自动发现所有 #[Command] 注解的类。
 * boot() 和 exit() 是生命周期钩子，路由在 run.php 中处理。
 */
#[CliApplication(name: "testrunner", version: "1.0.0", scan: __DIR__ . "/Command")]
class TestRunnerApp
{
    public static function boot(): void {

    }
    public static function exit(): void {}
}
