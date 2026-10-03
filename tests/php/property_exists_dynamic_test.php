<?php
function exists_check($ok, $message) {
    if (!$ok) { throw new Exception($message); }
}
$object = new stdClass();
$object->nullable = null;
exists_check(property_exists($object, 'nullable'), 'null dynamic property lost');
exists_check(!isset($object->nullable), 'null dynamic property isset');
unset($object->nullable);
exists_check(!property_exists($object, 'nullable'), 'unset dynamic property retained');
$decoded = json_decode('{"nullable":null}');
exists_check(property_exists($decoded, 'nullable'), 'decoded stdClass property lost');
class ExistsExample {
    private $hidden = null;
    public int $pending;
    public static $shared;
    public function containsHidden() { return property_exists($this, 'hidden'); }
    public function __get($name) { return 1; }
}
$declared = new ExistsExample();
exists_check(property_exists($declared, 'hidden') && $declared->containsHidden(), 'private declaration / this');
exists_check(property_exists($declared, 'pending') && !isset($declared->pending), 'uninitialized declaration');
exists_check(property_exists(ExistsExample::class, 'shared'), 'static declaration');
exists_check(!property_exists($declared, 'virtual'), 'magic get is not a property declaration');
echo "property exists dynamic OK\n";
