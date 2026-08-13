<?php

namespace tests\php;

function same_name_fn($x) { return 'fn:'.$x; }

function same_name_caller($same_name_fn) {
    // 形参与函数同名时，bare same_name_fn(...) 应调用函数而非变量
    return same_name_fn('ok');
}

$r = same_name_caller('var');
if ($r !== 'fn:ok') {
    Log::fatal('期望 fn:ok，实际: ' . var_export($r, true));
}
Log::info('func_vs_var_same_name 测试通过');
