<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Output\BufferedOutput;

$buf = new BufferedOutput();
$t = new Table($buf);
$t->setHeaders(['ID', 'Name']);
$t->setRows([['1', 'Alice'], ['2', 'Bob']]);
$t->render();
$out = $buf->fetch();
Log::info("OUT:\n$out");
Log::info('hex ends: ' . bin2hex(substr($out, -40)));

if (!str_contains($out, 'Alice') || !str_contains($out, 'Bob')) {
    Log::fatal('missing rows');
}
Log::info('symfony_table_only_test 测试通过');
