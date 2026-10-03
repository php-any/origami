<?php
function return_ref_assert($ok, $message) { if (!$ok) { throw new Exception($message); } }
function &nested_ref(&$array, &$counter) { if (true) { return $array[$counter++]; } }
$array = [3, 9]; $counter = 0;
$alias =& nested_ref($array, $counter);
$alias = 7;
return_ref_assert($counter === 1 && $array === [7, 9], 'nested index reference evaluated twice');
$copy = nested_ref($array, $counter);
$copy = 15;
return_ref_assert($counter === 2 && $array[1] === 9, 'ordinary reference return call aliased copy');
class RefReturnOwner {
    public int $value = 4;
    public function &get(): int { if (true) { return $this->value; } }
    public function &index(&$array) { return $array[0]; }
}
$owner = new RefReturnOwner;
$property =& $owner->get();
$property = 8;
return_ref_assert($owner->value === 8, 'method property reference');
$caught = false;
try { $property = []; } catch (TypeError $error) { $caught = true; }
return_ref_assert($caught && $owner->value === 8, 'typed reference return guard');
$method = 'get';
$dynamic =& $owner->$method();
$dynamic = 11;
return_ref_assert($owner->value === 11, 'dynamic method reference');
$closure = function &() use (&$array) { if (true) { return $array[1]; } };
$closed =& $closure(); $closed = 17;
return_ref_assert($array[1] === 17, 'closure nested reference return');
$arrow = fn &(&$input) => $input[0];
$arrowAlias =& $arrow($array); $arrowAlias = 23;
return_ref_assert($array[0] === 23, 'arrow reference return');
$staticArrow = static fn &(&$input) => $input[1];
$staticAlias =& $staticArrow($array); $staticAlias = 29;
return_ref_assert($array[1] === 29, 'static arrow reference return');
echo "reference return identity ok\n";
