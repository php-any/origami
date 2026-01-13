<?php
$iterator = new DirectoryIterator("tests/php");
$iterator->rewind();
while ($iterator->valid()) {
    echo $iterator->current() . "\n";
    $iterator->next();
}
