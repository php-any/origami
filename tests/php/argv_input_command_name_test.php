<?php

namespace tests\php;

require __DIR__.'/../../cli_test/vendor/autoload.php';

use Symfony\Component\Console\Input\ArgvInput;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputDefinition;

echo "=== argv_input_command_name_test.php ===\n";

echo "RAW \$_SERVER['argv']:\n";
var_dump($_SERVER['argv'] ?? null);

$input = new ArgvInput();

echo "ArgvInput firstArgument BEFORE bind:\n";
var_dump($input->getFirstArgument());

$definition = new InputDefinition([
    new InputArgument('command', InputArgument::REQUIRED, 'The command to execute'),
]);

$input->bind($definition);

echo "ArgvInput firstArgument AFTER bind:\n";
var_dump($input->getFirstArgument());

$ref = new \ReflectionClass($input);
$prop = $ref->getProperty('arguments');
$prop->setAccessible(true);

echo "ArgvInput internal arguments AFTER bind:\n";
var_dump($prop->getValue($input));

