<?php
namespace NominalCase;
class ÄClass {}
class äClass {}
interface ÄContract {}
interface äContract {}
class OnlyUpper implements ÄContract {}
$upper = new ÄClass();
$lower = new äClass();
if (is_a($upper, äClass::class) || is_a($lower, ÄClass::class)) { throw new \Exception('non-ASCII class names folded'); }
if (!is_a($upper, 'NominalCase\\ÄCLASS') || !is_a($lower, 'nominalcase\\äclass')) { throw new \Exception('ASCII case not folded'); }
if (is_a(new OnlyUpper(), äContract::class) || !is_a(new OnlyUpper(), 'nominalcase\\Äcontract')) { throw new \Exception('interface case rules'); }
echo "nominal identifier case OK\n";
