<?php

namespace tests\php;

/**
 * ArgvInput：array_shift($argv) 必须写回本地数组，否则命令名会停在 artisan。
 */

$a = ['artisan', 'inspire'];
$shifted = array_shift($a);
if ($shifted !== 'artisan') {
    Log::fatal('array_shift 返回值错误: '.var_export($shifted, true));
}
if (($a[0] ?? null) !== 'inspire' || count($a) !== 1) {
    Log::fatal('array_shift 未写回数组: '.var_export($a, true));
}

class ArgvOpt_Ctor
{
    public $tokens;

    public function __construct(?array $argv = null)
    {
        $argv ??= ['artisan', 'inspire'];
        array_shift($argv);
        $this->tokens = $argv;
    }
}

$i = new ArgvOpt_Ctor();
if (($i->tokens[0] ?? null) !== 'inspire') {
    Log::fatal('构造函数 ?array $argv=null 后 array_shift 失败: '.var_export($i->tokens, true));
}

$_SERVER['argv'] = ['artisan', 'serve', '--port=18080'];
class ArgvOpt_Server
{
    public $tokens;

    public function __construct(?array $argv = null)
    {
        $argv ??= $_SERVER['argv'] ?? [];
        array_shift($argv);
        $this->tokens = $argv;
    }
}
$s = new ArgvOpt_Server();
if (($s->tokens[0] ?? null) !== 'serve') {
    Log::fatal('从 $_SERVER[argv] 构造后命令名错误: '.var_export($s->tokens, true));
}

Log::info('argv_shift_command 测试通过');
