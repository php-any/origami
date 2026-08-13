<?php

namespace tests\php;

class ArrayAppendObjectHolder
{
    public $items = ['named' => 'value'];

    public function push($value)
    {
        $this->items[] = $value;
    }

    public function all()
    {
        return $this->items;
    }
}

$items = ['named' => 'value'];
$items[] = 'direct';

if (count($items) === 2 && $items[0] === 'direct') {
    Log::info('关联数组 [] 追加测试通过');
} else {
    Log::fatal('关联数组 [] 追加测试失败');
}

$holder = new ArrayAppendObjectHolder();
$holder->push('property');
$propertyItems = $holder->all();

if (count($propertyItems) === 2 && $propertyItems[0] === 'property') {
    Log::info('对象数组属性 [] 追加测试通过');
} else {
    Log::fatal('对象数组属性 [] 追加测试失败');
}
