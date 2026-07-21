<?php

namespace tests\php;

/**
 * 验证 array_merge_recursive 保留字符串键（Illuminate Validation 依赖）。
 */

$merged = array_merge_recursive([], ['title' => ['required', 'min:1']]);
if (!isset($merged['title']) || $merged['title'][0] !== 'required') {
    Log::fatal('array_merge_recursive 丢失字符串键');
}

Log::info('array_merge_recursive 测试通过');
