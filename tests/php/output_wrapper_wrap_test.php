<?php

namespace tests\php;
if (!is_file(dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php')) {
    Log::info("skip: 缺少 vendor 依赖，跳过测试");
    return;
}


require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Helper\OutputWrapper;

$w = new OutputWrapper();
$in = 'Formatted success message';
$out = $w->wrap($in, 70, "\n");
Log::info('in: ' . var_export($in, true));
Log::info('out: ' . var_export($out, true));
Log::info('hex: ' . bin2hex($out));

if ($out !== $in && !str_contains($out, 'Formatted')) {
    Log::fatal('OutputWrapper::wrap corrupted text: ' . var_export($out, true));
}

Log::info('output_wrapper_wrap_test 测试通过');
