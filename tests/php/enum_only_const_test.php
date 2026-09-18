<?php
namespace tests\php;
enum EnumOnlyConst_Size: string {
    case TwoExtraLarge = '2xl';
    public const ExtraExtraLarge = self::TwoExtraLarge;
}
$alias = EnumOnlyConst_Size::ExtraExtraLarge;
Log::info('ok');
