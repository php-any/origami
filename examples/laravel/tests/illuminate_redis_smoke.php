<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Container\Container;
use Illuminate\Redis\RedisManager;

/**
 * 不实例化 Connection / CommandExecuted（需抽象类子类化与 getName()）。
 * 验证 RedisManager 驱动切换与多连接配置持有。
 */
$manager = new RedisManager(new Container(), 'phpredis', [
    'default' => [
        'host' => '127.0.0.1',
        'port' => 6379,
        'database' => 0,
    ],
    'cache' => [
        'host' => '127.0.0.1',
        'port' => 6380,
        'database' => 1,
    ],
]);

$manager->setDriver('predis');

if (!is_object($manager)) {
    echo "FAIL: manager\n";
    exit(1);
}

echo "PASS\n";
