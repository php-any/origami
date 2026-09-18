<?php

namespace tests\php;

/**
 * PHP CLI：putenv / Dotenv 写入 $_SERVER 后仍应保留 argv。
 * Origami 若先被 putenv 建成空 $_SERVER，会挡住 argv 初始化，artisan serve 会变成 list。
 */

putenv('ORIGAMI_ARGV_PROBE=1');
if (!isset($_SERVER['argv']) || !is_array($_SERVER['argv'])) {
    Log::fatal('putenv 后 $_SERVER[argv] 丢失');
}
if (count($_SERVER['argv']) < 1) {
    Log::fatal('$_SERVER[argv] 不应为空');
}
if (!in_array('tests/php/server_argv_after_putenv_test.php', $_SERVER['argv'], true)
    && strpos(implode(' ', $_SERVER['argv']), 'server_argv_after_putenv_test.php') === false) {
    Log::fatal('$_SERVER[argv] 未包含脚本名: '.var_export($_SERVER['argv'], true));
}

Log::info('server_argv_after_putenv 测试通过');
