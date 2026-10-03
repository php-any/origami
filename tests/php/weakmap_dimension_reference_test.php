<?php
$map = new WeakMap;
$key = new stdClass;
$map[$key] = ['items' => []];
$map[$key]['items'][] = 3;
$map[$key]['items'][] = 4;
$copy = $map->offsetGet($key);
$copy['items'][] = 9;
if ($map[$key] !== ['items' => [3, 4]]) { throw new Exception('WeakMap dimension or method copy'); }
$reference =& $map[$key];
$reference['value'] = 6;
if ($map[$key]['value'] !== 6 || (new ReflectionMethod(WeakMap::class, 'offsetGet'))->returnsReference()) { throw new Exception('WeakMap dimension reference metadata'); }
unset($map[$key]);
$reference['value'] = 7;
if (isset($map[$key]) || $reference['value'] !== 7) { throw new Exception('WeakMap detached slot'); }
echo "WeakMap dimension reference OK\n";
