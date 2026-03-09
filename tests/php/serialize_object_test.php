<?php

namespace tests\php;

/**
 * serialize 对象最小语义测试：
 * - 对简单的用户类实例，输出 PHP 兼容的 O:...:...:{...} 格式
 * - 仅覆盖公共属性场景，确保类名和属性名/值按预期编码
 */

class SerializeObject_Foo
{
    public int $a = 1;
    public int $b = 2;
}

$obj = new SerializeObject_Foo();

$s = serialize($obj);

// 这里的期望值基于 PHP 官方序列化格式推导：
// O:<len>:"<FQN>":2:{s:1:"a";i:1;s:1:"b";i:2;}
// 其中类名为 tests\php\SerializeObject_Foo，长度为 29 个字节。
$expected = 'O:29:"tests\php\SerializeObject_Foo":2:{s:1:"a";i:1;s:1:"b";i:2;}';

if ($s !== $expected) {
    Log::fatal('serialize 对象失败，期望 ' . $expected . '，实际: ' . $s);
}

Log::info('serialize 对象 PHP 兼容格式测试通过，结果为: ' . $s);

