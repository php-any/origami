<?php

namespace tests\php;

/**
 * 函数内 extract 后 require 的文件必须看见已导入变量（Laravel View getRequire）。
 */
$dir = sys_get_temp_dir().DIRECTORY_SEPARATOR.'origami_extract_require_test';
@mkdir($dir, 0777, true);
$file = $dir.DIRECTORY_SEPARATOR.'view.php';
file_put_contents($file, '<?php return isset($trigger) ? $trigger : "MISSING";');

$out = (static function () use ($file) {
    $__data = ['trigger' => 'FROM_EXTRACT'];
    extract($__data, EXTR_SKIP);
    return require $file;
})();

if ($out !== 'FROM_EXTRACT') {
    \Log::fatal('extract+require 作用域断裂: '.var_export($out, true));
}

// 对象也应可见
file_put_contents($file, '<?php return is_object($trigger) ? $trigger->html : "MISSING";');
class ExtractReq_Slot { public $html = 'OBJ'; }
$out2 = (static function () use ($file) {
    $__data = ['trigger' => new \tests\php\ExtractReq_Slot()];
    extract($__data, EXTR_SKIP);
    return require $file;
})();
if ($out2 !== 'OBJ') {
    \Log::fatal('extract+require 对象断裂: '.var_export($out2, true));
}

\Log::info('extract_require_scope 测试通过');
