<?php

namespace tests\php;

/**
 * Context 池不得改写仍被 Generator 持有的局部变量（生成器跨多次普通调用仍可见）。
 */

function CallPool_gen()
{
    $n = 7;
    yield $n;
    yield $n + 1;
}

function CallPool_churn($x)
{
    return $x + 1;
}

$g = CallPool_gen();
$g->rewind();
if ($g->current() !== 7) {
    Log::fatal('生成器 rewind 后 current 应为 7: ' . var_export($g->current(), true));
}
for ($i = 0; $i < 200; $i++) {
    CallPool_churn($i);
}
if ($g->current() !== 7) {
    Log::fatal('生成器 current 被调用帧复用改写: ' . var_export($g->current(), true));
}
$g->next();
if ($g->current() !== 8) {
    Log::fatal('生成器 next 后值错误: ' . var_export($g->current(), true));
}

function CallPool_inner()
{
    return debug_backtrace(DEBUG_BACKTRACE_IGNORE_ARGS);
}

function CallPool_outer()
{
    return CallPool_inner();
}

$trace = CallPool_outer();
$names = [];
foreach ($trace as $frame) {
    if (isset($frame['function'])) {
        $names[] = $frame['function'];
    }
}
$okInner = false;
$okOuter = false;
foreach ($names as $n) {
    if ($n === 'CallPool_inner' || substr($n, -strlen('CallPool_inner')) === 'CallPool_inner') {
        $okInner = true;
    }
    if ($n === 'CallPool_outer' || substr($n, -strlen('CallPool_outer')) === 'CallPool_outer') {
        $okOuter = true;
    }
}
if (!$okInner || !$okOuter) {
    Log::fatal('debug_backtrace 缺少当前请求帧: ' . json_encode($names));
}

Log::info('调用帧池与生成器/backtrace 测试通过');
