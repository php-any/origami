<?php

class ArrayVisibilityParent {
    protected $inherited = 'parent';
    private $parentSecret = 'parent secret';
}
class ArrayVisibilityChild extends ArrayVisibilityParent {
    public $visible = 'public';
    protected $hidden = 'protected';
    private $secret = 'private';
    public $implicitNull;
    public int $uninitialized;
    public static $staticValue = 'static';
    public function castThis(): array { return (array) $this; }
}
$object = new ArrayVisibilityChild();
$properties = (array) $object;
if (($properties["\0*\0hidden"] ?? null) !== 'protected'
    || ($properties["\0ArrayVisibilityChild\0secret"] ?? null) !== 'private'
    || ($properties["\0*\0inherited"] ?? null) !== 'parent'
    || ($properties["\0ArrayVisibilityParent\0parentSecret"] ?? null) !== 'parent secret'
    || ($properties['visible'] ?? null) !== 'public') {
    throw new Exception('Object array cast lost property visibility');
}
if (!array_key_exists('implicitNull', $properties) || $properties['implicitNull'] !== null
    || array_key_exists('uninitialized', $properties) || array_key_exists('staticValue', $properties)
    || array_key_exists('hidden', $properties) || array_key_exists('secret', $properties)) {
    throw new Exception('Object array cast included incorrect properties');
}
if ($object->castThis() !== $properties) {
    throw new Exception('Casting $this differs from casting the instance');
}
$dynamic = new stdClass();
$dynamic->name = 'dynamic';
$dynamic->{'12'} = 'numeric key';
if ((array) $dynamic !== ['name' => 'dynamic', 12 => 'numeric key']) {
    throw new Exception('Dynamic properties must retain their keys');
}
$exception = new RuntimeException('message', 17);
$properties = (array) $exception;
if (($properties["\0*\0message"] ?? null) !== 'message'
    || ($properties["\0*\0code"] ?? null) !== 17
    || !is_int($properties["\0*\0line"] ?? null)
    || !is_string($properties["\0*\0file"] ?? null)) {
    throw new Exception('Exception cast must retain protected scalar values');
}
echo "object array visibility OK\n";
