<?php

namespace tests\pdo;

/**
 * pdo：PDO 类与 FETCH_ASSOC 常量。
 */

if (!class_exists('PDO')) {
    Log::fatal('PDO 类未注册');
}
$v = \PDO::FETCH_ASSOC;
if ($v === null || $v === false || $v === '') {
    Log::fatal('PDO::FETCH_ASSOC 缺失');
}

Log::info('pdo 类存在测试通过');
