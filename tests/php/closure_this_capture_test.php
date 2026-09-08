<?php

namespace tests\php;

/**
 * 实例方法内定义的闭包应捕获 $this，延迟调用时仍可访问实例属性。
 */
class ClosureThisCapture_Hook
{
    public $component = 'comp';

    public function boot()
    {
        $self = $this;
        return function () use ($self) {
            return $self->component;
        };
    }

    public function bootThis()
    {
        return function () {
            return $this->component;
        };
    }
}

$h = new ClosureThisCapture_Hook();
$fn1 = $h->boot();
if ($fn1() !== 'comp') {
    Log::fatal('use($self) 捕获失败: '.var_export($fn1(), true));
}

$fn2 = $h->bootThis();
try {
    $got = $fn2();
    if ($got !== 'comp') {
        Log::fatal('$this 闭包捕获失败: '.var_export($got, true));
    }
} catch (\Throwable $e) {
    Log::fatal('$this 闭包调用异常: '.$e->getMessage());
}

Log::info('closure_this_capture 测试通过');
