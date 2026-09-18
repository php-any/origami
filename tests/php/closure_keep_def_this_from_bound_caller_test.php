<?php

namespace tests\php;

/**
 * 方法内定义的闭包保留定义处 $this；在另一对象的 Closure::bind 上下文中
 * 普通调用不得被调用方 BoundThis 污染（Livewire EventBus / ExtendBlade）。
 */

class KeepDefThis_Host
{
    public string $name = 'host';

    public function makeListener(): \Closure
    {
        return function () {
            return $this->name;
        };
    }
}

class KeepDefThis_Component
{
    public string $name = 'component';
}

$host = new KeepDefThis_Host();
$fn = $host->makeListener();
$comp = new KeepDefThis_Component();

$fromBound = \Closure::bind(function ($fn) {
    return $fn();
}, $comp, $comp);

$got = $fromBound($fn);
if ($got !== 'host') {
    Log::fatal('调用方 BoundThis 不得覆盖监听器定义处 $this, got=' . var_export($got, true));
}

$rebound = \Closure::bind($fn, $comp, $comp);
$got2 = $rebound();
if ($got2 !== 'component') {
    Log::fatal('Closure::bind 显式重绑应生效, got=' . var_export($got2, true));
}

Log::info('闭包保留定义处 $this（调用方 Bound 不污染）测试通过');
