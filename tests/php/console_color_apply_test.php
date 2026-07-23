<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Color;

$c = new Color('black', 'green', []);
$out = $c->apply('OK');
Log::info('color apply: ' . var_export($out, true));
Log::info('hex: ' . bin2hex($out));

use Symfony\Component\Console\Formatter\OutputFormatter;
$f = new OutputFormatter(true);
$formatted = $f->format('<fg=black;bg=green> OK </>');
Log::info('formatter: ' . var_export($formatted, true));
Log::info('hex: ' . bin2hex($formatted));

Log::info('console_color_apply_test 探测完成');
