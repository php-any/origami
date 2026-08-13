<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Process\Factory;

/**
 * 真实子进程依赖 Symfony Process::getCommandLine() 等方法（Origami 子类方法解析待完善）。
 * 先验证 Factory::result() 假进程 API。
 */
$factory = new Factory();
$result = $factory->result('origami-process', '', 0);

if (!$result->successful()) {
    echo "FAIL: not successful\n";
    exit(1);
}

if (trim($result->output()) !== 'origami-process') {
    echo "FAIL: output ";
    var_export($result->output());
    echo "\n";
    exit(1);
}

echo "PASS\n";
