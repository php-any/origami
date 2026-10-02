<?php
function nullTypeCheck($actual, $expected, $label) {
    if ($actual !== $expected) { throw new Exception($label); }
}
function nullTypeError($callback, $label) {
    try { $callback(); } catch (TypeError $expected) { return; }
    throw new Exception('missing TypeError: ' . $label);
}
function requiredString(string $value) { return $value; }
function requiredReference(string &$value) { return $value; }
function optionalString(string $value = null) { return $value; }
function nullableString(?string $value) { return $value; }
function nullableWithDefault(?string $value = 'default') { return $value; }
function requiredReturn(): string { return null; }
function nullableReturn(): ?string { return null; }
function missingReturn(): ?string {}
nullTypeError(fn() => requiredString(null), 'non nullable argument');
$reference = null;
nullTypeError(function () use (&$reference) { requiredReference($reference); }, 'non nullable reference');
nullTypeCheck(optionalString(), null, 'implicit nullable default');
nullTypeCheck(optionalString(null), null, 'implicit nullable explicit argument');
nullTypeCheck(nullableString(null), null, 'explicit nullable argument');
nullTypeCheck(nullableWithDefault(null), null, 'explicit null does not use default');
nullTypeCheck(nullableWithDefault(), 'default', 'omitted argument uses default');
nullTypeCheck(nullableReturn(), null, 'explicit nullable return');
nullTypeError(fn() => requiredReturn(), 'non nullable return');
nullTypeError(fn() => missingReturn(), 'missing nullable return');
nullTypeError(function (): string { return null; }, 'closure non nullable return');
nullTypeError(function (): ?string {}, 'closure missing return');
class NullTypeObject {
    public string $value = 'before';
    public ?string $nullable = 'before';
    public function required(): string { return null; }
    public function missing(): ?string {}
}
$object = new NullTypeObject();
nullTypeError(function () use ($object) { $object->value = null; }, 'non nullable property');
nullTypeCheck($object->value, 'before', 'failed null property write');
$object->nullable = null;
nullTypeCheck($object->nullable, null, 'nullable property');
nullTypeError(fn() => $object->required(), 'method non nullable return');
nullTypeError(fn() => $object->missing(), 'method missing return');
echo "typed null return: PASS\n";
