<?php
$iterator = new DirectoryIterator("tests/php");
foreach ($iterator as $key => $value) {
    echo "key=$key, value=$value\n";
    break;
}
