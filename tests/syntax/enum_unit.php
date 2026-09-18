<?php

namespace tests\syntax;

/**
 * PHP 8.1：无 backing 的 unit enum。
 */

enum SyntaxUnit_Color
{
    case Red;
    case Blue;
}

$c = SyntaxUnit_Color::Red;
if ($c !== SyntaxUnit_Color::Red) {
    Log::fatal('unit enum 同一性失败');
}
if ($c->name !== 'Red') {
    Log::fatal('unit enum name 失败: ' . var_export($c->name, true));
}
if (!($c instanceof \UnitEnum)) {
    Log::fatal('UnitEnum instanceof 失败');
}

Log::info('syntax unit enum 测试通过');
