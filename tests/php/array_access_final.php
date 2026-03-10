<?php

namespace tests\php;

class ArrayAccessTest implements ArrayAccess
{
    private array $storage = [];

    public function offsetExists(mixed $offset): bool {
        Log::info("offsetExists called");
        return isset($this->storage[$offset]);
    }

    public function offsetGet(mixed $offset): mixed {
        Log::info("offsetGet called");
        return $this->storage[$offset] ?? null;
    }

    public function offsetSet(mixed $offset, mixed $value): void {
        Log::info("offsetSet called");
        $this->storage[$offset] = $value;
    }

    public function offsetUnset(mixed $offset): void {
        Log::info("offsetUnset called");
        unset($this->storage[$offset]);
    }

    public function testThis(): void {
        Log::info("在方法内部测试 \$this[0] = 'test'");
        $this[0] = 'test';
        
        Log::info("读取 \$this[0]");
        $val = $this[0];
        Log::info("值：" . (is_string($val) ? $val : 'non-string'));
    }
}

Log::info("=== 测试外部对象数组访问 ===");
$obj = new ArrayAccessTest();

Log::info("测试：\$obj[0] = 'zero'");
$obj[0] = 'zero';

Log::info("测试：\$obj[1] = 'one'");
$obj[1] = 'one';

Log::info("测试：\$obj['key'] = 'value'");
$obj['key'] = 'value';

Log::info("测试：读取 \$obj[0]");
$val0 = $obj[0];
Log::info("结果：" . (is_string($val0) ? $val0 : 'non-string'));

Log::info("测试：读取 \$obj[1]");
$val1 = $obj[1];
Log::info("结果：" . (is_string($val1) ? $val1 : 'non-string'));

Log::info("测试：读取 \$obj['key']");
$keyVal = $obj['key'];
Log::info("结果：" . (is_string($keyVal) ? $keyVal : 'non-string'));

Log::info("\n=== 测试 \$this 数组访问 ===");
$obj->testThis();

Log::info("\n=== 所有测试完成 ===");
