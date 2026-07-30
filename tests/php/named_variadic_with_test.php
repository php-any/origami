<?php

class Rec
{
    public function __construct(
        public readonly string $a,
        public readonly string $b = 'default',
        public array $extra = [],
    ) {}

    public function with(mixed ...$args): self
    {
        foreach (['a', 'b', 'extra'] as $prop) {
            $args[$prop] ??= $this->{$prop};
        }
        return new self(...$args);
    }
}

$r = new Rec('hello');
$r2 = $r->with(extra: ['x' => 1, 'y' => 2]);
echo 'a=' . $r2->a . "\n";
echo 'b=' . $r2->b . "\n";
echo 'extra=' . json_encode($r2->extra) . "\n";
if ($r2->a !== 'hello' || $r2->b !== 'default' || $r2->extra !== ['x' => 1, 'y' => 2]) {
    echo "FAIL\n";
    exit(1);
}
echo "OK named_variadic_with\n";
