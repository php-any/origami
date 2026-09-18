<?php

namespace tests\php;

/**
 * PHP 8.1 BackedEnum::tryFrom / from / cases（Filament IconSize::tryFrom 依赖）。
 */
enum EnumTryFrom_Size: string
{
    case Small = 'sm';
    case Large = 'lg';

    public const AliasLarge = self::Large;
}

$from = EnumTryFrom_Size::tryFrom('sm');
if ($from !== EnumTryFrom_Size::Small) {
    Log::fatal('tryFrom 已知值失败');
}
if (EnumTryFrom_Size::tryFrom('nope') !== null) {
    Log::fatal('tryFrom 未知值应为 null');
}

$viaFrom = EnumTryFrom_Size::from('lg');
if ($viaFrom !== EnumTryFrom_Size::Large) {
    Log::fatal('from 已知值失败');
}

$cases = EnumTryFrom_Size::cases();
if (!is_array($cases) || count($cases) !== 2) {
    Log::fatal('cases 不应包含 const 别名, count=' . count($cases));
}

$caught = false;
try {
    EnumTryFrom_Size::from('nope');
} catch (\ValueError $e) {
    $caught = true;
}
if (!$caught) {
    Log::fatal('from 未知值应抛 ValueError');
}

Log::info('enum tryFrom/from/cases 测试通过');
