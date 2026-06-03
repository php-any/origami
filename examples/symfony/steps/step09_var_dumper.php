<?php

use Symfony\Component\VarDumper\VarDumper;
use Symfony\Component\VarDumper\Cloner\VarCloner;
use Symfony\Component\VarDumper\Dumper\CliDumper;

// step09: symfony/var-dumper

$cloner = new VarCloner();
$data = $cloner->cloneVar(['key' => 'value', 'num' => 42]);
step_check('VarCloner cloneVar', $data !== null);

$buffer = '';
$dumper = new CliDumper(function ($line) use (&$buffer) {
    $buffer .= $line;
});
$dumper->dump($data);
step_check('CliDumper dump', str_contains($buffer, 'key'), substr($buffer, 0, 120));

step_check('step09_var_dumper', true, 'symfony/var-dumper');
