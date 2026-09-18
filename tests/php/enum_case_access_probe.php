<?php
enum IconSize: string {
    case TwoExtraLarge = '2xl';
}
echo IconSize::TwoExtraLarge->value, "\n";
echo IconSize::TwoExtraLarge->name, "\n";
