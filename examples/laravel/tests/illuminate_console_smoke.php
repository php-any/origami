<?php

require dirname(__DIR__) . '/vendor/autoload.php';
require __DIR__ . '/support/ConsoleSmokeHelloCommand.php';

use Symfony\Component\Console\Input\ArrayInput;
use Symfony\Component\Console\Output\BufferedOutput;

$command = new ConsoleSmokeHelloCommand();
$input = new ArrayInput(['name' => 'laravel']);
$output = new BufferedOutput();

$code = $command->run($input, $output);
if ($code !== 0) {
    echo "FAIL: exit code $code\n";
    exit(1);
}

$text = trim($output->fetch());
if ($text !== 'hello:laravel') {
    echo "FAIL: output ";
    var_export($text);
    echo "\n";
    exit(1);
}

echo "PASS\n";
