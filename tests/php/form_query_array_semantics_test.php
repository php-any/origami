<?php
parse_str('first=1&repeat=one&repeat=two&items[]=a&items[]=b&nested[name]=origami&nested[0]=zero&a.b=x&empty=&plus=a+b&encoded=%26', $result);
$expected = ['first' => '1', 'repeat' => 'two', 'items' => ['a', 'b'], 'nested' => ['name' => 'origami', 0 => 'zero'], 'a_b' => 'x', 'empty' => '', 'plus' => 'a b', 'encoded' => '&'];
if ($result !== $expected) { throw new Exception('query field semantics'); }
echo "form query arrays OK\n";
$parsed = ['old' => ['keep' => 1]];
$snapshot = $parsed;
parse_str('new=value', $parsed);
if ($parsed !== ['new' => 'value'] || $snapshot !== ['old' => ['keep' => 1]]) { throw new Exception('parse_str output COW'); }
