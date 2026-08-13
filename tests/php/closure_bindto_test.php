<?php

/**
 * Closure::bindTo — Laravel ClosureCommand / inspire 依赖。
 */
class ClosureBindTo_Host
{
    public string $tag = 'host';

    public function runViaBindTo(callable $cb): string
    {
    $bound = $cb->bindTo($this, static::class);
    return call_user_func($bound);
    }
}

$host = new ClosureBindTo_Host();
$result = $host->runViaBindTo(function () {
    return $this->tag;
});

if ($result !== 'host') {
    \Log::fatal('bindTo($this) 失败: ' . var_export($result, true));
}
\Log::info('closure_bindto_test OK');
