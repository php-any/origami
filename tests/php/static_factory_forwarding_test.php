<?php
class FactoryRoot {
    public static function create(): static { return self::build(); }
    private static function build(): static { return new static; }
}
class FactoryChild extends FactoryRoot {}
if (get_class(FactoryChild::create()) !== FactoryChild::class) { throw new Exception('late static factory'); }
echo "OK\n";
