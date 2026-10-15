<?php
trait FactoryClockDefaults {
    private ?string $clock = null;
    protected array $options = ['value' => 1];
    public int $uninitialized;
    public function clock() { return $this->clock; }
    public function options() { return $this->options; }
    public function changeOption() { $this->options['value'] = 2; }
}
class MutableFactoryDate extends DateTime {
    use FactoryClockDefaults;
    public function __construct() { throw new Exception('factory must not call user constructor'); }
}
class ImmutableFactoryDate extends DateTimeImmutable {
    use FactoryClockDefaults;
    public function __construct() { throw new Exception('factory must not call user constructor'); }
}
class InheritedFactoryDate extends MutableFactoryDate {
    public string $label = 'child';
}

foreach ([MutableFactoryDate::class, ImmutableFactoryDate::class, InheritedFactoryDate::class] as $class) {
    foreach (['format', 'interface'] as $factory) {
        $date = $factory === 'format'
            ? $class::createFromFormat('Y-m-d H:i:s', '2026-10-14 12:34:56')
            : $class::createFromInterface(new DateTime('2026-10-14 12:34:56'));
        if (get_class($date) !== $class || $date->clock() !== null || $date->options() !== ['value' => 1]) {
            throw new Exception('factory must initialize subclass/trait defaults: ' . $class . ' ' . $factory);
        }
        if ($date->format('Y-m-d H:i:s') !== '2026-10-14 12:34:56') {
            throw new Exception('factory must preserve date state');
        }
        if ($class === InheritedFactoryDate::class && $date->label !== 'child') {
            throw new Exception('inherited factory must initialize child defaults');
        }
        $caught = false;
        try { $unused = $date->uninitialized; } catch (Error $error) { $caught = true; }
        if (!$caught) { throw new Exception('typed property without default must remain uninitialized'); }
        $date->changeOption();
        $fresh = $class::createFromFormat('Y-m-d H:i:s', '2026-10-14 12:34:56');
        if ($fresh->options() !== ['value' => 1]) { throw new Exception('factory instances must isolate array defaults'); }
    }
}
echo "DateTime factory subclass defaults: PASS\n";
