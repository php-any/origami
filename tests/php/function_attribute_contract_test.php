<?php
function attributeCheck($ok, $message) { if (!$ok) throw new Exception($message); }
#[Attribute(Attribute::TARGET_FUNCTION | Attribute::IS_REPEATABLE)]
class LazyMarker {
    public static int $constructed = 0;
    public function __construct(public string $label) { self::$constructed++; }
}
#[LazyMarker(label: 'first')]
#[LazyMarker('second')]
function markedFunction() {}
attributeCheck(LazyMarker::$constructed === 0, 'lazy declaration');
$attributes = (new ReflectionFunction('markedFunction'))->getAttributes(LazyMarker::class);
attributeCheck(count($attributes) === 2 && LazyMarker::$constructed === 0, 'lazy reflection');
attributeCheck($attributes[0]->getName() === LazyMarker::class, 'name');
attributeCheck($attributes[0]->getArguments() === ['label' => 'first'], 'named arguments');
attributeCheck($attributes[0]->getTarget() === Attribute::TARGET_FUNCTION && $attributes[0]->isRepeated(), 'target and repeat metadata');
attributeCheck($attributes[0]->newInstance()->label === 'first' && LazyMarker::$constructed === 1, 'constructor runs on demand');
class NotAnAttribute {}
#[NotAnAttribute] function unmarkedFunction() {}
try { (new ReflectionFunction('unmarkedFunction'))->getAttributes()[0]->newInstance(); throw new Exception('missing marker validation'); } catch (Error $error) {}
#[Attribute(Attribute::TARGET_CLASS)] class WrongTarget {}
#[WrongTarget] function wrongFunction() {}
try { (new ReflectionFunction('wrongFunction'))->getAttributes()[0]->newInstance(); throw new Exception('missing target validation'); } catch (Error $error) {}
#[Attribute] class SingleMarker {}
#[SingleMarker, SingleMarker] function repeatedFunction() {}
try { (new ReflectionFunction('repeatedFunction'))->getAttributes()[0]->newInstance(); throw new Exception('missing repetition validation'); } catch (Error $error) {}
echo "function attributes OK\n";
