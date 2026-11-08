#!/usr/bin/env php
<?php
$source = file_get_contents(__FILE__);
$tail = substr($source, __COMPILER_HALT_OFFSET__);
// Worktree checkout may use CRLF on Windows; the halt offset is byte-accurate either way.
$normalized = str_replace("\r\n", "\n", $tail);
$expected = "ARCHIVE_PAYLOAD\n<?php declare(strict_types=1);\n";
if ($normalized !== $expected) {
    throw new \Exception('Incorrect compiler halt byte offset: ' . bin2hex($tail));
}
echo "halt compiler uppercase passed\n";
__HALT_COMPILER();ARCHIVE_PAYLOAD
<?php declare(strict_types=1);
