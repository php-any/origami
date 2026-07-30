<?php

namespace tests\php;

/**
 * 对象方法命名参数可跳过靠前可选参数（Laravel withRouting(web: ...)）
 */

class NamedMethod_Builder
{
    public function withRouting(
        $using = null,
        $web = null,
        $api = null,
        $commands = null,
        $health = null
    ) {
        return compact('using', 'web', 'api', 'commands', 'health');
    }
}

$b = new NamedMethod_Builder();
$web = '/path/to/web.php';
$commands = '/path/to/console.php';
$health = '/up';

$r = $b->withRouting(web: $web, commands: $commands, health: $health);

if ($r['using'] !== null) {
    Log::fatal('using 应为空, 实际: '.var_export($r['using'], true));
}
if ($r['web'] !== $web) {
    Log::fatal('web 错误: '.var_export($r['web'], true));
}
if ($r['api'] !== null) {
    Log::fatal('api 应为空, 实际: '.var_export($r['api'], true));
}
if ($r['commands'] !== $commands) {
    Log::fatal('commands 错误: '.var_export($r['commands'], true));
}
if ($r['health'] !== $health) {
    Log::fatal('health 错误: '.var_export($r['health'], true));
}

// method_exists 对不存在的类名应返回 false，不抛错
if (method_exists('/not/a/class.php', '__invoke') !== false) {
    Log::fatal('method_exists 未知类应返回 false');
}

Log::info('named_method_args_skip_first 测试通过');
