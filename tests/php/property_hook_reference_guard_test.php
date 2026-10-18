<?php
class HookReferenceGuard {
    public int $writes = 0;
    public int $value = 1 {
        set { $this->writes++; $this->value = $value * 2; }
    }
    public array $items = [1] {
        set { $this->writes++; $this->items = $value; }
    }
    public int $plain = 1;
    public array $storage = [1];
    public array $virtualItems { get => $this->storage; }
    public function updateInside(): void { $this->value = 4; }
}
$object = new HookReferenceGuard;
$caught = false;
try { $reference =& $object->value; } catch (Error $e) { $caught = true; }
if (!$caught || $object->value !== 1 || $object->writes !== 0) {
    throw new Exception('hooked property reference must be rejected before mutation');
}
$caught = false;
try { $object->items[] = 9; } catch (Error $e) { $caught = true; }
if (!$caught || $object->items !== [1] || $object->writes !== 0) {
    throw new Exception('indirect write must not bypass setter');
}
$object->virtualItems[] = 9;
if ($object->virtualItems !== [1]) { throw new Exception('virtual getter result must not create backing storage'); }
$object->storage = [[1]];
$object->virtualItems[0][0] = 9;
if ($object->storage !== [[1]]) { throw new Exception('virtual getter arrays must retain nested value semantics'); }
$external = 1;
$object->storage = [&$external];
$object->virtualItems[0] = 9;
if ($external !== 9) { throw new Exception('virtual getter array copy must preserve explicit element references'); }
$object->updateInside();
if ($object->value !== 8 || $object->writes !== 1) {
    throw new Exception('hook must still write its own backing value');
}
$reference =& $object->plain;
$reference = 5;
if ($object->plain !== 5) { throw new Exception('ordinary references must remain live'); }
echo "property hook reference guard OK\n";
