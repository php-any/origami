<?php

namespace tests\php;

/**
 * new static::$arr[$key](...) 动态类名下标（Laravel HasCollection::newCollection）。
 */

class NewStaticPropIndex_Collection
{
    public $items;

    public function __construct($items = [])
    {
        $this->items = $items;
    }
}

class NewStaticPropIndex_Model
{
    protected static array $resolvedCollectionClasses = [];
    protected static string $collectionClass = NewStaticPropIndex_Collection::class;

    public function newCollection(array $models = [])
    {
        static::$resolvedCollectionClasses[static::class] ??= static::$collectionClass;

        return new static::$resolvedCollectionClasses[static::class]($models);
    }
}

$m = new NewStaticPropIndex_Model();
$c = $m->newCollection(['a', 'b']);
if (!($c instanceof NewStaticPropIndex_Collection)) {
    \Log::fatal('new static::$arr[$key] 应得到 Collection 实例');
}
if ($c->items !== ['a', 'b']) {
    \Log::fatal('构造参数未正确传入');
}

\Log::info('new_static_prop_index_class 测试通过');
