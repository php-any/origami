<?php
class SyntaxVirtualHook {
    public string $first = 'Ada';
    public string $last = 'Lovelace';
    public string $fullName { get => $this->first . ' ' . $this->last; }
}
$object = new SyntaxVirtualHook;
if ($object->fullName !== 'Ada Lovelace') {
    throw new Exception('virtual get hook result');
}
$object->first = 'Grace';
if ($object->fullName !== 'Grace Lovelace') {
    throw new Exception('virtual get hook must reevaluate');
}
$caught = false;
try { $object->fullName = 'ignored'; } catch (Error $e) { $caught = true; }
if (!$caught) { throw new Exception('get-only virtual property must reject writes'); }
$caught = false;
try { unset($object->fullName); } catch (Error $e) { $caught = str_contains($e->getMessage(), 'Cannot unset hooked property'); }
if (!$caught || $object->fullName !== 'Grace Lovelace') {
    throw new Exception('hooked property must reject unset');
}
echo "virtual get hook OK\n";
