<?php

/**
 * 阶段 2：Bootstrap\Console\Command 继承 illuminate/console，Origami execute() 双轨冒烟。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();

use App\Console\Commands\InspireCommand;
use Bootstrap\Console\Command as BootstrapCommand;
use Illuminate\Console\Command as IlluminateCommand;

if (!is_subclass_of(BootstrapCommand::class, IlluminateCommand::class)) {
    echo "FAIL: Bootstrap\\Console\\Command should extend Illuminate\\Console\\Command\n";
    exit(1);
}

$command = new InspireCommand();
if (!$command instanceof IlluminateCommand) {
    echo "FAIL: InspireCommand is not an Illuminate\\Console\\Command\n";
    exit(1);
}

$_SERVER['argv'] = ['inspire', '--quiet'];
$code = $command->execute();
if ($code !== 0) {
    echo "FAIL: execute() exit code $code\n";
    exit(1);
}

echo "PASS\n";
