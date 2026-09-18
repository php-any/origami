<?php
require __DIR__.'/../../vendor/autoload.php';
echo \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::class, "\n";
\Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop();
echo "openLoop ok\n";
\Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop();
echo "closeLoop ok\n";
