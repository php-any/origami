<?php

class MethodExistsAnimal {
    public function speak() {}
    protected function breathe() {}
    private function digest() {}
}

class MethodExistsDog extends MethodExistsAnimal {
    public function fetch() {}
}

$dog = new MethodExistsDog();

// 测试 1：实例方法存在
var_dump(method_exists($dog, 'fetch'));     // true
var_dump(method_exists($dog, 'speak'));    // true（继承）
var_dump(method_exists($dog, 'breathe')); // true（继承 protected）

// 测试 2：不存在的方法
var_dump(method_exists($dog, 'fly'));      // false

// 测试 3：通过类名字符串
var_dump(method_exists('MethodExistsDog', 'fetch'));   // true
var_dump(method_exists('MethodExistsDog', 'speak'));   // true
var_dump(method_exists('MethodExistsAnimal', 'speak')); // true
var_dump(method_exists('MethodExistsAnimal', 'fetch')); // false

// 测试 4：不存在的类（先用 class_exists 保护）
var_dump(class_exists('NonExistent', false) && method_exists('NonExistent', 'method')); // false

echo "method_exists 测试完成\n";
