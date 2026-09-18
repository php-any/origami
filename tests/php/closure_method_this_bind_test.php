<?php

namespace tests\php;

/**
 * 方法内 function(){} 闭包应绑定 $this，供 Livewire SupportRedirects::boot 的 bind 工厂使用。
 */
class RedirectThis_BindHost
{
    public $component = 'COMP';

    public function boot()
    {
        $factory = function () {
            return $this->component;
        };
        return $factory;
    }
}

$h = new RedirectThis_BindHost();
$fn = $h->boot();
$got = $fn();
if ($got !== 'COMP') {
    \Log::fatal('方法内闭包未绑定 $this, got='.var_export($got, true));
}

\Log::info('closure_method_this_bind 测试通过');
