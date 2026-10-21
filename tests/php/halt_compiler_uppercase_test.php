#!/usr/bin/env php
<?php
$source = file_get_contents(__FILE__);
$tail = substr($source, __COMPILER_HALT_OFFSET__);
if ($tail !== "ARCHIVE_PAYLOAD\n<?php declare(strict_types=1);\n") {
    throw new \Exception('Incorrect compiler halt byte offset: ' . bin2hex($tail));
}
echo "halt compiler uppercase passed\n";
__HALT_COMPILER();ARCHIVE_PAYLOAD
<?php declare(strict_types=1);
