<?php
function checkSplice($label, $actual, $expected) {
    if ($actual !== $expected) {
        throw new RuntimeException($label . ': ' . json_encode($actual));
    }
}

// Livewire's morph-aware Blade compiler removes a matched stack entry with
// three arguments. Omitting replacement must remove, without inserting null.
$stack = [['type' => 'script'], ['type' => 'style']];
$removed = array_splice($stack, 1, 1);
checkSplice('three arguments', $stack, [['type' => 'script']]);
checkSplice('removed', $removed, [['type' => 'style']]);

$values = [10, 20, 30];
checkSplice('omitted length', array_splice($values, 1), [20, 30]);
checkSplice('remaining', $values, [10]);

$values = [10, 20, 30];
array_splice($values, 1, null, 'replacement');
checkSplice('nullable length and scalar', $values, [10, 'replacement']);

$values = [10, 20, 30];
array_splice($values, 1, 1, null);
checkSplice('explicit null replacement', $values, [10, 30]);

$values = [10, 20];
array_splice(array: $values, offset: 1, length: 0, replacement: ['x', 'y']);
checkSplice('named replacement', $values, [10, 'x', 'y', 20]);

echo "array_splice optional replacement: PASS\n";
