<?php
if (!is_file(__DIR__ . '/../../examples/laravel/vendor/autoload.php')) {
    Log::info("skip: 缺少 vendor 依赖，跳过测试");
    return;
}


require_once __DIR__ . '/../../examples/laravel/vendor/autoload.php';

use Illuminate\Support\Arr;

$data = ['view' => ['paths' => ['/tmp/views']]];
$v = Arr::get($data, 'view.paths.0');
if ($v !== '/tmp/views') {
    echo "FAIL Arr::get dot numeric: ";
    var_export($v);
    echo "\n";
    exit(1);
}

echo "PASS\n";
