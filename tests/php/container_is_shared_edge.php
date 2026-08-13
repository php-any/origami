<?php
class MiniContainerEdge {
    protected array $instances = [];
    protected array $bindings = [];

    public function isShared(string $abstract): bool {
        return isset($this->instances[$abstract]) ||
               (isset($this->bindings[$abstract]['shared']) &&
               $this->bindings[$abstract]['shared'] === true);
    }

    public function bindPartial(string $abstract): void {
        $this->bindings[$abstract] = ['concrete' => 'x'];
    }
}

$c = new MiniContainerEdge();
$abstract = 'App\\Http\\Kernel';
$c->bindPartial($abstract);
$r = $c->isShared($abstract);
echo "partial isShared=" . ($r ? '1' : '0') . "\n";
