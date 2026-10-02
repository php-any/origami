<?php

function bulkKeyCheck($condition, $message) {
    if (!$condition) {
        throw new Exception($message);
    }
}

$a = [];
$a[10] = 3;
$a[''] = 1;
usort($a, function ($left, $right) { return $left - $right; });
bulkKeyCheck($a === [1, 3], 'usort must drop every key, including empty-string keys');
$a[] = 4;
bulkKeyCheck($a === [1, 3, 4], 'usort must reset automatic-key state and lookup cache');

$a = [];
$a[99] = 2;
$b = $a;
usort($b, function ($left, $right) { return $left - $right; });
bulkKeyCheck(array_keys($a) === [99] && array_keys($b) === [0], 'usort of a single element must reindex a private copy');
$b[] = 3;
bulkKeyCheck($b === [2, 3], 'single-element sort must reset the next integer key');

$a = [];
$a[10] = 2;
$a[5] = 1;
array_multisort($a);
bulkKeyCheck($a === [1, 2], 'array_multisort must reorder and renumber integer keys');
$a[] = 3;
bulkKeyCheck($a === [1, 2, 3], 'array_multisort must reset automatic-key state');

$a = [];
$a['title'] = 'title';
$a[10] = 10;
$a[20] = 20;
$removed = array_splice($a, 1, 1, [30]);
bulkKeyCheck($removed === [10], 'splice returns removed elements');
bulkKeyCheck(array_keys($a) === ['title', 0, 1], 'splice preserves strings and renumbers integer keys');
bulkKeyCheck($a[0] === 30 && $a[1] === 20, 'splice replacement keeps its insertion order');
$a[] = 40;
bulkKeyCheck(array_keys($a) === ['title', 0, 1, 2], 'splice resets automatic-key state');

echo "array bulk API keys: OK\n";
