<?php
// Dump how origami tokenizes a blade-like FQN
$code = '<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop();';
$debug = __DIR__.'/../../examples/laravel13/storage/origami-debug';
file_put_contents($debug.'/tok_probe.php', $code);
echo "wrote\n";
