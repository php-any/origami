<?php
function destructure_key_check($value, $expected) {
    if ($value !== $expected) { throw new Exception('destructure numeric keys'); }
}
[$a, $b] = [1 => 'one', 0 => 'zero'];
destructure_key_check([$a, $b], ['zero', 'one']);
[$a, [$b, $c]] = [1 => [1 => 'inner-one', 0 => 'inner-zero'], 0 => 'outer'];
destructure_key_check([$a, $b, $c], ['outer', 'inner-zero', 'inner-one']);
$warnings = [];
set_error_handler(function ($level, $message) use (&$warnings) { $warnings[] = $message; return true; });
[$a, , $b] = [0 => 'zero', 2 => 'two'];
destructure_key_check([$a, $b, $warnings], ['zero', 'two', []]);
[$a, $b] = [1 => 'one'];
destructure_key_check([$a, $b, $warnings], [null, 'one', ['Undefined array key 0']]);
restore_error_handler();
echo "destructure numeric keys: PASS\n";
