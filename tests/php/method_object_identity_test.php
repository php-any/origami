<?php
class FluentIdentity {
    private $value = 0;
    public function selfObject() { return $this; }
    public function increment() { $this->value++; return $this; }
    public function value() { return $this->value; }
}
$original = new FluentIdentity;
$returned = $original->selfObject();
for ($i = 0; $i < 1000; $i++) { $returned->increment()->increment(); }
if ($returned !== $original || $original->value() !== 2000) {
    throw new RuntimeException('fluent object identity or private property scope failed');
}
echo "OK: method object identity\n";
