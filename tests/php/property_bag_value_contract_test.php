<?php
function property_bag_check($ok, $message) {
    if (!$ok) { throw new Exception($message); }
}
$object = (object)['' => null, 8 => 'eight', 'nested' => [1]];
property_bag_check($object instanceof stdClass && is_object($object) && !is_array($object), 'cast object identity');
property_bag_check(property_exists($object, '') && $object->{''} === null && $object->{'8'} === 'eight', 'cast property keys');
property_bag_check((array)(object)null === [] && (bool)(object)[] === true, 'empty object casts');
$closure = function () {};
property_bag_check((object)$closure === $closure, 'closure object cast identity');
foreach ([[$object, []], [[], $object], [[], 1]] as $operands) {
    $rejected = false;
    try { $sum = $operands[0] + $operands[1]; } catch (TypeError $error) { $rejected = true; }
    property_bag_check($rejected, 'object or scalar accepted by array union');
}
property_bag_check(array_map(null, ['keep' => 1, 8 => null]) === ['keep' => 1, 8 => null], 'map null callback keys');
property_bag_check(array_map(null, [1, 2], ['a']) === [[1, 'a'], [2, null]], 'map longest input');
property_bag_check(array_map(function ($a, $b) { return [$a, $b]; }, [1], ['a', 'b']) === [[1, 'a'], [null, 'b']], 'map callback padding');
property_bag_check(array_replace(['' => 1, 8 => 2], ['' => 3, '01' => 4]) === ['' => 3, 8 => 2, '01' => 4], 'replace key identity');
property_bag_check(array_replace_recursive(['' => ['a' => 1, 8 => 2]], ['' => ['a' => 3]]) === ['' => ['a' => 3, 8 => 2]], 'recursive replace key identity');
$keys = ['AA' => 1, 8 => 2, '' => 3];
property_bag_check(array_diff_ukey($keys, ['aa' => 0], function ($a, $b) { return strcasecmp((string)$a, (string)$b); }) === [8 => 2, '' => 3], 'callback key difference');
echo "Property bag value contract OK\n";
