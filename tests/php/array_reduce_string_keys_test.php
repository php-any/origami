<?php
namespace tests\php;

$a = iterator_to_array((function() { yield 'Laravel'; })());
echo "keys: "; var_export(array_keys($a)); echo "\n";

$sum = array_reduce($a, function ($acc, $v) { return $acc.$v; }, '');
if ($sum !== 'Laravel') {
    Log::fatal("reduce yield fail: ".var_export($sum, true)." a=".var_export($a, true));
}

// string key array
$b = ['0' => 'Laravel'];
$sum2 = array_reduce($b, function ($acc, $v) { return $acc.$v; }, '');
if ($sum2 !== 'Laravel') {
    Log::fatal("reduce string key fail: ".var_export($sum2, true));
}

Log::info('array_reduce_string_keys 测试通过');
