<?php
function ref_variadic_assert($ok, $message) { if (!$ok) { throw new Exception($message); } }
function change_refs(int &...$values) { foreach ($values as &$value) { $value += 2; } return count($values); }
$a = '3'; $b = 7;
ref_variadic_assert(change_refs($a, $b) === 2 && $a === 5 && $b === 9, 'typed variadic references');
$array = ['10', 20];
change_refs(...$array);
ref_variadic_assert($array === [12, 22], 'unpacked references');
class ReferenceStringOwner { public string $value = '4'; }
$owner = new ReferenceStringOwner;
function take_int(int &$value) {}
$caught = false;
try { take_int($owner->value); } catch (TypeError $error) { $caught = true; }
ref_variadic_assert($caught && $owner->value === '4', 'by-reference coercion bypassed property type');
$one = 1; $left = ['x' => &$one]; $right = $left; sort($right);
$id1 = ReflectionReference::fromArrayElement($left, 'x')->getId();
$id2 = ReflectionReference::fromArrayElement($right, 0)->getId();
ref_variadic_assert(strlen($id1) === 20 && $id1 === $id2, 'ReflectionReference identity');
echo "variadic reference constraints ok\n";
