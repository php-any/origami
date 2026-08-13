<?php

/**
 * 回归：外层闭包 use ($x) 后，若函数体未直接引用 $x、只再传给内层 use，
 * 捕获仍必须生效（Dotenv Loader::load 的 array_reduce + Option::map 路径）。
 */
$outer = 'captured';
$fn = static function () use ($outer) {
    $inner = static function () use ($outer) {
        return $outer;
    };
    return $inner();
};
$got = $fn();
if ($got !== 'captured') {
    echo "FAIL nested_use_only got=" . var_export($got, true) . "\n";
    exit(1);
}
echo "OK nested_use_only\n";
