<?php

/**
 * Illuminate\Support\Arr::from：关联数组原样返回（Collection::getArrayableItems）。
 */
$src = ['a' => 1, 'b' => 2];
$got = Illuminate\Support\Arr::from($src);
if (!is_array($got) || ($got['a'] ?? null) !== 1 || ($got['b'] ?? null) !== 2) {
    fwrite(STDERR, 'Arr::from 失败: '.var_export($got, true)."\n");
    exit(1);
}
echo "arr_from_ok\n";
