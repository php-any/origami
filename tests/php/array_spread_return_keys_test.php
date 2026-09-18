<?php

namespace tests\php;

/**
 * 方法返回的关联数组展开时必须保留键名。
 */

class SpreadReturn_KeysHost
{
    public static function props(): array
    {
        return ['lazy' => true, 'columnSpan' => 2];
    }
}

$props = SpreadReturn_KeysHost::props();
if (($props['lazy'] ?? null) !== true) {
    Log::fatal('返回关联数组读键失败: ' . var_export($props, true));
}

$merged = [...SpreadReturn_KeysHost::props(), 'x' => 1];
if (($merged['lazy'] ?? null) !== true) {
    Log::fatal('方法返回数组展开丢键: ' . var_export($merged, true));
}
if (($merged['columnSpan'] ?? null) !== 2) {
    Log::fatal('方法返回数组展开丢 columnSpan: ' . var_export($merged, true));
}

Log::info('方法返回关联数组展开保留键 测试通过');
