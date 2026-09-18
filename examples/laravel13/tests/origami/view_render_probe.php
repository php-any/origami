<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();
echo "boot\n";
try {
  $html = view('welcome')->render();
  echo "welcome len=".strlen($html)."\n";
} catch (Throwable $e) {
  echo "EX: ".$e->getMessage()."\n";
  echo $e->getFile().":".$e->getLine()."\n";
}
