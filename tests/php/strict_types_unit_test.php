<?php
declare(strict_types=1);
require __DIR__ . '/strict_types_fixtures/weak.php';
require __DIR__ . '/strict_types_fixtures/strict.php';
function strictCheck($actual, $expected, $label) {
    if ($actual !== $expected) { throw new Exception($label); }
}
function expectStrictError($callback, $label) {
    try { $callback(); } catch (TypeError $error) { return; }
    throw new Exception('missing TypeError: ' . $label);
}
expectStrictError(fn() => weakString(42), 'caller determines parameter mode');
strictCheck(weakCallStrict(), '42', 'weak caller to strict declaration');
expectStrictError(fn() => strictCallWeak(), 'callee body has declaration mode');
strictCheck(weakReturn(), '42', 'weak return conversion');
expectStrictError(fn() => strictReturn(), 'strict return rejection');
strictCheck(strictFloat(42), 42.0, 'strict integer widens to float');
strictCheck(strictNumber(42), 42, 'union preserves exact integer');
expectStrictError(fn() => strictFloat('42'), 'strict float rejects string');
expectStrictError(fn() => strictNumber('42'), 'strict union rejects string');
strictCheck(strictNullableFloat(42), 42.0, 'nullable float widening');
strictCheck(strictNullableString(null), null, 'nullable null');
expectStrictError(fn() => strictNullableString(42), 'nullable string rejects integer');
expectStrictError(fn() => weakString(...[42]), 'direct spread argument');
$stringCallable = 'weakString';
expectStrictError(fn() => $stringCallable(value: 42), 'variable function named argument');
$object = new WeakTypedObject();
expectStrictError(fn() => $object->accepts(42), 'method parameter');
strictCheck($object->converted(), '42', 'method return mode');
expectStrictError(fn() => strictWrite($object), 'strict property assignment');
strictCheck($object->text, 'before', 'failed property write preserves value');
strictCheck(weakWrite($object), '42', 'weak property assignment result');
expectStrictError(fn() => new StrictTypedObject(42), 'constructor');
$strictObject = new StrictTypedObject('ok');
expectStrictError(fn() => $strictObject->rejected(), 'strict method return');
expectStrictError(fn() => $strictObject->accepts(value: 42), 'method named argument');
expectStrictError(fn() => $strictObject->accepts(...[42]), 'method spread argument');
expectStrictError(fn() => $strictObject->variadic('ok', 42), 'method variadic argument');
expectStrictError(fn() => StrictTypedObject::acceptsStatic(value: 42), 'static named argument');
expectStrictError(fn() => StrictTypedObject::acceptsStatic(42), 'static argument');
$invokable = new StrictInvokable();
expectStrictError(fn() => $invokable(42), 'invokable object argument');
$weakClosure = weakMakeClosure();
$strictClosure = strictMakeClosure();
$strictArrow = strictMakeArrow();
strictCheck($weakClosure(), '42', 'closure declaration mode');
expectStrictError($strictClosure, 'strict closure return');
expectStrictError($strictArrow, 'strict arrow return');
$reference = 42;
expectStrictError(function () use (&$reference) { strictReference($reference); }, 'reference argument');
strictCheck($reference, 42, 'failed reference conversion preserves value');
expectStrictError(function () use ($strictObject, &$reference) { $strictObject->reference($reference); }, 'method reference argument');
expectStrictError(fn() => strictVariadic('ok', 42), 'variadic argument');
expectStrictError(fn() => call_user_func('strictString', 42), 'forwarded callback strict mode');
expectStrictError(fn() => call_user_func_array('strictString', [42]), 'forwarded array callback strict mode');
strictCheck(array_map('strictString', [42]), ['42'], 'internal callback weak conversion');
expectStrictError(fn() => array_map('strictString', [[]]), 'internal callback type rejection');
expectStrictError(fn() => array_map(fn($value): string => 42, [0]), 'callback return declaration mode');
expectStrictError(function () { foreach (strictMakeGenerator() as $value) {} }, 'generator declaration mode');
echo "strict types unit: PASS\n";
