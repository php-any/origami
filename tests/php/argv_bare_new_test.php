<?php

namespace tests\php;

/**
 * artisan: $app->handleCommand(new ArgvInput) 无括号。
 */

$_SERVER['argv'] = ['artisan', 'inspire'];

class ArgvBare_New
{
    public $tokens;

    public function __construct(?array $argv = null)
    {
        $argv ??= $_SERVER['argv'] ?? [];
        array_shift($argv);
        $this->tokens = $argv;
    }
}

function argv_bare_handle($input)
{
    return $input->tokens[0] ?? null;
}

$withParen = argv_bare_handle(new ArgvBare_New());
$noParen = argv_bare_handle(new ArgvBare_New);

if ($withParen !== 'inspire') {
    Log::fatal('new ArgvBare_New() 失败: '.var_export($withParen, true));
}
if ($noParen !== 'inspire') {
    Log::fatal('new ArgvBare_New 无括号失败: '.var_export($noParen, true));
}

Log::info('argv_bare_new 测试通过');
