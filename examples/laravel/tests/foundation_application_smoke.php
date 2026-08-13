<?php

/**
 * Foundation Application 冒烟：app() / environment / config / Dispatchable。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();

use Illuminate\Foundation\Application;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Config\Repository;

$illuminateApp = illuminate_container();
if (!($illuminateApp instanceof Application)) {
    echo "FAIL: illuminate_container not Application\n";
    exit(1);
}

if (!function_exists('app')) {
    echo "FAIL: app() helper missing\n";
    exit(1);
}

$viaHelper = \app();
if (!($viaHelper instanceof Application)) {
    echo "FAIL: app() did not return Application\n";
    exit(1);
}

if ($viaHelper !== $illuminateApp) {
    echo "FAIL: app() instance mismatch\n";
    exit(1);
}

$env = $illuminateApp->environment();
if (!is_string($env) || $env === '') {
    echo "FAIL: environment()\n";
    var_export($env);
    echo "\n";
    exit(1);
}

if (!$illuminateApp->environment($env)) {
    echo "FAIL: environment($env) match\n";
    exit(1);
}

$console = $illuminateApp->runningInConsole();
if (!is_bool($console)) {
    echo "FAIL: runningInConsole type\n";
    exit(1);
}

$config = $illuminateApp->make('config');
if (!($config instanceof Repository)) {
    echo "FAIL: config binding\n";
    exit(1);
}

$name = $config->get('app.name');
if (!is_string($name) || $name === '') {
    echo "FAIL: app.name\n";
    exit(1);
}

// Dispatchable 是 trait；Origami 对 trait 自动加载可能受限，先验收 PendingDispatch 类
if (!class_exists(\Illuminate\Foundation\Bus\PendingDispatch::class)) {
    echo "FAIL: PendingDispatch missing\n";
    exit(1);
}

$dispatchableFile = dirname(__DIR__) . '/vendor/laravel/framework/src/Illuminate/Foundation/Bus/Dispatchable.php';
if (!is_file($dispatchableFile)) {
    echo "FAIL: Dispatchable.php missing on disk\n";
    exit(1);
}

echo "PASS\n";
