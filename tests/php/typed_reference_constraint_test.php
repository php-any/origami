<?php
function referenceCheck($actual, $expected) {
    if ($actual !== $expected) { throw new Exception('typed reference value: ' . json_encode($actual)); }
}
class ReferenceProperties {
    public int $number = 1;
    public int|string $flexible = 1;
    public string $text = 'one';
}
$object = new ReferenceProperties();
$alias =& $object->number;
$alias = '2';
referenceCheck([$alias, $object->number], [2, 2]);
try { $alias = []; throw new Exception('typed alias accepted array'); }
catch (TypeError $expected) {}
referenceCheck($object->number, 2);
$object->flexible =& $alias;
try { $alias = '3'; throw new Exception('incompatible conversions accepted'); }
catch (TypeError $expected) {}
referenceCheck([$alias, $object->number, $object->flexible], [2, 2, 2]);
$alias = 3;
referenceCheck([$object->number, $object->flexible], [3, 3]);
unset($object->number);
$alias = 'text';
referenceCheck($object->flexible, 'text');
function writeReference(&$value) { $value = []; }
$text =& $object->text;
try { writeReference($text); throw new Exception('by-ref call removed property constraint'); }
catch (TypeError $expected) {}
referenceCheck($object->text, 'one');
class ReferenceSelf {
    public self $value;
}
class ReferenceChild extends ReferenceSelf {}
$self = new ReferenceChild();
$self->value = new ReferenceSelf();
$selfAlias =& $self->value;
try { $selfAlias = new stdClass(); throw new Exception('self reference accepted wrong class'); }
catch (TypeError $expected) {}
echo "typed reference constraints OK\n";
