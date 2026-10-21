<?php
class ReflectionStaticProbe { private static array $values = ['before']; public static function values() { return self::$values; } }
$p = new ReflectionProperty(ReflectionStaticProbe::class, 'values');
$p->setValue(null, ['after']);
if ($p->getValue() !== ['after'] || ReflectionStaticProbe::values() !== ['after']) { throw new Exception('Static reflection write failed'); }
try { $p->setValue(null, 123); throw new Exception('Property type lost'); } catch (TypeError $e) {}
echo "static property reflection passed\n";
