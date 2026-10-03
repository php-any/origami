<?php
function ref_metadata_assert($ok,$message) { if (!$ok) { throw new Exception($message); } }
function &ref_metadata_function(&$value) { return $value; }
class RefMetadataOwner { public function &byRef(&$value) { return $value; } public function byValue($value) { return $value; } }
ref_metadata_assert((new ReflectionFunction('ref_metadata_function'))->returnsReference(), 'function reference metadata');
$closure = function &(&$value) { return $value; };
ref_metadata_assert((new ReflectionFunction($closure))->returnsReference(), 'closure reference metadata');
ref_metadata_assert((new ReflectionMethod(RefMetadataOwner::class,'byRef'))->returnsReference(), 'method reference metadata');
ref_metadata_assert(!(new ReflectionMethod(RefMetadataOwner::class,'byValue'))->returnsReference(), 'value method metadata');
$value=4; $array=[''=>&$value, -3=>&$value, 0=>&$value];
$empty=ReflectionReference::fromArrayElement($array,'');
$negative=ReflectionReference::fromArrayElement($array,-3);
$zero=ReflectionReference::fromArrayElement($array,0);
ref_metadata_assert(strlen($empty->getId()) === 20 && $empty->getId() === $negative->getId() && $negative->getId() === $zero->getId(), 'reference key identities');
$caught=false;
try { ReflectionReference::fromArrayElement($array,'missing'); } catch (ReflectionException $error) { $caught=true; }
ref_metadata_assert($caught, 'missing reference key exception');
$caught=false;
try { ReflectionReference::fromArrayElement($array,'0'); } catch (ReflectionException $error) { $caught=true; }
ref_metadata_assert($caught, 'reflection must preserve literal key kind');
echo "reference reflection matrix ok\n";
interface InterfaceReferenceMetadata { public static function &reference(array &$value): array; }
$interfaceMethod = new ReflectionMethod(InterfaceReferenceMetadata::class, 'reference');
ref_metadata_assert($interfaceMethod->isStatic() && $interfaceMethod->isAbstract() && $interfaceMethod->returnsReference(), 'interface static reference flags');
