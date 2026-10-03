<?php
function frameCheck($ok, $message) { if (!$ok) throw new Exception($message); }
function frameChurn($a, $b, $c) { $local = [$a, $b, $c]; return $local[0]; }
class FrameReference {
    public function &local() { $value = 42; return $value; }
    public function captured() { $value = 7; return function () use (&$value) { return ++$value; }; }
    public function arrayReference() { $value = 5; return [&$value, &$value]; }
}
$object = new FrameReference();
$reference =& $object->local();
$closure = $object->captured();
$aliases = $object->arrayReference();
for ($i = 0; $i < 500; $i++) frameChurn($i, $i + 1, $i + 2);
frameCheck($reference === 42, 'reference return bucket survived frame reuse');
$reference = 99;
frameCheck($closure() === 8 && $closure() === 9, 'closure capture survived frame reuse');
$aliases[0] = 17;
frameCheck($aliases[1] === 17, 'array reference aliases survived frame reuse');
$value = 1;
$first = function () use (&$value) { return ++$value; };
$second = function () use (&$value) { return ++$value; };
$value = 20;
frameCheck($first() === 21 && $second() === 22 && $value === 22, 'separate capture buckets share one reference cell');
unset($value);
for ($i = 0; $i < 500; $i++) frameChurn($i, $i + 1, $i + 2);
frameCheck($first() === 23 && $second() === 24, 'captures survive unset and pooled calls');
echo "Pooled frame identity OK\n";
