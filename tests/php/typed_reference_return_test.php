<?php
function &typedReference(): int {
    static $value = '4';
    return $value;
}
$alias =& typedReference();
if ($alias !== 4) { throw new Exception('typed return coercion'); }
$alias = 7;
if (typedReference() !== 7) { throw new Exception('typed return reference identity'); }
function &invalidTypedReference(): int {
    static $value = [];
    return $value;
}
try { $bad =& invalidTypedReference(); throw new Exception('missing return TypeError'); } catch (TypeError $error) {}
echo "typed reference return: PASS\n";
