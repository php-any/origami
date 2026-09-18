<?php
namespace tests\php;

enum EnumCaseOnly_Size: string
{
    case TwoExtraLarge = '2xl';
}

$case = EnumCaseOnly_Size::TwoExtraLarge;
if ($case->value !== '2xl') {
    Log::fatal('case access failed');
}
Log::info('enum_case_access_only 测试通过');
