<?php
function typeRefCheck($value, $expected, $label) {
    if ($value !== $expected) { throw new Exception($label); }
}
function typeRefError($callback, $label) {
    try { $callback(); } catch (TypeError $error) { return; }
    throw new Exception('missing TypeError: ' . $label);
}
function typeRefAbsent() {}
function typeRefMixed(): mixed {}
function typeRefVoid(): void {}
function typeRefBareVoid(): void { return; }
function typeRefNever(): never {}
function typeRefNeverThrow(): never { throw new RuntimeException('never throws'); }
typeRefCheck(typeRefAbsent(), null, 'no declaration');
typeRefError(fn() => typeRefMixed(), 'mixed still requires explicit return');
typeRefError(fn() => typeRefNever(), 'never cannot fall through');
typeRefCheck(typeRefVoid(), null, 'implicit void return');
typeRefCheck(typeRefBareVoid(), null, 'bare void return');
$void = function (): void {};
typeRefCheck($void(), null, 'void closure');
try { typeRefNeverThrow(); } catch (RuntimeException $error) { typeRefCheck($error->getMessage(), 'never throws', 'never exception'); }
class TypeRefMethods {
    public function absent() {}
    public function mixed(): mixed {}
    public function nothing(): void { return; }
    public function optional(?string $value): ?string { return $value; }
    public function number(int|false $value): int|false { return $value; }
}
$object = new TypeRefMethods();
typeRefCheck($object->nothing(), null, 'void method');
typeRefError(fn() => $object->mixed(), 'mixed method');
typeRefCheck($object->optional(null), null, 'nullable compact declaration');
typeRefCheck($object->number(false), false, 'literal false union');
typeRefCheck($object->number('2'), 2, 'compact union coercion');
$mixed = (new ReflectionMethod(TypeRefMethods::class, 'mixed'))->getReturnType();
typeRefCheck($mixed->getName(), 'mixed', 'explicit mixed reflection');
typeRefCheck($mixed->allowsNull(), true, 'mixed allows null');
typeRefCheck((new ReflectionMethod(TypeRefMethods::class, 'absent'))->getReturnType(), null, 'absent reflection');
$nullable = (new ReflectionMethod(TypeRefMethods::class, 'optional'))->getReturnType();
typeRefCheck($nullable instanceof ReflectionNamedType, true, 'nullable named reflection');
typeRefCheck($nullable->getName(), 'string', 'nullable reflection base');
typeRefCheck($nullable->allowsNull(), true, 'nullable reflection allows null');
echo "type ref declarations: PASS\n";
