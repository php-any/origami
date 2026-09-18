<?php

namespace tests\php;

/**
 * include/require 必须能看到调用者 extract 的变量，即使被引入文件顶层未引用该名。
 * 对应 Blade @capture：get_defined_vars() 需包含仅在内层闭包使用的 $attributes。
 */

$dir = sys_get_temp_dir() . '/origami_include_inject_test';
@mkdir($dir, 0777, true);
$file = $dir . '/inner.php';
file_put_contents($file, '<?php
$content = (function ($args) {
    return function () use ($args) {
        extract($args, EXTR_SKIP);
        return isset($attributes) ? $attributes : "missing";
    };
})(get_defined_vars());
return $content();
');

$__path = $file;
$__data = ['attributes' => 'from-extract', 'other' => 1];
$got = (static function () use ($__path, $__data) {
    extract($__data, EXTR_SKIP);
    return require $__path;
})();

@unlink($file);
@rmdir($dir);

if ($got !== 'from-extract') {
    Log::fatal('include 未注入仅内层引用的 extract 变量: ' . var_export($got, true));
}

Log::info('include 注入 extract 变量测试通过');
