<?php
function ref_call_assert($ok, $message) { if (!$ok) { throw new Exception($message); } }
class ReferenceCallBinding {
    public function mutate(int $amount, int &...$values): array {
        foreach ($values as &$value) { $value += $amount; }
        return array_keys($values);
    }
    public static function bump(int &$value): void { $value++; }
}
$object = new ReferenceCallBinding;
$values = ['4', 8]; $index = 0;
$object->mutate(2, $values[$index++], $values[$index++]);
ref_call_assert($values === [6, 10] && $index === 2, 'method arguments evaluated twice');
$method = 'mutate';
$object->$method(3, ...$values);
ref_call_assert($values === [9, 13], 'dynamic method unpacked references');
$a = '5'; $b = 7;
$keys = $object->mutate(amount: 4, first: $a, second: $b);
ref_call_assert($keys === ['first', 'second'] && $a === 9 && $b === 11, 'named variadic references');
$unpack = ['amount' => 1, 'last' => 3];
$object->mutate(...$unpack);
ref_call_assert($unpack['last'] === 4, 'string unpacked reference keys');
$empty = $object->mutate(1);
ref_call_assert($empty === [], 'zero variadic arguments');
ReferenceCallBinding::bump($values[0]);
ref_call_assert($values[0] === 10, 'static reference parameter');
$closure = function (int &...$args) { foreach ($args as &$value) { $value++; } };
$closure(...$values);
ref_call_assert($values === [11, 14], 'closure unpacked references');
$fixed = function (int &$left, int &$right) { $left += 2; $right += 3; };
$fixed(right: $values[1], left: $values[0]);
ref_call_assert($values === [13, 17], 'closure named references');
class ReferenceConstructor { public function __construct(int &$value) { $value++; } }
$object = new ReferenceConstructor($values[0]);
ref_call_assert($values[0] === 14, 'constructor references');
$bad = false;
try { $fixed(1, $values[0]); } catch (Error $error) { $bad = true; }
ref_call_assert($bad, 'nonvariable reference argument accepted');
echo "reference call binding ok\n";
