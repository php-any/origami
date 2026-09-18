<?php
// Exact string from Livewire SupportCompiledWireKeys.php line 64
$s = "$1\n<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey(\$component); ?>\n";
echo "bs=", substr_count($s, "\\"), "\n";
echo $s;
