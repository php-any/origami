<?php

function checkIssetAccess($ok, $message) {
    if (!$ok) { throw new Exception($message); }
}

class IssetAccessMagic {
    public $events = [];
    public function __isset($name) {
        $this->events[] = 'isset:' . $name;
        if ($name === 'throw') { throw new Exception('magic'); }
        return $name !== 'missing';
    }
    public function __get($name) {
        $this->events[] = 'get:' . $name;
        return ['leaf' => 7];
    }
}

$o = new IssetAccessMagic();
$name = 'present';
checkIssetAccess(isset($o->$name), 'dynamic isset');
checkIssetAccess($o->events === ['isset:present'], 'dynamic isset must not fetch');
$o->events = [];
checkIssetAccess(isset($o->{$name}['leaf']), 'dynamic property dimension');
checkIssetAccess($o->events === ['isset:present', 'get:present'], 'dimension fetch once');
$name = 'missing';
$o->events = [];
checkIssetAccess(($o->$name ?? 'fallback') === 'fallback', 'dynamic coalesce');
checkIssetAccess($o->events === ['isset:missing'], 'missing coalesce must not fetch');

function issetAccessThrowReceiver() { throw new Exception('receiver'); }
function expectIssetAccessException($callback, $message) {
    $caught = false;
    try { $callback(); } catch (Exception $e) { $caught = $e->getMessage() === $message; }
    checkIssetAccess($caught, 'exception propagation: ' . $message);
}
expectIssetAccessException(function () use ($o) { return isset($o->throw); }, 'magic');
expectIssetAccessException(function () use ($o) { $k = 'throw'; return isset($o->$k); }, 'magic');
expectIssetAccessException(function () use ($o) { return $o->throw ?? null; }, 'magic');
expectIssetAccessException(function () use ($o) { return isset($o->throw['leaf']); }, 'magic');
expectIssetAccessException(function () { return isset(issetAccessThrowReceiver()->x); }, 'receiver');

$a = ['' => 42];
$b = [0 => 99];
checkIssetAccess(isset($a[null]) && !isset($b[null]), 'null key is empty string');
checkIssetAccess(($a[null] ?? 0) === 42 && ($b[null] ?? 0) === 0, 'coalesce null key');
checkIssetAccess(!empty($a[null]) && empty($b[null]), 'empty null key');
$a = ['x' => null, 'zero' => 0, 'false' => false];
checkIssetAccess(!isset($a['x']) && isset($a['zero'], $a['false']), 'non-null semantics');
checkIssetAccess(!isset($a['absent']['nested']), 'missing nested key');
checkIssetAccess(!isset($a['absent'], issetAccessThrowReceiver()->x), 'argument short circuit');

$plain = new stdClass();
$plain->x = 1;
$caught = false;
try { isset($plain['x']); } catch (Error $e) { $caught = true; }
checkIssetAccess($caught, 'ordinary object is not ArrayAccess');
foreach ([[], new stdClass()] as $key) {
    $caught = false;
    try { isset($a[$key]); } catch (TypeError $e) { $caught = true; }
    checkIssetAccess($caught, 'illegal array key must throw TypeError');
    $caught = false;
    try { $unused = $a[$key] ?? null; } catch (TypeError $e) { $caught = true; }
    checkIssetAccess($caught, 'coalesce illegal array key');
    $caught = false;
    try { empty($a[$key]); } catch (TypeError $e) { $caught = true; }
    checkIssetAccess($caught, 'empty illegal array key');
}

class IssetAccessOffsets implements ArrayAccess {
    public $events = [];
    public function offsetExists($key): bool { $this->events[] = 'exists'; return true; }
    public function offsetGet($key): mixed { $this->events[] = 'get'; return null; }
    public function offsetSet($key, $value): void {}
    public function offsetUnset($key): void {}
}
$offsets = new IssetAccessOffsets();
checkIssetAccess(isset($offsets[[]]), 'ArrayAccess owns key semantics');
checkIssetAccess($offsets->events === ['exists'], 'isset ArrayAccess must not fetch');
echo "isset access semantics OK\n";
