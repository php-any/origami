<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();
try {
  $panel = Filament\Facades\Filament::getPanel('admin');
  Filament\Facades\Filament::setCurrentPanel($panel);
  $login = app(Filament\Auth\Pages\Login::class);
  $login->mount();
  $view = $login->render();
  echo "view class=".get_class($view)."\n";
  $html = $view->render();
  echo "len=".strlen($html)."\n";
  echo "head=".substr(preg_replace('/\s+/',' ', $html),0,300)."\n";
  echo "has_div=".(str_contains($html,'<div')?'yes':'no')."\n";
} catch (Throwable $e) {
  echo "EX: ".get_class($e).": ".$e->getMessage()."\n";
  echo "at ".$e->getFile().":".$e->getLine()."\n";
}
echo "DONE\n";
