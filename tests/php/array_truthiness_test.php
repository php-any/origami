<?php

namespace tests\php;

/**
 * 非空数组在布尔上下文中应为 true（Symfony Table::render 的 if (!$row) 依赖此点）。
 */

$row = ['1', 'Alice'];
if (!$row) {
    Log::fatal('non-empty array must be truthy');
}
if ($row) {
    Log::info('non-empty array is truthy OK');
} else {
    Log::fatal('non-empty array evaluated false');
}

$empty = [];
if ($empty) {
    Log::fatal('empty array must be falsy');
}

$sep = new \stdClass();
// stdClass might not exist - use ArrayObject or simple
$obj = (object)[];
if (!$obj) {
    Log::info('empty object falsy?');
}

Log::info('array_truthiness_test 测试通过');
