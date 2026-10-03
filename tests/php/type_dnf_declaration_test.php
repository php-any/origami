<?php
interface DnfLeft {}
interface DnfRight {}
interface DnfOther {}
class DnfBoth implements DnfLeft, DnfRight {}
class DnfAlternative implements DnfOther {}
function dnfValue((DnfLeft&DnfRight)|DnfOther|null $value): (DnfLeft&DnfRight)|DnfOther|null { return $value; }
class DnfProperty {
    public (DnfLeft&DnfRight)|DnfOther|null $value;
}
$property = new DnfProperty();
foreach ([new DnfBoth(), new DnfAlternative(), null] as $value) {
    $property->value = dnfValue($value);
    if ($property->value !== $value) { throw new Exception('DNF identity'); }
}
try { dnfValue(new stdClass()); throw new Exception('DNF accepted invalid value'); }
catch (TypeError $expected) {}
try { $property->value = new stdClass(); throw new Exception('DNF property accepted invalid value'); }
catch (TypeError $expected) {}
echo "DNF declarations OK\n";
