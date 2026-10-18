<?php
function &syntaxReferenceGenerator(array &$items): Generator {
    foreach ($items as &$item) {
        yield $item;
    }
}
$items = [2, 4, 6];
$seen = [];
foreach (syntaxReferenceGenerator($items) as &$value) {
    $value *= 3;
    $seen[] = $value;
}
unset($value);
if ($items !== [6, 12, 18] || $seen !== [6, 12, 18]) {
    throw new Exception('by-reference generator must write back to yielded storage');
}
echo "generator reference yield OK\n";
