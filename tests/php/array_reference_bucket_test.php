<?php
function bucket_assert($ok, $message) {
    if (!$ok) { throw new Exception($message); }
}
$value = 7;
$original = ['left' => &$value, 'right' => 2];
$copy = $original;
sort($copy);
bucket_assert(array_keys($original) === ['left', 'right'], 'sort changed source keys');
$copy[1] = 19;
bucket_assert($value === 19 && $original['left'] === 19, 'sort lost reference value');
$value = 23;
bucket_assert($copy[1] === 23, 'source alias lost after renumber');
$shifted = $original;
array_unshift($shifted, 1);
bucket_assert(array_keys($original) === ['left', 'right'], 'unshift changed source keys');
$numeric = [8 => &$value, 20 => 4];
$numericCopy = $numeric;
array_shift($numericCopy);
bucket_assert(array_keys($numeric) === [8, 20], 'shift renumbered source keys');
$numericCopy[] = &$value;
$value = 31;
bucket_assert($numeric[8] === 31 && $numericCopy[1] === 31, 'detached reference alias');
$nested = ['inner' => [9 => &$value]];
$nestedCopy = $nested;
sort($nestedCopy['inner']);
bucket_assert(array_keys($nested['inner']) === [9], 'nested sort changed source keys');
$nestedCopy['inner'][0] = 37;
bucket_assert($nested['inner'][9] === 37, 'nested reference cell lost');
echo "array reference buckets ok\n";
