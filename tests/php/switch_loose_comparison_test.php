<?php
function switch_match($value, $label) {
    switch ($value) { case $label: return true; default: return false; }
}
function switch_check($value, $label, $expected) {
    if (switch_match($value, $label) !== $expected) { throw new Exception('switch loose comparison'); }
}
switch_check('01', '1', true);
switch_check(' 1 ', 1, true);
switch_check('1e2', '100', true);
switch_check(1.2, 1.8, false);
switch_check(1, 1.8, false);
switch_check(0, 'text', false);
switch_check(true, 'text', true);
switch_check(false, '0', true);
switch_check(null, '', true);
switch_check(null, '0', false);
switch_check('9007199254740993', '9007199254740992', false);
switch_check('0x10', 16, false);
echo "switch loose comparison: PASS\n";
