<?php

namespace tests\php;

/**
 * 对齐 Laravel Container::build：$this->buildStack[] = spl_object_hash($closure)
 */
class SplHashAppend_Holder
{
    protected $buildStack = [];

    public function build($concrete)
    {
        if ($concrete instanceof \Closure) {
            $this->buildStack[] = spl_object_hash($concrete);
            try {
                return $concrete($this, []);
            } finally {
                array_pop($this->buildStack);
            }
        }
        return $concrete;
    }
}

$h = new SplHashAppend_Holder();
$ret = $h->build(function ($app, $params) {
    return 'ok-from-closure';
});
if ($ret !== 'ok-from-closure') {
    Log::fatal('闭包 build 返回值错误: ' . var_export($ret, true));
}

Log::info('spl_object_hash append 测试通过');
