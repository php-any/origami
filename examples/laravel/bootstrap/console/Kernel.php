<?php

namespace Bootstrap\Console;

use Cli\Annotation\CliApplication;

/**
 * Artisan CLI 内核（扫描 app/Console/Commands）
 */
#[CliApplication(name: 'Laravel Framework', version: '1.0.0', scan: __DIR__ . '/../../app/Console/Commands')]
class Kernel
{
    public static function boot(): void {}

    public static function exit(): void {}
}
