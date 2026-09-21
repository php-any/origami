<?php

namespace tests\php;

/**
 * 函数内 include 与调用者共享变量（PHP 语义），被引入文件的赋值必须写回。
 */

$dir = sys_get_temp_dir().DIRECTORY_SEPARATOR.'origami_incl_fn_share';
@mkdir($dir, 0777, true);
$inc = $dir.DIRECTORY_SEPARATOR.'inc.php';
file_put_contents($inc, "<?php\n\$incl_fn_share_a = 2;\n\$incl_fn_share_arr['k'] = 9;\n");

$got = (static function () use ($inc) {
    $incl_fn_share_a = 1;
    $incl_fn_share_arr = ['k' => 1];
    include $inc;
    return [$incl_fn_share_a, $incl_fn_share_arr['k']];
})();

@unlink($inc);
@rmdir($dir);

if ($got !== [2, 9]) {
    \Log::fatal('函数内 include 未写回调用者变量: '.var_export($got, true));
}

\Log::info('函数内 include 共享 zval 测试通过');
