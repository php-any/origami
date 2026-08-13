<?php

namespace tests\php;

/**
 * 复现 Symfony OutputStyle::success / ProgressBar 在 Origami 下的格式化异常。
 */

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Input\ArrayInput;
use Symfony\Component\Console\Output\BufferedOutput;
use Illuminate\Console\OutputStyle;

$input = new ArrayInput([]);
$buffer = new BufferedOutput();
$style = new OutputStyle($input, $buffer);

$style->success('Formatted success message');
$successOut = $buffer->fetch();
Log::info('success out: ' . var_export($successOut, true));

if (str_contains($successOut, '\\1') || trim($successOut) === '\\1' || preg_match('/\\\\1/', $successOut)) {
    Log::info('DETECT: success output contains broken formatter marker');
}

$buffer2 = new BufferedOutput();
$style2 = new OutputStyle($input, $buffer2);
$style2->writeln('<info>hello</info>');
$tagOut = $buffer2->fetch();
Log::info('tag out: ' . var_export($tagOut, true));

$buffer3 = new BufferedOutput();
$style3 = new OutputStyle($input, $buffer3);
$bar = $style3->createProgressBar(5);
$bar->setFormat(' %current%/%max% [%bar%] %percent:3s%%');
$bar->start();
$bar->advance(2);
$bar->finish();
$barOut = $buffer3->fetch();
Log::info('bar out: ' . var_export($barOut, true));

if (str_contains($barOut, 'NOVERB') || str_contains($barOut, 'EXTRA')) {
    Log::info('DETECT: progress bar format broken (Go fmt verb leak)');
}

Log::info('symfony_output_format_test 探测完成');
