<?php
$notices = 0;
set_error_handler(function ($level, $message) use (&$notices) {
    if ($level !== E_NOTICE || $message !== 'Only variable references should be yielded by reference') {
        throw new Exception('unexpected yield diagnostic: ' . $message);
    }
    $notices++;
    return true;
});
function referenceYieldExpression(&$calls): int { $calls++; return 43; }
function &referenceYieldTemporaries(&$calls): Generator {
    yield 42;
    yield referenceYieldExpression($calls);
}
$calls = 0;
$seen = [];
foreach (referenceYieldTemporaries($calls) as $value) { $seen[] = $value; }
if ($seen !== [42, 43] || $calls !== 1 || $notices !== 2) {
    throw new Exception('temporary yield must evaluate once and continue after Notice');
}
$seen = [];
foreach (referenceYieldTemporaries($calls) as &$value) { $value++; $seen[] = $value; }
unset($value);
if ($seen !== [43, 44] || $calls !== 2 || $notices !== 4) {
    throw new Exception('temporary yield must have an independent reference slot');
}
restore_error_handler();
function &referenceYieldSource(&$value) { return $value; }
function &referenceYieldCall(&$value): Generator { yield referenceYieldSource($value); }
$source = 10;
foreach (referenceYieldCall($source) as &$value) { $value = 12; }
unset($value);
if ($source !== 12) { throw new Exception('reference-returning call must retain source identity'); }
set_error_handler(function ($level, $message) { throw new Exception('handler failure'); });
$caught = false;
try { foreach (referenceYieldTemporaries($calls) as $value) {} }
catch (Exception $e) { $caught = $e->getMessage() === 'handler failure'; }
restore_error_handler();
if (!$caught || $calls !== 2) { throw new Exception('yield must preserve error handler exceptions'); }
echo "generator reference temporary OK\n";
