<?php
namespace NominalReview;

function verify($actual, $expected, $label) {
    if ($actual !== $expected) { throw new \Exception($label . ': ' . json_encode($actual)); }
}
interface RootContract {}
interface MiddleContract extends RootContract {}
interface LeafContract extends MiddleContract {}
class Base implements LeafContract {
    public function selfChecks() {
        verify(is_a($this, 'nominalreview\\rootcontract'), true, 'this deep interface');
        verify($this instanceof RootContract, true, 'this instanceof');
        verify(is_subclass_of($this, RootContract::class), true, 'this subclass interface');
    }
}
class Child extends Base {}
class Iterator {}
class Generator {}
class Exception extends \Exception {}
interface Throwable {}

$object = new Child();
verify(is_a($object, 'NOMINALREVIEW\\ROOTCONTRACT'), true, 'deep interface');
verify(is_a($object, '\\nominalreview\\base'), true, 'case insensitive parent');
verify(is_subclass_of($object, RootContract::class), true, 'subclass implements');
verify(is_subclass_of(Child::class, Base::class), true, 'string subclass');
verify(is_subclass_of(Child::class, Base::class, false), false, 'disallow string');
verify(is_a(Child::class, RootContract::class, true), true, 'string deep interface');
verify(is_a(Child::class, RootContract::class), false, 'default disallow string');
verify(is_a(LeafContract::class, RootContract::class, true), true, 'interface source');
verify(is_subclass_of(LeafContract::class, RootContract::class), true, 'interface subclass');
verify(is_subclass_of(Child::class, Child::class), false, 'strict same class');
$dynamic = 'Base';
verify($object instanceof $dynamic, false, 'dynamic string is not namespace relative');
$dynamic = 'nominalreview\\base';
verify($object instanceof $dynamic, true, 'dynamic qualified class');
verify($object instanceof RootContract, true, 'deep instanceof');
$classObject = new Base();
verify($object instanceof $classObject, true, 'object RHS denotes its class');
$invalidClass = 1;
try {
    $object instanceof $invalidClass;
    throw new \Exception('invalid instanceof RHS accepted');
} catch (\Error $right) { verify(is_a($right, 'TypeError'), false, 'invalid RHS throws Error'); }
$object->selfChecks();
function requireRoot(RootContract $value): RootContract { return $value; }
verify(requireRoot($object) === $object, true, 'parameter and return interface');
try { requireRoot(new \stdClass()); throw new \Exception('wrong nominal parameter accepted'); }
catch (\TypeError $right) {}
function badReturn(): RootContract { return new \stdClass(); }
try { badReturn(); throw new \Exception('wrong nominal return accepted'); }
catch (\TypeError $right) {}
function requireClosure(\Closure $value) { return true; }
verify(requireClosure(function () {}), true, 'real closure accepted');
foreach (['strlen', [Base::class, 'selfChecks']] as $callable) {
    try { requireClosure($callable); throw new \Exception('callable accepted as Closure'); }
    catch (\TypeError $right) {}
}
verify(class_alias(Base::class, __NAMESPACE__ . '\\BaseAlias'), true, 'alias registration');
verify(is_a(new Base(), __NAMESPACE__ . '\\BaseAlias'), true, 'target alias');
verify(is_a(__NAMESPACE__ . '\\BaseAlias', Base::class, true), true, 'source alias');
verify(is_subclass_of(__NAMESPACE__ . '\\BaseAlias', Base::class), false, 'alias same identity');
verify(is_subclass_of(Base::class, __NAMESPACE__ . '\\BaseAlias'), false, 'target alias same identity');
verify(is_subclass_of($object, __NAMESPACE__ . '\\BaseAlias'), true, 'subclass through alias');

$autoload = [];
spl_autoload_register(function ($name) use (&$autoload) { $autoload[] = $name; });
$missing = 'NominalUnknownTarget';
verify($object instanceof $missing, false, 'unknown instanceof');
verify(is_a($object, $missing), false, 'unknown is_a target');
verify(is_subclass_of($object, $missing), false, 'unknown subclass target');
verify(is_a('NominalUnknownSource', Base::class), false, 'disallow autoload source');
verify(is_subclass_of('NominalUnknownSource', Base::class, false), false, 'disallow subclass autoload');
verify(count($autoload), 0, 'target never autoloads');
verify(is_a('NominalUnknownSource', Base::class, true), false, 'allowed missing source');
verify($autoload, ['NominalUnknownSource'], 'autoload only source');
spl_autoload_register(function ($name) {
    if ($name === 'NominalThrowingSource') { throw new \RuntimeException('autoload failure'); }
});
try {
    is_a('NominalThrowingSource', Base::class, true);
    throw new \Exception('autoload exception was swallowed');
} catch (\RuntimeException $right) { verify($right->getMessage(), 'autoload failure', 'autoload propagation'); }

function requireIterable(iterable $value) { return true; }
foreach ([new Iterator(), new Generator()] as $fake) {
    try {
        requireIterable($fake);
        throw new \Exception('namespaced iterator accepted as iterable');
    } catch (\TypeError $right) {}
}

try { usleep(-1); }
catch (Exception $wrong) { throw new \Exception('namespace exception captured Error'); }
catch (\Exception $wrong) { throw new \Exception('Exception captured Error'); }
catch (Throwable $wrong) { throw new \Exception('namespace throwable captured Error'); }
catch (\Error $right) {
    verify(is_a($right, 'valueerror'), true, 'internal valueerror');
    verify($right instanceof \Throwable, true, 'internal throwable');
    verify(is_a($right, '\\Exception'), false, 'error not exception');
}
try { throw new \BadMethodCallException('test'); }
catch (\BadFunctionCallException $right) { verify(true, true, 'exception parent'); }
echo "nominal type relations OK\n";
