<?php
class ReadonlyMatrixValue {
    public readonly int $value;
    public readonly array $items;
    public function __construct() { $this->value = 1; $this->items = ['count' => 1]; }
    public function change() { $this->value = 2; }
}
$object = new ReadonlyMatrixValue;
$errors = 0;
try { $object->value = 2; } catch (Error $e) { $errors++; }
try { $object->change(); } catch (Error $e) { $errors++; }
try { $object->items['count'] = 2; } catch (Error $e) { $errors++; }
try { $reference =& $object->value; } catch (Error $e) { $errors++; }
try { (new ReflectionProperty($object, 'value'))->setValue($object, 2); } catch (Error $e) { $errors++; }
if ($errors !== 5 || $object->value !== 1 || $object->items !== ['count' => 1]) { throw new Exception('readonly mutation'); }
readonly class ReadonlyImplicitMatrix {
    public int $value;
    public function __construct() { $this->value = 3; }
}
$implicit = new ReadonlyImplicitMatrix;
try { $implicit->value = 4; } catch (Error $e) { $errors++; }
try { $implicit->extra = 4; } catch (Error $e) { $errors++; }
if ($errors !== 7 || $implicit->value !== 3) { throw new Exception('readonly implicit or dynamic property'); }
echo "readonly matrix OK\n";
class ReadonlyCloneMatrix {
    public function __construct(public readonly int $value = 3) {}
    public function __clone() { $this->value = $this->value + 1; try { $this->value = 9; throw new Exception('second clone initialization allowed'); } catch (Error $e) {} }
}
$original = new ReadonlyCloneMatrix;
$copy = clone $original;
if ($original->value !== 3 || $copy->value !== 4) { throw new Exception('readonly clone initialization'); }
try { $original->__construct(8); throw new Exception('promoted readonly reinitialization allowed'); } catch (Error $e) {}
