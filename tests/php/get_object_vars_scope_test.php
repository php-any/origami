<?php
class ObjectVarsProbe {
    public $visible = 1;
    protected $protectedValue = 2;
    private $privateValue = 3;
    public int $uninitialized;
    public static $shared = 4;
    public function vars() { return get_object_vars($this); }
    public function __get($name) { throw new \Exception('Magic getter must not run'); }
}
$object = new ObjectVarsProbe();
if (get_object_vars($object) !== ['visible' => 1]) {
    throw new \Exception('External visibility or uninitialized property handling failed');
}
if ($object->vars() !== ['visible' => 1, 'protectedValue' => 2, 'privateValue' => 3]) {
    throw new \Exception('Caller scope not preserved');
}
$dynamic = json_decode('{"first":1,"second":2}');
if (get_object_vars($dynamic) !== ['first' => 1, 'second' => 2]) {
    throw new \Exception('Dynamic properties lost');
}
echo "object vars scope passed\n";
