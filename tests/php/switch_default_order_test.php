<?php
function switch_order_check($value, $expected) {
    if ($value !== $expected) { throw new Exception('switch default order'); }
}
$seen = [];
switch (1) { case 1: $seen[] = 'one'; default: $seen[] = 'default'; case 2: $seen[] = 'two'; }
switch_order_check($seen, ['one', 'default', 'two']);
$seen = [];
switch (99) { default: $seen[] = 'default'; case 2: $seen[] = 'two'; }
switch_order_check($seen, ['default', 'two']);
$seen = [];
switch (2) { default: $seen[] = 'default'; case 2: $seen[] = 'two'; break; }
switch_order_check($seen, ['two']);
$seen = [];
switch (99) { case 1: $seen[] = 'one'; default: case 2: $seen[] = 'two'; }
switch_order_check($seen, ['two']);
echo "switch default order: PASS\n";
