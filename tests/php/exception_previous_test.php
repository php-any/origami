<?php

$inner = new Exception('inner');
$outer = new RuntimeException('outer', 0, $inner);

try {
    throw $outer;
} catch (Throwable $caught) {
    $previous = $caught->getPrevious();

    if (!($previous instanceof Exception)) {
        echo "FAIL previous type\n";
        exit(1);
    }

    if ($previous->getMessage() !== 'inner') {
        echo "FAIL previous message\n";
        exit(1);
    }
}

echo "PASS\n";
