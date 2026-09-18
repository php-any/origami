<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();
$fqn = 'Livewire\\Features\\SupportCompiledWireKeys\\SupportCompiledWireKeys';
echo "exists: ".(class_exists($fqn)?'yes':'no')."\n";
echo "name: $fqn\n";
try {
  echo "call: ".Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::class."\n";
} catch (Throwable $e) {
  echo "EX: ".$e->getMessage()."\n";
}
