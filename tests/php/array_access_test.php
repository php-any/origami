<?php

namespace tests\php;

/**
 *ArrayAccess 接口测试
 * 
 *ArrayAccess 是 PHP 内置接口，允许对象像数组一样被访问
 * 接口定义：
 * interface ArrayAccess {
 *     public function offsetExists(mixed $offset): bool;
 *     public function offsetGet(mixed $offset): mixed;
 *     public function offsetSet(mixed $offset, mixed $value): void;
 *     public function offsetUnset(mixed $offset): void;
 * }
 */

// 实现 ArrayAccess 接口的类
class DataContainer implements ArrayAccess
{
    private array $data = [];

    public function __construct(array $initialData = [])
    {
        $this->data = $initialData;
    }

    // 检查偏移量是否存在
    public function offsetExists(mixed $offset): bool
    {
        return isset($this->data[$offset]);
    }

    // 获取指定偏移量的值
    public function offsetGet(mixed $offset): mixed
    {
        return $this->data[$offset] ?? null;
    }

    // 设置指定偏移量的值
    public function offsetSet(mixed $offset, mixed $value): void
    {
        if (is_null($offset)) {
            $this->data[] = $value;
        } else {
            $this->data[$offset] = $value;
        }
    }

    // 删除指定偏移量的元素
    public function offsetUnset(mixed $offset): void
    {
        unset($this->data[$offset]);
    }

    // 获取所有数据
    public function getData(): array
    {
        return $this->data;
    }
}

// 测试 1: 基本设置和获取
Log::info("=== ArrayAccess 基础测试 ===");
$container = new DataContainer();
$container['name'] = 'John';
$container['age'] = 25;
$container['city'] = 'New York';

if ($container['name'] === 'John') {
    Log::info("offsetSet/offsetGet 测试通过");
} else {
    Log::fatal("offsetSet/offsetGet 测试失败");
}

// 测试 2: offsetExists - 检查偏移量是否存在
Log::info("=== offsetExists 测试 ===");
if ($container->offsetExists('name')) {
    Log::info("offsetExists(存在) 测试通过");
} else {
    Log::fatal("offsetExists(存在) 测试失败");
}

if (!$container->offsetExists('country')) {
    Log::info("offsetExists(不存在) 测试通过");
} else {
    Log::fatal("offsetExists(不存在) 测试失败");
}

// 测试 3: offsetUnset - 删除元素
Log::info("=== offsetUnset 测试 ===");
$container['temp'] = 'temporary';
if ($container->offsetExists('temp')) {
    Log::info("临时值已设置");
}

$container->offsetUnset('temp');
if (!$container->offsetExists('temp')) {
    Log::info("offsetUnset 测试通过");
} else {
    Log::fatal("offsetUnset 测试失败");
}

// 测试 4: 数组语法糖访问（折言应该支持）
Log::info("=== 数组语法糖测试 ===");
$container['score'] = 95;
$score = $container['score'];
if ($score === 95) {
    Log::info("数组语法糖读取测试通过");
} else {
    Log::fatal("数组语法糖读取测试失败");
}

// 测试 5: 使用 isset 检测（PHP 特性）
Log::info("=== isset 检测测试 ===");
if (isset($container['age'])) {
    Log::info("isset 检测存在测试通过");
} else {
    Log::fatal("isset 检测存在测试失败");
}

if (!isset($container['salary'])) {
    Log::info("isset 检测不存在测试通过");
} else {
    Log::fatal("isset 检测不存在测试失败");
}

// 测试 6: 初始化数据测试
Log::info("=== 初始化数据测试 ===");
$initialData = ['a' => 1, 'b' => 2, 'c' => 3];
$container2 = new DataContainer($initialData);

if ($container2['a'] === 1 && $container2['b'] === 2 && $container2['c'] === 3) {
    Log::info("初始化数据测试通过");
} else {
    Log::fatal("初始化数据测试失败");
}

// 测试 7: 数字索引测试
Log::info("=== 数字索引测试 ===");
$container3 = new DataContainer();
$container3[0] = 'first';
$container3[1] = 'second';
$container3[2] = 'third';

if ($container3[0] === 'first' && $container3[1] === 'second' && $container3[2] === 'third') {
    Log::info("数字索引测试通过");
} else {
    Log::fatal("数字索引测试失败");
}

// 测试 8: 删除数字索引
Log::info("=== 删除数字索引测试 ===");
$container3->offsetUnset(1);
if (!$container3->offsetExists(1) && $container3->offsetExists(0) && $container3->offsetExists(2)) {
    Log::info("删除数字索引测试通过");
} else {
    Log::fatal("删除数字索引测试失败");
}

// 测试 9: instanceof 检查
Log::info("=== instanceof 检查 ===");
if ($container instanceof ArrayAccess) {
    Log::info("instanceof ArrayAccess 测试通过");
} else {
    Log::fatal("instanceof ArrayAccess 测试失败");
}

// 测试 10: 复杂数据类型
Log::info("=== 复杂数据类型测试 ===");
$complexData = [
    'user' => ['id' => 1, 'name' => 'Alice'],
    'settings' => ['theme' => 'dark', 'lang' => 'en'],
];
$container4 = new DataContainer($complexData);

$user = $container4['user'];
if (is_array($user) && $user['id'] === 1 && $user['name'] === 'Alice') {
    Log::info("复杂数据类型测试通过");
} else {
    Log::fatal("复杂数据类型测试失败");
}

Log::info("\n=== 所有 ArrayAccess 测试通过 ===");
