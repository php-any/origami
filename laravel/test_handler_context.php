<?php

require __DIR__.'/vendor/autoload.php';

$app = require_once __DIR__.'/bootstrap/app.php';

$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();

$handler = $app->make(Illuminate\Contracts\Debug\ExceptionHandler::class);

$ref = new ReflectionClass($handler);
$method = $ref->getMethod('context');
$method->setAccessible(true);

echo "calling context...\n";
$result = $method->invoke($handler);
echo "done\n";
var_export($result);
