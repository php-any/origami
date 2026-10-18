<?php
$ticks = 0;
function syntaxCountTick(int $increment): void { $GLOBALS['ticks'] += $increment; }
register_tick_function('syntaxCountTick', 1);
declare(ticks=1) {
    $values = [1, 2];
    $values[] = 3;
}
if ($ticks !== 2 || $values !== [1, 2, 3]) {
    throw new Exception('ticks count statements, not array elements or subexpressions');
}
unregister_tick_function('syntaxCountTick');
declare(ticks=1) { $values[] = 4; }
if ($ticks !== 2) { throw new Exception('unregistered callback executed'); }
echo "tick statement count OK\n";
