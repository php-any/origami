<?php

namespace tests\php;

/**
 * boolval() 基础语义（Filament BooleanStateCast 依赖）。
 */
if (boolval(1) !== true || boolval(0) !== false) {
    Log::fatal('boolval 数字失败');
}
if (boolval('1') !== true || boolval('') !== false || boolval('0') !== false) {
    Log::fatal('boolval 字符串失败');
}
if (boolval(null) !== false || boolval([]) !== false || boolval([1]) !== true) {
    Log::fatal('boolval null/array 失败');
}

Log::info('boolval 测试通过');
