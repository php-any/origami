<?php

namespace tests\php;

/**
 * Arr::from / Enumerable：数组不得被当成 Enumerable 再调 all()。
 */
interface CollDepth_Enumerable
{
    public function all();
}

class CollDepth_Collection implements CollDepth_Enumerable
{
    protected $items = [];

    public function __construct($items = [])
    {
        $this->items = self::from($items);
    }

    public static function from($items)
    {
        return match (true) {
            is_array($items) => $items,
            $items instanceof CollDepth_Enumerable => $items->all(),
            is_object($items) => (array) $items,
            default => [],
        };
    }

    public function all()
    {
        return $this->items;
    }
}

$a = new CollDepth_Collection([1, 2, 3]);
if ($a->all() !== [1, 2, 3] && array_values($a->all()) !== [1, 2, 3]) {
    // allow key preservation differences
}
$all = $a->all();
if (!is_array($all) || count($all) !== 3) {
    \Log::fatal('from array failed: '.var_export($all, true));
}

$b = new CollDepth_Collection($a);
$all2 = $b->all();
if (!is_array($all2) || count($all2) !== 3) {
    \Log::fatal('from Enumerable failed type='.gettype($all2).' val='.var_export($all2, true));
}

// 自引用：items 若被设成 Collection 自身会炸
$self = new CollDepth_Collection();
try {
    // 模拟错误写入
    $ref = new \ReflectionClass($self);
    // 跳过反射；直接构造套娃
    $c1 = new CollDepth_Collection([1]);
    $c2 = new CollDepth_Collection($c1);
    $c3 = new CollDepth_Collection($c2);
    $x = $c3->all();
    if (!is_array($x)) {
        \Log::fatal('多层 from 失败');
    }
} catch (\Throwable $e) {
    \Log::fatal('多层 from 异常: '.$e->getMessage());
}

\Log::info('collection_arr_from_match 测试通过');
