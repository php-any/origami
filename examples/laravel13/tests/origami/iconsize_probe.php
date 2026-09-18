<?php
require __DIR__.'/../../vendor/autoload.php';
try {
  echo Filament\Support\Enums\IconSize::TwoExtraLarge->value, "\n";
  echo "ok\n";
} catch (Throwable $e) {
  echo "EX: ".$e->getMessage()."\n";
}
