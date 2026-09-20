<?php

namespace tests\php;

/**
 * 顶层 include 仍应共享 $GLOBALS：a.php 赋值的变量，后续 b.php 能读到。
 */

$dir = sys_get_temp_dir().DIRECTORY_SEPARATOR.'origami_incl_top_share';
@mkdir($dir, 0777, true);
$a = $dir.DIRECTORY_SEPARATOR.'a.php';
$b = $dir.DIRECTORY_SEPARATOR.'b.php';
file_put_contents($a, "<?php\n\$incl_top_share_x = 7;\n");
file_put_contents($b, "<?php\nreturn \$incl_top_share_x ?? 'missing';\n");

include $a;
$got = include $b;

@unlink($a);
@unlink($b);
@rmdir($dir);

if ((int) $got !== 7) {
    \Log::fatal('顶层 include 未共享变量: '.var_export($got, true));
}

\Log::info('顶层 include 全局共享测试通过');
