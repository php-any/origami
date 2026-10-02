<?php
function check_store($actual, $expected, $label) {
    if ($actual !== $expected) {
        throw new Exception($label . ': ' . json_encode($actual));
    }
}
$sparse = [7 => 'seven', -2 => 'minus', '' => 'empty', '007' => 'text'];
check_store(array(null => 'null', false => 'false', '007' => 'text'), ['' => 'null', 0 => 'false', '007' => 'text'], 'array syntax key conversion');
check_store(count($sparse), 4, 'literal count');
check_store(array_keys($sparse), [7, -2, '', '007'], 'literal keys');
check_store(array_key_first($sparse), 7, 'first sparse key');
check_store(array_key_last(['x' => 1, '' => 2]), '', 'last empty string key');
check_store(array_combine([7, -2, '', '007'], [1, 2, 3, 4]), [7 => 1, -2 => 2, '' => 3, '007' => 4], 'combine keys');
check_store(array_is_list(['' => 1]), false, 'empty string is not a list');
check_store(array_slice(['x' => 1, '' => 2, 7 => 3], 1), ['' => 2, 0 => 3], 'slice empty string');
check_store([...$sparse, 'last'], [0 => 'seven', 1 => 'minus', '' => 'empty', '007' => 'text', 2 => 'last'], 'spread integer reindex');
check_store(['' => 1] + [0 => 2, '' => 3], ['' => 1, 0 => 2], 'union empty string identity');
check_store(array_diff_key([7 => 1, '' => 2, 'x' => 3], ['x' => 0]), [7 => 1, '' => 2], 'diff key identity');
check_store(array_intersect_key([7 => 1, '' => 2, 'x' => 3], ['' => 0]), ['' => 2], 'intersect empty string');
$splice = ['a' => 1, 'b' => 2, 7 => 3, '' => 4];
$replacement = ['x' => 9];
$removed = array_splice($splice, 1, -1, $replacement);
check_store($removed, ['b' => 2, 0 => 3], 'splice return keys');
check_store($splice, ['a' => 1, 0 => 9, '' => 4], 'splice replacement keys');
check_store($replacement, ['x' => 9], 'splice replacement isolation');
check_store(array_reverse(['' => 1, 7 => 2, 'x' => 3]), ['x' => 3, 0 => 2, '' => 1], 'reverse key retention');
check_store(array_reverse([7 => 1, '' => 2], true), ['' => 2, 7 => 1], 'reverse preserve keys');
$a = [30, 10, 20];
$copy = $a;
asort($copy);
check_store(array_keys($copy), [1, 2, 0], 'asort integer identity');
check_store($copy[0], 30, 'asort lookup');
check_store($a, [30, 10, 20], 'asort COW');
$copy[] = 40;
check_store(array_key_last($copy), 3, 'asort append counter');
$a = [9 => 'nine', 2 => 'two', 7 => 'seven'];
$copy = $a;
ksort($copy);
check_store(array_keys($copy), [2, 7, 9], 'ksort sparse');
check_store($copy[9], 'nine', 'ksort lookup');
krsort($copy);
check_store(array_keys($copy), [9, 7, 2], 'krsort sparse');
check_store(array_keys($a), [9, 2, 7], 'key sort COW');
$a = ['z' => 30, '' => 10, 9 => 20];
$copy = $a;
sort($copy);
check_store(array_keys($copy), [0, 1, 2], 'sort reindex');
check_store($copy, [10, 20, 30], 'sort values');
check_store(array_keys($a), ['z', '', 9], 'sort COW');
rsort($copy);
check_store($copy, [30, 20, 10], 'rsort values');
$numeric = ['10', '2', '1'];
sort($numeric, SORT_NUMERIC);
check_store($numeric, ['1', '2', '10'], 'numeric string sort');
$natural = ['A10', 'a2', 'a1'];
sort($natural, SORT_NATURAL | SORT_FLAG_CASE);
check_store($natural, ['a1', 'a2', 'A10'], 'natural case flag');
$callback = ['b' => 2, 'a' => 1];
$caught = false;
try {
    usort($callback, function ($left, $right) { throw new Exception('comparison failed'); });
} catch (Exception $error) { $caught = true; }
check_store($caught, true, 'comparison exception propagation');
check_store(array_keys($callback), [0, 1], 'failed usort reindexes keys');
$stable = ['first' => 1, 'second' => 1, 'third' => 0];
uasort($stable, fn($left, $right) => $left <=> $right);
check_store(array_keys($stable), ['third', 'first', 'second'], 'stable comparison ties');
$captured = [3, 1, 2];
$seen = [];
usort($captured, function ($left, $right) use (&$captured, &$seen) {
    $seen[] = implode(',', $captured);
    return $left <=> $right;
});
foreach ($seen as $state) { check_store($state, '3,1,2', 'comparison sees original array'); }
check_store($captured, [1, 2, 3], 'comparison result publishes after sort');
$single = ['z' => 8];
sort($single);
check_store(array_keys($single), [0], 'single sort');
$a = ['a10', 'a2', 'a1'];
natcasesort($a);
check_store(array_keys($a), [2, 1, 0], 'natural sort identity');
check_store($a[0], 'a10', 'natural sort lookup');
$heap = new SplMinHeap();
foreach ([8, 1, 5, 2] as $v) { $heap->insert($v); }
$values = [];
while (!$heap->isEmpty()) { $values[] = $heap->extract(); }
check_store($values, [1, 2, 5, 8], 'heap buffer mutation');
$queue = new SplQueue();
$queue->enqueue(1);
$queue->enqueue(2);
check_store($queue->dequeue(), 1, 'queue first');
$queue->enqueue(3);
check_store($queue->dequeue(), 2, 'queue second');
check_store($queue->dequeue(), 3, 'queue append after shift');
echo "array_store_mutations OK\n";
