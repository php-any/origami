<?php

namespace tests\php;

/**
 * 子类 isset($this->parentProp) / 读取父类声明属性不得 panic，应对齐 PHP。
 */

class ParentPropIsset_Parent
{
    public $x = 7;
}

class ParentPropIsset_Child extends ParentPropIsset_Parent
{
    public function check()
    {
        return isset($this->x) && $this->x === 7;
    }

    public function issetMissing()
    {
        return isset($this->noSuchProp);
    }
}

$c = new ParentPropIsset_Child();
if (!$c->check()) {
    Log::fatal('父类属性 isset/读取 失败');
}
if ($c->issetMissing()) {
    Log::fatal('不存在的属性 isset 应为 false');
}

Log::info('parent_property_isset 测试通过');
