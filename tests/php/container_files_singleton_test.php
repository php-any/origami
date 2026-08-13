<?php
namespace tests\php;

$base = dirname(__DIR__, 2).'/examples/laravel13';
require $base.'/vendor/autoload.php';

$app = \Illuminate\Foundation\Application::configure(basePath: $base)->create();

$app->singleton('files', function () {
    return new \Illuminate\Filesystem\Filesystem;
});

if (!$app->bound('files')) {
    Log::fatal('bound files 应为 true');
}
$f = $app->make('files');
if (!($f instanceof \Illuminate\Filesystem\Filesystem)) {
    Log::fatal('类型错误');
}

// Also register via real provider
$app2 = \Illuminate\Foundation\Application::configure(basePath: $base)->create();
$app2->register(new \Illuminate\Filesystem\FilesystemServiceProvider($app2));
if (!$app2->bound('files')) {
    Log::fatal('provider 后 bound files 应为 true');
}
$f2 = $app2->make('files');
if (!($f2 instanceof \Illuminate\Filesystem\Filesystem)) {
    Log::fatal('provider make 失败');
}

Log::info('container_files_singleton 测试通过');
