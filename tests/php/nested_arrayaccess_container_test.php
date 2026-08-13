<?php

class NestedArrayAccess_Container implements ArrayAccess
{
    private array $bindings = [];

    public function instance(string $key, mixed $value): void
    {
        $this->bindings[$key] = $value;
    }

    public function offsetExists(mixed $offset): bool
    {
        return array_key_exists($offset, $this->bindings);
    }

    public function offsetGet(mixed $offset): mixed
    {
        return $this->bindings[$offset] ?? null;
    }

    public function offsetSet(mixed $offset, mixed $value): void
    {
        $this->bindings[$offset] = $value;
    }

    public function offsetUnset(mixed $offset): void
    {
        unset($this->bindings[$offset]);
    }
}

$c = new NestedArrayAccess_Container();
$c->instance('config', [
    'database.default' => 'sqlite',
    'database.connections' => [
        'sqlite' => ['driver' => 'sqlite', 'database' => ':memory:'],
    ],
]);

$default = $c['config']['database.default'];
if ($default !== 'sqlite') {
    echo "FAIL read chained: ";
    var_export($default);
    echo "\n";
    exit(1);
}

$c['config']['database.default'] = 'mysql';
$after = $c['config']['database.default'];
if ($after !== 'mysql') {
    echo "FAIL write chained: ";
    var_export($after);
    echo "\n";
    exit(1);
}

echo "PASS\n";
