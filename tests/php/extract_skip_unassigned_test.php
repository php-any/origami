<?php

namespace tests\php;

/**
 * EXTR_SKIP：仅引用未赋值的变量不算「已存在」，应被 extract 导入。
 * 对应 Blade @capture 内 extract($args, EXTR_SKIP) 导入 $attributes。
 */

$args = ['attributes' => 'ok', 'name' => 'logo'];

$result = (function () use ($args) {
    // 体内引用 $attributes，符号表有槽但未赋值
    extract($args, EXTR_SKIP);
    return [$attributes ?? 'missing', $name ?? 'missing'];
})();

if ($result[0] !== 'ok' || $result[1] !== 'logo') {
    Log::fatal('EXTR_SKIP 未导入未赋值槽: ' . var_export($result, true));
}

// 已赋 null 的变量应被 EXTR_SKIP 跳过
$attrs2 = (function () {
    $attributes = null;
    extract(['attributes' => 'skip-me'], EXTR_SKIP);
    return $attributes;
})();
if ($attrs2 !== null) {
    Log::fatal('EXTR_SKIP 应对已赋 null 跳过: ' . var_export($attrs2, true));
}

Log::info('extract EXTR_SKIP 未赋值槽测试通过');
