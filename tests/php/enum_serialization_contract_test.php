<?php
error_reporting(0);
function checkEnumContract($condition, $message) { if (!$condition) { throw new Exception($message); } }
enum UnitProbe { case First; case Second; }
enum BackProbe: int { case Two = 2; case One = 1; }
checkEnumContract(UnitProbe::First instanceof UnitEnum, 'unit interface');
checkEnumContract(!(UnitProbe::First instanceof BackedEnum), 'unit is not backed');
checkEnumContract(BackProbe::cases() === [BackProbe::Two, BackProbe::One], 'case order');
checkEnumContract(BackProbe::from(1) === BackProbe::One, 'canonical from');
foreach ([UnitProbe::First, BackProbe::Two] as $case) {
    checkEnumContract(unserialize(serialize($case)) === $case, 'canonical E identity');
    checkEnumContract(unserialize(serialize($case), ['allowed_classes' => false]) === $case, 'enum ignores allowed classes');
    foreach ([fn() => new UnitProbe(), fn() => clone $case, fn() => (new ReflectionClass(UnitProbe::class))->newInstanceWithoutConstructor()] as $operation) {
        try { $operation(); throw new Exception('enum operation should fail'); } catch (Error $error) {}
    }
}
class CustomWire implements Serializable {
    public string $payload = 'a}b';
    public function serialize(): string { return $this->payload; }
    public function unserialize(string $payload): void { $this->payload = $payload; }
}
$object = new CustomWire();
$wire = serialize($object);
checkEnumContract($wire === 'C:10:"CustomWire":3:{a}b}', 'C wire byte length');
checkEnumContract(unserialize($wire)->payload === 'a}b', 'C unserialize hook');
echo "enum and custom serialization OK\n";
