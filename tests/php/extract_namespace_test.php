<?php

namespace tests\php;

/**
 * Application::extractNamespace 辅助逻辑测试：
 * - 直接拷贝 Symfony Application::extractNamespace 的实现
 * - 覆盖无命名空间、单级、多级命名空间等场景
 * - 期望行为与 PHP 原生一致，便于对比 Origami 结果
 */

class ExtractNamespace_Test
{
    public static function extractNamespace(string $name, ?int $limit = null): string
    {
        // 与 vendor/symfony/console/Application.php:1216 保持一致
        $parts = explode(':', $name, -1);

        return implode(':', null === $limit ? $parts : \array_slice($parts, 0, $limit));
    }
}

/**
 * 简单断言辅助
 *
 * @param string|null $expected
 * @param string      $actual
 * @param string      $label
 */
function assert_extract_ns($expected, $actual, string $label): void
{
    if ($expected !== $actual) {
        Log::fatal('extractNamespace 测试失败 ['.$label.']: 期望 ['.$expected.'] 实际 ['.$actual.']');
    }
}

// 用例集合（期望值按 PHP explode 负数 limit 语义推导）
$cases = [
    // 无命名空间：explode(':', 'list', -1) => []，最后 implode 为空字符串
    ['name' => 'list',          'limit' => 1, 'expected' => ''],

    // 单级命名空间：explode(':', 'foo:bar', -1) => ['foo']，limit=1 -> 'foo'
    ['name' => 'foo:bar',       'limit' => 1, 'expected' => 'foo'],

    // 多级命名空间：explode(':', 'foo:bar:baz', -1) => ['foo', 'bar']，limit=1 -> 'foo'
    ['name' => 'foo:bar:baz',   'limit' => 1, 'expected' => 'foo'],

    // 多级命名空间 + limit=2：保留前两级
    ['name' => 'foo:bar:baz',   'limit' => 2, 'expected' => 'foo:bar'],

    // limit=null：直接用 $parts
    ['name' => 'foo',           'limit' => null, 'expected' => ''],
    ['name' => 'foo:bar',       'limit' => null, 'expected' => 'foo'],
    ['name' => 'foo:bar:baz',   'limit' => null, 'expected' => 'foo:bar'],
];

foreach ($cases as $idx => $case) {
    $name     = $case['name'];
    $limit    = $case['limit'];
    $expected = $case['expected'];

    if ($limit === null) {
        $actual = ExtractNamespace_Test::extractNamespace($name);
        $label  = sprintf('%s, limit=null', $name);
    } else {
        $actual = ExtractNamespace_Test::extractNamespace($name, $limit);
        $label  = sprintf('%s, limit=%d', $name, $limit);
    }

    assert_extract_ns($expected, $actual, $label);
}

Log::info('extractNamespace 逻辑测试通过');

