<?php
function initializeCheck($condition, $message) {
    if (!$condition) throw new Exception($message);
}
class InitializationBase {
    public int $value = 1;
    public array $items = [1, 2];
    private int $privateValue = 10;
    public function baseValue() { return $this->privateValue; }
}
class InitializationMiddle extends InitializationBase {
    public int $value = 2;
    private int $privateValue = 20;
    public function middleValue() { return $this->privateValue; }
}
class InitializationLeaf extends InitializationMiddle {}
class InitializationTypedLeaf extends InitializationMiddle {
    public int $value;
}
class_alias(InitializationMiddle::class, 'InitializationAlias');
class InitializationAliasLeaf extends InitializationAlias {}
for ($i = 0; $i < 100; $i++) {
    $first = new InitializationLeaf;
    $second = new InitializationLeaf;
    initializeCheck($first->value === 2, 'nearest declaration default');
    initializeCheck($first->baseValue() === 10 && $first->middleValue() === 20, 'distinct private declarations');
    $first->items[0] = 99;
    initializeCheck($second->items === [1, 2], 'fresh array defaults on each object');
    initializeCheck((new InitializationAliasLeaf)->value === 2, 'parent class alias');
    $typed = new InitializationTypedLeaf;
    initializeCheck(!isset($typed->value), 'typed redeclaration hides parent default');
    try { $typed->value; throw new Exception('typed property must remain uninitialized'); }
    catch (Error $error) {}
}
echo "Inherited initialization OK\n";
