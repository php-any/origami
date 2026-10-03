<?php
function prop_ref_assert($ok, $message) { if (!$ok) { throw new Exception($message); } }
class StaticRefOwner {
    public static int $value = 4;
    public static array $stack = [];
    public static function push($value) { array_push(static::$stack, $value); }
}
StaticRefOwner::push('a');
prop_ref_assert(StaticRefOwner::$stack === ['a'], 'static array reference argument');
$static =& StaticRefOwner::$value;
$static = 7;
prop_ref_assert(StaticRefOwner::$value === 7, 'static property alias');
StaticRefOwner::$value = 9;
prop_ref_assert($static === 9, 'static direct write lost alias');
$caught = false;
try { $static = []; } catch (TypeError $error) { $caught = true; }
prop_ref_assert($caught && $static === 9, 'static alias type constraint');
$replacement = 12;
StaticRefOwner::$value =& $replacement;
$replacement = 16;
prop_ref_assert(StaticRefOwner::$value === 16 && $static === 9, 'static rebind lost old reference');
$static = [];
prop_ref_assert($static === [], 'static rebind retained old property constraint');
class MagicRefOwner {
    private array $values = ['x' => 3];
    public function &__get($name) { return $this->values[$name]; }
}
$magic = new MagicRefOwner;
$copy = $magic->x;
$copy = 8;
prop_ref_assert($magic->x === 3, 'magic ordinary read kept reference');
$alias =& $magic->x;
$alias = 11;
prop_ref_assert($magic->x === 11, 'magic reference read');
$key = 'x'; $dynamic =& $magic->$key; $dynamic = 14;
prop_ref_assert($magic->x === 14, 'dynamic magic reference read');
class OffsetRefOwner implements ArrayAccess {
    private array $values = [5];
    public function offsetExists($offset): bool { return isset($this->values[$offset]); }
    public function &offsetGet($offset): mixed { return $this->values[$offset]; }
    public function offsetSet($offset, $value): void { $this->values[$offset] = $value; }
    public function offsetUnset($offset): void { unset($this->values[$offset]); }
}
$offset = new OffsetRefOwner;
$copy = $offset[0]; $copy = 20;
prop_ref_assert($offset[0] === 5, 'offset ordinary read kept reference');
$index = 0; $alias =& $offset[$index++]; $alias = 24;
prop_ref_assert($index === 1 && $offset[0] === 24, 'offset reference or evaluation count');
echo "property reference identity ok\n";
