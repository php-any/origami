<?php
class StringScopeParent {
    private const LABEL = 'parent';
    private function suffix(): string { return 'private'; }
    public function __toString(): string {
        return self::LABEL . ':' . $this->suffix() . ':' . static::class;
    }
}
class StringScopeChild extends StringScopeParent {
    public function suffix(): string { return 'child'; }
}
function consumeScopeString(string $value): string { return $value; }
$object = new StringScopeChild();
$expected = 'parent:private:StringScopeChild';
if (consumeScopeString($object) !== $expected) { throw new RuntimeException('parameter string scope'); }
function returnScopeString(): string { return new StringScopeChild(); }
if (returnScopeString() !== $expected) { throw new RuntimeException('return string scope'); }
class StringScopeHolder { public string $value; }
$holder = new StringScopeHolder();
$holder->value = $object;
if ($holder->value !== $expected) { throw new RuntimeException('property string scope'); }
if (strval($object) !== $expected) { throw new RuntimeException('strval string scope'); }
if (strtr($object, ['parent' => 'base']) !== 'base:private:StringScopeChild') { throw new RuntimeException('strtr string scope'); }
class ThrowingScopeString {
    public function __toString(): string { throw new RuntimeException('string conversion failure'); }
}
try { consumeScopeString(new ThrowingScopeString()); throw new RuntimeException('missing conversion exception'); }
catch (RuntimeException $e) { if ($e->getMessage() !== 'string conversion failure') { throw $e; } }
echo "string coercion scope OK\n";
