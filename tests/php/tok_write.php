<?php
// Dump how origami tokenizes a blade-like FQN
$code = '<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop();';
file_put_contents(__DIR__.'/tok_probe.php', $code);
echo "wrote\n";
