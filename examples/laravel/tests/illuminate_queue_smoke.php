<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Container\Container;
use Illuminate\Queue\Jobs\JobName;
use Illuminate\Queue\SyncQueue;

/**
 * SyncQueue::push 依赖 ramsey/uuid（UuidInterface 解析仍为上游缺口）。
 * 验证 SyncQueue 容器绑定与 JobName 工具类。
 */
$container = new Container();
$queue = new SyncQueue();
$queue->setContainer($container);

if ($queue->getContainer() !== $container) {
    echo "FAIL: container binding\n";
    exit(1);
}

$parsed = JobName::parse('App\\Jobs\\Demo@handle');
if (($parsed[0] ?? null) !== 'App\\Jobs\\Demo' || ($parsed[1] ?? null) !== 'handle') {
    echo "FAIL: JobName::parse\n";
    var_export($parsed);
    echo "\n";
    exit(1);
}

echo "PASS\n";
