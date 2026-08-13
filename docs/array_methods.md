# 数组方法使用文档

本文档介绍了 Origami 语言中数组对象可用的所有方法，这些方法遵循 Node.js 的命名和签名风格。

## 基础操作方法

### push()

将一个或多个元素添加到数组的末尾，并返回新的数组长度。

```php
$arr = [1, 2, 3];
$length = $arr->push(4, 5); // 返回 5
echo $arr; // 输出: [1, 2, 3, 4, 5]
```

### pop()

移除并返回数组的最后一个元素，如果数组为空则返回 null。

```php
$arr = [1, 2, 3];
$last = $arr->pop(); // 返回 3
echo $arr; // 输出: [1, 2]
```

### shift()

移除并返回数组的第一个元素，如果数组为空则返回 null。

```php
$arr = [1, 2, 3];
$first = $arr->shift(); // 返回 1
echo $arr; // 输出: [2, 3]
```

### unshift()

将一个或多个元素添加到数组的开头，并返回新的数组长度。

```php
$arr = [1, 2, 3];
$length = $arr->unshift(0); // 返回 4
echo $arr; // 输出: [0, 1, 2, 3]
```

## 数组操作方法

### slice(start?, end?)

返回数组的一个浅拷贝，从 start 到 end（不包括 end）的元素组成的新数组。

```php
$arr = [1, 2, 3, 4, 5];
$slice1 = $arr->slice(1, 3); // 返回 [2, 3]
$slice2 = $arr->slice(2); // 返回 [3, 4, 5]
$slice3 = $arr->slice(-2); // 返回 [4, 5]
```

### splice(start, deleteCount?, ...items)

通过删除现有元素和/或添加新元素来更改数组的内容，返回被删除的元素数组。

```php
$arr = [1, 2, 3, 4, 5];
$deleted = $arr->splice(1, 2, 'a', 'b'); // 返回 [2, 3]
echo $arr; // 输出: [1, 'a', 'b', 4, 5]
```

### concat(...items)

合并两个或多个数组，返回一个新数组，包含所有数组的元素。

```php
$arr1 = [1, 2];
$arr2 = [3, 4];
$result = $arr1->concat($arr2, [5, 6]); // 返回 [1, 2, 3, 4, 5, 6]
```

### join(separator?)

将数组的所有元素转换为字符串并用指定的分隔符连接。

```php
$arr = ['apple', 'banana', 'orange'];
$str1 = $arr->join(); // 返回 "apple,banana,orange"
$str2 = $arr->join(' - '); // 返回 "apple - banana - orange"
```

### reverse()

反转数组中元素的顺序，并返回反转后的数组。

```php
$arr = [1, 2, 3, 4];
$reversed = $arr->reverse(); // 返回 [4, 3, 2, 1]
```

## 查找方法

### indexOf(searchElement, fromIndex?)

返回数组中第一个与指定元素相等的元素的索引，如果没找到则返回 -1。

```php
$arr = ['apple', 'banana', 'orange', 'banana'];
\\$index1 = $arr->indexOf('banana'); // 返回 1
\\$index2 = $arr->indexOf('banana', 2); // 返回 3
\\$index3 = $arr->indexOf('grape'); // 返回 -1
```

### includes(searchElement, fromIndex?)

判断数组是否包含指定的元素，如果包含则返回 true，否则返回 false。

```php
$arr = ['apple', 'banana', 'orange'];
$hasApple = $arr->includes('apple'); // 返回 true
$hasGrape = $arr->includes('grape'); // 返回 false
```

### find(callback)

返回数组中第一个满足回调函数条件的元素，如果没有找到则返回 null。

```php
$arr = [1, 2, 3, 4, 5];
$firstEven = $arr->find(function($element) {
    return $element % 2 == 0;
}); // 返回 2
```

### findIndex(callback)

返回数组中第一个满足回调函数条件的元素的索引，如果没有找到则返回 -1。

```php
$arr = [1, 2, 3, 4, 5];
$firstEvenIndex = $arr->findIndex(function($element) {
    return $element % 2 == 0;
}); // 返回 1
```

## 迭代方法

### forEach(callback)

对数组中的每个元素执行一次提供的回调函数。

```php
$arr = [1, 2, 3];
$arr->forEach(function($element, \\$index, $array) {
    echo "元素 $index: $element\n";
});
```

### map(callback)

创建一个新数组，其结果是该数组中的每个元素调用一次提供的回调函数后的返回值。

```php
$arr = [1, 2, 3];
$doubled = $arr->map(function($element) {
    return $element * 2;
}); // 返回 [2, 4, 6]
```

### filter(callback)

创建一个新数组，包含所有使回调函数返回 true 的元素。

```php
$arr = [1, 2, 3, 4, 5];
$evens = $arr->filter(function($element) {
    return $element % 2 == 0;
}); // 返回 [2, 4]
```

### reduce(callback, initialValue?)

将数组中的所有元素通过回调函数累积为单个值。

```php
$arr = [1, 2, 3, 4];
$sum = $arr->reduce(function($accumulator, $current) {
    return $accumulator + $current;
}, 0); // 返回 10

// 不提供初始值的情况
$sum2 = $arr->reduce(function($accumulator, $current) {
    return $accumulator + $current;
}); // 返回 10
```

### every(callback)

检查数组中的所有元素是否都满足回调函数的条件，如果所有元素都满足则返回 true，否则返回 false。

```php
$arr = [2, 4, 6, 8];
$allEven = $arr->every(function($element) {
    return $element % 2 == 0;
}); // 返回 true

$arr2 = [2, 4, 5, 8];
$allEven2 = $arr2->every(function($element) {
    return $element % 2 == 0;
}); // 返回 false
```

### some(callback)

检查数组中是否至少有一个元素满足回调函数的条件，如果有则返回 true，否则返回 false。

```php
$arr = [1, 3, 5, 7];
$hasEven = $arr->some(function($element) {
    return $element % 2 == 0;
}); // 返回 false

$arr2 = [1, 3, 4, 7];
$hasEven2 = $arr2->some(function($element) {
    return $element % 2 == 0;
}); // 返回 true
```

## 高级方法

### sort()

对数组元素进行排序，默认按字符串比较排序，并返回排序后的数组。

```php
$arr = [3, 1, 4, 1, 5];
$sorted = $arr->sort(); // 返回 [1, 1, 3, 4, 5]

$strArr = ['banana', 'apple', 'cherry'];
$sortedStr = $strArr->sort(); // 返回 ['apple', 'banana', 'cherry']
```

### flat(depth?)

将嵌套数组扁平化，返回一个新数组，其中所有子数组元素都被递归地连接到指定深度。

```php
$arr = [1, [2, 3], [4, [5, 6]]];
$flattened1 = $arr->flat(); // 返回 [1, 2, 3, 4, [5, 6]]
$flattened2 = $arr->flat(2); // 返回 [1, 2, 3, 4, 5, 6]
```

### flatMap(callback)

首先使用映射函数映射每个元素，然后将结果扁平化一层，返回一个新数组。

```php
$arr = [1, 2, 3];
$result = $arr->flatMap(function($element) {
    return [$element, $element * 2];
}); // 返回 [1, 2, 2, 4, 3, 6]
```

## 属性

### length

获取数组的长度。

```php
$arr = [1, 2, 3, 4, 5];
echo $arr->length; // 输出: 5
```

## 回调函数参数说明

对于需要回调函数的方法（如 `map`、`filter`、`reduce` 等），回调函数可以接收以下参数：

- **element**: 当前正在处理的数组元素
- **index**: 当前元素的索引
- **array**: 调用方法的数组

```php
$arr = [1, 2, 3];
$arr->map(function($element, $index, $array) {
    return $element + $index; // 元素值加上索引
}); // 返回 [1, 3, 5]
```

## 注意事项

1. 所有方法都遵循 Node.js 的命名和签名风格
2. 回调函数中的 `$this` 指向当前数组对象
3. 方法调用会修改原数组（如 `push`、`pop`、`shift`、`unshift`、`splice`、`reverse`、`sort`）
4. 其他方法返回新数组，不修改原数组
5. 字符串比较使用 `AsString()` 方法进行
6. 布尔值判断使用 `AsBool()` 方法进行

---

# 泛型列表 List<T> 使用文档

`List<T>` 是 Origami 语言提供的类型安全的泛型列表类，支持编译时类型检查，提供比普通数组更严格的类型约束。

## 创建 List 实例

### 基本语法

```php
// 创建指定类型的 List
$intList = new List<int>();        // 整数列表
$stringList = new List<string>();  // 字符串列表
$boolList = new List<bool>();      // 布尔值列表
```

### 类型安全

```php
$intList = new List<int>();

// ✓ 正确：添加整数
$intList->add(1);
$intList->add(2);
$intList->add(3);

// ✗ 错误：类型不匹配，会抛出异常
// $intList->add("hello");  // 类型错误
// $intList->add(3.14);     // 类型错误
```

## 基础操作方法

### add(item)

向列表末尾添加一个元素。

```php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);
echo $list->size(); // 输出: 3
```

### get(index)

根据索引获取元素，如果索引超出范围返回 null。

```php
$list = new List<string>();
$list->add("Hello");
$list->add("World");

echo $list->get(0); // 输出: Hello
echo $list->get(1); // 输出: World
echo $list->get(2); // 输出: null (索引超出范围)
```

### set(index, value)

设置指定索引位置的元素值。

```php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

$list->set(1, 10);
echo $list->get(1); // 输出: 10
```

### size()

获取列表的大小。

```php
$list = new List<string>();
$list->add("apple");
$list->add("banana");
echo $list->size(); // 输出: 2
```

### isEmpty()

检查列表是否为空。

```php
$list = new List<int>();
echo $list->isEmpty(); // 输出: true

$list->add(1);
echo $list->isEmpty(); // 输出: false
```

## 查找和移除方法

### contains(item)

检查列表是否包含指定元素。

```php
$list = new List<string>();
$list->add("apple");
$list->add("banana");

echo $list->contains("apple");  // 输出: true
echo $list->contains("orange"); // 输出: false
```

### indexOf(item)

获取指定元素在列表中的索引，如果不存在返回 -1。

```php
$list = new List<string>();
$list->add("apple");
$list->add("banana");
$list->add("orange");

echo $list->indexOf("banana"); // 输出: 1
echo $list->indexOf("grape"); // 输出: -1
```

### remove(item)

移除列表中第一个匹配的元素，返回是否成功移除。

```php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

$success = $list->remove(2);
echo $success; // 输出: true
echo $list->size(); // 输出: 2
```

### removeAt(index)

根据索引移除元素，返回是否成功移除。

```php
$list = new List<string>();
$list->add("apple");
$list->add("banana");
$list->add("orange");

$success = $list->removeAt(1);
echo $success; // 输出: true
echo $list->size(); // 输出: 2
```

### clear()

清空列表中的所有元素。

```php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

$list->clear();
echo $list->size(); // 输出: 0
echo $list->isEmpty(); // 输出: true
```

## 转换方法

### toArray()

将列表转换为普通数组。

```php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

$array = $list->toArray();
echo count($array); // 输出: 3
```

## 迭代器支持

`List<T>` 实现了 `Iterator` 接口，支持 `foreach` 循环和手动迭代。

### foreach 循环

```php
$list = new List<string>();
$list->add("Hello");
$list->add("World");
$list->add("Origami");

foreach ($list as $index => $value) {
    echo "索引 $index: $value\n";
}
```

### 手动迭代

```php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

$list->rewind();
while ($list->valid()) {
    echo "键: " . $list->key() . ", 值: " . $list->current() . "\n";
    $list->next();
}
```

### 迭代器方法详解

#### rewind()

重置迭代器到开始位置。

```php
$list = new List<int>();
$list->add(1);
$list->add(2);
$list->add(3);

$list->rewind();
echo "第一个元素: " . $list->current() . "\n"; // 输出: 1
```

#### current()

获取当前元素。

```php
$list->rewind();
$list->next(); // 移动到第二个元素
echo "当前元素: " . $list->current() . "\n"; // 输出: 2
```

#### key()

获取当前索引。

```php
$list->rewind();
$list->next();
echo "当前索引: " . $list->key() . "\n"; // 输出: 1
```

#### next()

移动到下一个元素。

```php
$list->rewind();
echo "第一个: " . $list->current() . "\n"; // 输出: 1
$list->next();
echo "第二个: " . $list->current() . "\n"; // 输出: 2
```

#### valid()

检查迭代器是否有效。

```php
$list->rewind();
while ($list->valid()) {
    echo "元素: " . $list->current() . "\n";
    $list->next();
}
```

## 类型约束示例

### 整数列表

```php
$intList = new List<int>();

// ✓ 正确操作
$intList->add(1);
$intList->add(2);
$intList->add(3);

// ✗ 类型错误
// $intList->add("hello");  // 编译时类型检查
// $intList->add(3.14);     // 编译时类型检查
```

### 字符串列表

```php
$stringList = new List<string>();

// ✓ 正确操作
$stringList->add("Hello");
$stringList->add("World");

// ✗ 类型错误
// $stringList->add(123);   // 编译时类型检查
// $stringList->add(true);  // 编译时类型检查
```

## 性能特点

1. **类型安全**: 编译时类型检查，避免运行时类型错误
2. **内存效率**: 使用 Go 的 slice 实现，动态扩容
3. **迭代器支持**: 支持 foreach 循环和手动迭代
4. **方法丰富**: 提供完整的列表操作方法

## 与普通数组的区别

| 特性     | 普通数组 | List<T> |
| -------- | -------- | ------- |
| 类型检查 | 运行时   | 编译时  |
| 类型安全 | 无       | 有      |
| 方法支持 | 丰富     | 基础    |
| 性能     | 高       | 高      |
| 内存管理 | 自动     | 自动    |

## 最佳实践

1. **必须指定类型**: `List<T>` 必须指定具体的泛型类型，不支持无类型约束
2. **类型安全优先**: 当需要类型安全时，优先使用 `List<T>` 而不是普通数组
3. **错误处理**: 注意处理类型不匹配的异常
4. **性能考虑**: 对于大量数据操作，考虑使用 `List<T>` 的类型安全优势

## 注意事项

1. **必须指定泛型类型**: `List<T>` 必须指定具体的泛型类型，不支持无类型约束的 `List()`
2. **编译时类型检查**: 泛型类型检查在编译时进行，类型不匹配会抛出异常
3. **类型安全**: 所有方法都是类型安全的，确保数据一致性
4. **迭代器状态**: 迭代器状态在方法调用间保持，注意重置迭代器

---

# 泛型哈希表 HashMap<K, V> 使用文档

`HashMap<K, V>` 是 Origami 语言提供的类型安全的泛型哈希表类，支持键值对存储，提供编译时类型检查。

## 创建 HashMap 实例

### 基本语法

```php
// 创建指定类型的 HashMap
$stringIntMap = new HashMap<string, int>();    // 字符串键，整数值
$intStringMap = new HashMap<int, string>();    // 整数键，字符串值
$stringStringMap = new HashMap<string, string>(); // 字符串键，字符串值
```

### 类型安全

```php
$map = new HashMap<string, int>();

// ✓ 正确：添加字符串键和整数值
$map->put("apple", 10);
$map->put("banana", 20);

// ✗ 错误：类型不匹配，会抛出异常
// $map->put(123, "hello");  // 键类型错误
// $map->put("grape", "world"); // 值类型错误
```

## 基础操作方法

### put(key, value)

添加或更新键值对。

```php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);
$map->put("apple", 15); // 更新现有键的值
```

### get(key)

根据键获取值，如果键不存在返回 null。

```php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);

echo $map->get("apple"); // 输出: 10
echo $map->get("grape"); // 输出: null (键不存在)
```

### remove(key)

根据键删除键值对，返回是否成功删除。

```php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);

$success = $map->remove("apple");
echo $success; // 输出: true
echo $map->size(); // 输出: 1
```

## 查询方法

### containsKey(key)

检查是否包含指定的键。

```php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);

echo $map->containsKey("apple"); // 输出: true
echo $map->containsKey("grape"); // 输出: false
```

### containsValue(value)

检查是否包含指定的值。

```php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);

echo $map->containsValue(10); // 输出: true
echo $map->containsValue(30); // 输出: false
```

### size()

获取哈希表的大小。

```php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);

echo $map->size(); // 输出: 2
```

### isEmpty()

检查哈希表是否为空。

```php
$map = new HashMap<string, int>();
echo $map->isEmpty(); // 输出: true

$map->put("apple", 10);
echo $map->isEmpty(); // 输出: false
```

## 批量操作方法

### clear()

清空哈希表中的所有键值对。

```php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);

$map->clear();
echo $map->size(); // 输出: 0
echo $map->isEmpty(); // 输出: true
```

### keys()

获取所有键的数组。

```php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);

$keys = $map->keys();
// $keys 包含 ["apple", "banana"]
```

### values()

获取所有值的数组。

```php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);

$values = $map->values();
// $values 包含 [10, 20]
```

## 迭代器支持

`HashMap<K, V>` 实现了 `Iterator` 接口，支持 `foreach` 循环和手动迭代。

### foreach 循环

```php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);
$map->put("orange", 30);

// 使用 foreach 循环
foreach ($map as $key => $value) {
    echo "$key: $value\n";
}
```

### 手动迭代

```php
$map = new HashMap<string, int>();
$map->put("apple", 10);
$map->put("banana", 20);

$map->rewind();
while ($map->valid()) {
    echo "键: " . $map->key() . ", 值: " . $map->current() . "\n";
    $map->next();
}
```

### 迭代器方法

- `rewind()`: 重置迭代器到开始位置
- `current()`: 获取当前值
- `key()`: 获取当前键
- `next()`: 移动到下一个键值对
- `valid()`: 检查迭代器是否有效

## 类型约束示例

### 字符串键，整数值

```php
$map = new HashMap<string, int>();

// ✓ 正确操作
$map->put("apple", 10);
$map->put("banana", 20);

// ✗ 类型错误
// $map->put(123, 10);     // 键类型错误
// $map->put("grape", "world"); // 值类型错误
```

### 整数键，字符串值

```php
$map = new HashMap<int, string>();

// ✓ 正确操作
$map->put(1, "one");
$map->put(2, "two");

// ✗ 类型错误
// $map->put("hello", "world"); // 键类型错误
// $map->put(3, 123);           // 值类型错误
```

## 性能特点

1. **类型安全**: 编译时类型检查，避免运行时类型错误
2. **哈希表性能**: O(1) 平均时间复杂度的键值查找
3. **迭代器支持**: 支持 foreach 循环和手动迭代
4. **内存效率**: 使用 Go 的 map 实现，动态扩容

## 与普通数组的区别

| 特性     | 普通数组      | HashMap<K, V> |
| -------- | ------------- | ------------- |
| 类型检查 | 运行时        | 编译时        |
| 键类型   | 整数索引      | 任意类型      |
| 访问方式 | 索引访问      | 键值访问      |
| 性能     | O(1) 索引访问 | O(1) 键值查找 |
| 用途     | 有序列表      | 键值映射      |

## 最佳实践

1. **必须指定类型**: `HashMap<K, V>` 必须指定键类型和值类型
2. **类型安全优先**: 当需要键值映射时，优先使用 `HashMap<K, V>`
3. **错误处理**: 注意处理类型不匹配的异常
4. **性能考虑**: 对于大量键值对操作，HashMap 比数组更高效

## 注意事项

1. **必须指定泛型类型**: `HashMap<K, V>` 必须指定键类型和值类型
2. **编译时类型检查**: 泛型类型检查在编译时进行，类型不匹配会抛出异常
3. **类型安全**: 所有方法都是类型安全的，确保数据一致性
4. **迭代器状态**: 迭代器状态在方法调用间保持，注意重置迭代器
