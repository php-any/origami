<?php
// Reproduce: if FQN string loses backslashes via wrong escape in double quotes (PHP-like)
// Origami keeps \L; real PHP 8.2 may strip. Simulate strip:
$fqn = 'Livewire\\Features\\SupportCompiledWireKeys\\SupportCompiledWireKeys';
$bad = str_replace('\\', '', $fqn);
echo "bad=$bad\n";
