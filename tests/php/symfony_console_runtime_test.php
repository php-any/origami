<?php

namespace tests\php;

/**
 * Symfony Console：Table / OutputWrapper / ProgressBar 在 Origami 下的冒烟（依赖 foreach/preg 修复）。
 */

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Helper\OutputWrapper;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\ArrayInput;
use Symfony\Component\Console\Output\BufferedOutput;
use Symfony\Component\Console\Style\SymfonyStyle;

$w = new OutputWrapper();
$wrapped = $w->wrap('Formatted success message', 70, "\n");
if (!str_contains($wrapped, 'Formatted')) {
    Log::fatal('OutputWrapper::wrap broken: ' . var_export($wrapped, true));
}

$buf = new BufferedOutput();
$t = new Table($buf);
$t->setHeaders(['ID', 'Name']);
$t->setRows([['1', 'Alice']]);
$t->render();
$tableOut = $buf->fetch();
if (!str_contains($tableOut, 'Alice')) {
    Log::fatal('Table::render failed: ' . var_export($tableOut, true));
}

$buf2 = new BufferedOutput();
$buf2->setDecorated(true);
$io = new SymfonyStyle(new ArrayInput([]), $buf2);
$io->success('Formatted success message');
$success = $buf2->fetch();
if (!str_contains($success, 'Formatted success message') || str_contains($success, '\\1')) {
    Log::fatal('SymfonyStyle::success broken: ' . var_export($success, true));
}

$buf3 = new BufferedOutput();
$io3 = new SymfonyStyle(new ArrayInput([]), $buf3);
$bar = $io3->createProgressBar(5);
$bar->setFormat(' %current%/%max% [%bar%] %percent:3s%%');
$bar->start();
$bar->advance(2);
$bar->finish();
$barOut = $buf3->fetch();
if (str_contains($barOut, 'NOVERB') || str_contains($barOut, 'EXTRA')) {
    Log::fatal('ProgressBar leaked Go fmt: ' . var_export($barOut, true));
}
if (!str_contains($barOut, '/5')) {
    Log::fatal('ProgressBar missing progress: ' . var_export($barOut, true));
}

Log::info('symfony_console_runtime_test 测试通过');
