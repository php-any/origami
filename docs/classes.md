# 类和对象

折言语言支持完整的面向对象编程，包括类、对象、继承、多态等特性。

## 基本类定义

### 简单类

```php
<?php
// 基本类定义
class Person {
    // 属性
    private string $name;
    private int $age;

    // 构造函数
    public function __construct(string $name, int $age) {
        $this->name = $name;
        $this->age = $age;
    }

    // 方法
    public function introduce(): string {
        return "I'm {$this->name}, {$this->age} years old.";
    }

    public function getAge(): int {
        return $this->age;
    }

    public function setAge(int $age): void {
        $this->age = $age;
    }
}

// 创建对象
object $person = new Person("Alice", 25);
echo $person->introduce() + "\n";
echo "Age: " + $person->getAge() + "\n";
```

### 访问修饰符

```php
<?php
class Example {
    // 公共属性：任何地方都可以访问
    public string $publicVar = "public";

    // 私有属性：只能在类内部访问
    private string $privateVar = "private";

    // 受保护属性：只能在类内部和子类中访问
    protected string $protectedVar = "protected";

    // 公共方法
    public function publicMethod(): string {
        return "This is public";
    }

    // 私有方法
    private function privateMethod(): string {
        return "This is private";
    }

    // 受保护方法
    protected function protectedMethod(): string {
        return "This is protected";
    }

    // 公共方法可以访问私有和受保护成员
    public function accessPrivate(): string {
        return $this->privateVar + " accessed from public method";
    }
}
```

## 构造函数和析构函数

### 构造函数

```php
<?php
class User {
    private string $name;
    private string $email;
    private int $createdAt;

    // 基本构造函数
    public function __construct(string $name, string $email) {
        $this->name = $name;
        $this->email = $email;
        $this->createdAt = time();
    }

    // 带默认值的构造函数
    public function __construct(string $name, string $email = "") {
        $this->name = $name;
        $this->email = $email;
        $this->createdAt = time();
    }

    public function getName(): string {
        return $this->name;
    }

    public function getEmail(): string {
        return $this->email;
    }
}

// 创建对象
object $user = new User("Alice", "alice@example.com");
echo "User: " + $user->getName() + "\n";
```

### 静态构造函数

```php
<?php
class Database {
    private static object $instance = null;
    private string $connectionString;

    private function __construct(string $connectionString) {
        $this->connectionString = $connectionString;
    }

    // 静态工厂方法
    public static function getInstance(string $connectionString): object {
        if (self::$instance === null) {
            self::$instance = new Database($connectionString);
        }
        return self::$instance;
    }

    public function getConnectionString(): string {
        return $this->connectionString;
    }
}

// 使用单例模式
object $db = Database::getInstance("mysql://localhost/db");
echo "Connection: " + $db->getConnectionString() + "\n";
```

## 继承

### 基本继承

```php
<?php
// 父类
class Animal {
    protected string $name;
    protected int $age;

    public function __construct(string $name, int $age) {
        $this->name = $name;
        $this->age = $age;
    }

    public function getName(): string {
        return $this->name;
    }

    public function getAge(): int {
        return $this->age;
    }

    // 虚方法，可以被重写
    public function makeSound(): string {
        return "Some sound";
    }
}

// 子类
class Dog extends Animal {
    private string $breed;

    public function __construct(string $name, int $age, string $breed) {
        // 调用父类构造函数
        parent::__construct($name, $age);
        $this->breed = $breed;
    }

    // 重写父类方法
    public function makeSound(): string {
        return "Woof!";
    }

    public function getBreed(): string {
        return $this->breed;
    }
}

// 使用继承
object $dog = new Dog("Buddy", 3, "Golden Retriever");
echo "Dog: " + $dog->getName() + "\n";
echo "Sound: " + $dog->makeSound() + "\n";
echo "Breed: " + $dog->getBreed() + "\n";
```

### 多重继承（通过接口）

```php
<?php
// 接口定义
interface Movable {
    public function move(): void;
}

interface Soundable {
    public function makeSound(): string;
}

// 实现多个接口
class Car implements Movable, Soundable {
    private string $model;

    public function __construct(string $model) {
        $this->model = $model;
    }

    public function move(): void {
        echo "Car {$this->model} is moving\n";
    }

    public function makeSound(): string {
        return "Vroom!";
    }

    public function getModel(): string {
        return $this->model;
    }
}

// 使用多接口实现
object $car = new Car("Tesla");
$car->move();
echo "Sound: " + $car->makeSound() + "\n";
```

## 抽象类和接口

### 抽象类

```php
<?php
// 抽象类
abstract class Shape {
    protected float $area;

    // 抽象方法，必须被子类实现
    abstract public function calculateArea(): float;

    // 具体方法
    public function getArea(): float {
        return $this->area;
    }

    public function displayArea(): void {
        echo "Area: " + $this->area + "\n";
    }
}

// 实现抽象类
class Circle extends Shape {
    private float $radius;

    public function __construct(float $radius) {
        $this->radius = $radius;
    }

    public function calculateArea(): float {
        $this->area = 3.14159 * $this->radius * $this->radius;
        return $this->area;
    }
}

class Rectangle extends Shape {
    private float $width;
    private float $height;

    public function __construct(float $width, float $height) {
        $this->width = $width;
        $this->height = $height;
    }

    public function calculateArea(): float {
        $this->area = $this->width * $this->height;
        return $this->area;
    }
}

// 使用抽象类
object $circle = new Circle(5);
$circle->calculateArea();
$circle->displayArea();

object $rectangle = new Rectangle(4, 6);
$rectangle->calculateArea();
$rectangle->displayArea();
```

### 接口

```php
<?php
// 接口定义
interface Logger {
    public function log(string $message): void;
    public function error(string $message): void;
}

interface Database {
    public function connect(): bool;
    public function query(string $sql): array;
    public function close(): void;
}

// 实现接口
class FileLogger implements Logger {
    private string $filename;

    public function __construct(string $filename) {
        $this->filename = $filename;
    }

    public function log(string $message): void {
        echo "LOG: " + $message + "\n";
    }

    public function error(string $message): void {
        echo "ERROR: " + $message + "\n";
    }
}

class MySQLDatabase implements Database {
    private string $host;
    private string $database;

    public function __construct(string $host, string $database) {
        $this->host = $host;
        $this->database = $database;
    }

    public function connect(): bool {
        echo "Connecting to MySQL at {$this->host}\n";
        return true;
    }

    public function query(string $sql): array {
        echo "Executing: " + $sql + "\n";
        return ["result" => "data"];
    }

    public function close(): void {
        echo "Closing MySQL connection\n";
    }
}

// 使用接口
object $logger = new FileLogger("app.log");
$logger->log("Application started");

object $db = new MySQLDatabase("localhost", "myapp");
$db->connect();
array $result = $db->query("SELECT * FROM users");
$db->close();
```

## 静态成员

### 静态属性和方法

```php
<?php
class MathUtils {
    // 静态属性
    public static float $PI = 3.14159;
    public static int $counter = 0;

    // 静态方法
    public static function add(int $a, int $b): int {
        return $a + $b;
    }

    public static function multiply(int $a, int $b): int {
        return $a * $b;
    }

    public static function incrementCounter(): int {
        self::$counter++;
        return self::$counter;
    }

    public static function getArea(float $radius): float {
        return self::$PI * $radius * $radius;
    }
}

// 使用静态成员
echo "PI: " + MathUtils::$PI + "\n";
echo "Sum: " + MathUtils::add(10, 20) + "\n";
echo "Product: " + MathUtils::multiply(5, 6) + "\n";
echo "Counter: " + MathUtils::incrementCounter() + "\n";
echo "Area: " + MathUtils::getArea(5) + "\n";
```

## 魔术方法

### 常用魔术方法

```php
<?php
class MagicExample {
    private array $data = [];

    // __get: 访问不存在的属性时调用
    public function __get(string $name): mixed {
        if (isset($this->data[$name])) {
            return $this->data[$name];
        }
        return null;
    }

    // __set: 设置不存在的属性时调用
    public function __set(string $name, mixed $value): void {
        $this->data[$name] = $value;
    }

    // __call: 调用不存在的方法时调用
    public function __call(string $name, array $arguments): mixed {
        echo "Calling undefined method: " + $name + "\n";
        return "Method {$name} not found";
    }

    // __toString: 对象转换为字符串时调用
    public function __toString(): string {
        return "MagicExample with " + count($this->data) + " properties";
    }
}

// 使用魔术方法
object $magic = new MagicExample();
$magic->name = "Alice";  // 调用 __set
echo "Name: " + $magic->name + "\n";  // 调用 __get
echo $magic->undefinedMethod() + "\n";  // 调用 __call
echo $magic + "\n";  // 调用 __toString
```
