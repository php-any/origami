<?php
$tickCount = 0;
function syntaxTickCallback(): void { $GLOBALS['tickCount']++; }
register_tick_function('syntaxTickCallback');
$before = $tickCount;
declare(ticks=1) {
    $tickProbe = 1;
    $tickProbe++;
}
$after = $tickCount;
unregister_tick_function('syntaxTickCallback');
if ($after <= $before || $tickProbe !== 2) {
    throw new Exception('declare ticks block must execute statements and callbacks');
}
$stable = $tickCount;
declare(ticks=1) {
    $tickProbe++;
}
if ($tickCount !== $stable || $tickProbe !== 3) {
    throw new Exception('unregistered tick callback must not execute');
}
echo "declare ticks block OK\n";
