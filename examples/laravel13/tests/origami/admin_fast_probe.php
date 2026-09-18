<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();
try {
  $http = $app->make(Illuminate\Contracts\Http\Kernel::class);
  $req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
  $response = $http->handle($req);
  $c = (string)$response->getContent();
  echo 'status='.$response->getStatusCode()."\n";
  echo 'len='.strlen($c)."\n";
  echo (str_contains(strtolower($c),'password')||str_contains(strtolower($c),'filament')||str_contains($c,'Email')) ? "markers:yes\n" : "markers:no\n";
} catch (Throwable $e) {
  echo 'EX: '.get_class($e).': '.$e->getMessage()."\n";
  echo 'at '.$e->getFile().':'.$e->getLine()."\n";
  $p=$e->getPrevious(); $i=0;
  while($p && $i++<3){ echo 'prev: '.$p->getMessage()."\n"; $p=$p->getPrevious(); }
}
echo "DONE\n";
