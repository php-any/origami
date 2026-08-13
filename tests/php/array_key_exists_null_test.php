<?php

namespace tests\php;

/**
 * array_key_exists：键存在且值为 null 时返回 true（与 isset 不同）
 */

$literal = ['x' => null];
if (!array_key_exists('x', $literal)) {
    Log::fatal("array_key_exists('x', ['x' => null]) 应为 true");
}
if (isset($literal['x'])) {
    Log::fatal("isset(['x' => null]['x']) 应为 false");
}
if (array_key_exists('missing', $literal)) {
    Log::fatal("array_key_exists('missing', ...) 应为 false");
}

$assigned = [];
$assigned['y'] = null;
if (!array_key_exists('y', $assigned)) {
    Log::fatal("赋值 null 后 array_key_exists 应为 true");
}

$mixed = ['a' => 1, 'b' => null, 'c' => false];
if (!array_key_exists('b', $mixed) || !array_key_exists('c', $mixed)) {
    Log::fatal("null/false 值的键仍应存在");
}

Log::info('array_key_exists null 测试通过');
