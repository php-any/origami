<?php
$events = [];
function dimensionKey($key) { global $events; $events[] = $key; return $key; }
function dimensionRhs() { global $events; $events[] = 'rhs'; return 7; }
$array = [];
$array[dimensionKey('outer')][dimensionKey('inner')] = dimensionRhs();
if ($events !== ['outer', 'inner', 'rhs'] || $array !== ['outer' => ['inner' => 7]]) { throw new Exception('complex evaluation'); }
class DimensionReferenceBox implements ArrayAccess {
    public array $value = ['inner' => 1];
    public array $calls = [];
    public function offsetExists($key): bool { $this->calls[] = 'exists'; return true; }
    public function &offsetGet($key): mixed { $this->calls[] = 'get'; return $this->value; }
    public function offsetSet($key, $value): void { $this->calls[] = 'set'; $this->value = $value; }
    public function offsetUnset($key): void { $this->calls[] = 'unset'; }
}
$box = new DimensionReferenceBox;
$snapshot = $box->value;
$box['outer']['inner'] = 9;
if ($box->value !== ['inner' => 9] || $snapshot !== ['inner' => 1] || $box->calls !== ['get']) { throw new Exception('reference dimension write'); }
$ref =& $box['outer']['inner'];
$ref = 11;
if ($box->value['inner'] !== 11 || $box->calls !== ['get', 'get']) { throw new Exception('nested dimension reference'); }
$events = [];
$array[dimensionKey('outer')][dimensionKey('inner')] += dimensionRhs();
if ($events !== ['outer', 'inner', 'rhs'] || $array['outer']['inner'] !== 14) { throw new Exception('compound evaluation'); }
$events = [];
$old = $array[dimensionKey('outer')][dimensionKey('inner')]++;
$new = ++$array[dimensionKey('outer')][dimensionKey('inner')];
if ($old !== 14 || $new !== 16 || $events !== ['outer', 'inner', 'outer', 'inner']) { throw new Exception('increment evaluation'); }
$box['outer']['inner'] += 2;
$old = $box['outer']['inner']++;
if ($old !== 13 || $box->value['inner'] !== 14 || $box->calls !== ['get', 'get', 'get', 'get']) { throw new Exception('reference compound dimension'); }
echo "complex dimensions OK\n";
$events = [];
if (!isset($array[dimensionKey('outer')][dimensionKey('inner')]) || $events !== ['outer', 'inner']) { throw new Exception('isset evaluation'); }
$events = [];
if (empty($array[dimensionKey('outer')][dimensionKey('inner')]) || $events !== ['outer', 'inner']) { throw new Exception('empty evaluation'); }
$events = [];
$value = $array[dimensionKey('outer')][dimensionKey('inner')] ?? dimensionRhs();
if ($value !== 16 || $events !== ['outer', 'inner']) { throw new Exception('coalesce evaluation'); }
$events = [];
$array[dimensionKey('outer')][dimensionKey('inner')] ??= dimensionRhs();
if ($events !== ['outer', 'inner'] || $array['outer']['inner'] !== 16) { throw new Exception('coalesce assign evaluation'); }
$events = [];
$array[dimensionKey('outer')][dimensionKey('missing')] ??= dimensionRhs();
if ($events !== ['outer', 'missing', 'rhs'] || $array['outer']['missing'] !== 7) { throw new Exception('coalesce assign missing'); }
$events = [];
$snapshot = $array;
unset($array[dimensionKey('outer')][dimensionKey('inner')]);
if ($events !== ['outer', 'inner'] || isset($array['outer']['inner']) || $snapshot['outer']['inner'] !== 16) { throw new Exception('unset evaluation COW'); }
$missing = [];
unset($missing[dimensionKey('none')][dimensionKey('child')]);
if ($missing !== []) { throw new Exception('unset must not vivify'); }
$box->calls = [];
unset($box['outer']['inner']);
if ($box->calls !== ['get'] || $box->value !== [] || $ref !== 14) { throw new Exception('reference dimension unset'); }
$events = [];
$plain = null;
$plain ??= dimensionRhs();
$plain ??= dimensionRhs();
if ($plain !== 7 || $events !== ['rhs']) { throw new Exception('plain coalesce assignment'); }
