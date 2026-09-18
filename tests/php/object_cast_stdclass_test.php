<?php

namespace tests\php;

/**
 * (object) 数组必须得到 stdClass，供 Blade $loop / Filament rowLoop(?stdClass) 使用。
 */

class ObjectCast_StdClassHolder
{
    protected ?\stdClass $loop = null;

    public function rowLoop(?\stdClass $loop): static
    {
        $this->loop = $loop;
        return $this;
    }

    public function getLoop(): ?\stdClass
    {
        return $this->loop;
    }
}

$a = (object) ['index' => 0, 'iteration' => 1, 'first' => true];
if (!is_object($a)) {
    \Log::fatal('(object) 结果应为 object, gettype='.gettype($a));
}
if (get_class($a) !== 'stdClass') {
    \Log::fatal('(object) 应为 stdClass, 实际 '.get_class($a));
}
if ($a->iteration !== 1) {
    \Log::fatal('stdClass 属性读取失败');
}
if (!($a instanceof \stdClass)) {
    \Log::fatal('instanceof stdClass 失败');
}

$h = new ObjectCast_StdClassHolder();
$h->rowLoop($a);
if ($h->getLoop()->index !== 0) {
    \Log::fatal('?stdClass 属性赋值后读取失败');
}

// 对齐 Laravel ManagesLoops::getLastLoop
$last = [
    'iteration' => 2,
    'index' => 1,
    'first' => false,
    'last' => true,
    'parent' => null,
];
$loop = (object) $last;
$h2 = new ObjectCast_StdClassHolder();
$h2->rowLoop($loop);
if ($h2->getLoop()->iteration !== 2) {
    \Log::fatal('getLastLoop 风格 (object) 赋值失败');
}

\Log::info('object_cast_stdclass 测试通过');
