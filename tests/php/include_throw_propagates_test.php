<?php

namespace tests\php;

/**
 * include/require 抛出的异常必须回到调用方 try/catch。
 * Program 若吞掉 throw，PhpEngine 会把残缺 HTML 交给 Livewire（RootTagMissing）。
 */

$caught = 'none';
try {
    require __DIR__.'/include_throw_target.php';
    $caught = 'continued';
} catch (\Throwable $e) {
    $caught = get_class($e).':'.$e->getMessage();
}

if ($caught === 'continued') {
    \Log::fatal('require 抛错后不应继续执行');
}
if ($caught === 'none') {
    \Log::fatal('require 抛错未被 try/catch 捕获');
}
if (!str_contains($caught, 'View [] not found.')) {
    \Log::fatal('require 异常消息不对: '.$caught);
}

$afterClosure = 'none';
try {
    (static function () {
        require __DIR__.'/include_throw_target.php';
    })();
    $afterClosure = 'continued';
} catch (\Throwable $e) {
    $afterClosure = 'caught:'.$e->getMessage();
}

if ($afterClosure === 'continued' || $afterClosure === 'none') {
    \Log::fatal('闭包内 require 抛错应向外抛, 实际: '.$afterClosure);
}

\Log::info('include_throw_propagates 测试通过');
