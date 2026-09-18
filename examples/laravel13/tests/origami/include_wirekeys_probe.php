<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();

// Simulate including a mini compiled blade snippet
$code = '<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); echo "ok\n"; \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop();';
$file = storage_path('framework/views/_wirekeys_probe.php');
file_put_contents($file, $code);
include $file;
echo "included ok\n";
