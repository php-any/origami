<?php

namespace tests\php;

/**
 * ?stdClass 属性应接受 Blade $loop（stdClass）与 (object) 转换结果。
 */

class StdClassProp_Holder
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

$h = new StdClassProp_Holder();
$obj = new \stdClass();
$obj->index = 0;
$obj->iteration = 1;
$h->rowLoop($obj);
if ($h->getLoop()->index !== 0) {
    \Log::fatal('stdClass 属性赋值失败');
}

// 模拟 Laravel getLastLoop：通常是 stdClass
$loop = (object) ['index' => 2, 'iteration' => 3, 'first' => false, 'last' => true];
$h2 = new StdClassProp_Holder();
try {
    $h2->rowLoop($loop);
    if ($h2->getLoop()->iteration !== 3) {
        \Log::fatal('(object) 数组转 stdClass 后属性读取失败');
    }
    \Log::info('stdclass_typed_property 测试通过');
} catch (\Throwable $e) {
    \Log::fatal('rowLoop 拒绝 loop: '.$e->getMessage().' type='.gettype($loop).' class='.(is_object($loop)?get_class($loop):''));
}
