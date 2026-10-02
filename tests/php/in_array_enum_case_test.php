<?php

enum FilterPlacement
{
    case Dropdown;
    case Above;
    case Collapsible;
}

enum OtherPlacement
{
    case Dropdown;
}

enum BackedPlacement: string
{
    case Dropdown = 'dropdown';
    case Above = 'above';
}

function checkEnumMembership(bool $condition, string $message): void
{
    if (! $condition) {
        throw new RuntimeException($message);
    }
}

checkEnumMembership(! in_array(FilterPlacement::Dropdown, [FilterPlacement::Above, FilterPlacement::Collapsible]), 'Different unit enum cases matched');
checkEnumMembership(in_array(FilterPlacement::Dropdown, [FilterPlacement::Above, FilterPlacement::Dropdown]), 'The same unit enum case did not match');
checkEnumMembership(! in_array(FilterPlacement::Dropdown, [OtherPlacement::Dropdown]), 'Cases from different enums matched');
checkEnumMembership(! in_array(BackedPlacement::Dropdown, [BackedPlacement::Above]), 'Different backed enum cases matched');
checkEnumMembership(! in_array(BackedPlacement::Dropdown, ['dropdown']), 'A backed case matched its scalar backing value');
checkEnumMembership(! in_array('Object(FilterPlacement)', [FilterPlacement::Dropdown]), 'Debug object text matched an enum case');
checkEnumMembership(in_array(FilterPlacement::Dropdown, [FilterPlacement::Dropdown], true), 'Strict case membership failed');
checkEnumMembership(! in_array(FilterPlacement::Dropdown, [FilterPlacement::Above], true), 'Strict case membership matched a different case');
echo "enum membership: PASS\n";
