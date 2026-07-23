<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Formatter\OutputFormatter;
use Symfony\Component\Console\Input\ArrayInput;
use Symfony\Component\Console\Output\BufferedOutput;
use Symfony\Component\Console\Style\SymfonyStyle;

// Direct format of what createBlock produces
$style = 'fg=black;bg=green';
$line = sprintf('<%s>%s</>', $style, ' [OK] Formatted success message');
Log::info('line: ' . var_export($line, true));

$f = new OutputFormatter(true);
$formatted = $f->format($line);
Log::info('formatted: ' . var_export($formatted, true));
Log::info('hex: ' . bin2hex($formatted));

// writeln array of two such lines
$buf = new BufferedOutput();
$io = new SymfonyStyle(new ArrayInput([]), $buf);
$io->writeln([$line, $line]);
Log::info('writeln: ' . var_export($buf->fetch(), true));

Log::info('done');
