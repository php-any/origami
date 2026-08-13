<?php
namespace tests\php;

/**
 * 命名参数跳过靠前可选参数（Laravel withRouting(web: ...) 场景）
 */

function withRouting_sim(
    $using = null,
    $web = null,
    $api = null,
    $commands = null,
    $channels = null,
    $pages = null,
    $health = null
) {
    return compact('using', 'web', 'api', 'commands', 'channels', 'pages', 'health');
}

$web = '/path/to/web.php';
$commands = '/path/to/console.php';
$health = '/up';

$r = withRouting_sim(web: $web, commands: $commands, health: $health);

if ($r['using'] !== null) {
    Log::fatal('using 应为空, 实际: '.var_export($r['using'], true));
}
if ($r['web'] !== $web) {
    Log::fatal('web 错误: '.var_export($r['web'], true));
}
if ($r['commands'] !== $commands) {
    Log::fatal('commands 错误: '.var_export($r['commands'], true));
}
if ($r['health'] !== $health) {
    Log::fatal('health 错误: '.var_export($r['health'], true));
}

// 方法形式
class NamedSkip_Builder {
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

$b = new NamedSkip_Builder();
$r2 = $b->withRouting(web: $web, commands: $commands, health: $health);
if ($r2['using'] !== null || $r2['web'] !== $web) {
    Log::fatal('方法命名参数失败: '.var_export($r2, true));
}

Log::info('named_arg_skip_first 测试通过');
