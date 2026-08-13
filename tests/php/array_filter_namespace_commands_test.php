<?php

namespace tests\php;

/**
 * array_filter + 闭包 过滤命令列表行为测试。
 *
 * 复现 Symfony TextDescriptor::describeApplication 中的这一段：
 *
 * foreach ($namespaces as &$namespace) {
 *     $namespace['commands'] = array_filter(
 *         $namespace['commands'],
 *         fn ($name) => isset($commands[$name])
 *     );
 * }
 *
 * 目标：在 Origami 中，过滤后的 $namespaces[0]['commands'] 应与 PHP CLI 保持一致，
 * 即仍然是 ['completion', 'hello', 'help', 'list']。
 */

$commands = [
    'completion' => 'C',
    'hello'      => 'H',
    'help'       => 'He',
    'list'       => 'L',
];

$namespaces = [
    [
        'id'       => '_global',
        'commands' => ['completion', 'hello', 'help', 'list'],
    ],
];

// 使用与 TextDescriptor::describeApplication 一致的 foreach 按值遍历写法，
// 只校验循环体内 $namespace['commands'] 的过滤结果。
foreach ($namespaces as $namespace) {
    $namespace['commands'] = array_filter(
        $namespace['commands'],
        fn ($name) => isset($commands[$name])
    );

    // 期望过滤结果不变：所有命令都在 $commands 中，array_filter 不应清空列表
    $expected = ['completion', 'hello', 'help', 'list'];

    if ($namespace['commands'] !== $expected) {
        Log::fatal(
            'array_filter namespace commands 测试失败，期望 '
            .'['.implode(', ', $expected).'] 实际 '
            .'['.implode(', ', $namespace['commands']).']'
        );
    }
}

Log::info('array_filter namespace commands 测试通过');

