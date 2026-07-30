<?php

$callback = strtoupper(...);
if (!($callback instanceof Closure)) {
    Log::fatal('一等函数可调用应返回 Closure');
}

$values = array_map($callback, ['a', 'b']);
if ($values !== ['A', 'B']) {
    Log::fatal('一等函数可调用执行失败');
}

Log::info('function_first_class_callable_test OK');
