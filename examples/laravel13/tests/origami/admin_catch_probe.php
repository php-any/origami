<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();

echo "boot ok\n";
$request = Illuminate\Http\Request::create('/admin/login', 'GET');
echo "request ok\n";
echo "hasMacro: ";
var_dump(method_exists($request, 'hasMacro'));
try {
  var_dump(Illuminate\Http\Request::hasMacro('foo'));
} catch (Throwable $e) {
  echo "EX hasMacro: ".$e->getMessage()."\n";
}

// trigger something that used hasMacro in rendering path
try {
  $httpKernel = $app->make(Illuminate\Contracts\Http\Kernel::class);
  echo "http kernel ok\n";
  $response = $httpKernel->handle($request);
  echo "status=".$response->getStatusCode()."\n";
  echo substr((string)$response->getContent(), 0, 200)."\n";
} catch (Throwable $e) {
  echo "EX: ".$e->getMessage()."\n";
  echo $e->getFile().":".$e->getLine()."\n";
  $c = $e->getPrevious();
  $i=0;
  while ($c && $i++<8) {
    echo "prev: ".$c->getMessage()." @ ".$c->getFile().":".$c->getLine()."\n";
    $c = $c->getPrevious();
  }
}
