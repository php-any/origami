<?php
class MiniContainer {
    protected array $instances = [];
    protected array $bindings = [];

    public function isShared(string $abstract): bool {
        return isset($this->instances[$abstract]) ||
               (isset($this->bindings[$abstract]['shared']) &&
               $this->bindings[$abstract]['shared'] === true);
    }
}

$c = new MiniContainer();
$abstract = 'App\\Http\\Kernel';
$r = $c->isShared($abstract);
echo "isShared=" . ($r ? '1' : '0') . "\n";
