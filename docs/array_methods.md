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

### 迭代器方法

- `rewind()`: 重置迭代器到开始位置
- `current()`: 获取当前元素
- `key()`: 获取当前键（索引）
- `next()`: 移动到下一个元素
- `valid()`: 检查迭代器是否有效

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
