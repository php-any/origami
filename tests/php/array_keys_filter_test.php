<?php
$a = ['int' => 1, 'string' => '1', 'null' => null, 'array' => [1]];
if (array_keys($a, 1, true) !== ['int']
    || array_keys($a, 1) !== ['int', 'string']
    || array_keys($a, null, true) !== ['null']
    || array_keys($a, [1], true) !== ['array']
    || array_keys($a) !== ['int', 'string', 'null', 'array']) {
    throw new \Exception('array_keys filter semantics failed');
}
echo "array keys filter passed\n";
