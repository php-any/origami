<?php
class DynamicConstructedProbe {
    public $owner;
    public function __construct($owner) { $this->owner = $owner; }
}
class DynamicFactoryProbe {
    private $map = ['key' => DynamicConstructedProbe::class];
    public function create() { return new $this->map['key']($this); }
}
$factory = new DynamicFactoryProbe();
$value = $factory->create();
if (!($value instanceof DynamicConstructedProbe) || $value->owner !== $factory) {
    throw new \Exception('Dynamic new consumed constructor arguments as function call');
}
echo "dynamic new property passed\n";
