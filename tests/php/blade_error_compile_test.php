<?php

namespace tests\php;

/**
 * Blade @error 编译产物：if (...): 与花括号 } 混用，且中途 ?> 插入 HTML。
 */
$html = <<<'HTML'
<div>
<?php if(true): ?>
err
<?php } ?>
<?php $__errorArgs = ['email'];
$__bag = true;
if ($__bag) :
if (isset($message)) { $__messageOriginal = $message; }
$message = 'bad'; ?> <span><?php echo $message; ?></span> <?php unset($message);
if (isset($__messageOriginal)) { $message = $__messageOriginal; }
}
unset($__errorArgs, $__bag); ?>
</div>
HTML;

file_put_contents(sys_get_temp_dir() . '/origami_blade_error_compile.php', $html);
ob_start();
include sys_get_temp_dir() . '/origami_blade_error_compile.php';
$out = ob_get_clean();
if (!str_contains($out, 'err') || !str_contains($out, 'bad')) {
    Log::fatal('Blade 混用 if:/} 未能渲染: ' . $out);
}
Log::info('Blade @error 编译产物测试通过');
