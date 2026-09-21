<?php

namespace tests\php;

/**
 * 无 use 的闭包在定义帧返回后仍可调用；实例方法内无 use 的闭包仍能用 $this。
 */

$fn = (static function () {
    return function () {
        return 41;
    };
})();
if ($fn() !== 41) {
    \Log::fatal('无 use 闭包在外层返回后调用失败: '.var_export($fn(), true));
}

class ClosureEmptyUse_Box
{
    public $n = 7;

    public function make()
    {
        return function () {
            return $this->n;
        };
    }
}

$cb = (new ClosureEmptyUse_Box())->make();
if ($cb() !== 7) {
    \Log::fatal('方法内无 use 闭包 $this 失败: '.var_export($cb(), true));
}

\Log::info('无 use 闭包生命周期测试通过');
