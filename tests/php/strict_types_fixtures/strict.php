<?php
declare(strict_types=1);
function strictString(string $value) { return $value; }
function strictFloat(float $value): float { return $value; }
function strictNumber(float|int $value) { return $value; }
function strictReturn(): string { return 42; }
function strictCallWeak() { return weakString(42); }
function strictWrite($object) { $object->text = 42; }
function strictMakeClosure() { return function (): string { return 42; }; }
function strictMakeArrow() { return fn(): string => 42; }
function strictReference(string &$value) { return $value; }
function strictVariadic(string ...$values) { return $values; }
function strictNullableFloat(?float $value) { return $value; }
function strictNullableString(?string $value) { return $value; }
function strictMakeGenerator() { yield strictString(42); }
class StrictInvokable {
    public function __invoke(string $value) { return $value; }
}
class StrictTypedObject {
    public function __construct(string $value) {}
    public function accepts(string $value) { return $value; }
    public function rejected(): string { return 42; }
    public function reference(string &$value) { return $value; }
    public function variadic(string ...$values) { return $values; }
    public static function acceptsStatic(string $value) { return $value; }
}
