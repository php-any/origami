<?php
function check_pointer($actual, $expected, $label) {
    if ($actual !== $expected) { throw new Exception($label . ': ' . json_encode($actual)); }
}
$a = [7 => 10, '' => 20, 'last' => 30];
check_pointer(key($a), 7, 'initial key');
check_pointer(current($a), 10, 'initial value');
check_pointer(next($a), 20, 'next');
check_pointer(key($a), '', 'empty string pointer');
$copy = $a;
check_pointer(next($copy), 30, 'copied pointer');
check_pointer(current($a), 20, 'pointer COW');
check_pointer(next($copy), false, 'past end');
check_pointer(key($copy), null, 'invalid key');
check_pointer(prev($copy), false, 'invalid pointer stays invalid');
check_pointer(end($copy), 30, 'end revives pointer');
check_pointer(prev($copy), 20, 'prev');
check_pointer(reset($copy), 10, 'reset');
check_pointer(prev($copy), false, 'before first');
check_pointer(next($copy), false, 'before first remains invalid');
$empty = [];
check_pointer(reset($empty), false, 'empty reset');
check_pointer(current([]), false, 'empty current');
check_pointer(key([]), null, 'empty key');
$first = [1, 2];
next($first);
unset($first[0]);
check_pointer(current($first), 2, 'unset preceding entry');
foreach ($a as $value) {}
check_pointer(current($a), 20, 'foreach does not move internal pointer');
check_pointer([7 => 1] === [0 => 1], false, 'strict comparison includes key');
check_pointer(['' => 1] === [0 => 1], false, 'strict comparison includes key type');
check_pointer(ARRAY_FILTER_USE_KEY, 2, 'filter key constant');
check_pointer(ARRAY_FILTER_USE_BOTH, 1, 'filter both constant');
check_pointer(array_filter([7 => 0, '' => 2, 9 => 3]), ['' => 2, 9 => 3], 'filter preserves keys');
check_pointer(array_filter([7 => 1, '' => 2], fn($key) => $key === '', ARRAY_FILTER_USE_KEY), ['' => 2], 'filter key type');
check_pointer(array_intersect([7 => 1, '' => 2, 9 => 3], [2, 3]), ['' => 2, 9 => 3], 'intersection keys');
check_pointer(array_flip([7 => 'x', '' => 'x', 9 => '']), ['x' => '', '' => 9], 'flip overwrite and empty key');
check_pointer(array_search(2, [7 => 1, '' => 2], true), '', 'search empty key');
$a = [7 => 1, '' => 2];
$copy = $a;
$keys = [];
array_walk($copy, function (&$value, $key) use (&$keys) { $value += 10; $keys[] = $key; });
check_pointer($keys, [7, ''], 'walk keys');
check_pointer($copy, [7 => 11, '' => 12], 'walk writes');
check_pointer($a, [7 => 1, '' => 2], 'walk COW');
echo "array_store_pointer OK\n";
