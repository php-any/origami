<?php

namespace tests\php;

/**
 * ReflectionFunction::getEndLine 应对齐 PHP（SerializableClosure / Livewire 依赖）。
 */
$fn = function () {
    return 1;
};

$r = new \ReflectionFunction($fn);
$start = $r->getStartLine();
$end = $r->getEndLine();

if ($start === false || $end === false) {
    Log::fatal('getStartLine/getEndLine 不应返回 false');
}
if (!is_int($start) || !is_int($end)) {
    Log::fatal('行号应为 int');
}
if ($end < $start) {
    Log::fatal("getEndLine($end) 应 >= getStartLine($start)");
}

Log::info("reflection_function_get_end_line 测试通过 start=$start end=$end");
