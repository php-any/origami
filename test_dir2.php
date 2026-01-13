<?php
$iterator = new DirectoryIterator("tests/php");
$iterator->rewind();
echo "First item:\n";
echo "  current(): " . $iterator->current() . "\n";
echo "  key(): " . $iterator->key() . "\n";
echo "  getFilename(): " . $iterator->getFilename() . "\n";
echo "  valid(): " . ($iterator->valid() ? "true" : "false") . "\n";
