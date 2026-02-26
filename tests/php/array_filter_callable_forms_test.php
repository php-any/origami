<?php

namespace tests\php;

/**
 * array_filter 各种 callable 形式测试：
 *
 * 1) 字符串函数名：'strlen'
 * 2) 静态方法回调：['tests\php\ArrayFilterHelper', 'is_even']
 */

class ArrayFilterHelper
{
    public static function is_even(int $n): bool
    {
        return $n % 2 === 0;
    }
}

// 1. 字符串函数名：strlen
$strings = ['', 'a', 'foo', ''];
$filteredStrings = array_filter($strings, 'strlen');

// strlen 返回 0/1/3，0 会被视为 false，因此只保留 'a' 和 'foo'
if (array_values($filteredStrings) !== ['a', 'foo']) {
    Log::fatal(
        'array_filter 使用字符串函数名 strlen 失败，实际='
        . json_encode(array_values($filteredStrings))
    );
}

// 2. 静态方法回调：['tests\php\ArrayFilterHelper', 'is_even']
$numbers = [1, 2, 3, 4, 5];
$filteredEvens = array_filter($numbers, [ArrayFilterHelper::class, 'is_even']);

if (array_values($filteredEvens) !== [2, 4]) {
    Log::fatal(
        'array_filter 使用静态方法回调失败，实际='
        . json_encode(array_values($filteredEvens))
    );
}

Log::info('array_filter 多种 callable 形式测试通过');

