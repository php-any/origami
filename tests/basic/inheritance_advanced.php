<?php

echo "=== 继承和接口高级测试 ===\n";

// 抽象类和方法
abstract class Shape {
    abstract public function area();
    public function describe() {
        return "I am a shape with area: " . $this->area();
    }
}

class Circle extends Shape {
    private $radius;
    public function __construct($radius) {
        $this->radius = $radius;
    }
    public function area() {
        return 3.14159 * $this->radius * $this->radius;
    }
}

class Rectangle extends Shape {
    private $width;
    private $height;
    public function __construct($width, $height) {
        $this->width = $width;
        $this->height = $height;
    }
    public function area() {
        return $this->width * $this->height;
    }
}

$circle = new Circle(5);
$rect = new Rectangle(4, 6);

if ($circle instanceof Shape && $rect instanceof Shape) {
    Log::info("抽象类继承测试通过");
} else {
    Log::fatal("抽象类继承测试失败");
}

// 多态
$shapes = [$circle, $rect];
$totalArea = 0;
for ($s in $shapes) {
    $totalArea = $totalArea + $s->area();
}
// 3.14159 * 25 + 24 = 78.53975 + 24 = 102.53975
if ($totalArea > 102 && $totalArea < 103) {
    Log::info("多态测试通过");
} else {
    Log::fatal("多态测试失败, 实际 {$totalArea}");
}

// 接口多实现
interface Loggable {
    public function toLog();
}

interface Storable {
    public function serialize();
}

class User implements Loggable, Storable {
    public $name;
    public function __construct($name) {
        $this->name = $name;
    }
    public function toLog() {
        return "User: " . $this->name;
    }
    public function serialize() {
        return "{\"name\":\"" . $this->name . "\"}";
    }
}

$user = new User("Alice");
if ($user instanceof Loggable && $user instanceof Storable && $user instanceof User) {
    Log::info("接口多实现测试通过");
} else {
    Log::fatal("接口多实现测试失败");
}

// final 方法
class Base {
    public function normal() {
        return "base normal";
    }
    final public function locked() {
        return "base locked";
    }
}
class Child extends Base {
    public function normal() {
        return "child normal";
    }
}

$child = new Child();
if ($child->normal() == "child normal" && $child->locked() == "base locked") {
    Log::info("final 方法测试通过");
} else {
    Log::fatal("final 方法测试失败");
}

// 接口继承接口
interface Readable {
    public function read();
}
interface Writable extends Readable {
    public function write();
}
class File implements Writable {
    public function read() {
        return "reading";
    }
    public function write() {
        return "writing";
    }
}
$file = new File();
if ($file instanceof Readable && $file instanceof Writable) {
    Log::info("接口继承测试通过");
} else {
    Log::fatal("接口继承测试失败");
}

// parent:: 调用
class BaseClass {
    public function greet() {
        return "Hello from base";
    }
}
class DerivedClass extends BaseClass {
    public function greet() {
        return parent::greet() . " + derived";
    }
}
$derived = new DerivedClass();
if ($derived->greet() == "Hello from base + derived") {
    Log::info("parent:: 调用测试通过");
} else {
    Log::fatal("parent:: 调用测试失败, 实际 '{$derived->greet()}'");
}

echo "=== 继承和接口高级测试完成 ===\n";
