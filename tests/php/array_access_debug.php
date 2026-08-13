<?php

namespace tests\php;

// 简化的 ArrayAccess 测试
class SimpleContainer implements ArrayAccess
{
    private array $data = [];

    public function offsetExists(mixed $offset): bool
    {
        // 调 试：打印$this->data 的内容
        Log::info("offsetExists: offset = " . (is_string($offset) ? $offset : 'non-string'));
        Log::info("offsetExists: data 内 容");
        foreach ($this->data as $key => $value) {
            Log::info("  Key: " . $key . " => " . $value);
        }
        
        // 使用 array_key_exists 代替 isset
        $result = array_key_exists($offset, $this->data);
        Log::info("offsetExists called with array_key_exists, result: " . ($result ? 'true' : 'false'));
        return $result;
    }

    public function offsetGet(mixed $offset): mixed
    {
        return $this->data[$offset] ?? null;
    }

    public function offsetSet(mixed $offset, mixed $value): void
    {
        Log::info("offsetSet called: offset=" . (is_string($offset) ? $offset : 'non-string') . ", value=" . (is_string($value) ? $value : 'non-string'));
        if (is_null($offset)) {
            $this->data[] = $value;
        } else {
            $this->data[$offset] = $value;
        }
        Log::info("offsetSet: 设置后的 data 内容");
        foreach ($this->data as $key => $val) {
            Log::info("  Key: " . $key. " => " . $val);
        }
    }

    public function offsetUnset(mixed $offset): void
    {
        unset($this->data[$offset]);
    }
}

Log::info("=== 简单测试 ===");
$container = new SimpleContainer();
$container['name'] = 'John';

Log::info("设置 name = John");
Log::info("直接访问 \$container->data: ");

// 测试直接访问数据
$data = $container['name'];
Log::info("通过 offsetGet 获取 name: " . ($data !== null ? $data : 'null'));

Log::info("检查 'name' 是否存在...");
$exists = $container->offsetExists('name');
Log::info("offsetExists 返回：" . ($exists ? 'true' : 'false'));

if ($exists === true) {
    Log::info("✓ offsetExists 测试通过");
} else {
    Log::fatal("✗ offsetExists 测试失败，期望 true");
}
