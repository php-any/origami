<?php

namespace tests\php;

/**
 * Closure::call($newThis, ...$args) — Laravel / Symfony 路由与视图会调用闭包实例方法 call。
 */
class ClosureCall_Host
{
    public string $tag = 'host';

    public function add(int $n): int
    {
        return strlen($this->tag) + $n;
    }
}

$host = new ClosureCall_Host();

$fn = function (int $n) {
    return $this->tag . ':' . $n;
};

$result = $fn->call($host, 7);
if ($result !== 'host:7') {
    Log::fatal('Closure::call($this, args) 失败: ' . var_export($result, true));
}

$noArg = function () {
    return $this->tag;
};
if ($noArg->call($host) !== 'host') {
    Log::fatal('Closure::call($this) 无参失败');
}

Log::info('closure_call 测试通过');
