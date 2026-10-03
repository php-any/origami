<?php
function ref_copy_assert($ok, $message) { if (!$ok) { throw new Exception($message); } }
$first = [1]; $alias =& $first[0];
$second = $first;
unset($alias);
$third = $second;
$third[0] = 7;
ref_copy_assert($first[0] === 1 && $second[0] === 1, 'single reference should unwrap after unset alias');
$alias =& $first[0]; $second = $first; $second[1] = 2;
unset($alias);
$spread = [...$first]; $spread[0] = 11;
ref_copy_assert($first[0] === 11 && $second[0] === 11, 'spread lost reference cell after separated copy');
$nested = [1]; $outer = [$nested];
$nested[0] = 3;
ref_copy_assert($outer[0][0] === 1, 'nested constructor array ownership');
$outer[0][0] = 5;
ref_copy_assert($nested[0] === 3, 'nested constructor write ownership');
class CloneReferenceOwner { public int $value = 4; }
$object = new CloneReferenceOwner;
$alias =& $object->value;
$clone = clone $object;
$clone->value = 9;
ref_copy_assert($object->value === 9 && $alias === 9, 'clone lost property reference');
echo "reference copy ownership ok\n";
$first = [1]; $alias =& $first[0];
$copy = $first; $copy[] = 2;
unset($copy[0]); unset($alias);
$next = $first; $next[0] = 8;
ref_copy_assert($first[0] === 1, 'removed bucket retained reference owner');
function borrow_reference_owner(&$value) {}
$first = [1]; borrow_reference_owner($first[0]);
$copy = $first; $copy[0] = 2;
ref_copy_assert($first[0] === 1, 'returned parameter retained reference owner');
$first = [1, 2]; foreach ($first as &$element) {} unset($element);
$copy = $first; $copy[0] = 3; $copy[1] = 4;
ref_copy_assert($first === [1, 2], 'foreach retained previous element owner');
function rebind_reference_owner(&$value, &$other) { $value =& $other; }
$first = [1]; $other = 2; rebind_reference_owner($first[0], $other);
$copy = $first; $copy[0] = 5;
ref_copy_assert($first[0] === 1, 'rebound parameter retained old owner');
