<?php
function scalar_check($actual, $expected, $label) {
    if ($actual !== $expected) { throw new Exception($label . ': ' . json_encode($actual)); }
}
function stringArgument(string $value) { return $value; }
function boolArgument(bool $value) { return $value; }
function floatArgument(float $value) { return $value; }
function intArgument(int $value) { return $value; }
function stringReturn(): string { return 42; }
function boolReturn(): bool { return 1; }
function floatReturn(): float { return 42; }
function unionIntString(string|int $value) { return $value; }
function unionNumbers(float|int $value) { return $value; }
scalar_check(stringArgument(42), '42', 'string parameter value');
scalar_check(stringArgument(true), '1', 'bool to string');
scalar_check(boolArgument('0'), false, 'string to bool');
scalar_check(boolArgument(2), true, 'integer to bool');
scalar_check(floatArgument(42), 42.0, 'float parameter value');
scalar_check(intArgument('2.0'), 2, 'numeric string to int');
scalar_check(stringReturn(), '42', 'string return value');
scalar_check(boolReturn(), true, 'bool return value');
scalar_check(floatReturn(), 42.0, 'float return value');
scalar_check(unionIntString(2.0), 2, 'union scalar priority');
scalar_check(unionIntString('2'), '2', 'union exact string');
scalar_check(unionNumbers('2'), 2, 'union integer string');
scalar_check(unionNumbers('2.0'), 2.0, 'union decimal string');
scalar_check(unionNumbers('2e2'), 200.0, 'union exponent string');
foreach (['NaN', 'INF', 'Infinity', '2tail', []] as $invalid) {
    try { floatArgument($invalid); throw new Exception('invalid float accepted'); }
    catch (TypeError $expected) {}
}
try { boolArgument([]); throw new Exception('array converted to bool declaration'); }
catch (TypeError $expected) {}
try { intArgument([]); throw new Exception('array converted to int declaration'); }
catch (TypeError $expected) {}
class ScalarProperties {
    public string $text;
    public bool $flag;
    public float $number;
    public function change() { $this->text = 7; $this->flag = '0'; $this->number = 7; }
    public function stringResult(): string { return 8; }
    public function floatResult(): float { return 8; }
}
$object = new ScalarProperties();
$object->text = 42;
scalar_check(($object->text = 42), '42', 'assignment returns converted value');
$object->flag = 1;
$object->number = 42;
scalar_check([$object->text, $object->flag, $object->number], ['42', true, 42.0], 'typed property values');
$object->change();
scalar_check([$object->text, $object->flag, $object->number], ['7', false, 7.0], 'this property values');
$name = 'text';
$object->$name = 9;
scalar_check(($object->$name = 9), '9', 'dynamic assignment returns converted value');
scalar_check($object->text, '9', 'dynamic property conversion');
scalar_check($object->stringResult(), '8', 'method string return');
scalar_check($object->floatResult(), 8.0, 'method float return');
class ScalarStringable {
    public function __toString(): string { return 'converted'; }
}
scalar_check(stringArgument(new ScalarStringable()), 'converted', 'stringable argument');
$object->text = new ScalarStringable();
scalar_check($object->text, 'converted', 'stringable property');
class ThrowingScalarStringable {
    public function __toString(): string { throw new RuntimeException('string conversion failed'); }
}
try { stringArgument(new ThrowingScalarStringable()); throw new Exception('conversion exception swallowed'); }
catch (RuntimeException $expected) { scalar_check($expected->getMessage(), 'string conversion failed', 'argument conversion exception'); }
try { $object->text = new ThrowingScalarStringable(); throw new Exception('property conversion exception swallowed'); }
catch (RuntimeException $expected) {}
scalar_check($object->text, 'converted', 'failed conversion preserves property');
function throwingStringReturn(): string { return new ThrowingScalarStringable(); }
try { throwingStringReturn(); throw new Exception('return conversion exception swallowed'); }
catch (TypeError $expected) {
    scalar_check($expected->getPrevious() instanceof RuntimeException, true, 'return conversion chains original exception');
}
echo "typed scalar storage OK\n";
