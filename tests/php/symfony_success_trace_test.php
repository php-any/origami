<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Input\ArrayInput;
use Symfony\Component\Console\Output\BufferedOutput;
use Symfony\Component\Console\Style\SymfonyStyle;
use Symfony\Component\Console\Formatter\OutputFormatter;

// Monkey: subclass to expose createBlock
class TraceStyle extends SymfonyStyle {
    public function publicCreateBlock($messages, $type, $style, $prefix, $padding, $escape) {
        // use reflection to call private createBlock
        $ref = new \ReflectionClass(\Symfony\Component\Console\Style\SymfonyStyle::class);
        $m = $ref->getMethod('createBlock');
        // may not work if setAccessible unsupported
        return 'skip';
    }
}

$buf = new BufferedOutput();
$buf->setDecorated(true);
$io = new SymfonyStyle(new ArrayInput([]), $buf);

// Inspect lineLength via dumping after constructing similar lines
$ref = new \ReflectionObject($io);
foreach ($ref->getProperties() as $p) {
    $name = $p->getName();
    if (str_contains($name, 'line') || str_contains($name, 'Length') || str_contains($name, 'buffer')) {
        Log::info('prop ' . $name);
    }
}

// Simulate padding path: isDecorated on style
Log::info('io decorated=' . ($io->isDecorated() ? '1' : '0'));

// Call info() which also uses block - simpler?
$buf3 = new BufferedOutput();
$buf3->setDecorated(true);
$io3 = new SymfonyStyle(new ArrayInput([]), $buf3);
$io3->writeln('<info>hello info</info>');
Log::info('info tag writeln: ' . var_export($buf3->fetch(), true));

$buf4 = new BufferedOutput();
$buf4->setDecorated(true);
$io4 = new SymfonyStyle(new ArrayInput([]), $buf4);
$io4->title('Hello Title');
Log::info('title: ' . var_export($buf4->fetch(), true));
Log::info('title hex: ' . bin2hex($buf4->fetch()));

Log::info('done');
