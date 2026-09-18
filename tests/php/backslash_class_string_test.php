<?php
// Simulate what happens if FQN is in a double-quoted string with single backslashes
$s = "\Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys";
echo "dq: [$s]\n";
$s2 = '\Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys';
echo "sq: [$s2]\n";
