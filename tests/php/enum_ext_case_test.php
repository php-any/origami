<?php
namespace tests\php;
enum EnumExtCase_Size: string {
    case TwoExtraLarge = '2xl';
    public const ExtraExtraLarge = self::TwoExtraLarge;
}
$case = EnumExtCase_Size::TwoExtraLarge;
Log::info('case=' . $case->value);
