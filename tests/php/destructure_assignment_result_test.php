<?php
function destructureCheck($actual, $expected, $label) {
    if ($actual !== $expected) { throw new Exception($label); }
}
$source = [3, 4];
destructureCheck(([$first, $second] = $source), $source, 'expression preserves array value');
destructureCheck($first, 3, 'first target');
destructureCheck($second, 4, 'second target');
destructureCheck(([$first, $second] = null), null, 'expression preserves null');
destructureCheck($first, null, 'null clears first target');
destructureCheck($second, null, 'null clears second target');
$rows = [['one', 1], ['two', 2]];
$seen = [];
while ([$name, $number] = array_shift($rows)) {
    $seen[] = $name;
    if (count($seen) > 2) { throw new Exception('destructuring loop failed to stop'); }
}
destructureCheck($seen, ['one', 'two'], 'loop stops at null');
destructureCheck($name, null, 'terminal assignment clears value');
echo "destructure assignment result: PASS\n";
