<?php
namespace tests\php;
function spread_func_take($a, $b) { return "$a-$b"; }
$arr = ['x', 'y'];
$r = spread_func_take(...$arr);
if ($r !== 'x-y') Log::fatal("got $r");
Log::info('spread_func_call 测试通过');
