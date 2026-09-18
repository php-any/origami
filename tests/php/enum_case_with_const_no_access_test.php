<?php
namespace tests\php;

enum EnumWithConst_Size: string
{
    case TwoExtraLarge = '2xl';
    public const ExtraExtraLarge = self::TwoExtraLarge;
}

// 只访问 case，不访问 const
$case = EnumWithConst_Size::TwoExtraLarge;
if ($case->value !== '2xl') {
    Log::fatal('case failed');
}
Log::info('enum_case_with_const_no_access 测试通过');
