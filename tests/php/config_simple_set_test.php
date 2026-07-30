<?php
namespace tests\php;
$base = dirname(__DIR__, 2).'/examples/laravel13';
require $base.'/vendor/autoload.php';
$cfg = new \Illuminate\Config\Repository(['x' => 1]);
$cfg->set('y', 2);
if ($cfg->get('y') !== 2) {
    Log::fatal('simple set fail: '.var_export($cfg->get('y'), true));
}
$cfg->set('app.providers', ['A']);
if ($cfg->get('app.providers') !== ['A']) {
    Log::fatal('nested fail');
}
Log::info('config_simple_set 测试通过');
