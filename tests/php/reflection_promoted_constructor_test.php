<?php
class ReflectionPromotedParent {
    public function __construct(private array $properties = ['default' => 7], public string $name = 'base') {}
    public function properties(): array { return $this->properties; }
}
class ReflectionPromotedChild extends ReflectionPromotedParent {}
$reflection = new ReflectionClass(ReflectionPromotedChild::class);
$a = $reflection->newInstance();
$b = $reflection->newInstanceArgs([]);
$c = $reflection->newInstanceArgs(['name' => 'named']);
$d = $reflection->newInstance(['supplied' => 8], 'args');
if ($a->properties() !== ['default' => 7] || $b->properties() !== ['default' => 7] || $c->properties() !== ['default' => 7] || $c->name !== 'named' || $d->properties() !== ['supplied' => 8] || $d->name !== 'args') { throw new Exception('reflection constructor binding'); }
echo "reflection promoted constructor OK\n";
class ReflectionPromotedExplicit extends ReflectionPromotedParent {
    public function __construct() { ReflectionPromotedParent::__construct(name: 'explicit'); }
}
$explicit = new ReflectionPromotedExplicit;
if ($explicit->name !== 'explicit' || $explicit->properties() !== ['default' => 7]) { throw new Exception('class qualified constructor promotion'); }
class ReflectionReferenceConstructor { public function __construct(&$value) { $value++; } }
$value = 4;
(new ReflectionClass(ReflectionReferenceConstructor::class))->newInstanceArgs(['value' => &$value]);
if ($value !== 5) { throw new Exception('reflection constructor reference'); }
class ReflectionHiddenConstructor { private function __construct() {} }
try { (new ReflectionClass(ReflectionHiddenConstructor::class))->newInstance(); throw new Exception('hidden constructor called'); } catch (ReflectionException $e) {}
