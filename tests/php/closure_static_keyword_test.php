<?php

namespace tests\php;

/**
 * 类方法中定义的闭包内 static:: 应绑定到定义时所在类（含静态方法）。
 */
class ClosureStaticKeyword_Utils
{
    static function escapeStringForHtml($subject)
    {
        return 'E:'.(string) $subject;
    }

    static function viaDirectClosure($v)
    {
        $fn = function ($value) {
            return static::escapeStringForHtml($value);
        };

        return $fn($v);
    }

    static function viaArrayMap($attrs)
    {
        $out = [];
        $fn = function ($value, $key) use (&$out) {
            $out[$key] = static::escapeStringForHtml($value);
        };
        foreach ($attrs as $k => $v) {
            $fn($v, $k);
        }

        return $out;
    }
}

$direct = ClosureStaticKeyword_Utils::viaDirectClosure('<x>');
if ($direct !== 'E:<x>') {
    Log::fatal("directClosure 失败: got=".var_export($direct, true));
}

$mapped = ClosureStaticKeyword_Utils::viaArrayMap(['a' => '1', 'b' => '2']);
if (($mapped['a'] ?? null) !== 'E:1' || ($mapped['b'] ?? null) !== 'E:2') {
    Log::fatal("arrayMap 失败: got=".var_export($mapped, true));
}

Log::info('closure_static_keyword 测试通过');
