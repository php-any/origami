<?php
trait SyntaxOverrideOne { public function value(): string { return 'one'; } }
trait SyntaxOverrideTwo { public function value(): string { return 'two'; } }
class SyntaxOverrideTarget {
    use SyntaxOverrideOne, SyntaxOverrideTwo {
        SyntaxOverrideTwo::value insteadof SyntaxOverrideOne;
        SyntaxOverrideOne::value as private;
        SyntaxOverrideOne::value as original;
    }
    public function value(): string { return 'class'; }
}
$object = new SyntaxOverrideTarget;
if ($object->value() !== 'class' || $object->original() !== 'one') {
    throw new Exception('class override wins; renamed alias keeps trait implementation');
}
echo "trait class override OK\n";
