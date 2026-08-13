<?php

echo "=== end() 函数测试 ===\n";

// 测试基本用法
$array = [1, 2, 3, 4, 5];
$result = end($array);
if($result == 5) {
    Log::info("基本 end 测试通过");
} else {
    Log::fatal("基本 end 测试失败");
}

// 测试空数组
$array = [];
$result = end($array);
if($result === null) {
    Log::info("空数组 end 测试通过");
} else {
    Log::fatal("空数组 end 测试失败");
}

// 测试单个元素
$array = [42];
$result = end($array);
if($result == 42) {
    Log::info("单个元素 end 测试通过");
} else {
    Log::fatal("单个元素 end 测试失败");
}

// 测试字符串数组
$array = ["a", "b", "c"];
$result = end($array);
if($result == "c") {
    Log::info("字符串数组 end 测试通过");
} else {
    Log::fatal("字符串数组 end 测试失败");
}

// 测试混合类型数组
$array = [1, "hello", 3.14, true];
$result = end($array);
if($result === true) {
    Log::info("混合类型数组 end 测试通过");
} else {
    Log::fatal("混合类型数组 end 测试失败");
}

// 测试 Iterator 对象
class EndArrayIterator implements Iterator {
    private $data;
    private int $pos;

    public function __construct($data) {
        $this->data = $data;
        $this->pos = 0;
    }

    public function current(): mixed {
        return $this->data[$this->pos];
    }

    public function key(): mixed {
        return $this->pos;
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

$nums = [10, 20, 30, 40];
$it = new EndArrayIterator($nums);
$result = end($it);
if($result == 40) {
    Log::info("Iterator end 测试通过");
} else {
    Log::fatal("Iterator end 测试失败");
}

// 测试空 Iterator
$emptyIt = new EndArrayIterator([]);
$result = end($emptyIt);
if($result === null) {
    Log::info("空 Iterator end 测试通过");
} else {
    Log::fatal("空 Iterator end 测试失败");
}

echo "=== end() 测试完成 ===\n";
