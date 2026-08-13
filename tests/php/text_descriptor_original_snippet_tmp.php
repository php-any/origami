<?php

namespace tests\php;

/**
 * 临时脚本：直接使用 TextDescriptor.php:197 原始 PHP 代码片段，
 * 在 Origami / PHP CLI 中分别运行，用于对比行为是否一致。
 *
 * 重点只还原这一行表达式的语义：
 *
 * $width = $this->getColumnWidth(
 *     array_merge(
 *         ...array_values(
 *             array_map(
 *                 fn ($namespace) => array_intersect($namespace['commands'], array_keys($commands)),
 *                 array_values($namespaces)
 *             )
 *         )
 *     )
 * );
 *
 * 这里用简单数组模拟 $namespaces / $commands，但保持原始表达式不拆分，方便你在 PHP CLI 里复制粘贴。
 */

// 模拟 TextDescriptor::describeApplication 中的 $commands / $namespaces 结构
$commands = [
    'list'    => 'list',
    'foo:bar' => 'foo:bar',
    'foo:baz' => 'foo:baz',
];

$namespaces = [
    [
        'id'       => '_global',
        'commands' => ['list'],
    ],
    [
        'id'       => 'foo',
        'commands' => ['foo:bar', 'foo:baz'],
    ],
];

// 完整保留原始 PHP 表达式（去掉 $this->getColumnWidth，只保留数组构造部分）
$flattened = array_merge(
    ...array_values(
        array_map(
            fn ($namespace) => array_intersect($namespace['commands'], array_keys($commands)),
            array_values($namespaces)
        )
    )
);

// 用 strlen 模拟 Helper::width，并用 max 计算最终宽度
$widths = array_map(static fn (string $name): int => \strlen($name), $flattened);
$width  = $widths ? \max($widths) + 2 : 0;

// 输出结果，方便在 PHP / Origami 两边直接对比
Log::info('TextDescriptor 原始片段 tmp: flattened='.json_encode($flattened));
Log::info('TextDescriptor 原始片段 tmp: width='.$width);

