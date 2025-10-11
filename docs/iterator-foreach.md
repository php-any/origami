# Iterator接口与foreach循环

折言语言支持自定义迭代器接口，允许开发者创建可迭代的对象，通过foreach循环进行遍历。

## Iterator接口

Iterator接口定义了迭代器必须实现的方法，使对象能够被foreach循环遍历。

### 接口定义

```php
interface Iterator {
    public function current(): mixed;
    public function key(): mixed;
    public function next(): void;
    public function rewind(): void;
    public function valid(): bool;
}
```

### 方法说明

1. **current()**: 返回当前迭代位置的元素值
2. **key()**: 返回当前迭代位置的键
3. **next()**: 将迭代器移动到下一个位置
4. **rewind()**: 将迭代器重置到初始位置
5. **valid()**: 检查当前位置是否有效

## 自定义Iterator实现

### 基本实现

以下是一个简单的数组迭代器实现：

```php
namespace tests\func;

// 实现语言级 Iterator 接口并进行 foreach 测试

class ArrayIterator implements Iterator {
    private array $data;
    private int $pos;

    public function __construct(array $data) {
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
        return $this->pos < $this->data->length;
    }
}
```

### 使用Iterator

创建Iterator实例并使用foreach循环遍历：

```php
array $nums = [1, 2, 3, 4];
object $it = new ArrayIterator($nums);

echo "-- foreach over Iterator with key --\n";
foreach ($it as $k => $v) {
    echo "k=" + $k + ", v=" + $v + "\n";
}

echo "-- foreach over Iterator value only --\n";
object $it2 = new ArrayIterator(["a", "b", "c"]);
foreach ($it2 as $v2) {
    echo $v2 + "\n";
}
```

## foreach循环

foreach循环是遍历数组和实现Iterator接口的对象的主要方式。

### 遍历数组

```php
// 索引数组
array $fruits = ["apple", "banana", "orange"];
foreach ($fruits as $fruit) {
    echo $fruit + "\n";
}

// 关联数组
array $person = [
    "name" => "Alice",
    "age" => 25,
    "city" => "Beijing"
];
foreach ($person as $key => $value) {
    echo $key + ": " + $value + "\n";
}
```

### 遍历Iterator对象

```php
// 带键值对的遍历
foreach ($iterator as $key => $value) {
    echo "Key: " + $key + ", Value: " + $value + "\n";
}

// 仅值的遍历
foreach ($iterator as $value) {
    echo "Value: " + $value + "\n";
}
```

## 高级Iterator实现

### 范围迭代器

```php
class RangeIterator implements Iterator {
    private int $start;
    private int $end;
    private int $current;

    public function __construct(int $start, int $end) {
        $this->start = $start;
        $this->end = $end;
        $this->current = $start;
    }

    public function current(): mixed {
        return $this->current;
    }

    public function key(): mixed {
        return $this->current - $this->start;
    }

    public function next(): void {
        $this->current++;
    }

    public function rewind(): void {
        $this->current = $this->start;
    }

    public function valid(): bool {
        return $this->current <= $this->end;
    }
}

// 使用范围迭代器
object $range = new RangeIterator(1, 5);
foreach ($range as $number) {
    echo $number + " ";
}
// 输出: 1 2 3 4 5
```

### 过滤迭代器

```php
class FilterIterator implements Iterator {
    private Iterator $iterator;
    private callable $filter;
    private mixed $currentValue;
    private mixed $currentKey;

    public function __construct(Iterator $iterator, callable $filter) {
        $this->iterator = $iterator;
        $this->filter = $filter;
        $this->rewind();
    }

    public function current(): mixed {
        return $this->currentValue;
    }

    public function key(): mixed {
        return $this->currentKey;
    }

    public function next(): void {
        do {
            $this->iterator->next();
            if ($this->iterator->valid()) {
                $this->currentValue = $this->iterator->current();
                $this->currentKey = $this->iterator->key();
            }
        } while ($this->iterator->valid() && !($this->filter)($this->currentValue));
    }

    public function rewind(): void {
        $this->iterator->rewind();
        if ($this->iterator->valid()) {
            $this->currentValue = $this->iterator->current();
            $this->currentKey = $this->iterator->key();
            
            // 检查第一个元素是否满足过滤条件
            if (!($this->filter)($this->currentValue)) {
                $this->next();
            }
        }
    }

    public function valid(): bool {
        return $this->iterator->valid() && ($this->filter)($this->currentValue);
    }
}

// 使用过滤迭代器
array $numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
object $arrayIterator = new ArrayIterator($numbers);

// 过滤出偶数
object $evenIterator = new FilterIterator($arrayIterator, function($value) {
    return $value % 2 == 0;
});

foreach ($evenIterator as $number) {
    echo $number + " ";
}
// 输出: 2 4 6 8 10
```

## 最佳实践

### 1. 实现所有必需方法

确保Iterator接口的所有方法都正确实现：

```php
class MyIterator implements Iterator {
    // 实现所有五个必需方法
    public function current(): mixed { /* ... */ }
    public function key(): mixed { /* ... */ }
    public function next(): void { /* ... */ }
    public function rewind(): void { /* ... */ }
    public function valid(): bool { /* ... */ }
}
```

### 2. 正确处理边界条件

```php
public function valid(): bool {
    // 确保正确检查迭代边界
    return $this->pos < $this->data->length && $this->pos >= 0;
}
```

### 3. 保持状态一致性

```php
public function rewind(): void {
    // 重置所有相关状态
    $this->pos = 0;
    // 可能还需要重置其他状态
}
```

## 注意事项

1. **性能考虑**: Iterator的每个方法都会在循环中频繁调用，确保实现高效
2. **内存管理**: 避免在Iterator中保存大量数据的副本
3. **异常处理**: 在Iterator方法中适当处理可能的异常情况
4. **类型安全**: 确保[current()](file:///Users/lvluo/Desktop/github.com/php-any/origami/data/value_array.go#L36-L38)和[key()](file:///Users/lvluo/Desktop/github.com/php-any/origami/data/value_array.go#L40-L42)方法返回正确的类型

通过实现Iterator接口，你可以创建自定义的可迭代对象，为折言语言提供更灵活和强大的数据遍历能力。