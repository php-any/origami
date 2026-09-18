<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();
try {
  $http = $app->make(Illuminate\Contracts\Http\Kernel::class);
  $req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
  echo "handling\n";
  $response = $http->handle($req);
  echo "status=".$response->getStatusCode()."\n";
  $c = (string)$response->getContent();
  echo "len=".strlen($c)."\n";
  echo (str_contains(strtolower($c),'password')||str_contains($c,'filament')||str_contains($c,'Email')) ? "markers:yes\n" : "markers:no\n";
  echo substr(preg_replace('/\s+/',' ',strip_tags($c)),0,200)."\n";
} catch (Throwable $e) {
  echo "EX: ".$e->getMessage()."\n";
  echo $e->getFile().":".$e->getLine()."\n";
}
