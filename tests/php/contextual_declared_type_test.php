<?php
function contextCheck($actual, $expected, $label) {
    if ($actual !== $expected) { throw new Exception($label); }
}
function contextError($callback, $label) {
    try { $callback(); } catch (TypeError $error) { return; }
    throw new Exception('missing contextual TypeError: ' . $label);
}
class TypeScopeParent {
    public function accepts(self $value): self { return $value; }
    public static function acceptsStatic(self $value): self { return $value; }
    public function parentResult(): self { return new TypeScopeParent(); }
    public function lateResult(): static { return new TypeScopeParent(); }
    public function callback() { return function (self $value): self { return $value; }; }
}
class TypeScopeChild extends TypeScopeParent {
    public function parentParameter(parent $value): parent { return $value; }
}
$parent = new TypeScopeParent();
$child = new TypeScopeChild();
contextCheck($child->accepts($parent), $parent, 'inherited self parameter and return');
contextCheck(TypeScopeChild::acceptsStatic($parent), $parent, 'static inherited self parameter');
contextCheck($child->parentResult() instanceof TypeScopeParent, true, 'self return uses declaration scope');
contextError(fn() => $child->lateResult(), 'static return requires called class');
contextCheck($child->parentParameter($parent), $parent, 'parent parameter and return');
$callback = $child->callback();
contextCheck($callback($parent), $parent, 'closure declaration scope');
contextCheck(array_map($callback, [$parent])[0], $parent, 'native callback declaration scope');
echo "contextual declared type: PASS\n";
