<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Helper\TableSeparator;
use Symfony\Component\Console\Output\BufferedOutput;

$t = new Table(new BufferedOutput());
$t->setHeaders(['ID', 'Name']);
$t->setRows([['1', 'Alice'], ['2', 'Bob']]);

// Probe via cloning render internals
$headers = [['ID', 'Name']];
$rows = [['1', 'Alice'], ['2', 'Bob']];
$divider = new TableSeparator();
$merged = array_merge($headers, [$divider], $rows);
Log::info('merged count=' . count($merged));
for ($i = 0; $i < count($merged); $i++) {
    Log::info("i=$i type=" . get_debug_type($merged[$i]) . ' is_array=' . (is_array($merged[$i]) ? '1' : '0') . ' is_sep=' . (($merged[$i] instanceof TableSeparator) ? '1' : '0'));
}

// Does setRows stick? Call getStyle or render step by step
$buf = new BufferedOutput();
$t2 = new Table($buf);
$t2->setHeaders(['A']);
$t2->addRow(['x']);
$t2->render();
Log::info('addRow out: ' . var_export($buf->fetch(), true));

Log::info('done');
