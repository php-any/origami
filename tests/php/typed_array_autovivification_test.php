<?php
class ArrayVivification {
    public array $colors;
    public function fill() { $this->colors['red'][500] = 'ok'; }
}
$object = new ArrayVivification();
$object->fill();
function actualArray(array $value): array { return $value; }
if (actualArray($object->colors) !== ['red' => [500 => 'ok']]) { throw new Exception('typed property array vivification'); }
$empty = null;
$empty['key'] = 'value';
if (actualArray($empty) !== ['key' => 'value']) { throw new Exception('null local array vivification'); }
echo "typed array autovivification: PASS\n";
