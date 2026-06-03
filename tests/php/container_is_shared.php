<?php
require __DIR__ . '/../../examples/laravel/vendor/autoload.php';

$app = require __DIR__ . '/../../examples/laravel/bootstrap/app.php';
(new \Illuminate\Foundation\Bootstrap\LoadConfiguration())->bootstrap($app);

$container = $app;
$abstract = \Illuminate\Contracts\Http\Kernel::class;
$r = $container->isShared($abstract);
echo "isShared=" . ($r ? '1' : '0') . "\n";
