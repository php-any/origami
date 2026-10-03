<?php
class TypedInitializationParent { private int $shadow; }
class TypedInitializationChild extends TypedInitializationParent { public $shadow = null; public int $value; public int $default = 7; }
$object = new TypedInitializationChild;
$property = new ReflectionProperty(TypedInitializationParent::class, 'shadow');
if ($property->isInitialized($object)) { throw new Exception('private initialization shadow'); }
$errors = 0;
try { $unused = $object->value; } catch (Error $e) { $errors++; }
unset($object->default);
try { $unused = $object->default; } catch (Error $e) { $errors++; }
if ($errors !== 2 || (new ReflectionProperty($object, 'default'))->isInitialized($object)) { throw new Exception('typed unset initialization'); }
$fresh = (new ReflectionClass(TypedInitializationChild::class))->newInstanceWithoutConstructor();
if ($fresh->default !== 7 || !(new ReflectionProperty($fresh, 'default'))->isInitialized($fresh)) { throw new Exception('constructorless defaults'); }
abstract class MethodFlagParent { abstract public function required(): int; final public function locked(): int { return 1; } }
class MethodFlagChild extends MethodFlagParent { public function required(): int { return 2; } }
$final = new ReflectionMethod(MethodFlagChild::class, 'locked');
$abstract = new ReflectionMethod(MethodFlagParent::class, 'required');
if (!$final->isFinal() || $final->isAbstract() || $final->getModifiers() !== 33 || !$abstract->isAbstract() || $abstract->getModifiers() !== 65) { throw new Exception('method declaration flags'); }
echo "typed property and method metadata OK\n";
class TypedDimensionInitialization { public array $array; public int $scalar; }
$dimension = new TypedDimensionInitialization;
$dimension->array['one']['two'] = 3;
if ($dimension->array !== ['one' => ['two' => 3]]) { throw new Exception('typed array dimension initialization'); }
try { $dimension->scalar['invalid'] = 1; throw new Exception('scalar dimension allowed'); } catch (TypeError $e) {}
