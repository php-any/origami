<?php

require __DIR__.'/../../vendor/autoload.php';

$_SERVER['argv'] = ['artisan', 'inspire'];

$input = new Symfony\Component\Console\Input\ArgvInput();
$ref = new ReflectionClass($input);
$tokens = $ref->getProperty('tokens');
$tokens->setAccessible(true);
echo 'TOKENS='.json_encode($tokens->getValue($input)).PHP_EOL;
echo 'FIRST_BEFORE='.json_encode($input->getFirstArgument()).PHP_EOL;

$def = new Symfony\Component\Console\Input\InputDefinition([
    new Symfony\Component\Console\Input\InputArgument('command', Symfony\Component\Console\Input\InputArgument::REQUIRED, 'The command to execute'),
    new Symfony\Component\Console\Input\InputOption('--help', '-h', Symfony\Component\Console\Input\InputOption::VALUE_NONE, 'help'),
    new Symfony\Component\Console\Input\InputOption('--quiet', '-q', Symfony\Component\Console\Input\InputOption::VALUE_NONE, 'quiet'),
    new Symfony\Component\Console\Input\InputOption('--verbose', '-v|vv|vvv', Symfony\Component\Console\Input\InputOption::VALUE_NONE, 'verbose'),
    new Symfony\Component\Console\Input\InputOption('--version', '-V', Symfony\Component\Console\Input\InputOption::VALUE_NONE, 'version'),
    new Symfony\Component\Console\Input\InputOption('--ansi', '', Symfony\Component\Console\Input\InputOption::VALUE_NEGATABLE, 'ansi', null),
    new Symfony\Component\Console\Input\InputOption('--no-interaction', '-n', Symfony\Component\Console\Input\InputOption::VALUE_NONE, 'no-interaction'),
    new Symfony\Component\Console\Input\InputOption('--env', null, Symfony\Component\Console\Input\InputOption::VALUE_OPTIONAL, 'env'),
    new Symfony\Component\Console\Input\InputOption('--silent', null, Symfony\Component\Console\Input\InputOption::VALUE_NONE, 'silent'),
]);

try {
    $input->bind($def);
    echo 'BIND=ok'.PHP_EOL;
} catch (Throwable $e) {
    echo 'BIND_ERR='.$e->getMessage().PHP_EOL;
}
echo 'FIRST_AFTER='.json_encode($input->getFirstArgument()).PHP_EOL;
echo 'TOKENS_AFTER='.json_encode($tokens->getValue($input)).PHP_EOL;
