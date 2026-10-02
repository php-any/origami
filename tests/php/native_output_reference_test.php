<?php
declare(strict_types=1);
function outputCheck($actual, $expected, $label) {
    if ($actual !== $expected) { throw new Exception($label); }
}
foreach ([null, 42, 'old', []] as $initial) {
    $matches = $initial;
    outputCheck(preg_match('/a/', 'a', $matches), 1, 'preg_match result');
    outputCheck($matches, ['a'], 'preg_match output');
    $matches = $initial;
    outputCheck(preg_match_all('/a/', 'aa', $matches), 2, 'preg_match_all result');
    outputCheck($matches, [['a', 'a']], 'preg_match_all output');
    $result = $initial;
    parse_str('a=1', $result);
    outputCheck($result, ['a' => '1'], 'parse_str output');
    $count = $initial;
    outputCheck(str_ireplace('x', 'y', 'Xx', $count), 'yy', 'str_ireplace result');
    outputCheck($count, 2, 'str_ireplace output');
}
outputCheck(preg_match('/a/', 'a'), 1, 'optional preg_match output');
outputCheck(preg_match_all('/a/', 'aa'), 2, 'optional preg_match_all output');
preg_match_all('/a/', 'aa', $undefinedMatches);
outputCheck($undefinedMatches, [['a', 'a']], 'undefined reference receives output');
echo "native output reference: PASS\n";
