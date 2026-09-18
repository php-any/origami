<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();
echo "boot\n";
try {
  var_dump(Filament\Support\View\ComponentAttributeBag::hasMacro('color'));
  echo "cab ok\n";
} catch (Throwable $e) {
  echo "CAB EX: ".$e->getMessage()."\n";
}
try {
  $http = $app->make(Illuminate\Contracts\Http\Kernel::class);
  $req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
  echo "handling\n";
  $response = $http->handle($req);
  echo "status=".$response->getStatusCode()."\n";
} catch (Throwable $e) {
  echo "EX: ".$e->getMessage()."\n";
  echo $e->getFile().":".$e->getLine()."\n";
}
