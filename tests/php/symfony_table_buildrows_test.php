<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Helper\TableSeparator;
use Symfony\Component\Console\Output\BufferedOutput;

// Minimal reproduction of render loop pieces
$headers = [['ID', 'Name']];
$rows = [['1', 'Alice']];
$divider = new TableSeparator();
$all = array_merge($headers, [$divider], $rows);

Log::info('=== identity check ===');
foreach ($all as $i => $row) {
    $same = ($divider === $row);
    $isSep = ($row instanceof TableSeparator);
    Log::info("i=$i same=$same isSep=$isSep type=" . get_debug_type($row));
}

// Use reflection-free: subclass Table? can't easily.
// Check Generator/TableRows by rendering with writeln debug - monkey patch via output

class CountingOutput extends BufferedOutput {
    public int $writes = 0;
    public function doWrite(string $message, bool $newline): void {
        $this->writes++;
        parent::doWrite($message, $newline);
        // can't Log easily from here in all cases
    }
}

$out = new CountingOutput();
$t = new Table($out);
$t->setHeaders(['ID']);
$t->setRows([['1']]);
$t->render();
Log::info('writes=' . $out->writes);
Log::info('content=' . var_export($out->fetch(), true));

Log::info('done');
