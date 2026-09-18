<?php
namespace tests\php;
enum IconSize: string {
    case TwoExtraLarge = '2xl';
    public const ExtraExtraLarge = self::TwoExtraLarge;
}
echo IconSize::TwoExtraLarge->value, "\n";
echo IconSize::ExtraExtraLarge->value, "\n";
