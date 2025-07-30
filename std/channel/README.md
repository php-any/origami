# Channel 模块

Channel 模块提供了对 Go 语言 `chan` 的封装，用于解决异步通讯问题。它允许在 Origami 脚本中进行并发通信，支持阻塞和非阻塞操作，并支持动态创建不同容量的 channel。

## 功能特性

- **真正的 Go Channel**: 底层使用 Go 的 `chan` 类型实现
- **异步通讯**: 支持生产者-消费者模式
- **阻塞/非阻塞操作**: 提供阻塞和非阻塞的发送方法，接收对齐 Go 的用法
- **动态容量**: 支持构造函数传参来指定 channel 容量
- **线程安全**: 使用 Go 的 `sync.Mutex` 确保线程安全
- **缓冲区支持**: 支持自定义缓冲容量的 channel
- **关闭状态管理**: 支持关闭 channel 并检查关闭状态

## 类定义

### Channel 类

```php
class Channel {
    public function __construct(int $capacity = 0)   // 构造函数，支持容量参数
    public function send($value): bool               // 非阻塞发送
    public function sendBlocking($value): bool       // 阻塞发送
    public function receive(): mixed                 // 接收数据（对齐 Go 的用法）
    public function close(): void                    // 关闭 channel
    public function isClosed(): bool                 // 检查是否关闭
    public function len(): int                       // 获取缓冲区长度
    public function cap(): int                       // 获取缓冲区容量
}
```

## 构造函数

### \_\_construct(int $capacity = 0)

创建 Channel 实例，支持指定容量。

**参数:**

- `$capacity`: channel 缓冲区容量，默认为 0（无缓冲）

**示例:**

```php
// 使用默认容量（0，无缓冲）
$channel1 = new Channel();

// 指定容量为 50
$channel2 = new Channel(50);

// 指定小容量用于测试
$channel3 = new Channel(5);
```

## 方法说明

### send($value): bool

非阻塞发送数据到 channel。

**参数:**

- `$value`: 要发送的数据

**返回值:**

- `bool`: 发送是否成功（如果缓冲区满则返回 false）

**示例:**

```php
$channel = new Channel(2); // 小容量
$result1 = $channel->send("数据1"); // 成功
$result2 = $channel->send("数据2"); // 成功
$result3 = $channel->send("数据3"); // 失败，缓冲区满
```

### sendBlocking($value): bool

阻塞发送数据到 channel。

**参数:**

- `$value`: 要发送的数据

**返回值:**

- `bool`: 发送是否成功（如果 channel 已关闭则返回 false）

**示例:**

```php
$channel = new Channel();
$result = $channel->sendBlocking("重要数据");
echo $result ? "发送成功" : "Channel 已关闭";
```

### receive(): mixed

接收数据（对齐 Go 的用法）。

**返回值:**

- `mixed`: 接收到的数据，如果 channel 已关闭且无数据则返回 null

**示例:**

```php
$value = $channel->receive();
if ($value !== null) {
    echo "接收到: " . $value;
} else {
    echo "没有数据可接收或 Channel 已关闭";
}
```

### close(): void

关闭 channel。

**示例:**

```php
$channel->close();
```

### isClosed(): bool

检查 channel 是否已关闭。

**返回值:**

- `bool`: 是否已关闭

**示例:**

```php
if ($channel->isClosed()) {
    echo "Channel 已关闭";
}
```

### len(): int

获取 channel 缓冲区中的数据数量。

**返回值:**

- `int`: 缓冲区中的数据数量

**示例:**

```php
echo "缓冲区中有 " . $channel->len() . " 个数据";
```

### cap(): int

获取 channel 缓冲区的容量。

**返回值:**

- `int`: 缓冲区容量

**示例:**

```php
echo "Channel 容量: " . $channel->cap();
```

## 使用示例

### 基本异步通讯

```php
<?php
// 创建默认容量的 channel（无缓冲）
$channel = new Channel();

// 非阻塞发送
$channel->send("数据1");
$channel->send("数据2");

// 接收数据（对齐 Go 的用法）
$data1 = $channel->receive(); // "数据1"
$data2 = $channel->receive(); // "数据2"

$channel->close();
```

### 指定容量的 Channel

```php
<?php
// 创建小容量的 channel
$channel = new Channel(3);

// 发送数据
for ($i = 1; $i <= 5; $i++) {
    $result = $channel->send("数据" . $i);
    echo "发送数据{$i}: " . ($result ? "成功" : "缓冲区满") . "\n";
}

// 接收数据
while ($channel->len() > 0) {
    $data = $channel->receive();
    echo "接收: " . $data . "\n";
}

$channel->close();
```

### 生产者-消费者模式

```php
<?php
// 创建指定容量的 channel
$channel = new Channel(10);

// 生产者
for ($i = 1; $i <= 5; $i++) {
    $channel->send("产品" . $i);
}

// 消费者
while ($channel->len() > 0) {
    $product = $channel->receive();
    if ($product !== null) {
        echo "消费: " . $product . "\n";
    }
}

$channel->close();
```

### 阻塞操作示例

```php
<?php
// 创建 channel
$channel = new Channel();

// 阻塞发送（会等待接收者）
$channel->sendBlocking("重要数据");

// 接收数据（Go 风格的接收）
$data = $channel->receive();
echo "接收到: " . $data;

$channel->close();
```

## 异步通讯特性

### 1. 真正的 Go Channel

- 底层使用 Go 的 `chan data.Value` 实现
- 支持 Go 的 channel 语义
- 线程安全，可在多个 goroutine 中使用

### 2. 动态容量支持

- 支持构造函数传参指定容量
- 默认容量为 0（无缓冲）
- 支持任意正数容量
- 容量为负数时使用默认容量 0

### 3. 对齐 Go 的用法

- **发送**: `send()`（非阻塞）和 `sendBlocking()`（阻塞）
- **接收**: `receive()`（对齐 Go 的 `<-chan` 用法）
- Go 的 channel 接收只有一种方式，我们保持一致

### 4. 缓冲区管理

- 支持自定义缓冲区容量
- 可以通过 `len()` 和 `cap()` 监控缓冲区状态
- 缓冲区满时，非阻塞发送会失败
- 无缓冲 channel 需要发送者和接收者同时准备好

### 5. 关闭语义

- 关闭后无法再发送数据
- 关闭后可以继续接收剩余数据
- 关闭后接收会返回 null

## 注意事项

1. **线程安全**: Channel 是线程安全的，可以在多个 goroutine 中安全使用
2. **阻塞行为**: `sendBlocking()` 会阻塞等待，`receive()` 对齐 Go 的用法
3. **容量限制**: 容量影响缓冲能力，小容量更容易测试缓冲区满的情况
4. **关闭后操作**: 向已关闭的 channel 发送数据会返回 false
5. **异步特性**: 真正的异步通讯，支持并发操作
6. **构造函数**: 支持传参指定容量，不传参时使用默认容量 0（无缓冲）
7. **无缓冲 channel**: 默认无缓冲，需要发送者和接收者同时准备好
8. **对齐 Go**: 接收方式对齐 Go 的 `<-chan` 用法，只有一种接收方式

## 实现细节

Channel 模块在底层使用 Go 的 channel 机制实现，提供了以下特性：

- 使用 `chan data.Value` 作为底层实现
- 使用 `sync.Mutex` 确保线程安全
- 支持动态创建不同容量的 channel
- 实现了 Go channel 的完整语义
- 提供了友好的 PHP 风格 API
- 支持构造函数传参来指定容量
- 接收方式对齐 Go 的用法
