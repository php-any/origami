<?php
namespace tests\php;

require dirname(__DIR__, 2).'/examples/laravel13/vendor/autoload.php';

$app = Illuminate\Foundation\Application::configure(basePath: dirname(__DIR__, 2).'/examples/laravel13')->create();

// Mimic FilesystemServiceProvider
$app->singleton('files', function () {
    return new Illuminate\Filesystem\Filesystem;
});

if (!$app->bound('files')) {
    Log::fatal('bound files 应为 true');
}
$f = $app->make('files');
if (!($f instanceof Illuminate\Filesystem\Filesystem)) {
    Log::fatal('make files 类型错误');
}
Log::info('container_string_binding 测试通过');
