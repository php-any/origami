<?php

$arr = ['/tmp/views'];
if (!array_key_exists('0', $arr)) {
    echo "FAIL array_key_exists\n";
    exit(1);
}
$v = $arr['0'];
if ($v !== '/tmp/views') {
    echo "FAIL read: ";
    var_export($v);
    echo "\n";
    exit(1);
}

echo "PASS\n";
