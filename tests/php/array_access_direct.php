<?php

namespace tests\php;

/**
 *ArrayAccess 接口测试 - 直接方法调用版本
 */
class TestContainer implements ArrayAccess
{
    private array $data = [];

    public function __construct(array $initialData = [])
    {
        $this->data = $initialData;
    }

    public function offsetExists(mixed $offset): bool
    {
        return isset($this->data[$offset]);
    }

    public function offsetGet(mixed $offset): mixed
    {
        return $this->data[$offset] ?? null;
    }

    public function offsetSet(mixed $offset, mixed $value): void
    {
        if (is_null($offset)) {
            $this->data[] = $value;
        } else {
            $this->data[$offset] = $value;
        }
    }

    public function offsetUnset(mixed $offset): void
    {
        unset($this->data[$offset]);
    }
}

Log::info("=== ArrayAccess 直接方法调用测试 ===");

// 测试 1: 使用 offsetSet 设置值
Log::info("测试 1: offsetSet");
$container= new TestContainer();
$container->offsetSet('name', 'John');
Log::info("✓ offsetSet 调用成功");

// 测试 2: 使用 offsetGet 获取值
Log::info("测试 2: offsetGet");
$value = $container->offsetGet('name');
if ($value === 'John') {
    Log::info("✓ offsetGet 测试通过，值：" . $value);
} else {
    Log::fatal("✗ offsetGet 测试失败");
}

// 测试 3: 使用 offsetExists 检查存在性
Log::info("测试 3: offsetExists");
if ($container->offsetExists('name')) {
    Log::info("✓ offsetExists(存在) 测试通过");
} else {
    Log::fatal("✗ offsetExists(存在) 测试失败");
}

if (!$container->offsetExists('age')) {
    Log::info("✓ offsetExists(不存在) 测试通过");
} else {
    Log::fatal("✗ offsetExists(不存在) 测试失败");
}

// 测试 4: 使用 offsetUnset 删除元素
Log::info("测试 4: offsetUnset");
$container->offsetSet('temp', 'temporary');
if ($container->offsetExists('temp')) {
    Log::info("临时值已设置");
}

$container->offsetUnset('temp');
if (!$container->offsetExists('temp')) {
    Log::info("✓ offsetUnset 测试通过");
} else {
    Log::fatal("✗ offsetUnset 测试失败");
}

// 测试 5: instanceof 检查
Log::info("测试 5: instanceof 检查");
if ($container instanceof ArrayAccess) {
    Log::info("✓ instanceof ArrayAccess 测试通过");
} else {
    Log::fatal("✗ instanceof ArrayAccess 测试失败");
}

// 测试 6: 初始化数据
Log::info("测试 6: 初始化数据");
$container2 = new TestContainer(['a' => 1, 'b' => 2]);
if ($container2->offsetGet('a') === 1 && $container2->offsetGet('b') === 2) {
    Log::info("✓ 初始化数据测试通过");
} else {
    Log::fatal("✗ 初始化数据测试失败");
}

Log::info("\n=== 所有 ArrayAccess 直接方法调用测试通过 ===");
