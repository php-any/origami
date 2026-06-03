<?php
$bindings = [];
$abstract = 'App\\Http\\Kernel';
// Must not print Warning when key missing on outer index
$r = isset($bindings[$abstract]['shared']);
echo "result=$r\n";
