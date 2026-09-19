<?php

namespace tests\php;

/**
 * trait 内 use { method as alias } 后，类再 use 该 trait 时必须保留别名方法。
 * Filament Notification 的 HasIconColor::getBaseIconColor 依赖此语义。
 */

trait TraitNestedAlias_Base
{
    public function getIconColor(): string
    {
        return 'base';
    }
}

trait TraitNestedAlias_Outer
{
    use TraitNestedAlias_Base {
        getIconColor as getBaseIconColor;
    }

    public function getIconColor(): string
    {
        return $this->getBaseIconColor() . '+outer';
    }
}

class TraitNestedAlias_Host
{
    use TraitNestedAlias_Outer;
}

$o = new TraitNestedAlias_Host();
$got = $o->getIconColor();
if ($got !== 'base+outer') {
    Log::fatal('嵌套 trait 别名调用失败: ' . var_export($got, true));
}
$base = $o->getBaseIconColor();
if ($base !== 'base') {
    Log::fatal('getBaseIconColor 未合并到类: ' . var_export($base, true));
}

Log::info('嵌套 trait 方法别名测试通过');
