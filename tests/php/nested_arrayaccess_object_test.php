<?php

class NestedArrayAccessObject_Store implements ArrayAccess
{
    private array $values;

    public function __construct(array $values)
    {
        $this->values = $values;
    }

    public function offsetExists(mixed $offset): bool
    {
        return array_key_exists($offset, $this->values);
    }

    public function offsetGet(mixed $offset): mixed
    {
        return $this->values[$offset] ?? null;
    }

    public function offsetSet(mixed $offset, mixed $value): void
    {
        $this->values[$offset] = $value;
    }

    public function offsetUnset(mixed $offset): void
    {
        unset($this->values[$offset]);
    }

    public function get(string $key): mixed
    {
        return $this->values[$key] ?? null;
    }
}

class NestedArrayAccessObject_Container implements ArrayAccess
{
    private array $values = [];

    public function offsetExists(mixed $offset): bool
    {
        return array_key_exists($offset, $this->values);
    }

    public function offsetGet(mixed $offset): mixed
    {
        return $this->values[$offset] ?? null;
    }

    public function offsetSet(mixed $offset, mixed $value): void
    {
        $this->values[$offset] = $value;
    }

    public function offsetUnset(mixed $offset): void
    {
        unset($this->values[$offset]);
    }
}

$container = new NestedArrayAccessObject_Container();
$container['store'] = new NestedArrayAccessObject_Store([
    'selected' => 'first',
    'preserved' => 'yes',
]);

$container['store']['selected'] = 'second';

if (!($container['store'] instanceof NestedArrayAccessObject_Store)) {
    echo "FAIL nested object replaced\n";
    exit(1);
}

if ($container['store']['selected'] !== 'second') {
    echo "FAIL nested write\n";
    exit(1);
}

if ($container['store']->get('preserved') !== 'yes') {
    echo "FAIL sibling value lost\n";
    exit(1);
}

echo "PASS\n";
