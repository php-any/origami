<?php
namespace tests\php;
require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Terminal;
use Symfony\Component\Console\Helper\ProgressBar;
use Symfony\Component\Console\Output\ConsoleOutput;

$t = new Terminal();
for ($i = 0; $i < 5; $i++) {
    try {
        $w = $t->getWidth();
        Log::info("iter $i width=$w");
    } catch (\Throwable $e) {
        Log::fatal("iter $i: " . $e->getMessage());
    }
}

$out = new ConsoleOutput();
$bar = new ProgressBar($out, 20);
$bar->start();
for ($i = 0; $i < 10; $i++) {
    try {
        $bar->advance(2);
        Log::info("advance $i ok");
    } catch (\Throwable $e) {
        Log::fatal("advance $i: " . $e->getMessage());
    }
}
$bar->finish();
Log::info('terminal_width_loop_test 测试通过');
