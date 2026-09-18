<?php

namespace tests\php;

/**
 * 构造传入另一 Enumerable/数组包装时，all() 应一次得到扁平 array，不递归死循环。
 * 对齐 Illuminate Collection::__construct → getArrayableItems。
 */
class CollectionNested_EnumerableLike
{
    protected $items = [];

    public function __construct($items = [])
    {
        $this->items = $this->getArrayableItems($items);
    }

    public function all()
    {
        return $this->items;
    }

    protected function getArrayableItems($items)
    {
        if (is_array($items)) {
            return $items;
        }
        if ($items instanceof self) {
            return $items->all();
        }
        if (is_object($items) && method_exists($items, 'toArray')) {
            return $items->toArray();
        }
        return (array) $items;
    }
}

$inner = new CollectionNested_EnumerableLike(['a' => 1, 'b' => 2]);
$outer = new CollectionNested_EnumerableLike($inner);
$all = $outer->all();
if (!is_array($all)) {
    \Log::fatal('嵌套 Enumerable all() 应为 array, got '.gettype($all));
}
if (($all['a'] ?? null) !== 1 || ($all['b'] ?? null) !== 2) {
    \Log::fatal('嵌套 Enumerable 未扁平化: '.var_export($all, true));
}

// 自引用不应在 all() 路径上无限递归：构造时已归一化
$selfLike = new CollectionNested_EnumerableLike([1, 2, 3]);
$again = new CollectionNested_EnumerableLike($selfLike->all());
if (count($again->all()) !== 3) {
    \Log::fatal('数组再包装失败');
}

\Log::info('collection_nested_enumerable 测试通过');
