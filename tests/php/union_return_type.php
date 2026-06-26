<?php

echo "=== 联合返回类型测试 ===\n";

// 测试类方法的 string|int 返回类型
class TestIterator implements Iterator {
    private $data;
    private int $pos;

    public function __construct($data) {
        $this->data = $data;
        $this->pos = 0;
    }

    public function current(): mixed {
        return $this->data[$this->pos];
    }

    public function key(): string|int {
        return $this->pos; // 返回 int
    }

    public function next(): void {
        $this->pos = $this->pos + 1;
    }

    public function rewind(): void {
        $this->pos = 0;
    }

    public function valid(): bool {
        return $this->pos < count($this->data);
    }
}

// 测试返回 int 类型
$it1 = new TestIterator([10, 20, 30]);
$key1 = $it1->key();
if(gettype($key1) == "int" && $key1 == 0) {
    Log::info("string|int 返回类型 - int 测试通过");
} else {
    Log::fatal("string|int 返回类型 - int 测试失败");
}

// 测试返回 string 类型的 Iterator
class StringKeyIterator implements Iterator {
    private $data;
    private int $pos;

    public function __construct($data) {
        $this->data = $data;
        $this->pos = 0;
    }

    public function current(): mixed {
        return $this->data[$this->pos];
    }

    public function key(): string|int {
        $keys = ["a", "b", "c"];
        return $keys[$this->pos]; // 返回 string
    }

    public function next(): void {
        $this->pos = $this->pos + 1;
    }

    public function rewind(): void {
        $this->pos = 0;
    }

    public function valid(): bool {
        return $this->pos < count($this->data);
    }
}

// 测试返回 string 类型
$it2 = new StringKeyIterator([1, 2, 3]);
$key2 = $it2->key();
if(gettype($key2) == "string" && $key2 == "a") {
    Log::info("string|int 返回类型 - string 测试通过");
} else {
    Log::fatal("string|int 返回类型 - string 测试失败");
}

// 测试动态返回类型（根据条件返回不同类型）
class DynamicKeyIterator implements Iterator {
    private $data;
    private int $pos;

    public function __construct($data) {
        $this->data = $data;
        $this->pos = 0;
    }

    public function current(): mixed {
        return $this->data[$this->pos];
    }

    public function key(): string|int {
        if ($this->pos % 2 == 0) {
            return $this->pos; // 返回 int
        } else {
            return "key" + $this->pos; // 返回 string
        }
    }

    public function next(): void {
        $this->pos = $this->pos + 1;
    }

    public function rewind(): void {
        $this->pos = 0;
    }

    public function valid(): bool {
        return $this->pos < count($this->data);
    }
}

// 测试动态返回类型
$it3 = new DynamicKeyIterator([10, 20, 30, 40]);
$key3 = $it3->key(); // 应该是 0 (int)
if(gettype($key3) == "int" && $key3 == 0) {
    Log::info("string|int 返回类型 - 动态返回 int 测试通过");
} else {
    Log::fatal("string|int 返回类型 - 动态返回 int 测试失败");
}

$it3->next();
$key4 = $it3->key(); // 应该是 "key1" (string)
if(gettype($key4) == "string" && $key4 == "key1") {
    Log::info("string|int 返回类型 - 动态返回 string 测试通过");
} else {
    Log::fatal("string|int 返回类型 - 动态返回 string 测试失败");
}

// 测试使用 key() 函数
$it4 = new TestIterator([100, 200, 300]);
$key5 = key($it4);
if(gettype($key5) == "int" && $key5 == 0) {
    Log::info("key() 函数处理 string|int 返回类型测试通过");
} else {
    Log::fatal("key() 函数处理 string|int 返回类型测试失败");
}

// 测试其他联合返回类型
class MixedReturnClass {
    public function getValue(): string|int|bool {
        return 42; // 返回 int
    }

    public function getString(): string|int|bool {
        return "hello"; // 返回 string
    }

    public function getBool(): string|int|bool {
        return true; // 返回 bool
    }
}

$obj = new MixedReturnClass();
$val1 = $obj->getValue();
if(gettype($val1) == "int" && $val1 == 42) {
    Log::info("string|int|bool 返回类型 - int 测试通过");
} else {
    Log::fatal("string|int|bool 返回类型 - int 测试失败");
}

$val2 = $obj->getString();
if(gettype($val2) == "string" && $val2 == "hello") {
    Log::info("string|int|bool 返回类型 - string 测试通过");
} else {
    Log::fatal("string|int|bool 返回类型 - string 测试失败");
}

$val3 = $obj->getBool();
if(gettype($val3) == "bool" && $val3 === true) {
    Log::info("string|int|bool 返回类型 - bool 测试通过");
} else {
    Log::fatal("string|int|bool 返回类型 - bool 测试失败");
}

echo "=== 联合返回类型测试完成 ===\n";
