<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Helper\Helper;
use Symfony\Component\Console\Helper\OutputWrapper;
use Symfony\Component\Console\Formatter\OutputFormatter;
use Symfony\Component\Console\Input\ArrayInput;
use Symfony\Component\Console\Output\BufferedOutput;
use Symfony\Component\Console\Style\SymfonyStyle;

function build_success_lines(OutputFormatter $formatter, int $lineLength): array {
    $messages = ['Formatted success message'];
    $type = '[OK] ';
    $style = 'fg=black;bg=green';
    $prefix = ' ';
    $padding = true;
    $escape = true;

    $indentLength = Helper::width($type);
    $prefixLength = Helper::width(Helper::removeDecoration($formatter, $prefix));
    $lineIndentation = str_repeat(' ', $indentLength);
    $lines = [];
    $outputWrapper = new OutputWrapper();
    foreach ($messages as $key => $message) {
        if ($escape) {
            $message = OutputFormatter::escape($message);
        }
        $message = str_replace("\r\n", "\n", $message);
        $lines = array_merge($lines, explode("\n", $outputWrapper->wrap($message, $lineLength - $prefixLength - $indentLength, "\n")));
    }
    $firstLineIndex = 0;
    if ($padding && $formatter->isDecorated()) {
        $firstLineIndex = 1;
        array_unshift($lines, '');
        $lines[] = '';
    }
    foreach ($lines as $i => &$line) {
        $line = ($firstLineIndex === $i ? $type.$line : $lineIndentation.$line);
        $line = $prefix.$line;
        $line .= str_repeat(' ', max($lineLength - Helper::width(Helper::removeDecoration($formatter, $line)), 0));
        $line = sprintf('<%s>%s</>', $style, $line);
    }
    return $lines;
}

$f = new OutputFormatter(true);
$lines = build_success_lines($f, 79);
Log::info('line count=' . count($lines));
foreach ($lines as $i => $line) {
    Log::info("line $i: " . var_export($line, true));
    Log::info("formatted $i: " . var_export($f->format($line), true));
}

$buf = new BufferedOutput();
$buf->setDecorated(true);
foreach ($lines as $line) {
    $buf->writeln($line);
}
Log::info('manual block out: ' . var_export($buf->fetch(), true));

Log::info('done');
