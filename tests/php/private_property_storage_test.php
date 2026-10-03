<?php
class PrivateStorageParent {
    public function __construct(private int $number = 3) {}
    public function value(): int { return $this->number; }
    public function set(int $n): void { $this->number = $n; }
    public function exists(): bool { return isset($this->number); }
    public function remove(): void { unset($this->number); }
}
class PrivateStorageChild extends PrivateStorageParent {
    private int $number = 9;
    public function __construct() { parent::__construct(4); }
    public function childValue(): int { return $this->number; }
}
$object = new PrivateStorageChild;
if ($object->value() !== 4 || $object->childValue() !== 9 || !$object->exists()) { throw new Exception('promoted private scope'); }
$parent = new ReflectionProperty(PrivateStorageParent::class, 'number');
$child = new ReflectionProperty(PrivateStorageChild::class, 'number');
$parent->setValue($object, 7);
if ($parent->getValue($object) !== 7 || $child->getValue($object) !== 9 || $object->value() !== 7) { throw new Exception('reflection private scope'); }
if (json_encode($object) !== '{}') { throw new Exception('private json exposure'); }
$object->remove();
if ($object->exists() || $parent->isInitialized($object) || !$child->isInitialized($object)) { throw new Exception('private unset identity'); }
echo "private property storage OK\n";
