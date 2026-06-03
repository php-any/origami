<?php

use Symfony\Component\Console\Application;
use App\Command\PingCommand;
use Symfony\Component\Console\Input\ArrayInput;
use Symfony\Component\Console\Output\BufferedOutput;

// step08: symfony/console

$app = new Application('Origami Symfony Test', '1.0.0');
$app->setAutoExit(false);
$app->add(new PingCommand());

$output = new BufferedOutput();
$exitCode = $app->run(new ArrayInput(['command' => 'app:ping']), $output);
$content = $output->fetch();

step_check('Console Application run', $exitCode === 0, "exit=$exitCode");
step_check('Console output pong', str_contains($content, 'pong'), "output=$content");

$app2 = new Application('Test', '1.0');
$app2->setAutoExit(false);
$names = $app2->all();
step_check('Console Application all', is_array($names));

step_check('step08_console', true, 'symfony/console');
