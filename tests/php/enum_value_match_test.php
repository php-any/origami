<?php

namespace tests\php;

function enum_value_probe($value, $default = null) {
    return match (true) {
        $value instanceof \BackedEnum => $value->value,
        $value instanceof \UnitEnum => $value->name,
        default => $value ?? value_probe_default($default),
    };
}

function value_probe_default($default = null) {
    return $default;
}

$r = enum_value_probe('hello');
if ($r !== 'hello') {
    Log::fatal('string 直通失败: ' . var_export($r, true));
}

$r2 = enum_value_probe(null, 'fallback');
if ($r2 !== 'fallback') {
    Log::fatal('null default 失败: ' . var_export($r2, true));
}

Log::info('enum_value_match 测试通过');
