<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Helper\Helper;
use Symfony\Component\Console\Formatter\OutputFormatter;
use Symfony\Component\Console\Terminal;

$t = new Terminal();
Log::info('terminal width=' . $t->getWidth());

$f = new OutputFormatter(true);
$s = ' [OK] Formatted success message';
Log::info('width=' . Helper::width($s));
Log::info('removeDecoration=' . var_export(Helper::removeDecoration($f, $s), true));

$styled = $f->format('<fg=black;bg=green>'.$s.'</>');
Log::info('styled=' . var_export($styled, true));
Log::info('removeDec styled=' . var_export(Helper::removeDecoration($f, $styled), true));
Log::info('width styled=' . Helper::width(Helper::removeDecoration($f, $styled)));

Log::info('done');
