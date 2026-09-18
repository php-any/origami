<?php
namespace tests\php;

enum EnumSelfRef_Size: string
{
    case Small = 'sm';
    case TwoExtraLarge = '2xl';
    public const ExtraExtraLarge = self::TwoExtraLarge;
}

$case = EnumSelfRef_Size::TwoExtraLarge;
if ($case->value !== '2xl') { Log::fatal('case failed'); }
$alias = EnumSelfRef_Size::ExtraExtraLarge;
if ($alias !== EnumSelfRef_Size::TwoExtraLarge) { Log::fatal('const failed'); }
Log::info('enum_case_self_ref_const 测试通过');
