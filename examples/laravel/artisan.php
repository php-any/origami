<?php

/**
 * Artisan CLI 引导（由 laravel <command> 调用）
 */

require __DIR__ . '/bootstrap/app.php';
bootstrap_app();
bootstrap_cli_container();

require __DIR__ . '/bootstrap/console/Kernel.php';
