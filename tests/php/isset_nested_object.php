<?php
class IssetNested_C {
    protected array $bindings = ['foo' => ['shared' => true]];
    public function check(string $abstract): bool {
        return isset($this->bindings[$abstract]['shared']);
    }
}
$c = new IssetNested_C();
echo $c->check('App\\Http\\Kernel') ? "1\n" : "0\n";
echo $c->check('foo') ? "1\n" : "0\n";
