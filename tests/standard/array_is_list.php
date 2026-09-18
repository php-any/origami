<?php

namespace tests\standard;

/**
 * PHP 8.1：array_is_list 键必须是 0..n-1。
 */

if (array_is_list([]) !== true) {
    Log::fatal('空数组是 list');
}
if (array_is_list([1, 2, 3]) !== true) {
    Log::fatal('连续数字键是 list');
}
if (array_is_list(['a' => 1]) !== false) {
    Log::fatal('关联数组不是 list');
}
if (array_is_list([1 => 'a', 2 => 'b']) !== false) {
    Log::fatal('非从 0 起的数字键不是 list');
}

Log::info('standard array_is_list 测试通过');
