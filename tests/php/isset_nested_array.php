<?php
// isset on nested keys must not emit Undefined array key (PHP 8+)
$bindings = ['foo' => ['shared' => true]];
$abstract = 'App\\Http\\Kernel';

$w1 = isset($bindings[$abstract]['shared']);
$w2 = isset($bindings['foo']['shared']);
$w3 = isset($bindings[$abstract]['missing']);

echo "missing outer: " . ($w1 ? '1' : '0') . "\n";
echo "exists nested: " . ($w2 ? '1' : '0') . "\n";
echo "missing inner: " . ($w3 ? '1' : '0') . "\n";
