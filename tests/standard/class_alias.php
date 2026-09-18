<?php

namespace tests\standard;

/**
 * class_alias。
 */

class StandardAlias_Src
{
    public function id(): string
    {
        return 'ok';
    }
}

if (class_alias(StandardAlias_Src::class, 'tests\\standard\\StandardAlias_Dst') !== true) {
    Log::fatal('class_alias 返回失败');
}
$o = new StandardAlias_Dst();
if ($o->id() !== 'ok') {
    Log::fatal('别名类实例失败');
}

Log::info('standard class_alias 测试通过');
