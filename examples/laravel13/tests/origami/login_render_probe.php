<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();
try {
  $panel = Filament\Facades\Filament::getPanel('admin');
  Filament\Facades\Filament::setCurrentPanel($panel);
  $login = app(Filament\Auth\Pages\Login::class);
  $login->mount();
  $html = $login->render()->render();
  echo "len=".strlen($html)."\n";
} catch (Throwable $e) {
  echo "EX: ".$e->getMessage()."\n";
  echo "at ".$e->getFile().":".$e->getLine()."\n";
  $c=$e->getPrevious(); $i=0;
  while($c && $i++<8){
    echo "prev: ".$c->getMessage()." @ ".$c->getFile().":".$c->getLine()."\n";
    $c=$c->getPrevious();
  }
}
