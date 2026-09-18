<?php

namespace tests\php;

/**
 * Closure::bind + include 应把 echo 写入 ob_start 缓冲区，且 $this 可用。
 */
class ObBindInclude_Comp
{
    public $content = 'hello-bind';
}

$path = sys_get_temp_dir() . DIRECTORY_SEPARATOR . 'origami_ob_bind_v.php';
file_put_contents($path, '<?php echo $this->content;');

$comp = new ObBindInclude_Comp();
ob_start();
\Closure::bind(function () use ($path) {
    include $path;
}, $comp, $comp)();
$out = ob_get_clean();

@unlink($path);

if ($out !== 'hello-bind') {
    Log::fatal('ob+bind+include 期望 hello-bind，实际: ' . var_export($out, true));
}
Log::info('ob_bind_include 测试通过');
