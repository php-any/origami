<?php

namespace tests\php;

function nc_helper($x) {
    return 'helper:' . $x;
}

// PHP: $a ?? f($b)  ===  $a ?? (f($b))，不应解析成 ($a ?? f)($b)
$a = 'keep';
$r = $a ?? nc_helper('x');
if ($r !== 'keep') {
    Log::fatal('非 null 左侧应短路，实际: ' . var_export($r, true));
}

$a = null;
$r = $a ?? nc_helper('y');
if ($r !== 'helper:y') {
    Log::fatal('null 左侧应调用右侧函数，实际: ' . var_export($r, true));
}

Log::info('null_coalesce_call_precedence 测试通过');
