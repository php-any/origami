<?php

namespace tests\php;

$autoload = __DIR__.'/../../cli_test/vendor/autoload.php';
if (!is_file($autoload)) {
    Log::info("skip: cli_test vendor 不存在，跳过 Symfony ArgvInput 测试");
    return;
}
require $autoload;

use Go\Test\Application;
use Symfony\Component\Console\Input\ArgvInput;

$input = new ArgvInput();
$app   = new Application();

echo "argv from _SERVER:\n";
var_dump($_SERVER['argv'] ?? null);

echo "ArgvInput tokens BEFORE bind:\n";
var_dump($input->getRawTokens(false));

// 模拟 Application::doRun 中的 bind 行为
$input->bind($app->getDefinition());

echo "ArgvInput tokens AFTER bind:\n";
var_dump($input->getRawTokens(false));

echo "ArgvInput first argument AFTER bind:\n";
var_dump($input->getFirstArgument());

