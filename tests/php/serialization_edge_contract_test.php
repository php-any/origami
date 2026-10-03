<?php
function edge_check($ok, $message) { if (!$ok) throw new Exception($message); }
$warnings = [];
set_error_handler(function ($code, $message) use (&$warnings) { $warnings[] = [$code, $message]; return true; });

class EdgeSleepChange {
    public $x = 1;
    public function __sleep() { $this->x = 2; $this->y = 3; return ['x', 'y']; }
}
edge_check(serialize(new EdgeSleepChange()) === 'O:15:"EdgeSleepChange":2:{s:1:"x";i:2;s:1:"y";i:3;}', 'sleep sees added properties');
class EdgeSleepBase { private $x = 1; public function __sleep() { return ['x']; } }
class EdgeSleepChild extends EdgeSleepBase { private $x = 2; }
edge_check(serialize(new EdgeSleepChild()) === "O:14:\"EdgeSleepChild\":1:{s:17:\"\0EdgeSleepChild\0x\";i:2;}", 'sleep chooses child private property');
class EdgeSleepBad { public function __sleep() { return null; } }
$warnings = [];
edge_check(serialize(new EdgeSleepBad()) === 'N;', 'invalid sleep result serializes null');
edge_check(count($warnings) === 1 && $warnings[0][0] === E_WARNING, 'invalid sleep emits warning');
class EdgeSleepDuplicate { public $x = 1; public function __sleep() { return ['x', 'x', 'missing']; } }
$warnings = [];
edge_check(serialize(new EdgeSleepDuplicate()) === 'O:18:"EdgeSleepDuplicate":1:{s:1:"x";i:1;}', 'sleep deduplicates properties');
edge_check(count($warnings) === 2 && $warnings[0][0] === E_WARNING && $warnings[1][0] === E_WARNING, 'sleep duplicate and missing warnings');

class EdgeVisibility { protected $x = 1; private $y = 2; }
$serialized = serialize(new EdgeVisibility());
edge_check(serialize(unserialize($serialized, ['allowed_classes' => false])) === $serialized, 'incomplete class preserves mangled property names');

edge_check(unserialize('a:1:{i:0;a:0:{}}', ['max_depth' => 1]) === [[]], 'empty nested arrays do not consume depth');
$warnings = [];
edge_check(unserialize('a:1:{i:0;a:1:{i:0;i:1;}}', ['max_depth' => 1]) === false, 'nested populated array exceeds depth');
edge_check(count($warnings) === 2 && str_contains($warnings[0][1], 'Maximum depth of 1 exceeded'), 'depth warning');
edge_check(unserialize('a:1:{i:0;O:8:"stdClass":0:{}}', ['max_depth' => 1]) === false, 'empty object consumes depth');
$oldDepth = ini_set('unserialize_max_depth', '1');
edge_check(unserialize('a:1:{i:0;a:1:{i:0;i:1;}}') === false, 'depth ini default');
edge_check(unserialize('a:1:{i:0;a:1:{i:0;i:1;}}', ['max_depth' => 0]) === [[1]], 'explicit unlimited depth');
ini_set('unserialize_max_depth', $oldDepth);

foreach (['d:nan;', 'd:NaN;', 'd:+INF;', 'd:Inf;', 'd:0x1p2;', 'd:1_0;', 'd: 1;', 'd:1e;'] as $bad) {
    edge_check(unserialize($bad) === false, 'reject non-PHP double syntax: ' . $bad);
}
edge_check(unserialize('d:1e999;') === INF && unserialize('d:-1e999;') === -INF, 'double overflow');
edge_check(unserialize('d:1e-999;') === 0.0, 'double underflow');
edge_check(unserialize('d:+.5;') === 0.5 && unserialize('d:1.;') === 1.0, 'PHP decimal forms');
$warnings = [];
edge_check(unserialize('i:9223372036854775808;') === PHP_INT_MAX, 'integer overflow saturates');
edge_check(unserialize('i:-9223372036854775809;') === PHP_INT_MIN, 'integer underflow saturates');
edge_check(count($warnings) === 2 && str_contains($warnings[0][1], 'Numerical result out of range'), 'integer overflow warnings');
$oldPrecision = ini_set('serialize_precision', '3');
edge_check(serialize([1.23456, 1000000.0, 1.0e-5]) === 'a:3:{i:0;d:1.23;i:1;d:1.0E+6;i:2;d:1.0E-5;}', 'configured precision');
ini_set('serialize_precision', '0');
edge_check(serialize(1.23456) === 'd:1.0E+0;', 'zero precision');
ini_set('serialize_precision', '-1');
edge_check(serialize([0.0, -0.0, 1.0e16, 1.0e17]) === 'a:4:{i:0;d:0;i:1;d:-0;i:2;d:10000000000000000;i:3;d:1.0E+17;}', 'shortest floats and negative zero');
ini_set('serialize_precision', $oldPrecision);

class EdgeTypedReference { public int $x; }
$rejected = false;
try { unserialize('a:2:{i:0;s:1:"4";i:1;O:18:"EdgeTypedReference":1:{s:1:"x";R:2;}}'); }
catch (TypeError $e) { $rejected = true; }
edge_check($rejected, 'typed reference rejects numeric strings without coercion');
$copy = unserialize('a:2:{i:0;i:4;i:1;O:18:"EdgeTypedReference":1:{s:1:"x";R:2;}}');
$copy[0] = 8;
edge_check($copy[1]->x === 8, 'typed property reference identity');
restore_error_handler();
edge_check(serialize(unserialize('a:1:{i:0;R:1;}')) === 'a:1:{i:0;N;}', 'exclusive temporary root recursion');
$rootCycle = unserialize('a:1:{i:0;R:1;}');
edge_check(serialize($rootCycle) === 'a:1:{i:0;N;}', 'sole root reference is released');
edge_check(serialize(unserialize('a:2:{i:0;R:1;i:1;R:1;}')) === 'a:2:{i:0;a:2:{i:0;R:2;i:1;R:2;}i:1;R:2;}', 'multiple root aliases survive');
echo "serialization edge contract OK\n";
