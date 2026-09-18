<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

// Compile a filament blade that might trigger wire keys
$view = 'filament-panels::components.layout.base';
try {
  echo "compiling $view\n";
  $c = app('blade.compiler')->compileString('@livewire(\'test\')');
  echo substr($c, 0, 500), "\n";
} catch (Throwable $e) {
  echo "EX: ".$e->getMessage()."\n";
}

// Show what generated openLoop looks like when compiling a foreach in livewire context
try {
  Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop();
  $html = app('blade.compiler')->compileString('@foreach(($x??[]) as $i){{ $i }}@endforeach');
  Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop();
  echo "compiled:\n", $html, "\n";
  // Check if backslashes present
  echo "has bs: ".(str_contains($html, 'Livewire\\Features') || str_contains($html, 'Livewire\Features') ? 'yes':'no')."\n";
  echo "mangled: ".(str_contains($html, 'LivewireFeatures') ? 'yes':'no')."\n";
} catch (Throwable $e) {
  echo "EX2: ".$e->getMessage()."\n";
}
