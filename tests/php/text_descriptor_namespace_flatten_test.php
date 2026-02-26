<?php

namespace tests\php;

/**
 * TextDescriptor::describeApplication 中 namespace/commands 扁平化逻辑单元测试：
 *
 * 目标：验证下面这一段在 Origami 中的行为是否与 PHP 一致：
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
 * 这里为了避免「数组下标转对象」这一层额外噪音，只保留与命令列表相关的核心逻辑：
 *
 * array_merge(
 *     ...array_values(
 *         array_map(
 *             fn ($nsCommands) => array_intersect($nsCommands, array_keys($commands)),
 *             array_values($namespacesCommands)
 *         )
 *     )
 * )
 *
 * 使用简单的 $namespacesCommands / $commands 数组，专注验证 array_map + array_values + ... 展开 +
 * array_intersect + array_merge 的组合语义，方便与 PHP CLI 直接对比。
 */

// 构造模拟的 $commands：
// - GLOBAL（_global）命名空间：命令 list
// - foo 命名空间：命令 foo:bar, foo:baz
$commands = [
    'list'    => 'list',
    'foo:bar' => 'foo:bar',
    'foo:baz' => 'foo:baz',
];

// 仅保留每个 namespace 下「命令名列表」这一层，避免额外的关联数组结构干扰
$namespacesCommands = [
    ['list'],                 // _global
    ['foo:bar', 'foo:baz'],   // foo
];

// 分步执行，便于在 Origami 中观察每一步的行为
$step1 = array_values($namespacesCommands);
Log::info('step1='.json_encode($step1));

$step2 = array_map(
    fn ($nsCommands) => array_intersect($nsCommands, array_keys($commands)),
    $step1
);
$step2 = array_map(
    fn ($nsCommands) => array_intersect($nsCommands, array_keys($commands)),
    $step1
);
Log::info('step2='.json_encode($step2));

$step3 = array_values($step2);
Log::info('step3='.json_encode($step3));

// 对应 TextDescriptor.php:197 的数组构造，去掉最外层 getColumnWidth
$flattened = array_merge(...$step3);

// 期望顺序：['list', 'foo:bar', 'foo:baz']
$expected = ['list', 'foo:bar', 'foo:baz'];

// 比较内容与顺序
if ($flattened !== $expected) {
    Log::fatal(
        'TextDescriptor namespace 扁平化逻辑测试失败，期望 '
        .'['.implode(', ', $expected).'] 实际 '
        .'['.implode(', ', $flattened).']'
    );
}

// 额外：用 TextDescriptor::getColumnWidth 的字符串分支语义，手动算一遍宽度（这里简单用 strlen 模拟 Helper::width）
// 这样方便在 PHP CLI 下对比数值是否一致。
$widths = array_map(static fn (string $name): int => \strlen($name), $flattened);
$width  = $widths ? \max($widths) + 2 : 0;

// 对我们构造的数据，最长命令是 'foo:baz'（长度 7），因此 width 应为 7 + 2 = 9。
if ($width !== 9) {
    Log::fatal('TextDescriptor getColumnWidth 简化逻辑测试失败，期望 9 实际 '.$width);
}

Log::info('TextDescriptor namespace 扁平化/宽度简化逻辑测试通过');

