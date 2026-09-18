<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();
Filament\Facades\Filament::setCurrentPanel(Filament\Facades\Filament::getPanel('admin'));
try {
  app(Filament\Auth\Pages\Login::class)->mount();
} catch (Throwable $e) {
  echo $e->getMessage()."\n";
  echo $e->getTraceAsString()."\n";
}
