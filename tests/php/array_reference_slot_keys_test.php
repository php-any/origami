<?php

function arraySlotCheck($condition, $message) {
    if (!$condition) {
        throw new Exception($message);
    }
}

$a = [10, 20];
$first = &$a[0];
$first = 11;
arraySlotCheck($a === [11, 20], 'packed reference must reuse key 0');
arraySlotCheck(count($a) === 2, 'reference must not insert a duplicate integer key');
unset($first);
arraySlotCheck($a[0] === 11, 'unsetting a reference must preserve its array entry');

$a = [];
$empty = &$a[''];
$empty = 'empty';
$a[0] = 'zero';
arraySlotCheck(array_keys($a) === ['', 0], 'empty string and integer 0 are distinct');
arraySlotCheck($a[''] === 'empty' && $a[0] === 'zero', 'empty-key reference must survive integer assignment');

$a = [10, 20];
$sparse = &$a[7];
$sparse = 70;
$appended = &$a[];
$appended = 80;
arraySlotCheck(array_keys($a) === [0, 1, 7, 8], 'append reference must use the next PHP integer key');
arraySlotCheck($a[7] === 70 && $a[8] === 80, 'sparse and appended references must write through');

$a = [10, 20];
$numeric = &$a['1'];
$numeric = 21;
$leadingZero = &$a['01'];
$leadingZero = 99;
arraySlotCheck(array_keys($a) === [0, 1, '01'], 'canonical numeric strings share integer keys');
arraySlotCheck($a[1] === 21 && $a['01'] === 99, 'noncanonical numeric string stays a string');

$a = [[1, 2]];
$b = $a;
$nested = &$b[0][0];
$nested = 9;
arraySlotCheck($a === [[1, 2]], 'nested reference binding must separate a value copy');
arraySlotCheck($b === [[9, 2]], 'nested reference must stay bound after writeback');

$a = ['left' => 1, 'right' => 2];
$keys = ['left', 'right'];
$index = 0;
$once = &$a[$keys[$index++]];
$once = 3;
arraySlotCheck($index === 1, 'reference key expression must be evaluated once');
arraySlotCheck($a['left'] === 3 && $a['right'] === 2, 'reference must bind the evaluated key');

$a = [];
$boolean = &$a[true];
$boolean = 1;
$float = &$a[2.0];
$float = 2;
arraySlotCheck(array_keys($a) === [1, 2], 'boolean and float references use integer keys');

$a = [];
$a[''] = 'keep';
$nullKey = &$a[null];
arraySlotCheck($nullKey === 'keep', 'null key must bind the existing empty-string entry');
$nullKey = 'changed';
arraySlotCheck($a[''] === 'changed' && count($a) === 1, 'null key reference writes through without appending');

echo "array reference slot keys: OK\n";
