<?php

// 验证 match 中对 `$value instanceof \BackedEnum => $value->value` 写法的解析与执行

enum Status: string {
    case OPEN = 'open';
    case CLOSED = 'closed';
}

function getStatusValue($value): string {
    return match (true) {
        $value instanceof \BackedEnum => $value->value,
        default => 'unknown',
    };
}

// 测试 1: 传入枚举实例，应返回其 value
$s = Status::OPEN;
echo "Status::OPEN => " . getStatusValue($s) . "\n"; // 期望: open

// 测试 2: 传入非枚举，应走 default 分支
$s2 = 123;
echo "int 123 => " . getStatusValue($s2) . "\n"; // 期望: unknown

