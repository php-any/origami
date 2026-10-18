<?php
class SyntaxAsymmetricBase {
    public private(set) string $name = 'initial';
    public protected(set) int $count = 1;
    public function rename(string $name): void { $this->name = $name; }
}
class SyntaxAsymmetricChild extends SyntaxAsymmetricBase {
    public function increment(): void { $this->count++; }
    public function forbiddenRename(): void { $this->name = 'child'; }
}
$object = new SyntaxAsymmetricChild;
$object->rename('owner');
$object->increment();
$errors = 0;
try { $object->name = 'outside'; } catch (Error $e) { $errors++; }
try { $object->count = 9; } catch (Error $e) { $errors++; }
try { $object->forbiddenRename(); } catch (Error $e) { $errors++; }
try { $reference =& $object->name; } catch (Error $e) { $errors++; }
if ($errors !== 4 || $object->name !== 'owner' || $object->count !== 2) {
    throw new Exception('asymmetric read/write and reference visibility');
}
echo "asymmetric visibility OK\n";
