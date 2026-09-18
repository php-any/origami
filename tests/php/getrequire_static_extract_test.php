<?php

namespace tests\php;

/**
 * Laravel Filesystem::getRequire：类方法内 static 闭包 use + extract + require。
 * extract() 不得把闭包符号表还回 context pool，否则随后 $__path 会变成 slots=0。
 */

$dir = sys_get_temp_dir() . DIRECTORY_SEPARATOR . 'origami_getrequire_static';
@mkdir($dir, 0777, true);
$file = $dir . DIRECTORY_SEPARATOR . 'loaded.php';
file_put_contents($file, '<?php return $title . "-" . $__path_ok;');

class GetRequireStaticExtract_FS
{
    public function getRequire($path, array $data = [])
    {
        $__path = $path;
        $__data = $data;
        return (static function () use ($__path, $__data) {
            extract($__data, EXTR_SKIP);
            return require $__path;
        })();
    }

    public function afterExtractUse()
    {
        $__path = 'kept';
        return (static function () use ($__path) {
            extract([]);
            return $__path;
        })();
    }
}

$fs = new GetRequireStaticExtract_FS();

$gotUse = $fs->afterExtractUse();
if ($gotUse !== 'kept') {
    Log::fatal('extract 后 use 变量丢失: ' . var_export($gotUse, true));
}

$got = $fs->getRequire($file, ['title' => 'Hello', '__path_ok' => 'yes']);
if ($got !== 'Hello-yes') {
    Log::fatal('getRequire 模拟失败: ' . var_export($got, true));
}

@unlink($file);
@rmdir($dir);

Log::info('getrequire_static_extract 测试通过');
