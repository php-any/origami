<?php

namespace tests\php;

// 模拟 Illuminate\Support\enum_value + 全局 value()
function value($value, ...$args) {
    return $value instanceof \Closure ? $value(...$args) : $value;
}

function enum_value($value, $default = null) {
    return match (true) {
        $value instanceof \BackedEnum => $value->value,
        $value instanceof \UnitEnum => $value->name,
        default => $value ?? value($default),
    };
}

$r = enum_value('hello');
if ($r !== 'hello') {
    Log::fatal('期望 hello，实际: ' . var_export($r, true));
}

$r = enum_value(null, 'd');
if ($r !== 'd') {
    Log::fatal('期望 d，实际: ' . var_export($r, true));
}

Log::info('enum_value_laravel 测试通过');
