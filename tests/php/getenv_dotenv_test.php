<?php

namespace tests\php;

/**
 * 验证 getenv 能读取 $_ENV / $_SERVER 中的变量（Symfony Dotenv 默认写入此处，不调用 putenv）。
 */

// 模拟 Dotenv::populate 行为
$_ENV['APP_URL'] = 'http://localhost:8080';
$_SERVER['APP_URL'] = 'http://localhost:8080';

$url = getenv('APP_URL');
if ($url !== 'http://localhost:8080') {
    Log::fatal('getenv(APP_URL) 失败: ' . var_export($url, true));
}

// $_ENV 优先于 OS 环境（若存在同名变量）
$_ENV['PATH'] = '/custom/path';
$path = getenv('PATH');
if ($path !== '/custom/path') {
    Log::fatal('getenv 应优先返回 $_ENV 中的值: ' . var_export($path, true));
}

Log::info('getenv $_ENV 回退测试通过');
