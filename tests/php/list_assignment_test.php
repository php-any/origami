<?php
$calls = 0;
function listSource() {
    global $calls;
    ++$calls;
    return ['left', 'skip', ['right', 42]];
}
$result = (list($a, , list($b, $c)) = listSource());
if ($a !== 'left' || $b !== 'right' || $c !== 42 || $calls !== 1 || $result !== ['left', 'skip', ['right', 42]]) {
    throw new \Exception('list assignment failed');
}
list($version, $mode) = explode('|', '3.1|off');
if ($version !== '3.1' || $mode !== 'off') {
    throw new \Exception('list explode failed');
}
echo "list assignment passed\n";
