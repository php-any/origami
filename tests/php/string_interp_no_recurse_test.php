<?php

namespace tests\php;

/**
 * 双引号插值：{$var} 的值若含 $xxx，不得二次展开。
 * Filament @capture 依赖 "<?php {$name} = ..." 其中 $name='$content'。
 */
$name = '$content';
$arguments = '$logo, $isDarkMode = false';
$out = "
                <?php {$name} = (function (\$args) {
                    return function ({$arguments}) use (\$args) {
                        extract(\$args, EXTR_SKIP);
                        ob_start(); ?>
            ";

if (!str_contains($out, '<?php $content = (function')) {
    Log::fatal('二次插值或丢失 $content: ' . json_encode($out));
}
if (str_contains($out, '<?php  = (function') || str_contains($out, "<?php   = (function")) {
    Log::fatal('出现空左值赋值: ' . json_encode($out));
}

// 单层插值：值本身带 $
$a = '$b';
$b = 'WRONG';
$s = "X{$a}Y";
if ($s !== 'X$bY') {
    Log::fatal('{$a} 不应再展开 $b，实际: ' . json_encode($s));
}

Log::info('string_interp_no_recurse 测试通过');
