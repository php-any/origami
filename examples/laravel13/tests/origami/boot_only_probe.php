<?php
echo "start\n";
require __DIR__.'/../../vendor/autoload.php';
echo "autoload\n";
$app = require __DIR__.'/../../bootstrap/app.php';
echo "app\n";
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
echo "kernel\n";
$kernel->bootstrap();
echo "boot ok\n";
