# 迭代器使用文档

本文档介绍了 Origami 语言中迭代器的使用方法，包括 `List<T>` 和 `HashMap<K, V>` 的迭代器功能。

## 迭代器接口

Origami 语言中的迭代器实现了标准的 `Iterator` 接口，提供统一的遍历方式。

### 迭代器方法

- `rewind()`: 重置迭代器到开始位置
- `current()`: 获取当前元素
- `key()`: 获取当前键（索引或键名）
- `next()`: 移动到下一个元素
- `valid()`: 检查迭代器是否有效

## List<T> 迭代器

### 基本用法

```php
<?php
// 创建 List<int>
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

// 手动迭代
$list->rewind();
while ($list->valid()) {
    echo "索引: " . $list->key() . ", 值: " . $list->current() . "\n";
    $list->next();
}
```

### foreach 循环支持

```php
<?php
$list = new List<string>();
$list->add("Hello");
$list->add("World");
$list->add("Origami");

// 使用 foreach 循环
foreach ($list as $index => $value) {
    echo "索引 $index: $value\n";
}
```

### 迭代器方法详解

#### rewind() - 重置迭代器

```php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

// 重置到开始位置
$list->rewind();
echo "第一个元素: " . $list->current() . "\n"; // 输出: 1
```

#### current() - 获取当前元素

```php
$list->rewind();
$list->next(); // 移动到第二个元素
echo "当前元素: " . $list->current() . "\n"; // 输出: 2
```

#### key() - 获取当前索引

```php
$list->rewind();
$list->next();
echo "当前索引: " . $list->key() . "\n"; // 输出: 1
```

#### next() - 移动到下一个元素

```php
$list->rewind();
echo "第一个: " . $list->current() . "\n"; // 输出: 1
$list->next();
echo "第二个: " . $list->current() . "\n"; // 输出: 2
```

#### valid() - 检查迭代器是否有效

```php
$list->rewind();
while ($list->valid()) {
    echo "元素: " . $list->current() . "\n";
    $list->next();
}
```

## HashMap<K, V> 迭代器

### 基本用法

```php
<?php
// 创建 HashMap<string, int>
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);
$map->put("orange", 30);

// 手动迭代
$map->rewind();
while ($map->valid()) {
    echo "键: " . $map->key() . ", 值: " . $map->current() . "\n";
    $map->next();
}
```

### foreach 循环支持

```php
<?php
$map = new HashMap<string, string>();
$map->put("name", "Origami");
$map->put("version", "1.0");
$map->put("type", "language");

// 使用 foreach 循环
foreach ($map as $key => $value) {
    echo "$key: $value\n";
}
```

### 迭代器方法详解

#### rewind() - 重置迭代器

```php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);

// 重置到开始位置
$map->rewind();
echo "第一个键: " . $map->key() . "\n"; // 输出: apple
echo "第一个值: " . $map->current() . "\n"; // 输出: 10
```

#### current() - 获取当前值

```php
$map->rewind();
$map->next(); // 移动到下一个键值对
echo "当前值: " . $map->current() . "\n"; // 输出: 20
```

#### key() - 获取当前键

```php
$map->rewind();
$map->next();
echo "当前键: " . $map->key() . "\n"; // 输出: banana
```

#### next() - 移动到下一个键值对

```php
$map->rewind();
echo "第一个键: " . $map->key() . "\n"; // 输出: apple
$map->next();
echo "第二个键: " . $map->key() . "\n"; // 输出: banana
```

#### valid() - 检查迭代器是否有效

```php
$map->rewind();
while ($map->valid()) {
    echo "键: " . $map->key() . ", 值: " . $map->current() . "\n";
    $map->next();
}
```

## 集合操作方法

### List<T> 操作方法

#### 添加元素

```php
<?php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

// 在指定位置插入元素
$list->set(1, 10); // 将索引1的元素设置为10
```

#### 删除元素

```php
<?php
$list = new List<string>();
$list->add("apple");
$list->add("banana");
$list->add("orange");

// 按值删除元素
$success = $list->remove("banana");
echo "删除结果: " . ($success ? "成功" : "失败") . "\n";

// 按索引删除元素
$removed = $list->removeAt(0);
echo "删除的元素: " . $removed . "\n";

// 清空所有元素
$list->clear();
echo "清空后大小: " . $list->size() . "\n";
```

#### 查找元素

```php
<?php
$list = new List<string>();
$list->add("apple");
$list->add("banana");
$list->add("orange");

// 检查是否包含元素
$contains = $list->contains("banana");
echo "是否包含 banana: " . ($contains ? "是" : "否") . "\n";

// 获取元素索引
$index = $list->indexOf("orange");
echo "orange 的索引: " . $index . "\n";

// 获取元素
$value = $list->get(1);
echo "索引1的元素: " . $value . "\n";
```

#### 集合信息

```php
<?php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

// 获取集合大小
echo "集合大小: " . $list->size() . "\n";

// 检查是否为空
echo "是否为空: " . ($list->isEmpty() ? "是" : "否") . "\n";

// 转换为数组
$array = $list->toArray();
echo "数组长度: " . count($array) . "\n";
```

### HashMap<K, V> 操作方法

#### 添加和更新键值对

```php
<?php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);
$map->put("orange", 30);

// 更新现有键的值
$map->put("apple", 15);
echo "更新后 apple 的值: " . $map->get("apple") . "\n";
```

#### 删除键值对

```php
<?php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);
$map->put("orange", 30);

// 删除指定键
$success = $map->remove("banana");
echo "删除 banana: " . ($success ? "成功" : "失败") . "\n";

// 清空所有键值对
$map->clear();
echo "清空后大小: " . $map->size() . "\n";
```

#### 查找键值对

```php
<?php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);
$map->put("orange", 30);

// 检查是否包含键
$hasKey = $map->containsKey("apple");
echo "是否包含键 apple: " . ($hasKey ? "是" : "否") . "\n";

// 检查是否包含值
$hasValue = $map->containsValue(20);
echo "是否包含值 20: " . ($hasValue ? "是" : "否") . "\n";

// 获取值
$value = $map->get("banana");
echo "banana 的值: " . $value . "\n";
```

#### 获取所有键和值

```php
<?php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);
$map->put("orange", 30);

// 获取所有键
$keys = $map->keys();
echo "所有键: " . implode(", ", $keys) . "\n";

// 获取所有值
$values = $map->values();
echo "所有值: " . implode(", ", $values) . "\n";
```

#### 集合信息

```php
<?php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);

// 获取集合大小
echo "集合大小: " . $map->size() . "\n";

// 检查是否为空
echo "是否为空: " . ($map->isEmpty() ? "是" : "否") . "\n";
```

## 高级迭代器用法

### 条件迭代

```php
<?php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);
$list->add(4);
$list->add(5);

// 只处理偶数
$list->rewind();
while ($list->valid()) {
    $value = $list->current();
    if ($value % 2 == 0) {
        echo "偶数: $value\n";
    }
    $list->next();
}
```

### 提前退出

```php
<?php
$list = new List<string>();
$list->add("apple");
$list->add("banana");
$list->add("orange");
$list->add("grape");

// 找到特定元素后退出
$list->rewind();
while ($list->valid()) {
    $value = $list->current();
    echo "检查: $value\n";
    if ($value == "orange") {
        echo "找到目标元素!\n";
        break;
    }
    $list->next();
}
```

### 反向遍历

```php
<?php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

// 反向遍历（需要先获取所有元素）
$values = $list->toArray();
for ($i = count($values) - 1; $i >= 0; $i--) {
    echo "反向: " . $values[$i] . "\n";
}
```

## 迭代器状态管理

### 状态重置

```php
<?php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

// 第一次遍历
$list->rewind();
while ($list->valid()) {
    echo "第一次: " . $list->current() . "\n";
    $list->next();
}

// 重置后再次遍历
$list->rewind();
while ($list->valid()) {
    echo "第二次: " . $list->current() . "\n";
    $list->next();
}
```

### 嵌套迭代

```php
<?php
$list1 = new List<int>();
$list1->add(1);
$list1->add(2);

$list2 = new List<string>();
$list2->add("a");
$list2->add("b");

// 嵌套迭代
$list1->rewind();
while ($list1->valid()) {
    $num = $list1->current();

    $list2->rewind();
    while ($list2->valid()) {
        $str = $list2->current();
        echo "$num$str\n";
        $list2->next();
    }

    $list1->next();
}
```

## 性能考虑

### 迭代器 vs 数组访问

```php
<?php
$list = new List<int>();
// ... 添加大量元素

// 推荐：使用迭代器
$list->rewind();
while ($list->valid()) {
    $value = $list->current();
    // 处理元素
    $list->next();
}

// 不推荐：频繁调用 get()
for ($i = 0; $i < $list->size(); $i++) {
    $value = $list->get($i); // 每次调用都有开销
    // 处理元素
}
```

### 内存效率

```php
<?php
// 对于大型集合，迭代器更节省内存
$largeList = new List<string>();
// ... 添加大量元素

// 使用迭代器，一次只处理一个元素
$largeList->rewind();
while ($largeList->valid()) {
    $value = $largeList->current();
    // 处理当前元素
    $largeList->next();
}
```

## 实用操作示例

### 批量删除元素

```php
<?php
$list = new List<string>();
$list->add("apple");
$list->add("banana");
$list->add("orange");
$list->add("grape");
$list->add("kiwi");

// 删除所有包含 'a' 的元素
$toRemove = [];
$list->rewind();
while ($list->valid()) {
    $value = $list->current();
    if (strpos($value, 'a') !== false) {
        $toRemove[] = $value;
    }
    $list->next();
}

// 删除收集到的元素
foreach ($toRemove as $item) {
    $list->remove($item);
}

echo "删除后剩余元素: " . $list->size() . "\n";
```

### 条件更新元素

```php
<?php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);
$map->put("orange", 30);

// 将所有值大于15的键值对的值翻倍
$map->rewind();
while ($map->valid()) {
    $key = $map->key();
    $value = $map->current();
    if ($value > 15) {
        $map->put($key, $value * 2);
    }
    $map->next();
}

// 显示更新后的结果
$map->rewind();
while ($map->valid()) {
    echo $map->key() . ": " . $map->current() . "\n";
    $map->next();
}
```

### 查找和替换

```php
<?php
$list = new List<string>();
$list->add("Hello");
$list->add("World");
$list->add("Hello");
$list->add("Origami");

// 查找并替换所有 "Hello" 为 "Hi"
for ($i = 0; $i < $list->size(); $i++) {
    $value = $list->get($i);
    if ($value == "Hello") {
        $list->set($i, "Hi");
    }
}

// 显示结果
$list->rewind();
while ($list->valid()) {
    echo $list->current() . " ";
    $list->next();
}
echo "\n";
```

### 集合合并

```php
<?php
$list1 = new List<int>();
$list1->add(1);
$list1->add(2);

$list2 = new List<int>();
$list2->add(3);
$list2->add(4);

// 合并两个列表
$list2->rewind();
while ($list2->valid()) {
    $list1->add($list2->current());
    $list2->next();
}

echo "合并后大小: " . $list1->size() . "\n";
```

### 键值对转换

```php
<?php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);
$map->put("orange", 30);

// 创建反向映射 (值 -> 键)
$reverseMap = new HashMap<int, string>();
$map->rewind();
while ($map->valid()) {
    $key = $map->key();
    $value = $map->current();
    $reverseMap->put($value, $key);
    $map->next();
}

// 显示反向映射
$reverseMap->rewind();
while ($reverseMap->valid()) {
    echo $reverseMap->current() . " -> " . $reverseMap->key() . "\n";
    $reverseMap->next();
}
```

## 错误处理和异常

### 类型安全错误

```php
<?php
try {
    $list = new List<int>();
    $list->add("hello"); // 类型错误：期望 int 但得到 string
} catch (Exception $e) {
    echo "类型错误: " . $e->getMessage() . "\n";
}
```

### 索引越界处理

```php
<?php
$list = new List<string>();
$list->add("apple");
$list->add("banana");

// 安全获取元素
$index = 5; // 超出范围
if ($index >= 0 && $index < $list->size()) {
    $value = $list->get($index);
    echo "元素: " . $value . "\n";
} else {
    echo "索引超出范围\n";
}
```

### 空集合处理

```php
<?php
$list = new List<int>();

// 检查集合是否为空
if (!$list->isEmpty()) {
    $list->rewind();
    while ($list->valid()) {
        echo $list->current() . "\n";
        $list->next();
    }
} else {
    echo "集合为空\n";
}
```

## 常见错误和解决方案

### 忘记重置迭代器

```php
<?php
$list = new List<int>();
$list->add(1);
$list->add(2);

// 错误：没有重置迭代器
while ($list->valid()) { // 可能不会执行
    echo $list->current() . "\n";
    $list->next();
}

// 正确：先重置迭代器
$list->rewind();
while ($list->valid()) {
    echo $list->current() . "\n";
    $list->next();
}
```

### 在迭代过程中修改集合

```php
<?php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

// 错误：在迭代过程中修改集合
$list->rewind();
while ($list->valid()) {
    $value = $list->current();
    if ($value == 2) {
        $list->remove(2); // 这可能导致迭代器状态不一致
    }
    $list->next();
}

// 正确：先收集要删除的元素，再删除
$toRemove = [];
$list->rewind();
while ($list->valid()) {
    $value = $list->current();
    if ($value == 2) {
        $toRemove[] = $value;
    }
    $list->next();
}

// 然后删除
foreach ($toRemove as $value) {
    $list->remove($value);
}
```

## 最佳实践

1. **总是重置迭代器**: 在开始迭代前调用 `rewind()`
2. **使用 foreach 循环**: 对于简单遍历，优先使用 `foreach`
3. **避免在迭代中修改**: 不要在迭代过程中修改集合
4. **及时释放资源**: 迭代完成后不需要特殊清理
5. **考虑性能**: 对于大型集合，迭代器比索引访问更高效

## 注意事项

1. **迭代器状态**: 迭代器状态在方法调用间保持
2. **并发安全**: 迭代器不是线程安全的
3. **内存使用**: 迭代器本身占用很少内存
4. **性能**: 迭代器访问比随机访问更高效
5. **类型安全**: 迭代器返回的值保持原始类型
