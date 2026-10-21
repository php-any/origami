<?php
declare(strict_types=1);
function negationInt(int $n): int { return $n; }
$n = 28;
if (negationInt(-$n) !== -28 || -(-$n) !== 28) { throw new Exception('Integer negation changed type'); }
$min = -PHP_INT_MAX - 1;
if (!is_float(-$min) || -$min <= 0) { throw new Exception('Minimum integer overflow failed'); }
echo "integer negation passed\n";
