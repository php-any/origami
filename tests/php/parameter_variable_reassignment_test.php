<?php
function reassignment(int $value): array { $value = ['ok']; return $value; }
function referenceReassignment(int &$value) { $value = 'changed'; }
class ReassignmentCase {
    public function __construct(public int $value) { $value = ['local']; }
    public function change(int $value): string { $value = 'changed'; return $value; }
}
if (reassignment(4) !== ['ok']) { throw new Exception('ordinary parameter local'); }
$value = 1;
referenceReassignment($value);
if ($value !== 'changed') { throw new Exception('reference parameter local'); }
$object = new ReassignmentCase(3);
if ($object->value !== 3 || $object->change(2) !== 'changed') { throw new Exception('promoted parameter local'); }
$callback = function (int $value): string { $value = 'changed'; return $value; };
if ($callback(2) !== 'changed') { throw new Exception('closure parameter local'); }
try { reassignment([]); throw new Exception('missing argument TypeError'); } catch (TypeError $error) {}
echo "parameter variable reassignment: PASS\n";
