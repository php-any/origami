<?php
namespace tests\php;

/** Filament IconSize 模式：string backed enum + const 引用 case */
enum EnumIconSize_Size: string
{
    case ExtraLarge = 'xl';
    case TwoExtraLarge = '2xl';
    public const ExtraExtraLarge = self::TwoExtraLarge;
}

$direct = EnumIconSize_Size::TwoExtraLarge;
$viaConst = EnumIconSize_Size::ExtraExtraLarge;
if ($direct !== $viaConst || $direct->value !== '2xl') {
    Log::fatal('IconSize pattern failed');
}
Log::info('enum_icon_size_pattern 测试通过');
