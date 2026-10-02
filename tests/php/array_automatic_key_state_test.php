<?php

function autoKeyCheck($condition, $message) {
    if (!$condition) {
        throw new Exception($message);
    }
}

$a = [10, 20];
unset($a[1]);
$a[] = 30;
autoKeyCheck(array_keys($a) === [0, 2], 'unset must preserve the next integer key');

$a = [];
$a[100] = 1;
unset($a[100]);
$b = $a;
$b[] = 2;
autoKeyCheck(array_keys($b) === [101], 'value copy must preserve automatic-key history');
autoKeyCheck(count($a) === 0, 'append must stay within the value copy');

$a = [];
$a['title'] = 'title';
$a[] = 10;
$a[] = 20;
unset($a['title']);
autoKeyCheck(array_keys($a) === [0, 1], 'removing a string key must not renumber integer keys');

$a = [];
$a[10] = 10;
$a[2] = 2;
autoKeyCheck(array_pop($a) === 2, 'array_pop must follow insertion order');
$a[] = 11;
autoKeyCheck(array_keys($a) === [10, 11], 'popping a smaller key must not reset automatic-key history');
autoKeyCheck(array_pop($a) === 11, 'pop must return the last inserted integer');
$a[] = 12;
autoKeyCheck(array_keys($a) === [10, 11], 'pop of the preceding automatic key permits reuse');

$a = [];
$a[100] = 1;
unset($a[100]);
$a[2] = 2;
array_pop($a);
$a[] = 3;
autoKeyCheck(array_keys($a) === [101], 'pop must preserve an earlier deleted maximum');

$a = [];
$a[9] = 'first';
$a['name'] = 'name';
$a[7] = 'last';
autoKeyCheck(array_shift($a) === 'first', 'shift returns the first inserted value');
autoKeyCheck(array_keys($a) === ['name', 0], 'shift resets integer keys and preserves strings');
$a[] = 'tail';
autoKeyCheck(array_keys($a) === ['name', 0, 1], 'shift resets automatic-key state');

$a = [];
$a[7] = 7;
$a['name'] = 1;
$a[9] = 9;
array_unshift($a, 2, 3);
autoKeyCheck(array_keys($a) === [0, 1, 2, 'name', 3], 'unshift renumbers only integer keys');
$a[] = 10;
autoKeyCheck(array_keys($a) === [0, 1, 2, 'name', 3, 4], 'unshift resets automatic-key state');

$a = [];
$a[-5] = 'negative';
$negative = &$a[-6];
$negative = 'reference';
autoKeyCheck($a[-5] === 'negative' && $a[-6] === 'reference', 'negative integer keys are valid');
$a[] = 'tail';
$expectedTail = PHP_VERSION_ID >= 80300 ? -4 : 0;
autoKeyCheck(array_keys($a) === [-5, -6, $expectedTail], 'negative automatic-key rule follows PHP version');

$a = [];
$a[PHP_INT_MAX] = 'max';
$caught = false;
try {
    $a[] = 'overflow';
} catch (Error $error) {
    $caught = true;
}
autoKeyCheck($caught && $a[PHP_INT_MAX] === 'max', 'automatic-key overflow must throw without overwriting');
$caught = false;
try {
    array_push($a, 'overflow');
} catch (Error $error) {
    $caught = true;
}
autoKeyCheck($caught && $a[PHP_INT_MAX] === 'max', 'array_push must report occupied maximum key');

echo "array automatic key state: OK\n";
