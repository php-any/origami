<?php

/**
 * PSR-4 自动加载器
 *
 * 与 composer.json 中的 autoload.psr-4 保持一致：
 *   Agent\   => lib/   （框架代码）
 *   App\     => src/   （业务代码）
 *   Console\ => cli/   （命令行）
 *
 * 采用标准 spl_autoload_register，运行时按需加载类文件，
 * 无需 Go 侧硬编码命名空间，也无需手动逐个 include。
 */

spl_autoload_register(function (string $class): void {
    $prefixes = [
        "Agent\\" => __DIR__ . "/lib/",
        "App\\" => __DIR__ . "/src/",
        "Console\\" => __DIR__ . "/cli/",
    ];

    foreach ($prefixes as $prefix => $baseDir) {
        if (!str_starts_with($class, $prefix)) {
            continue;
        }
        $relative = substr($class, strlen($prefix));
        $file = $baseDir . str_replace("\\", "/", $relative) . ".php";
        if (is_file($file)) {
            require $file;
        }
        return;
    }
});
