<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();

echo "class Login exists: ".(class_exists(Filament\Auth\Pages\Login::class)?'yes':'no')."\n";

try {
  $panel = Filament\Facades\Filament::getPanel('admin');
  Filament\Facades\Filament::setCurrentPanel($panel);
  $login = app(Filament\Auth\Pages\Login::class);
  echo "login instance ok\n";
  $login->mount();
  echo "mount ok\n";
} catch (Throwable $e) {
  echo "EX: ".$e->getMessage()."\n";
  echo $e->getFile().":".$e->getLine()."\n";
  $c = $e->getPrevious();
  $i=0;
  while ($c && $i++<5) {
    echo "prev: ".$c->getMessage()." @ ".$c->getFile().":".$c->getLine()."\n";
    $c = $c->getPrevious();
  }
}
