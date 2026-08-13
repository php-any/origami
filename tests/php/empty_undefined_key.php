<?php
$abstractAliases = [];
$abstract = 'App\\Http\\Kernel';
$r = empty($abstractAliases[$abstract]);
echo "empty=" . ($r ? '1' : '0') . "\n";
