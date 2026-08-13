<?php

// bootstrap_smoke.php
// 阶段 2：验证 bootstrap/app.php 能正常创建 Application 并完成引导。

require __DIR__.'/../../vendor/autoload.php';

/** @var Illuminate\Foundation\Application $app */
$app = require __DIR__.'/../../bootstrap/app.php';

$checks = 0;

// 1. Application 实例正确
if (!$app instanceof Illuminate\Foundation\Application) {
    echo "FAIL: 不是 Application 实例\n";
    exit(1);
}
$checks++;

// 2. basePath 正确
if ($app->basePath() !== dirname(__DIR__, 2)) {
    echo "FAIL: basePath 不正确\n";
    exit(1);
}
$checks++;

// 3. 执行 bootstrap
$app->bootstrapWith([
    Illuminate\Foundation\Bootstrap\LoadEnvironmentVariables::class,
    Illuminate\Foundation\Bootstrap\LoadConfiguration::class,
]);

// 4. config 服务已绑定
if (!$app->bound('config')) {
    echo "FAIL: config 服务未绑定\n";
    exit(1);
}
$checks++;

// 5. config 能正常读取
$name = $app['config']->get('app.name');
if ($name !== 'Laravel') {
    echo "FAIL: app.name 不正确，got: $name\n";
    exit(1);
}
$checks++;

echo "OK: bootstrap smoke passed ($checks checks)\n";
