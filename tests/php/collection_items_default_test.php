<?php

namespace tests\php;

/**
 * 对齐 Illuminate Collection：protected $items = [] 未显式赋值时 all() 应为 []，可展开。
 */
class CollectionItemsDefault_Mini
{
    protected $items = [];

    public function __construct($items = [])
    {
        $this->items = $items;
    }

    public function all()
    {
        return $this->items;
    }
}

$c = new CollectionItemsDefault_Mini();
$all = $c->all();
if (!is_array($all)) {
    \Log::fatal('all() 应为 array, 实际 '.gettype($all));
}
if (count($all) !== 0) {
    \Log::fatal('默认 all() 应为空数组');
}

try {
    $spread = [...$c->all()];
} catch (\Throwable $e) {
    \Log::fatal('[...$c->all()] 不应失败: '.$e->getMessage());
}
if (!is_array($spread) || count($spread) !== 0) {
    \Log::fatal('展开空 all() 失败');
}

$c2 = new CollectionItemsDefault_Mini(['a' => 1, 'b' => 2]);
$s2 = [...$c2->all()];
if (($s2['a'] ?? null) !== 1 || ($s2['b'] ?? null) !== 2) {
    \Log::fatal('展开带键 all() 失败: '.var_export($s2, true));
}

\Log::info('collection_items_default 测试通过');
