<?php

namespace tests\php;

/**
 * extract + require 作用域共享（Laravel Filesystem::getRequire / PhpEngine）。
 */

$dir = __DIR__ . '/_extract_require_tmp';
@mkdir($dir, 0777, true);
$view = $dir . '/view.php';
file_put_contents($view, '<?php return "Hi ".$name;');

$__path = $view;
$__data = ['name' => 'Origami'];

$out = (static function () use ($__path, $__data) {
    extract($__data, EXTR_SKIP);
    return require $__path;
})();

@unlink($view);
@rmdir($dir);

if (trim($out) !== 'Hi Origami') {
    Log::fatal('extract+require 作用域共享失败: ' . var_export($out, true));
}

Log::info('extract+require 作用域共享测试通过');
