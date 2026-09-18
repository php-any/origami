<?php
namespace tests\php;

enum EnumSelfVs_Size: string
{
    case TwoExtraLarge = '2xl';
    public const ExtraExtraLarge = self::TwoExtraLarge;
}

// 外部 Class::Case 应可用
$ext = EnumSelfVs_Size::TwoExtraLarge;
if ($ext->value !== '2xl') {
    Log::fatal('external case access failed');
}

// const 引用 self::Case
$alias = EnumSelfVs_Size::ExtraExtraLarge;
if ($alias !== $ext) {
    Log::fatal('const alias failed');
}

Log::info('enum_self_vs_class_case 测试通过');
