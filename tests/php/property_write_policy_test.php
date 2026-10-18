<?php
class WritePolicyPlain { public int $value = 1; }
class WritePolicyHooked {
    public int $writes = 0;
    public int $value {
        set { $this->writes++; $this->value = $value * 2; }
    }
}
trait WritePolicyRestricted {
    public private(set) int $value = 1;
    public function update(int $value): void { $this->value = $value; }
}
class WritePolicyOwner { use WritePolicyRestricted; }
class WritePolicyChild extends WritePolicyOwner {}
function writePolicyAssign(object $object, int $value): void { $object->value = $value; }
$plain = new WritePolicyPlain;
$hooked = new WritePolicyHooked;
$restricted = new WritePolicyChild;
for ($i = 2; $i < 5; $i++) {
    writePolicyAssign($plain, $i);
    writePolicyAssign($hooked, $i);
    $restricted->update($i);
    $caught = false;
    try { writePolicyAssign($restricted, 100); } catch (Error $e) { $caught = true; }
    if (!$caught || $plain->value !== $i || $hooked->value !== $i * 2 || $restricted->value !== $i) {
        throw new Exception('write policy must follow the actual class and declaring scope');
    }
}
if ($hooked->writes !== 3) { throw new Exception('setter must run exactly once per assignment'); }
$reference =& $plain->value;
$reference = 9;
if ($plain->value !== 9) { throw new Exception('ordinary property references must remain live'); }
$caught = false;
try { $reference =& $restricted->value; } catch (Error $e) { $caught = true; }
if (!$caught) { throw new Exception('private setter must reject an external reference'); }
$caught = false;
try { unset($restricted->value); } catch (Error $e) { $caught = true; }
if (!$caught || $restricted->value !== 4) { throw new Exception('private setter must reject external unset: ' . var_export([$caught, $restricted->value], true)); }
echo "property write policy OK\n";
