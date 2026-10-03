<?php
function reference_named_keys(&$first, &...$rest) { $first++; foreach ($rest as &$item) { $item++; } return array_keys($rest); }
$first = 1;
$empty = 2;
$arguments = ['first' => &$first, '' => &$empty];
if (reference_named_keys(...$arguments) !== [''] || $first !== 2 || $empty !== 3) { throw new Exception('empty named key lost'); }
try { reference_named_keys(...['first' => &$first, '' => &$empty, 0 => &$empty]); throw new Exception('positional argument accepted'); } catch (Error $e) {}
function reference_named_fixed(&$value) {}
try { reference_named_fixed(...['' => &$empty]); throw new Exception('unknown empty parameter accepted'); } catch (Error $e) {}
echo "reference named keys OK\n";
