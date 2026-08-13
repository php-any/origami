<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Formatter\OutputFormatter;

$msg = 'Formatted success message';
$escaped = OutputFormatter::escape($msg);
Log::info('escaped: ' . var_export($escaped, true));

$style = 'fg=black;bg=green';
$line = sprintf('<%s>%s</>', $style, ' [OK] ' . $escaped);
Log::info('tagged: ' . var_export($line, true));

$f = new OutputFormatter(true);
Log::info('format: ' . var_export($f->format($line), true));

// What about empty padded lines like createBlock padding
$emptyStyled = sprintf('<%s>%s</>', $style, ' ' . str_repeat(' ', 100));
Log::info('empty styled format: ' . var_export($f->format($emptyStyled), true));
Log::info('hex: ' . bin2hex($f->format($emptyStyled)));

Log::info('done');
