<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Helper\Helper;
use Symfony\Component\Console\Formatter\OutputFormatter;
use Symfony\Component\Console\Input\ArrayInput;
use Symfony\Component\Console\Output\BufferedOutput;
use Symfony\Component\Console\Style\SymfonyStyle;

$buf = new BufferedOutput();
$io = new SymfonyStyle(new ArrayInput([]), $buf);

// Reflect lineLength
$ref = new \ReflectionClass($io);
// try public methods
Log::info('decorated=' . ($buf->isDecorated() ? '1' : '0'));
Log::info('formatter decorated=' . ($buf->getFormatter()->isDecorated() ? '1' : '0'));

// Force decorate
$buf->setDecorated(true);
$buf->getFormatter()->setDecorated(true);

$style = 'fg=black;bg=green';
$type = '[OK] ';
$prefix = ' ';
$msg = 'Formatted success message';
$lineLength = 120; // guess

$formatter = $buf->getFormatter();
$prefixLength = Helper::width(Helper::removeDecoration($formatter, $prefix));
$indentLength = Helper::width($type);
Log::info("prefixLen=$prefixLength indentLen=$indentLength");

$line = $type . $msg;
$line = $prefix . $line;
$pad = max($lineLength - Helper::width(Helper::removeDecoration($formatter, $line)), 0);
$line .= str_repeat(' ', $pad);
$line = sprintf('<%s>%s</>', $style, $line);
Log::info('built line: ' . var_export($line, true));
Log::info('formatted: ' . var_export($formatter->format($line), true));

// Now real success with decorated buffer
$buf2 = new BufferedOutput();
$buf2->setDecorated(true);
$io2 = new SymfonyStyle(new ArrayInput([]), $buf2);
$io2->success('Formatted success message');
$out = $buf2->fetch();
Log::info('success decorated: ' . var_export($out, true));
Log::info('hex: ' . bin2hex($out));

Log::info('done');
