<?php
$seen = [];
foreach ([1, 2] as $outer) {
    foreach ([1, 2] as $inner) { $seen[] = $outer; continue 0x2; }
    throw new Exception('continue did not leave inner loop');
}
if ($seen !== [1, 2]) { throw new Exception('continue levels'); }
$seen = [];
foreach ([1, 2] as $outer) {
    foreach ([1, 2] as $inner) { $seen[] = $outer; break 0b10; }
    throw new Exception('break did not leave both loops');
}
if ($seen !== [1]) { throw new Exception('break levels'); }
echo "loop control levels: PASS\n";
